package log

import "fmt"

var (
	errChan = make(chan error)
	msgCHan = make(chan string)
	err     error
	msg     string
)

func init() {
	startLog()
}

func startLog() {
	go func(errChan <-chan error, msgCHan <-chan string) {
		for {
			select {
			case msg = <-msgCHan:
				fmt.Println(msg)
			case err = <-errChan:
				fmt.Println(err)
			}
		}
	}(errChan, msgCHan)
}

func GetLogChan() (chan<- error, chan<- string) {
	return errChan, msgCHan
}
