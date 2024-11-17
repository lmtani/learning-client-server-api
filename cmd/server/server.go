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

	"github.com/lmtani/learning-client-server-api/internal/database"
	"github.com/lmtani/learning-client-server-api/internal/entities"
	_ "github.com/mattn/go-sqlite3"
)

const (
	CotacaoRoute      = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	ApiCotacaoTimeout = 200 * time.Millisecond
	DatabaseTimeout   = 10 * time.Millisecond
	SqliteDbPath      = "./quotes.db"
	ServerPort        = ":8080"
)

func main() {

	db, err := initDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	server := &Server{Database: db}
	http.HandleFunc("/cotacao", server.HandleCotacao)

	fmt.Println("Server running on port", ServerPort)
	if err := http.ListenAndServe(ServerPort, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func initDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", SqliteDbPath)
	if err != nil {
		return nil, err
	}

	if err := database.CreateTable(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

type Server struct {
	Database *sql.DB
}

func (s *Server) HandleCotacao(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	logger := &RequestLogger{StartTime: startTime, RequestId: fmt.Sprintf("%x", startTime.UnixNano())}
	// create a short hash to identify request
	logger.Log("Request received")
	ctx, cancel := context.WithTimeout(r.Context(), ApiCotacaoTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, CotacaoRoute, nil)
	if err != nil {
		msg := fmt.Sprintf("Error creating request: %s", err)
		logger.Log(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			msg := fmt.Sprintf("Request to API timed out after %s", ApiCotacaoTimeout)
			logger.Log(msg)
			http.Error(w, msg, http.StatusRequestTimeout)
			return
		}

		http.Error(w, "Error requesting USD-BRL quote", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprintf("Error reading response body: %s", err)
		logger.Log(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	var exchange entities.CurrencyExchange
	if err := json.Unmarshal(data, &exchange); err != nil {
		msg := fmt.Sprintf("Error unmarshalling JSON: %s", err)
		logger.Log(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	if err := database.AddQuote(s.Database, &exchange.UsdBrl, DatabaseTimeout); err != nil {
		logger.Log(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(entities.Cotacao{Bid: exchange.UsdBrl.Bid}); err != nil {
		msg := fmt.Sprintf("Error encoding exchange data: %s", err)
		logger.Log(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	logger.Log("Request handled")
}

type RequestLogger struct {
	StartTime time.Time
	RequestId string
}

func (l *RequestLogger) Log(msg string) {
	ts := time.Since(l.StartTime).Round(time.Millisecond)
	fmt.Printf("[%s] %s %s\n", l.RequestId, ts, msg)
}
