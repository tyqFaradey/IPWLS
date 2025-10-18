from sqlalchemy.orm import Session

from . import models

def _commit(connection: Session, value) -> None:
    
    connection.commit()
    connection.refresh(value)


def get_ip_by_value(connection: Session, ip_value: str) -> models.IPAddress | None:
    return connection.query(models.IPAddress).filter(models.IPAddress.ip == ip_value).first()


def create_ip(connection: Session, ip_value: str) -> models.IPAddress:
    ip_record = models.IPAddress(ip=ip_value)

    connection.add(ip_record)
    _commit(connection, ip_record)

    return ip_record

def delete_ip(connection: Session, ip_value: str) -> models.IPAddress | None:
    ip_record = get_ip_by_value(connection, ip_value)
    if ip_record:
        connection.delete(ip_record)
        connection.commit()
    
    return ip_record
