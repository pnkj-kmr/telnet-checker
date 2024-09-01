package internal

// TelnetChecker interface help to freeze
// the input and outout function definations
type TelnetChecker interface {
	GetInput() []Input
	ProduceOutput(<-chan Output, chan<- struct{})
}

type CmdParams struct {
	Expect  string `json:"expect,omitempty"`
	Command string `json:"command"`
	Eof     string `json:"eof,omitempty"`
}

// Input represent the input
type Input struct {
	Tag      string      `json:"tag,omitempty"`
	Host     string      `json:"host"`
	Port     int         `json:"port,omitempty"`
	Timeout  int         `json:"timeout,omitempty"`
	Commands []CmdParams `json:"commands"`
}

// Output represent the output
type Output struct {
	I   Input    `json:"input"`
	Err []string `json:"error,omitempty"`
	O   []string `json:"output,omitempty"`
}

// CmdPipe helps to define the Command line Agrs
type CmdPipe struct {
	Ifile         string
	Ofile         string
	Workers       int
	Port          int
	Timeout       int
	CircleTimeout int
}
