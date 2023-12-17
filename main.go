package main

import (
	"context"
	"fmt"
	"github.com/Master0fMagic/wotb-auction-bot/bot"
	"github.com/Master0fMagic/wotb-auction-bot/config"
	"github.com/Master0fMagic/wotb-auction-bot/dto"
	"github.com/Master0fMagic/wotb-auction-bot/provider"
	"github.com/Master0fMagic/wotb-auction-bot/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		log.WithError(err).Fatal("error parsing config")
	}

	lvl, err := initLogger(cfg.LogLevel)
	if err != nil {
		log.WithError(err).Fatal("error config logger")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	errorGroup, ctx := errgroup.WithContext(ctx)

	monitoringStorage, err := storage.NewSQLiteMonitoringStorage(cfg.DbPath)
	if err != nil {
		log.WithError(err).Fatal("error initializing monitoring storage")
	}

	flowStorage := storage.NewRuntimeAddMonitoringFlowStorage()
	actionProvider := provider.NewCachedAuctionDataProvider(provider.NewHTTPActionProvider(cfg.AuctionAPI), cfg.AuctionCacheLifetime)

	tgBot, err := bot.New(cfg.BotToken, lvl)
	if err != nil {
		log.WithError(err).Error("error initializing tg bot")
	}

	initBot(monitoringStorage, flowStorage, actionProvider, tgBot)
	errorGroup.Go(func() error {
		return tgBot.Run(ctx)
	})
	errorGroup.Go(func() error {
		return runVehiclesChecks(ctx, actionProvider, monitoringStorage, tgBot, cfg.AuctionCacheLifetime)
	})
	errorGroup.Go(func() error {
		return actionProvider.Run(ctx)
	})
	if err := errorGroup.Wait(); err != nil {
		log.WithError(err).Error("error awaiting error group")
	}
}

func runVehiclesChecks(ctx context.Context, dataProvider provider.AuctionDataProvider,
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
					if err := monitoringStorage.Remove(ctx, u.ChatID, v.Name); err != nil {
						return err
					}
				}
			}
		}
	}
}

func initLogger(logLevel string) (log.Level, error) {
	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		return lvl, err
	}
	log.SetLevel(lvl)

	return lvl, nil
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
		bot.GetCommandNamePredicate("all_data_short"),
		bot.GetAllDataShortCommandHandler(dataProvider),
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
