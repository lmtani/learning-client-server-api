package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/lmtani/learning-client-server-api/internal/database"
	"github.com/lmtani/learning-client-server-api/internal/entities"
)

const (
	CotacaoRoute      = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	ApiCotacaoTimeout = 200 * time.Millisecond
	DatabaseTimeout   = 10 * time.Millisecond
)

func main() {

	err := createDatabase()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), ApiCotacaoTimeout)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, CotacaoRoute, nil)
		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {

			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				http.Error(w, "Request timeout", http.StatusRequestTimeout)
				return
			}

			http.Error(w, "Error requesting USD-BRL quote", http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading response body", http.StatusInternalServerError)
			return
		}

		var exchange entities.CurrencyExchange
		if err := json.Unmarshal(data, &exchange); err != nil {
			http.Error(w, "Error unmarshalling JSON", http.StatusInternalServerError)
			return
		}

		db, err := sql.Open("sqlite3", "./quotes.db")
		if err != nil {
			http.Error(w, "Error opening database", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		if err := database.AddQuote(db, &exchange.UsdBrl, DatabaseTimeout); err != nil {
			http.Error(w, fmt.Sprintf("Database error: %s", err), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(entities.Cotacao{Bid: exchange.UsdBrl.Bid}); err != nil {
			http.Error(w, "Error encoding exchange data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	})

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func createDatabase() error {
	db, err := sql.Open("sqlite3", "./quotes.db")
	if err != nil {
		return err
	}
	defer db.Close()

	if err := database.CreateTable(db); err != nil {
		return err
	}
	return nil
}
