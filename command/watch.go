package command

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/gotify/cli/v2/config"
	"github.com/gotify/cli/v2/utils"
	"github.com/gotify/go-api-client/v2/models"
	"gopkg.in/urfave/cli.v1"
)

func Watch() cli.Command {
	return cli.Command{
		Name:      "watch",
		Usage:     "watch the result of a command and pushes output difference",
		ArgsUsage: "<cmd>",
		Flags: []cli.Flag{
			cli.Float64Flag{Name: "interval,n", Usage: "watch interval (sec)", Value: 2},
			cli.IntFlag{Name: "priority,p", Usage: "Set the priority"},
			cli.StringFlag{Name: "exec,x", Usage: "Pass command to exec (default to \"sh -c\")", Value: "sh -c"},
			cli.StringFlag{Name: "title,t", Usage: "Set the title (empty for command)"},
			cli.StringFlag{Name: "token", Usage: "Override the app token"},
			cli.StringFlag{Name: "url", Usage: "Override the Gotify URL"},
			cli.StringFlag{Name: "output,o", Usage: "Output verbosity (short|default|long)", Value: "default"},
		},
		Action: doWatch,
	}
}

func doWatch(ctx *cli.Context) {
	conf, confErr := config.ReadConfig(config.GetLocations())

	cmdArgs := ctx.Args()
	cmdStringNotation := strings.Join(cmdArgs, " ")
	execArgs := strings.Split(ctx.String("exec"), " ")
	cmdArgs = append(execArgs[1:], cmdStringNotation)
	execCmd := execArgs[0]

	outputMode := ctx.String("output")
	if !(outputMode == "default" || outputMode == "long" || outputMode == "short") {
		utils.Exit1With("output mode should be short|default|long")
		return
	}
	interval := ctx.Float64("interval")
	priority := ctx.Int("priority")
	title := ctx.String("title")
	if title == "" {
		title = cmdStringNotation
	}
	token := ctx.String("token")
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
	parsedURL, err := url.Parse(stringURL)
	if err != nil {
		utils.Exit1With("invalid url", stringURL)
		return
	}

	watchInterval := time.Duration(interval*1000) * time.Millisecond

	evalCmdOutput := func() (string, error) {
		cmd := exec.Command(execCmd, cmdArgs...)
		timeOut := time.After(watchInterval)
		outputBuf := bytes.NewBuffer([]byte{})
		cmd.Stdout = outputBuf
		cmd.Stderr = outputBuf
		err := cmd.Start()
		if err != nil {
			return "", fmt.Errorf("command failed to invoke: %v", err)
		}
		done := make(chan error)
		go func() {
			err := cmd.Wait()
			if err != nil {
				done <- fmt.Errorf("command failed to invoke: %v", err)
			}
			done <- nil
		}()
		select {
		case err := <-done:
			return outputBuf.String(), err
		case <-timeOut:
			cmd.Process.Kill()
			return outputBuf.String(), errors.New("command timed out")
		}
	}

	lastOutput, err := evalCmdOutput()
	if err != nil {
		utils.Exit1With("first run failed", err)
	}
	for range time.NewTicker(watchInterval).C {
		output, err := evalCmdOutput()
		if err != nil {
			output += fmt.Sprintf("\n!== <%v> ==!", err)
		}
		if output != lastOutput {
			msgData := bytes.NewBuffer([]byte{})

			switch outputMode {
			case "long":
				fmt.Fprintf(msgData, "command output for \"%s\" changed:\n\n", cmdStringNotation)
				fmt.Fprintln(msgData, "== BEGIN OLD OUTPUT ==")
				fmt.Fprint(msgData, lastOutput)
				fmt.Fprintln(msgData, "== END OLD OUTPUT ==")
				fmt.Fprintln(msgData, "== BEGIN NEW OUTPUT ==")
				fmt.Fprint(msgData, output)
				fmt.Fprintln(msgData, "== END NEW OUTPUT ==")
			case "default":
				fmt.Fprintf(msgData, "command output for \"%s\" changed:\n\n", cmdStringNotation)
				fmt.Fprintln(msgData, "== BEGIN NEW OUTPUT ==")
				fmt.Fprint(msgData, output)
				fmt.Fprintln(msgData, "== END NEW OUTPUT ==")
			case "short":
				fmt.Fprintf(msgData, output)
			}

			msgString := msgData.String()
			fmt.Println(msgString)
			pushMessage(parsedURL, token, models.MessageExternal{
				Title:    title,
				Message:  msgString,
				Priority: priority,
			}, true)
			lastOutput = output
		}
	}

}
