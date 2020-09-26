package main

import (
	"github.com/binance-exchange/go-binance"
	"strings"
	"time"
)

func getPrice(coin string, b binance.Binance) float64{
	println("Getting price for: " + coin)
	coinPrice, coinPriceError := b.Ticker24(binance.TickerRequest{Symbol: coin})
	if coinPriceError != nil{
		println(coin)
		panic(coinPriceError)
	}
	return coinPrice.AskPrice
}

func Amount(b binance.Binance, key string) float64 {
	accountInfo, accountInfoErr := b.Account(binance.AccountRequest{
		RecvWindow: 5 * time.Second,
		Timestamp:  time.Now(),
	})

	for _, coin := range accountInfo.Balances{
		if coin.Asset == key{
			return coin.Free
		}
	}
	if accountInfoErr != nil{
		panic(accountInfoErr)
	}
	return 0
}
//Gets all Current Coins on Spot Account
//wether locked or free
func getBalance(b binance.Binance) []string{
	balanceRes, balanceErr := b.Account(binance.AccountRequest{
		RecvWindow: 5 * time.Second,
		Timestamp:  time.Now(),
	})
	if balanceErr != nil{
		panic(balanceErr)
	}
	var balanceRet[] string
	for _, el := range balanceRes.Balances{
		if el.Free > 0 || el.Locked > 0{
			balanceRet = append(balanceRet, el.Asset + " Free: " + Float64ToString(el.Free) + " Locked: " + Float64ToString(el.Locked))
		}
	}
	return balanceRet
}

func getCashBalance(b binance.Binance) float64{
	balanceRes, balanceErr := b.Account(binance.AccountRequest{
		RecvWindow: 5 * time.Second,
		Timestamp:  time.Now(),
	})
	if balanceErr != nil{
		panic(balanceErr)
	}
	var cashBalance float64
	for _, el := range balanceRes.Balances{
		if !(el.Free > 0 || el.Locked > 0){
			continue
		}
		if !strings.Contains(el.Asset, "BTC") && !strings.Contains(el.Asset, "USDT"){
			var price = getPrice(el.Asset + "USDT", b)
			println("Price for " + el.Asset + " => " + Float64ToString(price))
			cashBalance += el.Free * price
		}
	}
	var btcPrice = getPrice("BTCUSDT", b)
	return(btcPrice)
}

//Return wether a Coin exists in the Binance API or not
//Coin Parameter without specification e.g. USDT/BTC -> UDST
func coinExists(coin string, b binance.Binance) bool{
	coinRequest, emptyError := b.Account(binance.AccountRequest{
		RecvWindow: 5 * time.Second,
		Timestamp:  time.Now(),
	})

	if emptyError != nil{
		panic(emptyError)
	}
	for _, coinR := range coinRequest.Balances{
		if coinR.Asset == coin{
			return true
		}
	}
	return false
}