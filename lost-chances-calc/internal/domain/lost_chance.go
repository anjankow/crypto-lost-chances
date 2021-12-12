package domain

type LostChance struct {
	FiatName      Fiat           `json:"fiat"`
	CryptocurName Cryptocurrency `json:"cryptocurrency"`

	InputFiatAmount  float64 `json:"inputFiatAmount"`
	OutputFiatAmount float64 `json:"outputFiatAmount"`
}
