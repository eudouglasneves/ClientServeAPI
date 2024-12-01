package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/go-resty/resty/v2"
)

// Estrutura para mapear a resposta da API de cotação
type Cotacao struct {
	USD struct {
		Bid string `json:"bid"`
	} `json:"USD"`
}

// Função para consultar a cotação do dólar
func obterCotacao(ctx context.Context) (string, error) {
	client := resty.New()
	client.SetTimeout(200 * time.Millisecond)

	var cotacao Cotacao
	resp, err := client.R().SetContext(ctx).SetResult(&cotacao).Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", fmt.Errorf("erro ao obter cotação")
	}

	return cotacao.USD.Bid, nil
}

// Função para salvar a cotação no banco de dados
func salvarCotacaoNoDB(ctx context.Context, db *sql.DB, cotacao string) error {
	// Definir o timeout máximo de 10ms
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	query := "INSERT INTO cotacoes (valor) VALUES (?)"
	_, err := db.ExecContext(ctx, query, cotacao)
	return err
}

// Função para lidar com a requisição de cotação
func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	// Criar um contexto para o servidor
	ctx := r.Context()

	// Obter cotação
	cotacao, err := obterCotacao(ctx)
	if err != nil {
		http.Error(w, "Erro ao obter cotação", http.StatusInternalServerError)
		return
	}

	// Conectar ao banco de dados SQLite
	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Salvar cotação no banco
	err = salvarCotacaoNoDB(ctx, db, cotacao)
	if err != nil {
		http.Error(w, "Erro ao salvar cotação no banco", http.StatusInternalServerError)
		return
	}

	// Retornar cotação ao cliente
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"bid": cotacao})
}

func main() {
	// Criar banco de dados se não existir
	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS cotacoes (id INTEGER PRIMARY KEY AUTOINCREMENT, valor TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	// Criar servidor HTTP
	http.HandleFunc("/cotacao", cotacaoHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
