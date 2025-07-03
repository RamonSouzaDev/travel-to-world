FROM golang:1.21-alpine AS builder

WORKDIR /app

# Instalar dependências do sistema
RUN apk add --no-cache git ca-certificates gcc musl-dev

# Copiar arquivos go
COPY go.mod ./

# Baixar dependências primeiro
RUN go mod download

# Copiar todo o código DEPOIS de baixar dependências
COPY . .

# Executar go mod tidy DEPOIS de copiar o código
RUN go mod tidy

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Estágio final - imagem minimalista
FROM alpine:latest

# Instalar certificados SSL
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binário
COPY --from=builder /app/main .

# Expor porta
EXPOSE 8080

# Comando para iniciar
CMD ["./main"]
