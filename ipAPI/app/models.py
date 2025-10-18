from sqlalchemy import Column, Integer, String, DateTime
from sqlalchemy.sql import func

from .dataBase import Base


class IPAddress(Base):
    __tablename__ = "IPAddress"

    id = Column(Integer, primary_key=True, autoincrement=True, index=True)
    ip = Column(String(45), unique=True, nullable=False, index=True)
    added_at = Column(DateTime(timezone=True), server_default=func.now())
