### Dockerfile ###
# Usando uma imagem oficial do Go como base
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev postgresql-client git
RUN go install github.com/cosmtrek/air@latest

# Copia os arquivos para o container
COPY go.mod go.sum ./
RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4

COPY . .

# Compila a aplicação
RUN go build -o main ./cmd/server

FROM alpine:3.20
WORKDIR /app/

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main .

USER appuser

EXPOSE 8080

CMD ["/app/main"]