package modem

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

var (
	ErrCreatingSMS = errors.New("sms: there was an error while creating the sms")
	ErrSendingSMS  = errors.New("sms: there was an error while sending the sms")
)

func GetAllModem() (*map[int]string, error) {
	cmd := exec.Command("mmcli", "-L")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}
	modems := strings.Split(out.String(), "\n")
	modems = modems[:len(modems)-1]
	var op = make(map[int]string)
	for _, v := range modems {
		infos := strings.Split(strings.TrimSpace(v), " ")
		paths := strings.Split(infos[0], "/")
		number, _ := strconv.Atoi(paths[len(paths)-1])
		op[number] = infos[len(infos)-1]
	}
	fmt.Println(op)
	return &op, nil
}

type Modem struct {
	Numero   int
	Operator string
}

func getSMSNumber(out *bytes.Buffer) (int, error) {
	s := strings.TrimSpace(out.String())
	s = strings.Replace(s, " ", "", -1)
	splitted := strings.Split(s, "/")
	number, _ := strconv.Atoi(splitted[len(splitted)-1])
	return number, nil
}

func (mod *Modem) SendSMS(sms string, num string) error {
	initiateCommand := fmt.Sprintf("--messaging-create-sms=number='%s',text='%s'", num, sms)
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("mmcli", "-m", strconv.Itoa(mod.Numero), initiateCommand)
	cmd.Stdin = strings.NewReader(initiateCommand)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return ErrCreatingSMS
	}
	fmt.Println(stdout.String())
	number, _ := getSMSNumber(&stdout)
	cmd = exec.Command("mmcli", "--sms", strconv.Itoa(number), "--send")
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(stderr.String())
		return ErrSendingSMS
	}
	return nil
}
