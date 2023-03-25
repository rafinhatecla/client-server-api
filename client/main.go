package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Cambio struct {
	Bid string `json:"bid"`
}

func main() {
	cambio, err := buscaCotacao()
	if err != nil {
		panic(err)
	}

	err = saveCotacao(cambio)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Dólar: %s\n", cambio.Bid)
}

func buscaCotacao() (*Cambio, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var cambio Cambio
	if err := json.NewDecoder(resp.Body).Decode(&cambio); err != nil {
		return nil, err
	}

	return &cambio, nil
}

func saveCotacao(cambio *Cambio) error {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s\n", cambio.Bid))
	if err != nil {
		return err
	}

	return nil
}
