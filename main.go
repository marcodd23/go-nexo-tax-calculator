package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	coingecko "github.com/marcodd23/internal/coingecko_api"
)

const (
	//Csv columns:
	Transaction = iota
	Type
	InputCurrency
	InputAmount
	OutputCurrency
	OutputAmount
	USDEquivalent
	Details
	DateTime
)

const exchangeTransaction = "Exchange"
const transferToProWalletTransaction = "TransferToProWallet"
const transferFromProWalletTransaction = "TransferFromProWallet"
const exchangeDepositedOn = "ExchangeDepositedOn"
const unlockingTermDeposit = "UnlockingTermDeposit"
const lockingTermDeposit = "LockingTermDeposit"
const dateLayoutCsv = "2006-01-02 15:04:05"
const dateLayoutInput = "02-01-2006"

var coinsWallet map[string]float64

var coinsNameMap = map[string]string{
	"BTC":   "bitcoin",
	"ETH":   "ethereum",
	"NEXO":  "nexo",
	"USDC":  "usd-coin",
	"AVAX":  "avalanche-2",
	"UST":   "terrausd-wormhole",
	"LINK":  "chainlink",
	"MATIC": "matic-network",
	"SOL":   "solana",
	"ETHW":  "ethereum-pow-iou",
	"DOT":   "polkadot",
	"LUNA":  "terra-luna",
	"LUNA2": "terra-luna-2",
}

type recordData struct {
	Type      string
	InCoin    string
	InAmount  float64
	OutCoin   string
	OutAmount float64
	Date      time.Time
}

func main() {
	var filePath string
	var targetDate string

	flag.StringVar(&filePath, "path", "nexo_transactions.csv", "Path to the transaction file")
	flag.StringVar(&targetDate, "targetDate", "01-01-2023", "Date when to calculate the balance: Ex. 01-01-2023")

	// Parse the flags from the command line
	flag.Parse()

	run(filePath, targetDate)
}

func run(filePath string, targetDateStr string) {
	coinsWallet = make(map[string]float64)

	var totalInvestmentEuro float64

	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the first row and discard it
	_, err = reader.Read()
	if err != nil {
		panic(err)
	}

	targetDate, err := time.Parse(dateLayoutInput, targetDateStr)
	if err != nil {
		log.Panicf("error parsing target date paramenter: %v", err)
	}

	fmt.Printf("Target Date: %v\n", targetDate)

	// Read the CSV records one by one
	for {

		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		record := mapRecord(row)

		//targetDate := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)

		if record.Date.After(targetDate) || record.Type == exchangeDepositedOn || record.Type == unlockingTermDeposit || record.Type == lockingTermDeposit {
			continue
		} else {
			if record.Type == exchangeTransaction {
				addAmount(coinsWallet, record.InCoin, record.InAmount)
				addAmount(coinsWallet, record.OutCoin, record.OutAmount)
			} else if record.Type == transferToProWalletTransaction || record.Type == transferFromProWalletTransaction {
				addAmount(coinsWallet, record.InCoin, record.InAmount)
			} else {
				addAmount(coinsWallet, record.OutCoin, record.OutAmount)
			}
		}
	}

	for coin, value := range coinsWallet {
		if coin != "EURX" {
			coinSymbol := coinsNameMap[coin]
			price := coingecko.QueryCoingeckoApi(coinSymbol, targetDateStr)
			valueInEur := value * price.Eur
			totalInvestmentEuro += valueInEur
			fmt.Printf("#%s: %f, Price: %f ### Value Euro: %f Eur\n", coin, value, valueInEur, value*price.Eur)
		} else {
			totalInvestmentEuro += value
			fmt.Printf("#%s: %f Eur\n", coin, value)
		}
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Printf("Total investment in Euro the %s: %f Euro\n", targetDateStr, totalInvestmentEuro)
}

func mapRecord(row []string) recordData {
	inAmount, err := strconv.ParseFloat(row[InputAmount], 64)
	if err != nil {
		fmt.Println("Error converting string to float64:", err)
		inAmount = 0.0
	}
	outAmount, err := strconv.ParseFloat(row[OutputAmount], 64)
	if err != nil {
		fmt.Println("Error converting string to float64:", err)
		outAmount = 0.0
	}

	date, err := time.Parse(dateLayoutCsv, row[DateTime])
	if err != nil {
		log.Panicf("Error Parsing Date: %v", err)
	}

	return recordData{
		Type:      row[Type],
		InCoin:    row[InputCurrency],
		InAmount:  inAmount,
		OutCoin:   row[OutputCurrency],
		OutAmount: outAmount,
		Date:      date,
	}
}

func addAmount(coinsWallet map[string]float64, coin string, amount float64) {
	if _, ok := coinsWallet[coin]; ok {
		coinsWallet[coin] = coinsWallet[coin] + amount
	} else {
		coinsWallet[coin] = amount
	}
}
