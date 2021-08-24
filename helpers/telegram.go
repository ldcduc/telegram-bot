package helpers

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/urfave/cli"
)

var (
	botAPITokenFlag = []cli.StringFlag{
	}
	botAPITokenFlagTest = []cli.StringFlag {
	}
	chatIDFlag = cli.Int64Flag {
		Name:  "chatId",
		Usage: "The ID of group/chanel",
		Value: -, // Real alert group
	}
	chatIDFlagTest = cli.Int64Flag {
		Name:  "chatId",
		Usage: "The ID of group/chanel",
		Value: -, // Group for testing
	}
)
type Telegram struct {
	ChatId  int64
	Bot     *tgbotapi.BotAPI
	IsDebug bool
}

/* Test case
Test bots: 5 bots
	Lower edge
./telegram_bot start test 0
./telegram_bot start test 0 'message'
	First bot
./telegram_bot start test 1 
./telegram_bot start test 1 'message'
	Normal bot
./telegram_bot start test 2 
./telegram_bot start test 2 'message'
	Last bot 
./telegram_bot start test 5
./telegram_bot start test 5 'message'
	Higher edge
./telegram_bot start test 6
./telegram_bot start test 6 'message'
Alert bots: 3 bots
	Lower edge
./telegram_bot start 0
./telegram_bot start 0 'message'
	First bot
./telegram_bot start 1
./telegram_bot start 1 'message'
	Normal bot
./telegram_bot start 2 
./telegram_bot start 2 'message'
	Last bot
./telegram_bot start 3
./telegram_bot start 3 'message'
	Higher edge
./telegram_bot start 4
./telegram_bot start 4 'message'

./telegram_bot start id 'message'
*/

func NewTeleClientFlag() []cli.Flag {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), os.Args)
	if len(os.Args) > 2 {
		if os.Args[2] == "test" {
			id, err := strconv.Atoi(os.Args[3])
			if err != nil { // No id supplied -> Choose the first bot
				id = 1
			} 
			if id > len(botAPITokenFlagTest) { // id > number of bots -> Choose the last bot
				id = len(botAPITokenFlagTest)
			}
			return []cli.Flag{botAPITokenFlagTest[int(math.Max(0, float64(id - 1)))], chatIDFlagTest}
		} else {
			id, err := strconv.Atoi(os.Args[2])
			if err != nil { // No id supplied -> Choose the first bot
				id = 1
			} 
			if id > len(botAPITokenFlag) { // id > number of bots -> Choose the last bot
				id = len(botAPITokenFlag) 
			}
			return []cli.Flag{botAPITokenFlag[int(math.Max(0, float64(id - 1)))], chatIDFlag}
		}
	}
	// ./telegram_bot start
	return []cli.Flag{botAPITokenFlag[0], chatIDFlag}
}

func NewTeleClientFromFlag(ctx *cli.Context, MaxNodes int) (*Telegram, error, string) {
	var (
		botAPIToken = ctx.String(botAPITokenFlag[0].Name)
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

	if len(os.Args) >= 3 { 
		_, err := strconv.Atoi(os.Args[len(os.Args) - 1]);
		if (err != nil) {
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
