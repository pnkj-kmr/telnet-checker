package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"time"
	"unicode"
)

// Conn implements net.Conn interface for Telnet protocol plus some set of
// Telnet specific methods.
type Conn struct {
	net.Conn
	r *bufio.Reader

	unixWriteMode bool

	cliSuppressGoAhead bool
	cliEcho            bool
}

func NewConn(conn net.Conn) (*Conn, error) {
	c := Conn{
		Conn: conn,
		r:    bufio.NewReaderSize(conn, 256),
	}
	return &c, nil
}

func Dial(network, addr string) (*Conn, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return NewConn(conn)
}

func DialTimeout(network, addr string, timeout time.Duration) (*Conn, error) {
	conn, err := net.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, err
	}
	return NewConn(conn)
}

// SetUnixWriteMode sets flag that applies only to the Write method.
// If set, Write converts any '\n' (LF) to '\r\n' (CR LF).
func (c *Conn) SetUnixWriteMode(uwm bool) {
	c.unixWriteMode = uwm
}

func (c *Conn) Close() error {
	return c.Conn.Close()
}

func (c *Conn) do(option byte) error {
	// log.Println("do:", option)
	_, err := c.Conn.Write([]byte{cmdIAC, cmdDo, option})
	return err
}

func (c *Conn) dont(option byte) error {
	// log.Println("dont:", option)
	_, err := c.Conn.Write([]byte{cmdIAC, cmdDont, option})
	return err
}

func (c *Conn) will(option byte) error {
	// log.Println("will:", option)
	_, err := c.Conn.Write([]byte{cmdIAC, cmdWill, option})
	return err
}

func (c *Conn) wont(option byte) error {
	// log.Println("wont:", option)
	_, err := c.Conn.Write([]byte{cmdIAC, cmdWont, option})
	return err
}

func (c *Conn) sub(option byte, data ...byte) error {
	// log.Println("sub:", data, option)
	if _, err := c.Conn.Write([]byte{cmdIAC, cmdSB, option}); err != nil {
		return err
	}
	if _, err := c.Conn.Write(data); err != nil {
		return err
	}
	_, err := c.Conn.Write([]byte{cmdIAC, cmdSE})
	return err
}

func (c *Conn) deny(cmd, option byte) (err error) {
	// log.Println("deny:", cmd, option)
	switch cmd {
	case cmdDo:
		err = c.wont(option)
	case cmdDont:
		// nop
	case cmdWill, cmdWont:
		err = c.dont(option)
	}
	return
}

func (c *Conn) skipSubnego() error {
	// log.Println("skipSubnego --")
	for {
		if b, err := c.r.ReadByte(); err != nil {
			return err
		} else if b == cmdIAC {
			if b, err = c.r.ReadByte(); err != nil {
				return err
			} else if b == cmdSE {
				return nil
			}
		}
	}
	// var data []byte
	// o, err := c.r.ReadByte()
	// if err != nil {
	// 	return err
	// }
	// for {
	// 	b, err := c.r.ReadByte()
	// 	if err != nil {
	// 		return errors.New("read IAC SE of IAC SB '" + strconv.FormatInt(int64(b), 10) + "' fail, " + err.Error())
	// 	}
	// 	if b != cmdIAC {
	// 		data = append(data, b)
	// 		continue
	// 	}

	// 	b, err = c.r.ReadByte()
	// 	if err != nil {
	// 		return errors.New("read IAC SE of IAC SB '" + strconv.FormatInt(int64(b), 10) + "' fail, " + err.Error())
	// 	}

	// 	if b == cmdSE {
	// 		break
	// 	}

	// 	data = append(data, cmdIAC)
	// 	data = append(data, b)
	// }

	// switch o {
	// case optTermType:
	// 	if len(data) == 1 && data[0] == 1 {
	// 		// IAC SB TERMINAL-TYPE IS xterm IAC SE
	// 		_, err = c.Conn.Write([]byte{cmdIAC, cmdSB, optTermType, 0, 'x', 't', 'e', 'r', 'm', cmdIAC, cmdSE})
	// 		return err
	// 	}
	// }
	// return nil
}

