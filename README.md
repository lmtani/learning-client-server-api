# Prática da Pós "Go Expert": Client-Server-API

## Resumo

Ao receber uma requisição, o backend:

- Consulta uma API externa (Timeout em 200ms)
- Grava em um banco local (Sqlite3) os dados recebidos (Timeout em 10ms)
- Retorna os dados para o client apenas o campo com a cotação (bid)

Já o frontend:

- Consulta o recurso do backend (Timeout em 100ms)
- Grava a cotação em um arquivo chamado "cotacao.txt"


## Para rodar

```bash
# Em uma sessão do terminal
go run cmd/server/server.go

# Em outra sessão do terminal
go run cmd/client/client.go
```

Exemplo de registros no banco de dados:

```tsv
id  code  codein  name                             high    low     varBid  pctChange  bid     ask     timestamp   createDate
--  ----  ------  -------------------------------  ------  ------  ------  ---------  ------  ------  ----------  -------------------
1   USD   BRL     Dólar Americano/Real Brasileiro  5.7959  5.7947  0.0057  0.17       5.7947  5.7965  1731715198  2024-11-15 20:59:58
2   USD   BRL     Dólar Americano/Real Brasileiro  5.7959  5.7947  0.0057  0.17       5.7947  5.7965  1731715198  2024-11-15 20:59:58
3   USD   BRL     Dólar Americano/Real Brasileiro  5.7959  5.7947  0.0057  0.17       5.7947  5.7965  1731715198  2024-11-15 20:59:58
4   USD   BRL     Dólar Americano/Real Brasileiro  5.7959  5.7947  0.0057  0.17       5.7947  5.7965  1731715198  2024-11-15 20:59:58
5   USD   BRL     Dólar Americano/Real Brasileiro  5.7959  5.7947  0.0057  0.17       5.7947  5.7965  1731715198  2024-11-15 20:59:58
```
