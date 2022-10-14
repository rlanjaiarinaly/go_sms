package modem

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Get the SMS number from the output of the sms creating command
func getSMSNumber(out *bytes.Buffer) (int, error) {
	s := strings.TrimSpace(out.String())
	s = strings.Replace(s, " ", "", -1)
	splitted := strings.Split(s, "/")
	number, err := strconv.Atoi(splitted[len(splitted)-1])
	if err != nil {
		return -1, ErrFetchingSMS
	}
	return number, nil
}

// The SendSMS function takes as argument two strings and returns an error if there is at least one
// Otherwise, it sends the message to the corresponding recipients
func (mod *Modem) SendSMSToRecipient(sms string, num string) error {
	sms = strings.Replace(sms, "\"", "''", -1)
	initiateCommand := fmt.Sprintf("--messaging-create-sms=number=%s,text=\"%s\"", num, sms)
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("mmcli", "-m", strconv.Itoa(mod.Numero), initiateCommand)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return ErrCreatingSMS
	}
	number, err := getSMSNumber(&stdout)
	if err != nil {
		return err
	}
	err = mod.SendSMS(number)
	if err != nil {
		return ErrSendingSMS
	}
	return nil
}

// Send the index specified SMS
func (mod *Modem) SendSMS(index int) error {
	cmd := exec.Command("mmcli", "--sms", strconv.Itoa(index), "--send")
	return cmd.Run()
}
