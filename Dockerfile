### Dockerfile ###
# Usando uma imagem oficial do Go como base
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Instala dependências do sistema
RUN apk add --no-cache gcc musl-dev postgresql-client git

# Copia os arquivos para o container
COPY go.mod go.sum ./
RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4

COPY . .

# Compila a aplicação
RUN go build -o main ./cmd/server

# Imagem final
FROM alpine:latest
WORKDIR /root/

# Instala dependências
RUN apk add --no-cache ca-certificates

# Copia o binário da aplicação
COPY --from=builder /app/main .

# Define a porta que o container expõe
EXPOSE 8080

# Comando de execução da aplicação
CMD ["/root/main"]
