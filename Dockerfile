FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o scheduler .

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/scheduler .
COPY --from=builder /app/web ./web

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

CMD ["./scheduler"]

#Пример командной строки для запуска
#docker build -t scheduler .
#docker run --name scheduler -p 7540:7540 -v "D:\Dev\final_project:/data" -e TODO_PASSWORD=12345 scheduler