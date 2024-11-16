package main

import (
	"encoding/json"
	"fmt"
	"github.com/lmtani/learning-client-server-api/internal/entities"
	"io"
	"net/http"
	"os"
)

const ServerResourceRoute = "http://localhost:8080/cotacao"

func main() {
	resp, err := http.Get(ServerResourceRoute)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
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
