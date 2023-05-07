# Gotify-CLI [![travus-badge][travis-badge]][travis] [![badge-release][badge-release]][release]

Gotify-CLI is a command line client for pushing messages to [gotify/server][gotify/server]. It is **not** required to push messages. See [alternatives](#alternatives).

<p align="center">
    <img src="gotify_cli.gif"/>
</p>

## Features

* stores token/url in a config file
* initialization wizard
* piping support (`echo message | gotify push`)
* simple to use
* watch and push script result changes (`gotify watch "curl http://example.com/api | jq '.data'"`)

## Alternatives

You can simply use [curl](https://curl.haxx.se/), [HTTPie](https://httpie.org/) or any other http-client to push messages.

```bash
$ curl -X POST "https://push.example.de/message?token=<apptoken>" -F "title=my title" -F "message=my message"
$ http -f POST "https://push.example.de/message?token=<apptoken>" title="my title" message="my message"
```

## Installation

### Manually

Download the [latest release][release] for your os: (this example uses version `v2.2.0`)

```bash
$ wget -O gotify https://github.com/gotify/cli/releases/download/v2.2.0/gotify-cli-linux-amd64
# or
$ curl -Lo gotify https://github.com/gotify/cli/releases/download/v2.2.0/gotify-cli-linux-amd64
```

Make `gotify` executable:

```bash
$ chmod +x gotify
```

Test if the Gotify-CLI works: *(When it doesn't work, you may have downloaded the wrong file or your device/os isn't supported)*

```bash
$ gotify version
```

It should output something like this:

```bash
Version:   2.2.0
Commit:    95cc125e39f805a1369b0c746c8bf3af15797a57
BuildDate: 2020-09-25-16:27:04
```

*(optional)* Move the executable to a folder on your `$PATH`:

```bash
$ mv gotify /usr/bin/gotify
```

### Arch Linux

Gotify-CLI is available [in the AUR][AUR], so you can install it with
[any AUR helper][AUR-helpers], i.e. [`yay`][yay]:

```bash
$ yay -S gotify-cli
```

> Maintenance of the AUR package is not performed by the Gotify team.

### macOS

Gotify-CLI is available [in Homebrew][brew], so you can install it with:

```bash
$ brew install gotify
```

> Maintenance of the Brew package is not performed by the Gotify team.

## Usage

Now you can either run the initialization wizard or [create a config manually](#Configuration). This tutorial uses the wizard.
```bash
$ gotify init
```
When you've finished initializing Gotify-CLI, you are ready to push messages to [gotify/server][gotify/server].

Here are some examples commands, you can view the "push help" via `gotify help push` (or have a look at [push help](#push-help)).
```json
$ gotify push my message
$ gotify push "my message"
$ echo my message | gotify push
$ gotify push < somefile
$ gotify push -t "my title" -p 10 "my message"
$ gotify watch "curl http://example.com/api | jq '.data'"
```

## Help

**Uses version `v2.2.0`**

```bash
NAME:
   Gotify - The official Gotify-CLI

USAGE:
   gotify [global options] command [command options] [arguments...]

VERSION:
   2.2.0

COMMANDS:
     init        Initializes the Gotify-CLI
     version, v  Shows the version
     config      Shows the config
     push, p     Pushes a message
     watch       watch the result of a command and pushes output difference
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Watch help

```
NAME:
   gotify watch - watch the result of a command and pushes output difference

USAGE:
   gotify watch [command options] <cmd>

OPTIONS:
   --interval value, -n value  watch interval (sec) (default: 2)
   --priority value, -p value  Set the priority (default: 0)
   --exec value, -x value      Pass command to exec (default: "sh -c")
   --title value, -t value     Set the title (empty for command)
   --token value               Override the app token [$GOTIFY_TOKEN]
   --url value                 Override the Gotify URL
   --output value, -o value    Output verbosity (short|default|long) (default: "default")
```

### Push help

```bash
$ gotify help push
NAME:
   gotify push - Pushes a message

USAGE:
   gotify push [command options] <message-text>

DESCRIPTION:
   the message can also provided in stdin f.ex:
   echo my text | gotify push

OPTIONS:
   --priority value, -p value    Set the priority (default: 0)
   --title value, -t value       Set the title (empty for app name)
   --token value                 Override the app token [$GOTIFY_TOKEN]
   --url value                   Override the Gotify URL
   --quiet, -q                   Do not output anything (on success)
   --contentType value           The content type of the message. See https://gotify.net/docs/msgextras#client-display
   --disable-unescape-backslash  Disable evaluating \n and \t (if set, \n and \t will be seen as a string)
```

## Configuration

**Note: The config can be created by `gotify init`.**

Gotify-CLI will search the following paths for a config file:
* `/etc/gotify/cli.json`
* `$XDG_CONFIG_HOME/gotify/cli.json`
* `~/.gotify/cli.json`
* `./cli.json`

### Structure

| name  | description | example |
| ----- | ----------- | ------- |
| token | an application token (a client token will not work) | `A4ZudDRdLT40L5X` |
| url   | the URL to your [gotify/server][gotify/server]      | `https://gotify.example.com` |
| defaultPriority   | Default priority ( set to 0 if not present)       | `6` |

### Config example

```json
{
  "token": "A4ZudDRdLT40L5X",
  "url": "https://gotify.example.com",
  "defaultPriority": 6
}
```

## Configuration option

If needed, you can disable SSL handcheck validation using an environment variable:
```
export GOTIFY_SKIP_VERIFY_TLS=True
```


### Dockerfile
The Dockerfile contains the steps necessary to build a new version of the CLI and then run it in
a minimal Alpine container.

**Build:**

```bash
docker build -t gotify/gotify-cli .
```

**Run (this assumes your `cli.json` file is in the current working directory):**

```bash
docker run -it -v "$PWD/cli.json:/home/app/cli.json" gotify/gotify-cli:latest push -p 5 "Test from Gotify CLI"
```

 [gotify/server]: https://github.com/gotify/server
 [travis-badge]: https://app.travis-ci.com/gotify/cli.svg?branch=master
 [travis]: https://app.travis-ci.com/gotify/cli
 [badge-release]: https://img.shields.io/github/release/gotify/cli.svg
 [release]: https://github.com/gotify/cli/releases/latest
 [AUR]: https://aur.archlinux.org/packages/gotify-cli
 [AUR-helpers]: https://wiki.archlinux.org/title/AUR_helpers
 [yay]: https://aur.archlinux.org/packages/yay/
 [brew]: https://brew.sh
