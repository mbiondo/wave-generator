FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o wave-generator main.go

ENV PORT=1155
EXPOSE ${PORT}

CMD ["/app/wave-generator", "--port", "${PORT}"]