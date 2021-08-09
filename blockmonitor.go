package main

import (
	"fmt"
	"log"
	"time"

	"github.com/urfave/cli"

	"telegram_bot/helpers"
	"telegram_bot/node"
)

const (
	// MaxTimes is a max try to sending messages
	MaxTimes = 3 // the times try sends message
	MaxNodes = 7 // number of nodes to check
	MinutePeriod = 5 // check every MinutePeriod minutes
)

func blockHeightMonitor(ctx *cli.Context) {
	client := &Client{
		SendCount: 0,
	}

	teleClient, err, startingInfo := helpers.NewTeleClientFromFlag(ctx, MaxNodes)
	if err != nil {
		log.Printf("can not init telegram bot %s", err.Error())
		return
	}
	log.Print("Connected to Telegram")
	client.TeleClient = teleClient
	
	var last_consensused_block_height = []int {}
	for i := 0; i < MaxNodes; i ++ {
		last_consensused_block_height = append(last_consensused_block_height, 0)
	}

	isLastTimeChanged := true
	second := MinutePeriod * 60
	fmt.Printf("Check for every %d seconds\n", second)
	caption := fmt.Sprintf("Monitoring bot started, checks for every %d seconds\nSend %d messages in maximum if failed", second, MaxTimes)
	sendAlert(client, startingInfo, caption, true)
	for {
		var tmp = []int {}
		for i := 0; i < MaxNodes; i ++ {
			tmp = append(tmp, last_consensused_block_height[i])
		}

		var failedNode = []int {}
		failedNode = node.Check_last_consensused_height(MaxNodes, last_consensused_block_height, failedNode)
		fmt.Println("failedNode len", len(failedNode))
		caption := ""
		if len(failedNode) > 0 {
			isLastTimeChanged = false
			for i := 0; i < len(failedNode); i ++ {
				caption = caption + fmt.Sprintf("node=%d: last_consensused_block_height=%d, current=%d\n", failedNode[i], tmp[failedNode[i] - 1], last_consensused_block_height[failedNode[i] - 1])
			}
			sendAlert(client, "Found no change in last_consensused_block_height - Alert sent", caption, isLastTimeChanged)
			fmt.Println("Alert sent: {:?}", caption)
		} else {
			isLastTimeChanged = true
			fmt.Println("Last consensused height changed")
		}

		time.Sleep(time.Second * time.Duration(second))
	}
}
func sendAlert(client *Client, msg string, caption string, forceSend bool) {
	if forceSend {
		// send message not increase counter
		log.Printf("================send msg: %s", msg)
		client.TeleClient.SendMessage(msg, caption)
		client.SendCount = 1
		return
	}
	if client.SendCount > MaxTimes {
		return
	}
	log.Printf("================send msg: %s", msg)
	client.TeleClient.SendMessage(msg, caption)
	client.SendCount++
}
