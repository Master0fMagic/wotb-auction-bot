package bot

import (
	"context"
	"github.com/Master0fMagic/wotb-auction-bot/provider"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

func GetCommandNamePredicate(name string) Predicate {
	return func(update tgbotapi.Update) bool {
		return update.Message.IsCommand() &&
			update.Message.Command() == name
	}
}

func GetStaticTextResponseHandler(response string) HandlerFunc {
	return func(update tgbotapi.Update) ([]tgbotapi.Chattable, error) {
		return []tgbotapi.Chattable{tgbotapi.NewMessage(update.Message.Chat.ID, response)}, nil
	}
}

func GetAddMonitoringCommandHandler(dataProvider provider.AuctionDataProvider) HandlerFunc {
	return func(update tgbotapi.Update) ([]tgbotapi.Chattable, error) {
		data, err := dataProvider.GetData(context.TODO())
		if err != nil {
			return nil, err
		}

		var row []tgbotapi.KeyboardButton

		for _, entity := range data {
			btn := tgbotapi.NewKeyboardButton(entity.Name)
			row = append(row, btn)
		}

		msgReply := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose an entity:")
		msgReply.ReplyMarkup = tgbotapi.NewReplyKeyboard(row)

		return []tgbotapi.Chattable{msgReply}, nil
	}
}

func GetDataCommandHandler(dataProvider provider.AuctionDataProvider) HandlerFunc {
	return func(update tgbotapi.Update) ([]tgbotapi.Chattable, error) {
		data, err := dataProvider.GetData(context.TODO())
		if err != nil {
			return nil, err
		}

		var msgs []tgbotapi.Chattable
		for _, v := range data {
			photoConfig := tgbotapi.NewPhotoShare(update.Message.Chat.ID, v.Img)
			photoConfig.Caption = v.String()
			msgs = append(msgs, photoConfig)
		}

		return msgs, nil
	}
}

func GetDataShortCommandHandler(dataProvider provider.AuctionDataProvider) HandlerFunc {
	return func(update tgbotapi.Update) ([]tgbotapi.Chattable, error) {
		data, err := dataProvider.GetData(context.TODO())
		if err != nil {
			return nil, err
		}

		var stringData []string
		for _, v := range data {
			stringData = append(stringData, v.String())
		}

		response := tgbotapi.NewMessage(update.Message.Chat.ID,
			strings.Join(stringData, "\n\n"))

		return []tgbotapi.Chattable{response}, nil
	}
}
