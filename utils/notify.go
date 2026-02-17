package utils

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/ahmedsat/ebda-cli/config"
)

/*
$ notify-send
Usage:
  notify-send [OPTION…] <SUMMARY> [BODY] - create a notification

Help Options:
  -?, --help                        Show help options

Application Options:
  -u, --urgency=LEVEL               Specifies the urgency level (low, normal, critical).
  -t, --expire-time=TIME            Specifies the timeout in milliseconds at which to expire the notification.
  -a, --app-name=APP_NAME           Specifies the app name for the notification
  -i, --icon=ICON                   Specifies an icon filename or stock icon to display.
  -n, --app-icon=ICON               Specifies an application icon filename or app icon name. The server may or may not display it.
  -c, --category=TYPE[,TYPE...]     Specifies the notification category.
  -e, --transient                   Create a transient notification
  -h, --hint=TYPE:NAME:VALUE        Specifies basic extra data to pass. Valid types are boolean, int, double, string, byte and variant.
  -p, --print-id                    Print the notification ID.
  --id-fd                           File descriptor where to write the notification ID.
  -r, --replace-id=REPLACE_ID       The ID of the notification to replace.
  -w, --wait                        Wait for the notification to be closed before exiting.
  -A, --action=[NAME=]Text...       Specifies the actions to display to the user. Implies --wait to wait for user input. May be set multiple times. The name of the action is output to stdout. If NAME is not specified, the numerical index of the option is used (starting with 0).
  --selected-action-fd              File descriptor where to write the action chosen by the user.
  --activation-token-fd             File descriptor where to write the action activation token. The daemon must support it.
  -v, --version                     Version of the package.
*/

type UrgencyLevel int

const (
	UrgencyLow UrgencyLevel = iota
	UrgencyNormal
	UrgencyCritical
)

type CategoryType int

type Hint struct {
	Type  string
	Name  string
	Value string
}

type Notification struct {
	Summary    string
	Body       string
	Help       bool
	Urgency    UrgencyLevel
	ExpireTime int
	AppName    string
	Icon       string
	AppIcon    string
	// Category          CategoryType // TODO: (AhmedSat) search for CategoryType
	Transient bool
	Hints     []Hint
	// PrintID           string // TODO: (AhmedSat) search for PrintID
	// IdFd              int // TODO: (AhmedSat) search for IdFd
	// ReplaceID         int // TODO: (AhmedSat) search for ReplaceID
	Wait bool
	// Actions           []string  // TODO: (AhmedSat) lookup for a good way to handle actions
	// SelectedActionFd  int // TODO: (AhmedSat) search for SelectedActionFd
	// ActivationTokenFd int // TODO: (AhmedSat) search for ActivationTokenFd
	Version string
}

func (n *Notification) Run() (stdout, stderr string, err error) {

	if config.DisableNotify {
		return
	}

	name := "notify-send"
	args := []string{}

	if n.Help {
		name = "notify-send"
		args = append(args, "--help")
		stdout, stderr, err = runCmd(name, args)
		return
	}

	if n.Version != "" {
		name = "notify-send"
		args = append(args, "--version")
		stdout, stderr, err = runCmd(name, args)
		return
	}

	switch n.Urgency {
	case UrgencyLow:
		args = append(args, "--urgency", "low")
	case UrgencyNormal:
		args = append(args, "--urgency", "normal")
	case UrgencyCritical:
		args = append(args, "--urgency", "critical")
	}

	if n.ExpireTime > 0 {
		args = append(args, "--expire-time", fmt.Sprintf("%d", n.ExpireTime))
	}

	if n.AppName != "" {
		args = append(args, "--app-name", n.AppName)
	}

	if n.Icon != "" {
		args = append(args, "--icon", n.Icon)
	}

	if n.AppIcon != "" {
		args = append(args, "--app-icon", n.AppIcon)
	}

	// TODO: (AhmedSat) search for CategoryType

	if n.Transient {
		args = append(args, "--transient")
	}

	for _, hint := range n.Hints {
		args = append(args, "--hint", fmt.Sprintf("%s:%s:%s", hint.Type, hint.Name, hint.Value))
	}

	// TODO: (AhmedSat) search for PrintID

	// TODO: (AhmedSat) search for IdFd

	// TODO: (AhmedSat) search for ReplaceID

	if n.Wait {
		args = append(args, "--wait")
	}

	// TODO: (AhmedSat) lookup for a good way to handle actions

	// TODO: (AhmedSat) search for SelectedActionFd

	// TODO: (AhmedSat) search for ActivationTokenFd

	if n.Summary == "" {
		err = fmt.Errorf("summary is required")
		return
	}
	args = append(args, n.Summary)
	if n.Body != "" {
		args = append(args, n.Body)
	}

	stdout, stderr, err = runCmd(name, args)
	if err != nil {
		return
	}

	return
}

func NewProgressNotification(summary string, body string, id string, progress int) (*Notification, error) {

	fmt.Fprintf(os.Stderr, "\r%s = > (%d%%)", body, progress)

	if progress < 0 || progress > 100 {
		return nil, fmt.Errorf("progress must be between 0 and 100 inclusive got %d", progress)
	}

	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	if summary == "" {
		return nil, fmt.Errorf("summary is required")
	}

	timeout := 1000 * 60
	if progress == 100 {
		timeout = 0
	}

	return &Notification{
		Summary:    fmt.Sprintf("%q", summary),
		Body:       fmt.Sprintf("%q", body),
		Urgency:    UrgencyLow,
		ExpireTime: timeout,
		Hints: []Hint{
			{
				Type:  "int",
				Name:  "value",
				Value: fmt.Sprintf("%d", progress),
			},
			{
				Type:  "string",
				Name:  "synchronous",
				Value: id,
			},
		},
	}, nil
}

func runCmd(name string, args []string) (stdout, stderr string, err error) {

	ctx, canceled := context.WithTimeout(context.Background(), time.Microsecond)
	defer canceled()

	cmd := exec.CommandContext(ctx, name, args...)
	outBytes := []byte{}
	bufStdout := bytes.NewBuffer(outBytes)
	errBytes := []byte{}
	bufStderr := bytes.NewBuffer(errBytes)
	cmd.Stdout = bufStdout
	cmd.Stderr = bufStderr
	err = cmd.Run()
	stdout = bufStdout.String()
	stderr = bufStderr.String()
	if err != nil {
		return
	}
	return
}