func (c *Conn) cmd(cmd byte) error {
	switch cmd {
	case cmdGA:
		return nil
	case cmdDo, cmdDont, cmdWill, cmdWont:
		// Process cmd after this switch.
	case cmdSB:
		return c.skipSubnego()
	default:
		return fmt.Errorf("unknown command: %d", cmd)
	}
	// Read an option
	o, err := c.r.ReadByte()
	if err != nil {
		return err
	}
	//log.Println("received cmd:", cmd, o)
	switch o {
	case optEcho:
		// Accept any echo configuration.
		switch cmd {
		case cmdDo:
			if !c.cliEcho {
				c.cliEcho = true
				err = c.will(o)
			}
		case cmdDont:
			if c.cliEcho {
				c.cliEcho = false
				err = c.wont(o)
			}
		case cmdWill:
			if !c.cliEcho {
				c.cliEcho = true
				err = c.do(o)
			}
		case cmdWont:
			if c.cliEcho {
				c.cliEcho = false
				err = c.dont(o)
			}
		}
	case optSuppressGoAhead:
		// We don't use GA so can allways accept every configuration
		switch cmd {
		case cmdDo:
			if !c.cliSuppressGoAhead {
				c.cliSuppressGoAhead = true
				err = c.will(o)
			}
		case cmdDont:
			if c.cliSuppressGoAhead {
				c.cliSuppressGoAhead = false
				err = c.wont(o)
			}
		case cmdWill:
			if !c.cliSuppressGoAhead {
				c.cliSuppressGoAhead = true
				err = c.do(o)
			}
		case cmdWont:
			if c.cliSuppressGoAhead {
				c.cliSuppressGoAhead = false
				err = c.dont(o)
			}
		}
	// case optTermType:
	// 	// Accept any echo configuration.
	// 	switch cmd {
	// 	case cmdDo:
	// 		err = c.will(o)
	// 	case cmdDont:
	// 	case cmdWill, cmdWont:
	// 		err = c.dont(o)
	// 	}
	case optWndSize:
		if cmd != cmdDo {
			err = c.deny(cmd, o)
			break
		}
		if err = c.will(o); err != nil {
			break
		}
		// Reply with max window size: 65535x65535
		err = c.sub(o, 255, 255, 255, 255)
		// // if custom term size needed
		// if cmd == cmdDo {
		// 	_, err = c.Conn.Write([]byte{cmdIAC, cmdSB, optWndSize, 0, c.columns, 0, c.rows, cmdIAC, cmdSE})
		// }
	default:
		// Deny any other option
		err = c.deny(cmd, o)
	}
	return err
}

func (c *Conn) tryReadByte() (b byte, retry bool, err error) {
	b, err = c.r.ReadByte()
	if err != nil || b != cmdIAC {
		return
	}
	b, err = c.r.ReadByte()
	if err != nil {
		return
	}
	if b != cmdIAC {
		err = c.cmd(b)
		if err != nil {
			return
		}
		retry = true
	}
	return
}

// // SetEcho tries to enable/disable echo on server side. Typically telnet
// // servers doesn't support this.
// func (c *Conn) SetEcho(echo bool) error {
// 	if echo {
// 		return c.do(optEcho)
// 	}
// 	return c.dont(optEcho)
// }

// ReadByte works like bufio.ReadByte
func (c *Conn) ReadByte() (b byte, err error) {
	retry := true
	for retry && err == nil {
		b, retry, err = c.tryReadByte()
	}
	return
}

// ReadRune works like bufio.ReadRune
func (c *Conn) ReadRune() (r rune, size int, err error) {
loop:
	r, size, err = c.r.ReadRune()
	if err != nil {
		return
	}
	if r != unicode.ReplacementChar || size != 1 {
		// Properly readed rune
		return
	}
	// Bad rune
	err = c.r.UnreadRune()
	if err != nil {
		return
	}
	// Read telnet command or escaped IAC
	_, retry, err := c.tryReadByte()
	if err != nil {
		return
	}
	if retry {
		// This bad rune was a begining of telnet command. Try read next rune.
		goto loop
	}
	// Return escaped IAC as unicode.ReplacementChar
	return
}

// Read is for implement an io.Reader interface
func (c *Conn) Read(buf []byte) (int, error) {
	var n int
	for n < len(buf) {
		b, retry, err := c.tryReadByte()
		if err != nil {
			return n, err
		}
		if !retry {
			buf[n] = b
			n++
		}
		if n > 0 && c.r.Buffered() == 0 {
			// Don't block if can't return more data.
			return n, err
		}
	}
	return n, nil
}

