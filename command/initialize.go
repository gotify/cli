package command

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"syscall"

	apiruntime "github.com/go-openapi/runtime"
	"github.com/gotify/cli/v2/config"
	"github.com/gotify/cli/v2/utils"
	"github.com/gotify/go-api-client/v2/auth"
	api "github.com/gotify/go-api-client/v2/client"
	"github.com/gotify/go-api-client/v2/client/application"
	"github.com/gotify/go-api-client/v2/client/message"
	"github.com/gotify/go-api-client/v2/gotify"
	"github.com/gotify/go-api-client/v2/models"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/urfave/cli.v1"
)

func Init() cli.Command {
	return cli.Command{
		Name:   "init",
		Usage:  "Initializes the Gotify-CLI",
		Action: doInit,
	}
}

func doInit(ctx *cli.Context) {
	serverURL := inputServerURL()
	hr()
	token := inputToken(gotify.NewClient(serverURL, utils.CreateHTTPClient()))
	hr()
	defaultPriority := inputDefaultPriority()
	hr()

	conf := &config.Config{
		URL:             serverURL.String(),
		Token:           token,
		DefaultPriority: defaultPriority,
	}

	pathToWrite, err := config.ExistingConfig(config.GetLocations())

	var writeErr error
	if err == config.ErrNoneSet {
		pathToWrite = inputConfigLocation()
		writeErr = config.WriteConfig(pathToWrite, conf)
	} else {
		writeErr = config.WriteConfig(pathToWrite, conf)
	}

	if writeErr == nil {
		fmt.Println("Written config to:", pathToWrite)
	} else {
		fmt.Println("Something went wrong: ", writeErr)
	}
}

func inputConfigLocation() string {
	locations := config.GetLocations()

	if len(locations) == 1 {
		return locations[0]
	}

	for {
		fmt.Println("Where to put the config file?")
		for i, location := range locations {
			fmt.Println(fmt.Sprintf("%d. %s", i+1, location))
		}
		value := inputString("Enter a number: ")
		hr()

		choice, err := strconv.Atoi(value)
		if err != nil {
			continue
		}
		if choice > 0 && choice <= len(locations) {
			indexedChoice := choice - 1
			return locations[indexedChoice]
		}
	}
}

func inputString(text string) string {
	fmt.Print(text)
	reader := bufio.NewReader(os.Stdin)
	readString, err := reader.ReadString('\n')
	if err != nil {
		utils.Exit1With(err)
	}
	return strings.TrimSpace(readString)
}

func inputToken(gotify *api.GotifyREST) string {
	for {
		fmt.Println("Configure an application token")
		fmt.Println("1. Enter an application-token")
		fmt.Println("2. Create an application token (with user/pass)")
		value := inputString("Enter 1 or 2: ")
		hr()

		switch value {
		case "1":
			return inputRawToken(gotify)
		case "2":
			return inputCredentialsAndCreateToken(gotify)
		}
	}
}

func inputCredentialsAndCreateToken(gotify *api.GotifyREST) string {
	for {
		fmt.Println("Enter Credentials (only used for creating the token not saved afterwards)")
		username := inputString("Username: ")
		fmt.Print("Password: ")
		password := readPassword()
		basicAuth := auth.BasicAuth(username, password)
		_, err := utils.SpinLoader("Authenticating", func(success chan interface{}, failure chan error) {
			user, err := gotify.User.CurrentUser(nil, basicAuth)
			if err == nil {
				success <- *user.Payload
			} else {
				failure <- err
			}
		})

		if err == nil {
			hr()
			return createToken(gotify, basicAuth)
		}
	}
}

func readPassword() string {
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return inputString("")
	}
	fmt.Println()
	return string(password)
}

func createToken(gotify *api.GotifyREST, auth apiruntime.ClientAuthInfoWriter) string {
	for {
		name := inputString("Application name: ")
		if name == "" {
			erred("Name may not be empty")
			continue
		}
		description := inputString("Application description (can be empty): ")

		resp, err := utils.SpinLoader("Creating", func(success chan interface{}, failure chan error) {
			params := application.NewCreateAppParams()
			params.Body = &models.Application{Name: name, Description: description}
			resp, err := gotify.Application.CreateApp(params, auth)
			if err == nil {
				success <- *resp.Payload
			} else {
				failure <- err
			}
		})

		if err == nil {
			app := resp.(models.Application)
			return app.Token
		}
	}

}

func inputRawToken(gotify *api.GotifyREST) string {
	for {
		enteredToken := inputString("Application Token: ")

		if len(enteredToken) != 15 {
			fmt.Println("A application token must have a length of 15 characters")
			hr()
			continue
		}
		_, err := utils.SpinLoader("Validating", func(success chan interface{}, failure chan error) {
			params := message.NewCreateMessageParams()
			params.Body = &models.MessageExternal{
				Title:    "Test message",
				Message:  "Test message from Gotify CLI",
				Priority: 0,
			}

			resp, err := gotify.Message.CreateMessage(params, auth.TokenAuth(enteredToken))

			if err == nil {
				success <- *resp.Payload
			} else {
				failure <- err
			}
		})
		if err == nil {
			return enteredToken
		}
		hr()
	}

}

func inputDefaultPriority() int {
	for {
		defaultPriorityStr := inputString("Default Priority [0-10]: ")
		defaultPriority, err := strconv.Atoi(defaultPriorityStr)
		if err != nil || (defaultPriority > 10 || defaultPriority < 0) {
			erred("Priority needs to be a number between 0 and 10.")
			continue
		} else {
			return defaultPriority
		}
		hr()
	}
}

func inputServerURL() *url.URL {
	for {
		rawURL := inputString("Gotify URL: ")
		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			erred("Could not parse URL:", err)
			continue
		}
		if parsedURL.Scheme == "" {
			erred("Add a scheme to the url (http:// or https://)")
			continue
		}

		if parsedURL.Host == "" {
			erred("The host part of the url may not be empty")
			continue
		}

		version, err := utils.SpinLoader("Connecting", func(success chan interface{}, failure chan error) {
			client := gotify.NewClient(parsedURL, utils.CreateHTTPClient())

			ver, e := client.Version.GetVersion(nil)
			if e == nil {
				success <- *ver.Payload
			} else {
				failure <- e
			}
		})
		if err == nil {
			info := version.(models.VersionInfo)
			fmt.Println(fmt.Sprintf("Gotify v%s@%s", info.Version, info.BuildDate))
			return parsedURL
		}
		hr()
	}
}

func erred(data ...interface{}) {
	fmt.Println(append([]interface{}{"Error"}, data...)...)
	hr()
}

func hr() {
	fmt.Println()
}
