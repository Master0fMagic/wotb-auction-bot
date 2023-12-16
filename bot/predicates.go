package bot

import (
	"context"
	"github.com/Master0fMagic/wotb-auction-bot/dto"
	"github.com/Master0fMagic/wotb-auction-bot/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetCommandNamePredicate(name string) Predicate {
	return func(update tgbotapi.Update) bool {
		if update.Message == nil {
			return false
		}
		return update.Message.IsCommand() &&
			update.Message.Command() == name
	}
}

func GetSetVehicleNameFlowCallbackPredicate(flowStorage storage.AddMonitoringFlowStorage, step dto.MonitoringStep) Predicate {
	return func(update tgbotapi.Update) bool {
		if update.CallbackQuery == nil {
			return false
		}

		ctx := context.TODO()
		flowStep, err := flowStorage.Get(ctx, update.CallbackQuery.Message.Chat.ID)
		return err == nil && flowStep != nil && flowStep.Step == step
	}
}

func GetSetVehicleMinimalCountFlowPredicate(flowStorage storage.AddMonitoringFlowStorage, step dto.MonitoringStep) Predicate {
	return func(update tgbotapi.Update) bool {
		if update.Message == nil {
			return false
		}

		ctx := context.TODO()
		flowStep, err := flowStorage.Get(ctx, update.Message.Chat.ID)
		return err == nil && flowStep != nil && flowStep.Step == step
	}
}
