package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
	"os"
	"bytes"
)

const timeout = 5 * time.Second

var (
	listenAddr   = os.Getenv("PROXY_LISTEN_ADDR")
	targetServer = os.Getenv("PROXY_TARGET_ADDR")
	apiURL       = os.Getenv("API_URL")
)

type APIResponse struct {
	IP      string `json:"ip"`
	Allowed bool   `json:"allowed"`
}

func checkWhitelist(ip string) bool {
	client := http.Client{Timeout: timeout}
	resp, err := client.Get(fmt.Sprintf("%s?ip=%s", apiURL, ip))
	if err != nil {
		fmt.Println("[API ERROR]", err)
		return false
	}
	defer resp.Body.Close()

	var result APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("[JSON ERROR]", err)
		return false
	}

	fmt.Printf("[API] IP=%s allowed=%v\n", result.IP, result.Allowed)
	return result.Allowed
}

func sendDisconnectMsg(conn net.Conn, msg string) {
	conn.Write([]byte(msg + "\n"))
	conn.Close()
}

////

func writeVarInt(buf *bytes.Buffer, value int32) {
	for {
		temp := byte(value & 0b01111111)
		value >>= 7
		if value != 0 {
			temp |= 0b10000000
		}
		buf.WriteByte(temp)
		if value == 0 {
			break
		}
	}
}

func writeString(buf *bytes.Buffer, s string) {
	writeVarInt(buf, int32(len(s)))
	buf.WriteString(s)
}

func sendMCDisconnect(conn net.Conn, reason string) {
	msg := map[string]string{"text": reason}
	jsonMsg, _ := json.Marshal(msg)

	var buf bytes.Buffer
	packetID := byte(0x0B)

	writeVarInt(&buf, int32(packetID))
	writeString(&buf, string(jsonMsg))

	packetData := buf.Bytes()
	var packetBuf bytes.Buffer
	writeVarInt(&packetBuf, int32(len(packetData)))
	packetBuf.Write(packetData)

	conn.Write(packetBuf.Bytes())
	conn.Close()
}

////
 

func handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	clientAddr := clientConn.RemoteAddr().(*net.TCPAddr)
	clientIP := clientAddr.IP.String()
	fmt.Printf("[CONNECT] %s\n", clientIP)

	// Проверяем IP через API
	if !checkWhitelist(clientIP) {
		fmt.Printf("[BLOCKED] %s — not in whitelist\n", clientIP)
		sendMCDisconnect(clientConn, "You are not whitelisted.")
		return
	}

	serverConn, err := net.Dial("tcp", targetServer)
	if err != nil {
		fmt.Println("[ERROR] Cannot connect to target server:", err)
		sendMCDisconnect(clientConn, "Internal proxy error.")
		return
	}
	defer serverConn.Close()

	fmt.Printf("[ALLOWED] %s connected to %s\n", clientIP, targetServer)

	go io.Copy(serverConn, clientConn)
	io.Copy(clientConn, serverConn)
}

func main() {
	fmt.Println("Minecraft Proxy started on", listenAddr)

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[ERROR] Accept failed:", err)
			continue
		}
		go handleConnection(conn)
	}
}
