package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lmtani/learning-client-server-api/internal/entities"
)

const CotacaoRoute = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

// CurrencyExchange represents the structure of the exchange rate data
type CurrencyExchange struct {
	UsdBrl UsdBrl `json:"USDBRL"`
}

// UsdBrl represents the details of the USD/BRL exchange rate
type UsdBrl struct {
	Code       string `json:"code"`
	CodeIn     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func main() {
	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get(CotacaoRoute)
		if err != nil {
			http.Error(w, "Error getting USD-BRL quote", http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Error reading response body", http.StatusInternalServerError)
			return
		}

		var exchange CurrencyExchange
		if err := json.Unmarshal(data, &exchange); err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		if err := json.NewEncoder(w).Encode(entities.Cotacao{Bid: exchange.UsdBrl.Bid}); err != nil {
			http.Error(w, "Error encoding exchange data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		return
	})

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
