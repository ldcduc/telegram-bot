package main

import (
	"fmt"
	"os"
	"telegram_bot/helpers"

	"github.com/urfave/cli"
)

type Client struct {
	TeleClient *helpers.Telegram
	SendCount  int
}

func main() {
	app := cli.NewApp()
	app.Name = "blcMonitor"
	app.Usage = "sends messages to telegram when node dont increase blocks"
	app.Version = "0.0.1"
	app.Commands = healthCheckCommand()

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func healthCheckCommand() []cli.Command {
	healthCheckCmd := cli.Command{
		// Action:      blcMonitor,
		Action:      blockHeightMonitor,
		Name:        "start",
		Usage:       "Alert to telegram when block is stuck",
		Description: `Alert to telegram when block is stuck`,
	}
	healthCheckCmd.Flags = helpers.NewTeleClientFlag()

	return []cli.Command{healthCheckCmd}
}