package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API chamada com sucesso"))
}
