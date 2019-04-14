package command

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gotify/cli/v2/config"
	"github.com/gotify/cli/v2/utils"
	"github.com/gotify/go-api-client/v2/auth"
	"github.com/gotify/go-api-client/v2/client/message"
	"github.com/gotify/go-api-client/v2/gotify"
	"github.com/gotify/go-api-client/v2/models"
	"gopkg.in/urfave/cli.v1"
)

func Push() cli.Command {
	return cli.Command{
		Name:        "push",
		Aliases:     []string{"p"},
		Usage:       "Pushes a message",
		ArgsUsage:   "<message-text>",
		Description: "the message can also provided in stdin f.ex:\n   echo my text | gotify push",
		Flags: []cli.Flag{
			cli.IntFlag{Name: "priority,p", Usage: "Set the priority"},
			cli.StringFlag{Name: "title,t", Usage: "Set the title (empty for app name)"},
			cli.StringFlag{Name: "token", Usage: "Override the app token"},
			cli.StringFlag{Name: "url", Usage: "Override the Gotify URL"},
			cli.BoolFlag{Name: "quiet,q", Usage: "Do not output anything (on success)"},
			cli.StringFlag{Name: "contentType", Usage: "The content type of the message. See https://gotify.net/docs/msgextras#client-display"},
		},
		Action: doPush,
	}
}

func doPush(ctx *cli.Context) {
	conf, confErr := config.ReadConfig(config.GetLocations())

	msgText := readMessage(ctx)

	priority := ctx.Int("priority")
	title := ctx.String("title")
	token := ctx.String("token")
	quiet := ctx.Bool("quiet")
	contentType := ctx.String("contentType")
	if token == "" {
		if confErr != nil {
			utils.Exit1With("token is not configured, run 'gotify init'")
			return
		}
		token = conf.Token
	}
	stringURL := ctx.String("url")
	if stringURL == "" {
		if confErr != nil {
			utils.Exit1With("url is not configured, run 'gotify init'")
			return
		}
		stringURL = conf.URL
	}

	msg := models.MessageExternal{
		Message:  msgText,
		Title:    title,
		Priority: priority,
	}

	if contentType != "" {
		msg.Extras = map[string]interface{}{
			"client::display": map[string]interface{}{
				"contentType": contentType,
			},
		}
	}

	parsedURL, err := url.Parse(stringURL)
	if err != nil {
		utils.Exit1With("invalid url", stringURL)
		return
	}

	pushMessage(parsedURL, token, msg, quiet)
}

func pushMessage(parsedURL *url.URL, token string, msg models.MessageExternal, quiet bool) {
	client := gotify.NewClient(parsedURL, &http.Client{})

	params := message.NewCreateMessageParams()
	params.Body = &msg
	_, err := client.Message.CreateMessage(params, auth.TokenAuth(token))
	if err == nil {
		if !quiet {
			fmt.Println("message created")
		}
	} else {
		utils.Exit1With(err)
	}
}

func readMessage(ctx *cli.Context) string {
	msgArgs := strings.Join(ctx.Args(), " ")

	msgStdin := utils.ReadFrom(os.Stdin)

	if msgArgs == "" && msgStdin == "" {
		utils.Exit1With("a message must be set, either as argument or via stdin")
	}

	if msgArgs != "" && msgStdin != "" {
		utils.Exit1With("a message is set via stdin and arguments, use only one of them")
	}

	if msgArgs == "" {
		return msgStdin
	} else {
		return msgArgs
	}
}
