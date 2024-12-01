package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	// Criar contexto com timeout de 300ms
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Realizar requisição HTTP para o servidor
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatalf("Erro ao criar requisição: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Erro na resposta do servidor: %v", resp.Status)
	}

	// Processar a resposta JSON
	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalf("Erro ao decodificar a resposta: %v", err)
	}

	// Salvar cotação no arquivo
	cotacao := result["bid"]
	err = ioutil.WriteFile("cotacao.txt", []byte(fmt.Sprintf("Dólar: %s", cotacao)), 0644)
	if err != nil {
		log.Fatalf("Erro ao salvar cotação em arquivo: %v", err)
	}

	fmt.Printf("Cotação do Dólar: %s\n", cotacao)
}
