package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cambio struct {
	Bid string `json:"bid"`
}

func main() {
	criaBanco()

	http.HandleFunc("/cotacao", cotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	cotacao, err := buscaCotacao()
	if err != nil {
		fmt.Fprintln(w, err)
	}

	err = saveCotacao(cotacao)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	json.NewEncoder(w).Encode(cotacao)
}

func criaBanco() {
	db, err := sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		panic(err)
	}

	cotacaoTable := `CREATE TABLE IF NOT EXISTS cambio(id INTEGER PRIMARY KEY AUTOINCREMENT, bid TEXT NOT NULL);`
	_, err = db.Exec(cotacaoTable)
	if err != nil {
		panic(err)
	}

	db.Close()
}

// https://economia.awesomeapi.com.br/json/last/USD-BRL
func buscaCotacao() (*Cambio, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var payload map[string]Cambio
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	cambio := payload["USDBRL"]

	return &cambio, nil
}

func saveCotacao(cambio *Cambio) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10)
	defer cancel()

	db, err := sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		return err
	}

	defer db.Close()

	stmt, err := db.PrepareContext(ctx, `INSERT INTO cambio(bid) VALUES(?);`)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, cambio.Bid)
	if err != nil {
		return err
	}

	return nil
}
