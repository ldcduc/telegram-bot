package helpers

import (
	"fmt"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/urfave/cli"
)

var (
	botAPITokenFlag = cli.StringFlag{
		Name:  "apiToken",
		Usage: "The API token",
		Value: "", // Bot in real alert group
	}
	botAPITokenFlagTest = []cli.StringFlag {
	}
	chatIDFlag = cli.Int64Flag{
		Name:  "chatId",
		Usage: "The ID of group/chanel",
		Value: -572689309, // Real alert group
	}
	chatIDFlagTest = cli.Int64Flag{
		Name:  "chatId",
		Usage: "The ID of group/chanel",
		Value: -533307840, // Group for testing
	}
)
type Telegram struct {
	ChatId  int64
	Bot     *tgbotapi.BotAPI
	IsDebug bool
}

// ./telegram_bot start test id 'message'
// ./telegram_bot start id 'message'

func NewTeleClientFlag() []cli.Flag {
	if len(os.Args) > 2 {
		if os.Args[2] == "test" {
			id, err := strconv.Atoi(os.Args[3]);
			if err != nil {
				id = 1
			} 
			if id >= len(botAPITokenFlagTest) {
				id = len(botAPITokenFlagTest)
			}
			return []cli.Flag{botAPITokenFlagTest[id - 1], chatIDFlagTest}
		}
	}	
	return []cli.Flag{botAPITokenFlag, chatIDFlag}
}


func NewTeleClientFromFlag(ctx *cli.Context, MaxNodes int) (*Telegram, error, string) {
	var (
		botAPIToken = ctx.String(botAPITokenFlag.Name)
		chatID      = ctx.Int64(chatIDFlag.Name)
	)

	startingInfo := fmt.Sprintf("%d node(s) scenario", MaxNodes)

	telegram := &Telegram{
		ChatId:  chatID,
		IsDebug: false,
	}
	bot, err := tgbotapi.NewBotAPI(botAPIToken)
	if err != nil {
		return nil, err, startingInfo
	}
	bot.Debug = telegram.IsDebug

	if len(os.Args) > 2 {
		if os.Args[len(os.Args) - 1] != "test" {
			startingInfo = fmt.Sprintf("%d node(s) scenario: %s", MaxNodes, os.Args[len(os.Args) - 1])
		}
	}	

	telegram.Bot = bot
	return telegram, nil, startingInfo
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
