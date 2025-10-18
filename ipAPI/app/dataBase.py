import pymysql
pymysql.install_as_MySQLdb() 

from sqlalchemy import Engine, create_engine, text
from sqlalchemy.orm import sessionmaker, declarative_base
from sqlalchemy.exc import OperationalError
import os, time

Base = declarative_base()
DATABASE_URL = os.getenv("DB_URL", "")

for attempt in range(10):
    try:
        engine = create_engine(DATABASE_URL)

        with engine.connect() as conn:
            conn.execute(text("SELECT 1"))

        print("Connected to database")

        break
    except OperationalError as e:
        print(f"Database not ready (attempt {attempt + 1}/10), retrying in 3s...")
        time.sleep(3)
else:
    raise Exception("Could not connect to database after 10 attempts")

_SessionLocal = sessionmaker(bind=engine, autoflush=False, autocommit=False)

def check_connection():
    try:
        with engine.connect() as conn:
            conn.execute(text("SELECT 1"))
        return True
    except OperationalError:
        return False

def get_connection():
    connection = _SessionLocal()
    try:
        yield connection
    finally:
        connection.close()
