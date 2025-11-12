# Build stage
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Copia os arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Build do binário
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o billing ./cmd/billing

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Copia o binário do stage anterior
COPY --from=builder /app/billing .

EXPOSE 8082

CMD ["./billing"]