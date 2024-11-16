package entities

type Cotacao struct {
	Bid string `json:"bid"`
}

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
