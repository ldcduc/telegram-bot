package main

import (
	"fmt"
	"os"
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
		fmt.Printf("can not init telegram bot %s", err.Error())
		return
	}
	fmt.Print("Connected to Telegram")
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
	var deadNode = []int {}
	for {
		var tmp = []int {}
		for i := 0; i < MaxNodes; i ++ {
			tmp = append(tmp, last_consensused_block_height[i])
		}

		var failedNode = []int {}
		failedNode = node.Check_last_consensused_height(MaxNodes, last_consensused_block_height, failedNode, deadNode)
		fmt.Println("failedNode len", len(failedNode))
		caption := ""
		if len(failedNode) > 0 {
			isLastTimeChanged = false
			for i := 0; i < len(failedNode); i ++ {
				caption = caption + fmt.Sprintf("node=%d: last_consensused_block_height=%d, current=%d\n", failedNode[i], tmp[failedNode[i] - 1], last_consensused_block_height[failedNode[i] - 1])
			}
			sendAlert(client, caption, fmt.Sprintf("Found no change in last_consensused_block_height - Alert sent - Previous dead nodes = %v", deadNode), isLastTimeChanged)
			fmt.Println("Alert sent:", caption)
			if !isLastTimeChanged && client.SendCount > MaxTimes { // not forceSend + sent enough messages
				client.SendCount = 0
				// Append failedNode into deadNode
				for _, fnode := range failedNode {
					existed := false
					for _, dnode := range deadNode {
						if fnode == dnode {
							existed = true
						}
					}
					if !existed {
						deadNode = append(deadNode, fnode)
					}
				}
				//
				if len(deadNode) == MaxNodes {
					sendAlert(client, "All nodes have been dead", "Monitoring bot stopped", isLastTimeChanged)
					fmt.Println("Monitoring bot stopped: All nodes have been dead")
					os.Exit(0)
				}
				fmt.Println("client.SendCount is reset to 0 - deadNode =", deadNode)
			}
		} else {
			isLastTimeChanged = true
			currentTime := time.Now()
			fmt.Println(currentTime.Format("2006-01-02 15:04:05"), "Last consensused height changed")
		}

		time.Sleep(time.Second * time.Duration(second))
	}
}
func sendAlert(client *Client, msg string, caption string, forceSend bool) {
	if forceSend {
		// send message not increase counter
		// fmt.Printf("================send msg: %s\n", msg)
		client.TeleClient.SendMessage(msg, caption)
		client.SendCount = 1
		return
	}
	if client.SendCount > MaxTimes {
		return
	}
	// fmt.Printf("================send msg: %s\n", msg)
	client.TeleClient.SendMessage(msg, caption)
	client.SendCount++
}
