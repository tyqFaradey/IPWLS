from fastapi import FastAPI, HTTPException, Query, Depends
from sqlalchemy.orm import Session

from .dataBase import engine, get_connection, check_connection
from . import models, queries


models.Base.metadata.create_all(bind=engine)


app = FastAPI(title="IP Check API")


@app.get("/health")
def health():
    if not check_connection():
        raise HTTPException(status_code=500, detail="Database connection error")
    return {"status": "ok"}

@app.get("/check/")
def check_ip(ip: str = Query(...), connection: Session = Depends(get_connection)):
    record = queries.get_ip_by_value(connection, ip)

    return {"ip": ip, "allowed": bool(record)}

@app.post("/add/")
def add_ip(ip: str = Query(...), connection: Session = Depends(get_connection)):
    if queries.get_ip_by_value(connection, ip): raise HTTPException(status_code=400, detail="IP already exists") # existing check
    added_ip_record = queries.create_ip(connection, ip)

    return {"ip": added_ip_record.ip, "time": added_ip_record.added_at}

@app.post("/delete/")
def remove_ip(ip: str = Query(...), connection: Session = Depends(get_connection)):
    removed_ip_record = queries.delete_ip(connection, ip)
    if not removed_ip_record: raise HTTPException(status_code=400, detail="IP does not exist")
    
    return {"ip": removed_ip_record.ip}
