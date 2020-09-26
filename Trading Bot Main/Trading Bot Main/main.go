package main

import (
	"context"
	"github.com/binance-exchange/go-binance"
	"github.com/go-kit/kit/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"strconv"
	"strings"
)

func Float64ToString(inputNum float64) string {
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}

func binanceMain(binanceApiKey string, binanceApiSecret string) binance.Binance{
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	hmacSigner := &binance.HmacSigner{
		Key: []byte(binanceApiSecret),
	}
	ctx, _ := context.WithCancel(context.Background())

	binanceService := binance.NewAPIService(
		"https://www.binance.com",
		binanceApiKey,
		hmacSigner,
		logger,
		ctx,
	)
	b := binance.NewBinance(binanceService)
	println(" => Finished Binance Authorization")
	return b
}

func tgMain(b binance.Binance, telegramApiKey string){
	bot, tgbotErr := tgbotapi.NewBotAPI(telegramApiKey)
	if tgbotErr != nil{
		panic(tgbotErr)
	}

	println("Authorized on account " + bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if strings.HasPrefix(update.Message.Text, "/help") {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Trading Bot Help: \n/price <Coin>\n/balance\n/cashbalance")
			bot.Send(msg)
		}

		if strings.HasPrefix(update.Message.Text, "/price "){
			coin := strings.Split(update.Message.Text, " ")[1]
			txtAnswer := "Current Price for " + coin + ": " + Float64ToString(getPrice(coin, b))
			answerPrice := tgbotapi.NewMessage(update.Message.Chat.ID, txtAnswer)
			answerPrice.ReplyToMessageID = update.Message.MessageID
			bot.Send(answerPrice)
		}
		if strings.HasPrefix(update.Message.Text, "/balance"){
			var coinBalance = getBalance(b)
			for _, el := range coinBalance{
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, el)
				bot.Send(msg)
			}
		}
		if strings.HasPrefix(update.Message.Text, "/cashbalance"){
			msg:= tgbotapi.NewMessage(update.Message.Chat.ID, "Current Balance: " + Float64ToString(getCashBalance(b)))
			bot.Send(msg)
		}
	}
}

func main() {
	if !pathExists("config.lev"){
		createFiles()
		println("Please update your Data in config.lev")
		os.Exit(0)
	}

	var (
		BapiKey = getData(0, "config.lev")
		bsecret = getData(1, "config.lev")
		tgKey = getData(2, "config.lev")
		)

	if BapiKey == "123456789" || bsecret == "987654321" || tgKey == "abcdefghiklmnopqrstuvwxyz" {
		println("Please update your Data in config.lev")
		os.Exit(0)
	}


	tgMain(binanceMain(BapiKey, bsecret), tgKey)
}