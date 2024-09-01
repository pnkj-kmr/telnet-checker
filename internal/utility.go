package internal

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"
)

func GetCommandLineInputs() CmdPipe {
	iFile := flag.String("f", "input.json", "input file name")
	oFile := flag.String("o", "output.json", "output file name")
	noWorkers := flag.Int("w", 4, "number of workers")
	port := flag.Int("p", 23, "generic default port")
	timeout := flag.Int("t", 20, "timeout [secs]")
	circleTimeout := flag.Int("ct", 120, "timeout [secs]")

	flag.Parse()

	log.Println("File accepted:", *iFile, "| output file:", *oFile)
	log.Println("workers:", *noWorkers, "| port:", *port)
	log.Println("timeout:", *timeout, "| circle timeout:", *circleTimeout)

	return CmdPipe{
		Ifile: *iFile, Ofile: *oFile, Workers: *noWorkers, Port: *port,
		Timeout: *timeout, CircleTimeout: *circleTimeout,
	}
}

func run(input Input, cmdAgrs CmdPipe, ind int) (out Output) {
	// binding input if empty from json
	if input.Port == 0 {
		input.Port = cmdAgrs.Port
	}
	if input.Timeout == 0 {
		input.Timeout = cmdAgrs.Timeout
	}

	var timeout = time.Second * time.Duration(input.Timeout)
	var errors []string
	var result []string
	// doing the telnet to end device
	// log.Println("dialing...", fmt.Sprintf("%s:%d", input.Host, input.Port), timeout)
	conn, err := DialTimeout("tcp", fmt.Sprintf("%s:%d", input.Host, input.Port), timeout)
	if err != nil {
		// log.Println("err ...", err)
		errors = append(errors, err.Error())
	} else {
		// log.Println("connection created ...", input.Host)
		defer conn.Close()
		conn.SetUnixWriteMode(true)

		for _, cmd := range input.Commands {
			// log.Println("running ...", cmd)
			if cmd.Expect != "" {
				conn.Expect(timeout, cmd.Expect)
				conn.Sendln(nil, timeout, []byte(cmd.Command))
			} else {
				conn.Sendln(nil, timeout, []byte(cmd.Command))
			}
			if cmd.Eof != "" {
				data, err := conn.ReadUntil(cmd.Eof)
				if err != nil {
					errors = append(errors, err.Error())
				}
				result = append(result, string(data))
			}
		}
	}

	// forming result Output
	// log.Println("result ...", result)
	out.I = input
	out.Err = errors
	out.O = result
	log.Println("Completed with", input.Host, ind)
	return
}

func Execute(inputs []Input, ch chan<- Output, cmdAgrs CmdPipe) (err error) {
	var wg sync.WaitGroup
	c := make(chan int, cmdAgrs.Workers)
	for i := 0; i < len(inputs); i++ {
		wg.Add(1)
		c <- 1
		go func(input Input, cmdAgrs CmdPipe, ind int) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cmdAgrs.CircleTimeout))
			defer func() { wg.Done(); cancel(); <-c }()

			chInner := make(chan Output)
			go func() {
				defer close(chInner)
				_output := run(input, cmdAgrs, ind)
				chInner <- _output
			}()

			select {
			case _output := <-chInner:
				ch <- _output
			case <-ctx.Done():
				// close(chInner)
				var _output = Output{I: input, Err: []string{ctx.Err().Error()}}
				ch <- _output
				log.Println(ind, input.Host, "Terminate Error", ctx.Err())
			}

		}(inputs[i], cmdAgrs, i+1)
		log.Println("Host sent for telnet -- ", inputs[i], i+1)
	}
	wg.Wait()
	return nil
}
