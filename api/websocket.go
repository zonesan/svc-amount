package api

import (
	"crypto/tls"
	"io"
	"strings"
	"sync"

	"github.com/zonesan/clog"
	"golang.org/x/net/websocket"
)

var done = make(chan bool)
var wg sync.WaitGroup

func inLoop(ws *websocket.Conn, errors chan<- error, in chan<- []byte) {
	var msg = make([]byte, 512)

	for {
		var n int
		var err error

		n, err = ws.Read(msg)

		if err != nil {
			errors <- err
			if err == io.EOF {
				clog.Debug("inloop end here.")
				wg.Done()
				return
			} else {
				continue
			}
		}

		in <- msg[:n]
	}
}

func processErrors(errors <-chan error) {
	for err := range errors {
		if err == io.EOF {
			clog.Warn("connection closed by remote caused by:", err)
			wg.Done()
			return
			// os.Exit(0)
		} else {
			clog.Error(err)
		}
	}
	clog.Debug("return here.")
}

func processReceivedMessages(in <-chan []byte) *svcAmount {
	for msg := range in {
		// I have no idea why msg[0] is 0x01, I just ignore it.
		msgStr := string(msg[1:])
		if len(msgStr) > 0 {
			ss := strings.Split(msgStr, "\n")
			clog.Debugf("%#v", ss)
			for _, s := range ss {
				if strings.Contains(s, "/run/secrets") {
					result := strings.Fields(s)
					clog.Trace(len(result), result)

					clog.Debugf("size: %v, used: %v, available: %v", result[1], result[2], result[3])
					amount := &svcAmount{Name: "volume", Size: result[1], Used: result[2], Available: result[3]}
					defer wg.Done()
					return amount
				}
			}
		}
	}
	return nil
}

// func outLoop(ws *websocket.Conn, out <-chan []byte, errors chan<- error) {
// 	for msg := range out {
// 		_, err := ws.Write(msg)
// 		if err != nil {
// 			errors <- err
// 		}
// 	}
// 	clog.Debug("out here.")
// }

func dial(url, protocol, origin string) (ws *websocket.Conn, err error) {
	config, err := websocket.NewConfig(url, origin)
	if err != nil {
		return nil, err
	}
	if protocol != "" {
		config.Protocol = []string{protocol}
	}
	config.TlsConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	return websocket.DialConfig(config)
}

func ws(url, origin string) (interface{}, error) {

	protocol := ""
	ws, err := dial(url, protocol, origin)

	if protocol != "" {
		clog.Debugf("connecting to %s via %s from %s...", url, protocol, origin)
	} else {
		clog.Debugf("connecting to %s from %s...", url, origin)
	}

	if err != nil {
		clog.Error(err)
		return nil, err
	}
	defer ws.Close()

	clog.Debugf("successfully connected to %s", url)

	wg.Add(3)

	errors := make(chan error)
	in := make(chan []byte)
	// out := make(chan []byte)

	defer close(errors)
	// defer close(out)
	defer close(in)

	go inLoop(ws, errors, in)
	// go processReceivedMessages(in)
	go processErrors(errors)
	// go outLoop(ws, out, errors)

	// scanner := bufio.NewScanner(os.Stdin)

	// for scanner.Scan() {
	// 	out <- []byte(scanner.Text())
	// }

	amount := processReceivedMessages(in)

	wg.Wait()

	return amount, nil
}
