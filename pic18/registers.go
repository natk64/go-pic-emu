package pic18

type RegisterTable struct {
	PCL     uint16
	PCLATH  uint16
	PCLATU  uint16
	STKPTR  uint16
	TOSL    uint16
	TOSH    uint16
	TOSU    uint16
	WREG    uint16
	STATUS  uint16
	FSR0H   uint16
	FSR0L   uint16
	FSR1H   uint16
	FSR1L   uint16
	FSR2H   uint16
	FSR2L   uint16
	BSR     uint16
	INTCON  uint16
	INTCON2 uint16
	INTCON3 uint16
	RCON    uint16
}

var Registers = RegisterTable{
	TOSU:    0xFFF,
	TOSH:    0xFFE,
	TOSL:    0xFFD,
	STKPTR:  0xFFC,
	PCLATU:  0xFFB,
	PCLATH:  0xFFA,
	PCL:     0xFF9,
	WREG:    0xFE8,
	STATUS:  0xFD8,
	FSR0H:   0xFEA,
	FSR0L:   0xFE9,
	FSR1H:   0xFE2,
	FSR1L:   0xFE1,
	FSR2H:   0xFDA,
	FSR2L:   0xFD9,
	BSR:     0xFE0,
	INTCON:  0xFF2,
	INTCON2: 0xFF1,
	INTCON3: 0xFF0,
	RCON:    0xFD0,
}
