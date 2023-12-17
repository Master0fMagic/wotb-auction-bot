package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"sync"
)

type HandlerFunc func(tgbotapi.Update, *tgbotapi.BotAPI) error
type Predicate func(tgbotapi.Update) bool

type messageHandler struct {
	predicate Predicate
	handler   HandlerFunc
}

type Bot struct {
	tgBot *tgbotapi.BotAPI
	token string
	mtx   sync.RWMutex

	msgHandlers []messageHandler
}

func New(token string) (*Bot, error) {
	tgBot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	tgBot.Debug = true

	return &Bot{
		tgBot:       tgBot,
		token:       token,
		msgHandlers: make([]messageHandler, 0),
	}, err
}

func (b *Bot) AddHandler(predicate Predicate, handlerFunc HandlerFunc) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	b.msgHandlers = append(b.msgHandlers, messageHandler{handler: handlerFunc, predicate: predicate})
}

func (b *Bot) Send(msg tgbotapi.Chattable) error {
	if _, err := b.tgBot.Send(msg); err != nil {
		// todo log error
		return err
	}

	return nil
}

func (b *Bot) Run(ctx context.Context) error {
	b.tgBot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.tgBot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-ctx.Done():
			// todo add log for canceling via ctx
			return nil
		case update := <-updates:
			b.mtx.RLock()
			for _, handler := range b.msgHandlers {
				if !handler.predicate(update) {
					continue
				}

				err := handler.handler(update, b.tgBot)
				if err != nil {
					// todo log with error
				}
				break
			}

			b.mtx.RUnlock()
		}
	}
}
