package utils

import (
	"bufio"
	"os"

	"github.com/dnote/cli/log"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

// GenerateUUID returns a uid
func GenerateUUID() string {
	return uuid.NewV4().String()
}

func getInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", errors.Wrap(err, "reading stdin")
	}

	return input, nil
}

// AskConfirmation prompts for user input to confirm a choice
func AskConfirmation(question string, optimistic bool) (bool, error) {
	var choices string
	if optimistic {
		choices = "(Y/n)"
	} else {
		choices = "(y/N)"
	}

	log.Printf("%s %s: ", question, choices)

	res, err := getInput()
	if err != nil {
		return false, errors.Wrap(err, "Failed to get user input")
	}

	confirmed := res == "y\n" || res == "y\r\n"

	if optimistic {
		confirmed = confirmed || res == "\n" || res == "\r\n"
	}

	return confirmed, nil
}