// ReadBytes works like bufio.ReadBytes
func (c *Conn) ReadBytes(delim byte) ([]byte, error) {
	var line []byte
	for {
		b, err := c.ReadByte()
		if err != nil {
			return nil, err
		}
		line = append(line, b)
		if b == delim {
			break
		}
	}
	return line, nil
}

// SkipBytes works like ReadBytes but skips all read data.
func (c *Conn) SkipBytes(delim byte) error {
	for {
		b, err := c.ReadByte()
		if err != nil {
			return err
		}
		if b == delim {
			break
		}
	}
	return nil
}

// ReadString works like bufio.ReadString
func (c *Conn) ReadString(delim byte) (string, error) {
	bytes, err := c.ReadBytes(delim)
	return string(bytes), err
}

func (c *Conn) readUntil(read bool, delims ...string) ([]byte, int, error) {
	if len(delims) == 0 {
		return nil, 0, nil
	}
	p := make([]string, len(delims))
	for i, s := range delims {
		if len(s) == 0 {
			return nil, 0, nil
		}
		p[i] = s
	}
	var line []byte
	for {
		b, err := c.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		if read {
			line = append(line, b)
		}
		for i, s := range p {
			if s[0] == b {
				if len(s) == 1 {
					return line, i, nil
				}
				p[i] = s[1:]
			} else {
				p[i] = delims[i]
			}
		}
	}
	// panic(nil)
}

// ReadUntilIndex reads from connection until one of delimiters occurs. Returns
// read data and an index of delimiter or error.
func (c *Conn) ReadUntilIndex(delims ...string) ([]byte, int, error) {
	return c.readUntil(true, delims...)
}

// ReadUntil works like ReadUntilIndex but don't return a delimiter index.
func (c *Conn) ReadUntil(delims ...string) ([]byte, error) {
	d, _, err := c.readUntil(true, delims...)
	return d, err
}

// SkipUntilIndex works like ReadUntilIndex but skips all read data.
func (c *Conn) SkipUntilIndex(delims ...string) (int, error) {
	_, i, err := c.readUntil(false, delims...)
	return i, err
}

// SkipUntil works like ReadUntil but skips all read data.
func (c *Conn) SkipUntil(delims ...string) error {
	_, _, err := c.readUntil(false, delims...)
	return err
}

func (c *Conn) Expect(timeout time.Duration, delims ...string) ([]byte, int, error) {
	if e := c.SetReadDeadline(time.Now().Add(timeout)); nil != e {
		return nil, 0, e
	}
	return c.readUntil(true, delims...)
}

func (c *Conn) Sendln(buf *bytes.Buffer, timeout time.Duration, s []byte) error {
	if e := c.SetWriteDeadline(time.Now().Add(timeout)); nil != e {
		return e
	}

	copy_buffer := s
	if !bytes.HasSuffix(s, []byte("\n")) {
		copy_buffer = make([]byte, len(s)+1)
		copy(copy_buffer, s)
		copy_buffer[len(s)] = '\n'
	}

	if nil != buf {
		buf.Write(copy_buffer)
	}
	_, err := c.Write(copy_buffer)
	return err
}

func (c *Conn) Send(buf *bytes.Buffer, timeout time.Duration, s []byte) error {
	if e := c.SetWriteDeadline(time.Now().Add(timeout)); nil != e {
		return e
	}

	if nil != buf {
		buf.Write(s)
	}
	_, err := c.Write(s)
	return err
}

// Write is for implement an io.Writer interface
func (c *Conn) Write(buf []byte) (int, error) {
	search := "\xff"
	if c.unixWriteMode {
		search = "\xff\n"
	}
	var (
		n   int
		err error
	)
	for len(buf) > 0 {
		var k int
		i := bytes.IndexAny(buf, search)
		if i == -1 {
			k, err = c.Conn.Write(buf)
			n += k
			break
		}
		k, err = c.Conn.Write(buf[:i])
		n += k
		if err != nil {
			break
		}
		switch buf[i] {
		case LF:
			k, err = c.Conn.Write([]byte{CR, LF})
		case cmdIAC:
			k, err = c.Conn.Write([]byte{cmdIAC, cmdIAC})
		}
		n += k
		if err != nil {
			break
		}
		buf = buf[i+1:]
	}
	return n, err
}
