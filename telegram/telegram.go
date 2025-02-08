package telegram

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TelegramBot struct {
	bot        *bot.Bot
	chatId     string
	choiceChan chan string // Channel to capture user choice
}

func (tgBot *TelegramBot) callbackQueryHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	query := update.CallbackQuery

	// Send user choice to the channel
	tgBot.choiceChan <- query.Data

	// Close the channel after receiving the response
	close(tgBot.choiceChan)
}

func Initialize(token, chatId string) (TelegramBot, error) {
	b, err := bot.New(token)
	if err != nil {
		return TelegramBot{}, err
	}

	tgBot := &TelegramBot{
		bot:    b,
		chatId: chatId,
	}

	// Register callback handler
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "", bot.MatchTypePrefix, tgBot.callbackQueryHandler)

	return *tgBot, nil
}

// SendQuery waits for a user response using a channel
func (b TelegramBot) SendQuery(ctx context.Context, text string, options []string) (string, error) {
	// Create a new channel for capturing the user’s choice
	b.choiceChan = make(chan string, 1)

	// Build inline keyboard
	var keyboard [][]models.InlineKeyboardButton
	for _, option := range options {
		keyboard = append(keyboard, []models.InlineKeyboardButton{
			{Text: option, CallbackData: option},
		})
	}

	// Send a message with inline keyboard
	_, err := b.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      b.chatId,
		Text:        text,
		ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: keyboard},
	})
	if err != nil {
		return "", err
	}

	// Wait for response
	choice := <-b.choiceChan
	return choice, nil
}

// SendQuery waits for a user response using a channel
func (b TelegramBot) SendMessage(ctx context.Context, text string) error {
	// Send a message with inline keyboard
	_, err := b.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: b.chatId,
		Text:   text,
	})
	if err != nil {
		return err
	}

	return nil
}
