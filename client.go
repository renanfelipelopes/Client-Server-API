package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"http://localhost:8080/cotacao",
		nil,
	)
	if err != nil {
		fmt.Println("Erro ao criar requisição:", err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("Timeout ao chamar servidor")
			return
		}
		fmt.Println("Erro ao chamar servidor:", err)
		return
	}
	defer res.Body.Close()

	var cotacao Cotacao

	err = json.NewDecoder(res.Body).Decode(&cotacao)
	if err != nil {
		fmt.Println("Erro ao decodificar JSON:", err)
		return
	}

	fmt.Println("Valor do dólar:", cotacao.Bid)
}
