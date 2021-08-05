package helpers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/urfave/cli"
)

var (
	botAPITokenFlag = cli.StringFlag{
		Name:  "apiToken",
		Usage: "The API token",
	}
	chatIDFlag = cli.Int64Flag{
		Name:  "chatId",
		Usage: "The ID of group/chanel",
	}
)
type Telegram struct {
	ChatId  int64
	Bot     *tgbotapi.BotAPI
	IsDebug bool
}

func NewTeleClientFlag() []cli.Flag {
	return []cli.Flag{botAPITokenFlag, chatIDFlag}
}


func NewTeleClientFromFlag(ctx *cli.Context) (*Telegram, error) {
	var (
		botAPIToken = ctx.String(botAPITokenFlag.Name)
		chatID      = ctx.Int64(chatIDFlag.Name)
	)

	telegram := &Telegram{
		ChatId:  chatID,
		IsDebug: false,
	}
	bot, err := tgbotapi.NewBotAPI(botAPIToken)
	if err != nil {
		return nil, err
	}
	bot.Debug = telegram.IsDebug

	telegram.Bot = bot
	return telegram, nil
}

func (t *Telegram) SendMessage(content string, caption string) error {
	text := fmt.Sprintf("<b>%s</b>: %s", caption, content)
	msg := tgbotapi.NewMessage(t.ChatId, text)
	msg.ParseMode = "html"
	_, err := t.Bot.Send(msg)

	if err != nil {
		return err
	}
	return nil
}