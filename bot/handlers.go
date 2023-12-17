package bot

import (
	"context"
	"fmt"
	"github.com/Master0fMagic/wotb-auction-bot/dto"
	"github.com/Master0fMagic/wotb-auction-bot/provider"
	"github.com/Master0fMagic/wotb-auction-bot/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
)

func GetStaticTextResponseHandler(response string) HandlerFunc {
	return func(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
		_, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, response))
		return err
	}
}

func GetAddMonitoringCommandHandler(dataProvider provider.AuctionDataProvider, flowStorage storage.AddMonitoringFlowStorage) HandlerFunc {
	return func(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
		ctx := context.TODO()
		data, err := dataProvider.GetData(ctx, true)
		if err != nil {
			return err
		}

		var rows [][]tgbotapi.InlineKeyboardButton

		for _, entity := range data {
			btn := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(entity.Name, entity.Name))
			rows = append(rows, btn)
		}

		msgReply := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose an entity:")
		msgReply.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)

		if err := flowStorage.Set(ctx, dto.AddMonitoringStep{
			Data: dto.MonitoringData{
				ChatID: update.Message.Chat.ID,
			},
			Step: dto.StepSelectVehicle,
		}); err != nil {
			return err
		}

		_, err = bot.Send(msgReply)
		return err
	}
}

func GetAddMonitoringVehicleStepHandler(flowStorage storage.AddMonitoringFlowStorage) HandlerFunc {
	return func(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
		ctx := context.TODO()
		chatID := update.CallbackQuery.Message.Chat.ID

		flowData, err := flowStorage.Get(ctx, chatID)
		if err != nil {
			return err
		}
		flowData.Data.VehicleName = update.CallbackQuery.Data
		flowData.Step = dto.StepEnterMinimalCount

		if err := flowStorage.Set(ctx, *flowData); err != nil {
			return err
		}

		msgId := update.CallbackQuery.Message.MessageID
		editTextQuery := tgbotapi.NewEditMessageText(chatID, msgId,
			fmt.Sprintf("You chosed %s\nEnter minimal count:\nor /cancel", flowData.Data.VehicleName))
		editLineupQuery := tgbotapi.NewEditMessageReplyMarkup(chatID, msgId, tgbotapi.NewInlineKeyboardMarkup())

		_, err = bot.Send(editTextQuery)
		if err != nil {
			return err
		}

		_, err = bot.Send(editLineupQuery)
		if err != nil {
			return err
		}

		return nil
	}
}

func GetAddMonitoringMinimalCountStepHandler(flowStorage storage.AddMonitoringFlowStorage,
	monitoringStorage storage.MonitoringStorage) HandlerFunc {
	return func(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
		ctx := context.TODO()
		chatID := update.Message.Chat.ID

		flowData, err := flowStorage.Get(ctx, chatID)
		if err != nil {
			return err
		}
		minimalCount, err := strconv.Atoi(update.Message.Text)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "You have entered invalid value. Please enter integer number")
			_, err := bot.Send(msg)
			return err
		}

		flowData.Data.MinimalCount = minimalCount
		if err := flowStorage.Remove(ctx, chatID); err != nil {
			return err
		}
		if err := monitoringStorage.Save(ctx, flowData.Data); err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Monitoring for %s and %d count saved!",
			flowData.Data.VehicleName,
			flowData.Data.MinimalCount))

		_, err = bot.Send(msg)
		return err
	}
}

func GetDataCommandHandler(dataProvider provider.AuctionDataProvider) HandlerFunc {
	return func(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
		data, err := dataProvider.GetData(context.TODO(), true)
		if err != nil {
			return err
		}

		for _, v := range data {
			photoConfig := tgbotapi.NewPhotoShare(update.Message.Chat.ID, v.Img)
			photoConfig.Caption = v.String()

			_, err = bot.Send(photoConfig)
			if err != nil {
				return nil
			}
		}

		return err
	}
}

func GetDataShortCommandHandler(dataProvider provider.AuctionDataProvider) HandlerFunc {
	return func(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
		data, err := dataProvider.GetData(context.TODO(), true)
		if err != nil {
			return err
		}

		var stringData []string
		for _, v := range data {
			stringData = append(stringData, v.String())
		}

		response := tgbotapi.NewMessage(update.Message.Chat.ID,
			strings.Join(stringData, "\n\n"))

		_, err = bot.Send(response)
		return err
	}
}

func GetAllDataShortCommandHandler(dataProvider provider.AuctionDataProvider) HandlerFunc {
	return func(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
		data, err := dataProvider.GetData(context.TODO(), false)
		if err != nil {
			return err
		}

		var stringData []string
		for _, v := range data {
			stringData = append(stringData, v.String())
		}

		response := tgbotapi.NewMessage(update.Message.Chat.ID,
			strings.Join(stringData, "\n\n"))

		_, err = bot.Send(response)
		return err
	}
}

func GetMonitoringCommandHandler(monitoringStorage storage.MonitoringStorage) HandlerFunc {
	return func(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
		data, err := monitoringStorage.GetAll(context.TODO())
		if err != nil {
			return err
		}

		var responseData string
		if len(data) == 0 {
			responseData = "empty monitoring"
		} else {
			var stringData []string
			for _, v := range data {
				stringData = append(stringData, fmt.Sprintf("vehicle: %s, chatID: %d, minCount: %d",
					v.VehicleName, v.ChatID, v.MinimalCount))
			}
			responseData = strings.Join(stringData, "\n\n")
		}

		response := tgbotapi.NewMessage(update.Message.Chat.ID, responseData)

		_, err = bot.Send(response)
		return err
	}
}

func GetCancelCommandHandler(flowStorage storage.AddMonitoringFlowStorage) HandlerFunc {
	return func(update tgbotapi.Update, bot *tgbotapi.BotAPI) error {
		ctx := context.TODO()
		chatID := update.Message.Chat.ID

		if err := flowStorage.Remove(ctx, chatID); err != nil {
			return err
		}

		return nil
	}
}
