# Projeto de Cotação de Dólar em Go com SQLite

Este é um projeto em Go que implementa um servidor HTTP que consulta a cotação do dólar e persiste as informações em um banco de dados SQLite. Ele é executado dentro de um contêiner Docker para facilitar a instalação e a execução.

## Requisitos

- [Docker](https://www.docker.com/get-started) deve estar instalado e funcionando corretamente no seu sistema.

## Estrutura do Projeto

- **server.go**: Código do servidor Go que consome a API de câmbio do dólar, persiste as informações no SQLite e expõe um endpoint `/cotacao`.
- **Dockerfile**: Arquivo que define a configuração do ambiente Docker.
- **go.mod** e **go.sum**: Arquivos de dependências do Go.

## Passos para Rodar o Projeto

### 1. Clonar o Repositório

Se ainda não tiver o projeto, faça o clone para o seu diretório de trabalho:

```bash
git https://github.com/eudouglasneves/ClientServeAPI
cd ClientServeAPI
docker build -t go-sqlite-app
docker run -p 8080:8080 go-sqlite-app
curl http://localhost:8080/cotacao
