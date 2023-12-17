package main

import (
	"context"
	"fmt"
	"github.com/Master0fMagic/wotb-auction-bot/bot"
	"github.com/Master0fMagic/wotb-auction-bot/dto"
	"github.com/Master0fMagic/wotb-auction-bot/provider"
	"github.com/Master0fMagic/wotb-auction-bot/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const botToken = ""

const apiUrl = "https://eu.wotblitz.com/en/api/events/items/auction/?page_size=80&type[]=vehicle&saleable=true"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	errorGroup, ctx := errgroup.WithContext(ctx)

	monitoringStorage := storage.NewRuntimeMonitoringStorage()
	flowStorage := storage.NewRuntimeAddMonitoringFlowStorage()
	actionProvider := provider.NewCachedAuctionDataProvider(provider.NewHttpActionProvider(apiUrl),
		time.Minute*10) // todo move to config

	tgBot, err := bot.New(botToken) // todo move to cfg
	if err != nil {
		log.Fatal(err)
	}

	initBot(monitoringStorage, flowStorage, actionProvider, tgBot)
	errorGroup.Go(func() error {
		return tgBot.Run(ctx)
	})
	errorGroup.Go(func() error {
		return run_vehicles_checks(ctx, actionProvider, monitoringStorage, tgBot, time.Second*10) // todo move to config
	})

	errorGroup.Wait()
}

func run_vehicles_checks(ctx context.Context, dataProvider provider.AuctionDataProvider,
	monitoringStorage storage.MonitoringStorage,
	bot *bot.Bot, interval time.Duration) error {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			data, err := dataProvider.GetData(ctx, true)
			if err != nil {
				return err
			}

			for _, v := range data {
				users, err := monitoringStorage.GetAllByVehicleAndCountGte(ctx, v.Name, v.CurrentCount)
				if err != nil {
					return err
				}

				for _, u := range users {
					photo := tgbotapi.NewPhotoShare(u.ChatID, v.Img)
					photo.Caption = fmt.Sprintf("Attention! Only %d %s`s left. Current price is %d gold",
						v.CurrentCount, v.Name, v.Price)

					if err := bot.Send(photo); err != nil {
						return err
					}
				}
			}
		}
	}
}

func initBot(monitoringStorage storage.MonitoringStorage,
	flowStorage storage.AddMonitoringFlowStorage,
	dataProvider provider.AuctionDataProvider, tgBot *bot.Bot) {
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
	tgBot.AddHandler(
		bot.GetCommandNamePredicate("add_monitoring"),
		bot.GetAddMonitoringCommandHandler(dataProvider, flowStorage),
	)
	tgBot.AddHandler(
		bot.GetCommandNamePredicate("monitoring"),
		bot.GetMonitoringCommandHandler(monitoringStorage),
	)
	tgBot.AddHandler(
		bot.GetSetVehicleNameFlowCallbackPredicate(flowStorage, dto.StepSelectVehicle),
		bot.GetAddMonitoringVehicleStepHandler(flowStorage),
	)
	tgBot.AddHandler(
		bot.GetSetVehicleMinimalCountFlowPredicate(flowStorage, dto.StepEnterMinimalCount),
		bot.GetAddMonitoringMinimalCountStepHandler(flowStorage, monitoringStorage),
	)
	tgBot.AddHandler(
		bot.GetCommandNamePredicate("cancel"),
		bot.GetCancelCommandHandler(flowStorage),
	)
}
