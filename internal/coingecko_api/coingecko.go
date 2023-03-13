package coingecko

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Price struct {
	Usd float64 `json:"usd"`
	Eur float64 `json:"eur"`
}

type MarketData struct {
	CurrentPrice Price `json:"current_price"`
}

type CoinGeckoResponse struct {
	MarketData MarketData `json:"market_data"`
}

func QueryCoingeckoApi(coin string, date string,) Price{
	// Set up the request URL with the date and Bitcoin symbol\
	//date := "01-01-2023" // Change this to the date you want
	//symbol := "bitcoin"
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/history?date=%s", coin, date)

	// Send the HTTP GET request to the Coingecko API
	//fmt.Printf("Quering: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Parse the JSON response
	var cgResponse CoinGeckoResponse
	err = json.NewDecoder(resp.Body).Decode(&cgResponse)
	if err != nil {
		panic(err)
	}

	// Print the price of Bitcoin on the given day
	// fmt.Printf("Bitcoin price on %s: $%.2f USD\n", date, cgResponse.MarketData.CurrentPrice.Usd)
	// fmt.Printf("Bitcoin price on %s: $%.2f EUR\n", date, cgResponse.MarketData.CurrentPrice.Eur)

	return cgResponse.MarketData.CurrentPrice
}