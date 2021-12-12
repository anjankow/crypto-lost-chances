package domain

type Cryptocurrency string

const (
	Bitcoin  Cryptocurrency = "BTC"
	Ada      Cryptocurrency = "ADA"
	Ethereum Cryptocurrency = "ETH"
	Doge     Cryptocurrency = "DOGE"
	Litecoin Cryptocurrency = "LTC"
	Stellar  Cryptocurrency = "XLM"
	Xrp      Cryptocurrency = "XRP"
)

var Cryptocurrencies []Cryptocurrency = []Cryptocurrency{Bitcoin, Ethereum, Ada, Doge, Litecoin, Stellar, Xrp}
