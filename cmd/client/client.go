package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/lmtani/learning-client-server-api/internal/entities"
)

const (
	ServerResourceRoute = "http://localhost:8080/cotacao"
	ServerTimeout       = 300 * time.Millisecond
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), ServerTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ServerResourceRoute, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			fmt.Println("Request to server timed out")
			return
		}
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(string(body))
		return
	}

	var c entities.Cotacao
	err = json.Unmarshal(body, &c)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "Dólar: %s\n", c.Bid)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
