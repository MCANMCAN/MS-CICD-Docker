from fastapi import FastAPI, HTTPException, Depends
from sqlalchemy import create_engine, Column, Integer, String
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, Session
from dotenv import load_dotenv
import os

# Load environment variables from .env file
load_dotenv()

app = FastAPI()

# Database configuration
DATABASE_URL = os.getenv("DATABASE_URL")
engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base = declarative_base()

# Define Task table schema
class Task(Base):
    __tablename__ = 'tasks'

    id = Column(Integer, primary_key=True, index=True)
    task_description = Column(String)
    task_status = Column(String, default='pending')

# Create tables
Base.metadata.create_all(engine)

# Dependency to get database session
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

# Create a new task
@app.post('/task/')
def create_task(task_data: dict, db: Session = Depends(get_db)):
    task_description = task_data.get('task_description')
    if not task_description:
        raise HTTPException(status_code=422, detail={"detail": [{"loc": ["body", "task_description"], "msg": "field required", "type": "value_error.missing"}]})
    new_task = Task(task_description=task_description)
    db.add(new_task)
    db.commit()
    db.refresh(new_task)
    return new_task

# Example of running the FastAPI application
if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)