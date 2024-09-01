package internal

const (
	CR = byte('\r')
	LF = byte('\n')
)

const (
	// SE                  240    End of subnegotiation parameters.
	cmdSE = 240
	// NOP                 241    No operation.
	cmdNOP = 241
	// Data Mark           242    The data stream portion of a Synch.
	//                            This should always be accompanied
	//                            by a TCP Urgent notification.
	cmdData = 242

	// Break               243    NVT character BRK.
	cmdBreak = 243
	// Interrupt Process   244    The function IP.
	cmdIP = 244
	// Abort output        245    The function AO.
	cmdAO = 245
	// Are You There       246    The function AYT.
	cmdAYT = 246
	// Erase character     247    The function EC.
	cmdEC = 247
	// Erase Line          248    The function EL.
	cmdEL = 248
	// Go ahead            249    The GA signal.
	cmdGA = 249
	// SB                  250    Indicates that what follows is
	//                            subnegotiation of the indicated
	//                            option.
	cmdSB = 250 // FA

	// WILL (option code)  251    Indicates the desire to begin
	//                            performing, or confirmation that
	//                            you are now performing, the
	//                            indicated option.
	cmdWill = 251 // FB
	// WON'T (option code) 252    Indicates the refusal to perform,
	//                            or continue performing, the
	//                            indicated option.
	cmdWont = 252 // FC
	// DO (option code)    253    Indicates the request that the
	//                            other party perform, or
	//                            confirmation that you are expecting
	//                            the other party to perform, the
	//                            indicated option.
	cmdDo = 253 // FD
	// DON'T (option code) 254    Indicates the demand that the
	//                            other party stop performing,
	//                            or confirmation that you are no
	//                            longer expecting the other party
	//                            to perform, the indicated option.
	cmdDont = 254 // FE

	// IAC                 255    Data Byte 255.
	cmdIAC = 255 //FF

)

const (
	// 1(0x01)    echo
	optEcho = 1
	// 3(0x03)   Suppress continuation (this option can be selected for transmission one character at a time)
	optSuppressGoAhead = 3
	// // 24(0x18)  terminal type
	// optTermType = 24
	// 31(0x1F)   window size
	optWndSize = 31
	// // 32(0x20)   terminal rate
	// optRate = 32
	// 33(0x21)   Remote flow control
	// 34(0x22)   way
	// 36(0x24)   environment variables
)
