package command

import (
	"strings"
	"twitch-bot/channels"
	"twitch-bot/config"
	"twitch-bot/core/command/messages"
)

func HandleMessage(inMsg *messages.Incoming) error {
	if inMsg == nil {
		return nil
	}
	args, _, err := prepareMessage(*inMsg.Message)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return nil
	}
	command := args[0]
	args = args[1:]

	switch command {
	case "ping":
		resp := "Pong! " + *inMsg.User + " in #" + *inMsg.Channel
		channels.MessagesOUT <- &messages.Outgoing{
			Platform:        messages.Twitch,
			Message:         &resp,
			Channel:         inMsg.Channel,
			User:            inMsg.User,
			IncomingMessage: inMsg,
		}
	}
	return nil
}

func prepareMessage(messageText string) ([]string, string, error) {
	conf, err := config.GetPublic()
	if err != nil {
		return nil, "", err
	}

	if !strings.HasPrefix(messageText, conf.Global.CommandPrefix) {
		return nil, "no-prefix", nil
	}

	messageText = strings.Trim(messageText, " \t\n󠀀⠀")
	args := strings.Split(strings.TrimPrefix(messageText, conf.Global.CommandPrefix), " ")
	return args, "", nil
}
