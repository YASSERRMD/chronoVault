FROM golang:1.21-alpine

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

EXPOSE 8080

CMD ["go", "run", "cmd/server/main.go"]
