package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	expect "github.com/google/goexpect"
)

const (
	timeOut = time.Millisecond * 500
	modeRW  = "rw"
	modeRO  = "ro"
)

var (
	loginVal      = fmt.Sprintln("user")
	passVal       = fmt.Sprintln("password")
	roCommand     = fmt.Sprintln("mount -o ro,remount /mnt/sys")
	rwCommand     = fmt.Sprintln("mount -o rw,remount /mnt/sys")
	exitCommand   = fmt.Sprintln("exit")
	rebootCommand = fmt.Sprintln("reboot")
	userRE        = regexp.MustCompile("username:")
	passRE        = regexp.MustCompile("password:")
	promptRE      = regexp.MustCompile("%")
)

// RO структура, для работы в случае запроса для установления режима "read only"
type RO struct {
	ip string
}

// RW структура, для работы в случае запроса для установления режима "read write"
type RW struct {
	ip string
}

// newRO функция для генерации объектра RO
func newRO(ip string) *RO {
	return &RO{ip}
}

// newRO функция для генерации объектра Rw
func newRW(ip string) *RW {
	return &RW{ip}
}

type modeSetter interface {
	setmode() int
}

// setmode метод для объекта RO, устанавливает соединение, выставляет режим, перезагружает контроллер
func (r *RO) setmode() int {
	conn, _, err := expect.Spawn(fmt.Sprintf("telnet %s", r.ip), -1)
	if err != nil {
		log.Println(err)
		return 0
	}
	defer conn.Close()

	conn.Expect(userRE, timeOut)
	conn.Send(loginVal)
	conn.Expect(passRE, timeOut)
	conn.Send(passVal)
	conn.Expect(promptRE, timeOut)
	conn.Send(roCommand)
	conn.Expect(promptRE, timeOut)
	conn.Send(rebootCommand)
	conn.Expect(promptRE, timeOut)
	conn.Send(exitCommand)
	return 1
}

// setmode метод для объекта RW, устанавливает соединение, выставляет режим, отключается
func (r *RW) setmode() int {
	conn, _, err := expect.Spawn(fmt.Sprintf("telnet %s", r.ip), -1)
	if err != nil {
		log.Println(err)
		return 0
	}
	defer conn.Close()

	conn.Expect(userRE, timeOut)
	conn.Send(loginVal)
	conn.Expect(passRE, timeOut)
	conn.Send(passVal)
	conn.Expect(promptRE, timeOut)
	conn.Send(rwCommand)
	conn.Expect(promptRE, timeOut)
	conn.Send(exitCommand)

	return 1
}

func main() {
	ip := os.Args[2]
	mode := os.Args[1]

	switch mode {
	case modeRW:
		operate := newRW(ip)
		fmt.Fprintln(os.Stdout, operate.setmode())
	case modeRO:
		operate := newRO(ip)
		fmt.Fprintln(os.Stdout, operate.setmode())
	}
}
