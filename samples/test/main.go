package main

import (
	"log"
	"os"
	"telnet-checker/internal"

	"time"
)

const timeout = 10 * time.Second

func main() {
	if len(os.Args) != 4 {
		log.Printf("Usage: %s {HOST:PORT USER PASSWD", os.Args[0])
		return
	}
	dst, user, passwd := os.Args[1], os.Args[2], os.Args[3]

	t, err := internal.Dial("tcp", dst)
	// checkErr(err)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	t.SetUnixWriteMode(true)

	var data []byte
	var data2 []byte
	t.Expect(timeout, "name: ")
	t.Sendln(nil, timeout, []byte(user))
	t.Expect(timeout, "ssword: ")
	t.Sendln(nil, timeout, []byte(passwd))
	t.Expect(timeout, "#")
	t.Sendln(nil, timeout, []byte("sh ver"))
	data, err = t.ReadUntil("19#")

	// // get configuration from device
	// t.Expect(timeout, "#")
	// t.Sendln(nil, timeout, []byte("show run"))
	// data2, _ = t.ReadUntil("019#")

	// checkErr(err)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	os.Stdout.Write(data)
	os.Stdout.WriteString("\n")
	os.Stdout.Write(data2)
	t.Close()
}
