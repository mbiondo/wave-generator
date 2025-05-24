FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o wave-generator main.go

# Instala redis-tools para debug opcional (no el server)
RUN apt-get update && apt-get install -y redis-tools && rm -rf /var/lib/apt/lists/*

ENV PORT=1155
ENV REDIS_ADDR=redis:6379
EXPOSE ${PORT}

CMD ["/app/wave-generator"]