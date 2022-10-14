package modem

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Modem struct {
	Numero   int
	Operator string
}

var (
	ErrCreatingSMS = errors.New("sms: there was an error while creating the sms")
	ErrSendingSMS  = errors.New("sms: there was an error while sending the sms")
	ErrNoOperator  = errors.New("sim: the operator can't be found")
	ErrFetchingSMS = errors.New("sms: failed fetching the sms")
)

func GetAllModem() (map[string]Modem, error) {
	cmd := exec.Command("mmcli", "-L")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}
	modems := strings.Split(out.String(), "\n")
	modems = modems[:len(modems)-1]
	var op = make(map[string]Modem)
	for _, v := range modems {
		infos := strings.Split(strings.TrimSpace(v), " ")
		paths := strings.Split(infos[0], "/")
		number, _ := strconv.Atoi(paths[len(paths)-1])
		cmd := exec.Command("mmcli", "-m", strconv.Itoa(number))
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		if err := cmd.Run(); err != nil {
			continue
		}
		operator, err := getOperatorName(&stdout)
		if err != nil {
			continue
		}
		op[strings.ToUpper(operator)] = Modem{Numero: number, Operator: operator}
	}
	return op, nil
}

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

func getOperatorName(out *bytes.Buffer) (string, error) {
	regex := regexp.MustCompile("operator name: (?P<deb>[a-zA-Z]+)")
	match := regex.FindStringSubmatch(out.String())
	if len(match) < 2 {
		return "", ErrNoOperator
	}
	return match[1], nil
}

func (mod *Modem) SendSMS(sms string, num string) error {
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
	cmd = exec.Command("mmcli", "--sms", strconv.Itoa(number), "--send")
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return ErrSendingSMS
	}
	return nil
}
