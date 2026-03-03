package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type CotacaoResponse struct {
	USDBRL USDBRL `json:"USDBRL"`
}

type USDBRL struct {
	Bid string `json:"bid"`
}

func main() {
	http.HandleFunc("/cotacao", handler)

	log.Println("Servidor rodando na porta 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://economia.awesomeapi.com.br/json/last/USD-BRL",
		nil,
	)
	if err != nil {
		log.Println("Erro ao criar requisição:", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Timeout ao chamar API")
			http.Error(w, "Timeout na API externa", http.StatusGatewayTimeout)
			return
		}

		log.Println("Erro ao chamar API externa:", err)
		http.Error(w, "Erro ao chamar API externa:", http.StatusBadGateway)
		return
	}
	defer res.Body.Close()

	var cotacao CotacaoResponse
	err = json.NewDecoder(res.Body).Decode(&cotacao)
	if err != nil {
		log.Println("Erro ao decodificar JSON:", err)
		http.Error(w, "Erro ao processar resposta da API", http.StatusInternalServerError)
		return
	}

	bid := cotacao.USDBRL.Bid

	db, err := sql.Open("sqlite", "./cotacao.db")
	if err != nil {
		log.Println("Erro ao abrir banco:", err)
		return
	}
	defer db.Close()

	createTable := `
	CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		created_at DATETIME
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Println("Erro ao criar tabela:", err)
		return
	}

	ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelDB()

	insert := `INSERT INTO cotacoes (bid, created_at) VALUES (?, ?)`

	_, err = db.ExecContext(ctxDB, insert, bid, time.Now())
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Timeout ao salvar no banco")
			return
		}
		log.Println("Erro ao inserir no banco:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"bid": bid,
	})
}
