
FROM python:3.8-slim-buster

WORKDIR /home/mouyuan/go/src/github.com/mouyuan4598/py

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY . .
 
CMD ["python", "./hello.py"]