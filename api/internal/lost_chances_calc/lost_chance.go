package lostchancescalc

type LostChance struct {
	FiatName      string `json:"fiat"`
	CryptocurName string `json:"cryptocurrency"`

	InputFiatAmount  float64 `json:"inputFiatAmount"`
	OutputFiatAmount float64 `json:"outputFiatAmount"`
}
