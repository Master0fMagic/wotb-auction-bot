package main

import (
	"context"
	"github.com/Master0fMagic/wotb-auction-bot/bot"
	"github.com/Master0fMagic/wotb-auction-bot/provider"
	"github.com/Master0fMagic/wotb-auction-bot/storage"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const botToken = "6872455980:AAHpJ9IYerYYSM8FkRUqfToFtdR-YNNw05Y"

const apiUrl = "https://eu.wotblitz.com/en/api/events/items/auction/?page_size=80&type[]=vehicle&saleable=true"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	errorGroup, ctx := errgroup.WithContext(ctx)

	strg := storage.NewRuntimeMonitoringStorage()
	provdr := provider.NewHttpActionProvider(apiUrl) // todo move to config

	tgBot, err := bot.New(botToken) // todo move to cfg
	if err != nil {
		log.Fatal(err)
	}

	initBot(strg, provdr, tgBot)
	errorGroup.Go(func() error {
		return tgBot.Run(ctx)
	})

	//if err := gocron.Every(1).Minute().Do(func() {
	//	checkEntities()
	//}); err != nil {
	//	log.Fatal(err)
	//}

	errorGroup.Wait()
}

func initBot(monitoringStorage storage.MonitoringStorage, dataProvider provider.AuctionDataProvider, tgBot *bot.Bot) {
	tgBot.AddHandler(
		bot.GetCommandNamePredicate("start"),
		bot.GetStaticTextResponseHandler("welcome to wotb auction bot"),
	)
	tgBot.AddHandler(
		bot.GetCommandNamePredicate("data"),
		bot.GetDataCommandHandler(dataProvider),
	)
	tgBot.AddHandler(
		bot.GetCommandNamePredicate("data_short"),
		bot.GetDataShortCommandHandler(dataProvider),
	)

}
