package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lmtani/learning-client-server-api/internal/entities"
)

const CotacaoRoute = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

func main() {
	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
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
