# Usar uma imagem base com o Go já instalado
FROM golang:1.20-alpine

# Instalar dependências necessárias (GCC, SQLite, etc)
RUN apk update && apk add --no-cache \
    gcc \
    libc-dev \
    sqlite-dev \
    make \  
    bash

# Definir o diretório de trabalho
WORKDIR /app

# Copiar o código-fonte Go para o diretório de trabalho no contêiner
COPY . .

# Rodar o comando para instalar as dependências do Go
RUN go mod tidy

# Expor a porta 8080 para a aplicação (caso seu servidor esteja rodando nela)
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["go", "run", "server.go"]
