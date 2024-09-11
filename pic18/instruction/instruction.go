package instruction

type Opcode uint8
type Instruction uint16

const (
	ILLEGAL Opcode = iota
	ADDWF
	ADDWFC
	ANDWF
	CLRF
	COMF
	CPFSEQ
	CPFSGT
	CPFSLT
	DECF
	DECFSZ
	DCFSNZ
	INCF
	INCFSZ
	INFSNZ
	IORWF
	MOVF
	MOVFF
	MOVWF
	MULWF
	NEGF
	RLCF
	RLNCF
	RRCF
	RRNCF
	SETF
	SUBFWB
	SUBWF
	SUBWFB
	SWAPF
	TSTFSZ
	XORWF
	BCF
	BSF
	BTFSC
	BTFSS
	BTG
	BC
	BN
	BNC
	BNN
	BNOV
	BNZ
	BOV
	BRA
	BZ
	CALL
	CLRWDT
	DAW
	GOTO
	NOP
	// NOP1 is executed the same as NOP.
	// This instruction is a result of executing the second word in a 2 word instruction.
	NOP1
	POP
	PUSH
	RCALL
	RESET
	RETFIE
	RETLW
	RETURN
	SLEEP
	ADDLW
	ANDLW
	IORLW
	LFSR
	MOVLB
	MOVLW
	MULLW
	SUBLW
	XORLW
	TBLRD
	TBLWT
	ADDFSR
	ADDULNK
	CALLW
	MOVSF
	MOVSS
	PUSHL
	SUBFSR
	SUBULNK
)

func (raw Instruction) Opcode() Opcode {
	if (raw & 0b1111111111111111) == 0b0000000000000100 {
		return CLRWDT
	} else if (raw & 0b1111111111111111) == 0b0000000000000111 {
		return DAW
	} else if (raw & 0b1111111111111111) == 0b0000000000000000 {
		return NOP
	} else if (raw & 0b1111111111111111) == 0b0000000000000110 {
		return POP
	} else if (raw & 0b1111111111111111) == 0b0000000000000101 {
		return PUSH
	} else if (raw & 0b1111111111111111) == 0b0000000011111111 {
		return RESET
	} else if (raw & 0b1111111111111111) == 0b0000000000000011 {
		return SLEEP
	} else if (raw & 0b1111111111111111) == 0b0000000000010100 {
		return CALLW
	} else if (raw & 0b1111111111111110) == 0b0000000000010000 {
		return RETFIE
	} else if (raw & 0b1111111111111110) == 0b0000000000010010 {
		return RETURN
	} else if (raw & 0b1111111111111100) == 0b0000000000001000 {
		return TBLRD
	} else if (raw & 0b1111111111111100) == 0b0000000000001100 {
		return TBLWT
	} else if (raw & 0b1111111111110000) == 0b0000000100000000 {
		return MOVLB
	} else if (raw & 0b1111111111000000) == 0b1110111000000000 {
		return LFSR
	} else if (raw & 0b1111111111000000) == 0b1110100011000000 {
		return ADDULNK
	} else if (raw & 0b1111111111000000) == 0b1110100111000000 {
		return SUBULNK
	} else if (raw & 0b1111111110000000) == 0b1110101100000000 {
		return MOVSF
	} else if (raw & 0b1111111110000000) == 0b1110101110000000 {
		return MOVSS
	} else if (raw & 0b1111111100000000) == 0b1110001000000000 {
		return BC
	} else if (raw & 0b1111111100000000) == 0b1110011000000000 {
		return BN
	} else if (raw & 0b1111111100000000) == 0b1110001100000000 {
		return BNC
	} else if (raw & 0b1111111100000000) == 0b1110011100000000 {
		return BNN
	} else if (raw & 0b1111111100000000) == 0b1110010100000000 {
		return BNOV
	} else if (raw & 0b1111111100000000) == 0b1110000100000000 {
		return BNZ
	} else if (raw & 0b1111111100000000) == 0b1110010000000000 {
		return BOV
	} else if (raw & 0b1111111100000000) == 0b1110000000000000 {
		return BZ
	} else if (raw & 0b1111111100000000) == 0b1110111100000000 {
		return GOTO
	} else if (raw & 0b1111111100000000) == 0b0000110000000000 {
		return RETLW
	} else if (raw & 0b1111111100000000) == 0b0000111100000000 {
		return ADDLW
	} else if (raw & 0b1111111100000000) == 0b0000101100000000 {
		return ANDLW
	} else if (raw & 0b1111111100000000) == 0b0000100100000000 {
		return IORLW
	} else if (raw & 0b1111111100000000) == 0b0000111000000000 {
		return MOVLW
	} else if (raw & 0b1111111100000000) == 0b0000110100000000 {
		return MULLW
	} else if (raw & 0b1111111100000000) == 0b0000100000000000 {
		return SUBLW
	} else if (raw & 0b1111111100000000) == 0b0000101000000000 {
		return XORLW
	} else if (raw & 0b1111111100000000) == 0b1110100000000000 {
		return ADDFSR
	} else if (raw & 0b1111111100000000) == 0b1110101000000000 {
		return PUSHL
	} else if (raw & 0b1111111100000000) == 0b1110100100000000 {
		return SUBFSR
	} else if (raw & 0b1111111000000000) == 0b0110001000000000 {
		return CPFSEQ
	} else if (raw & 0b1111111000000000) == 0b0110010000000000 {
		return CPFSGT
	} else if (raw & 0b1111111000000000) == 0b0110000000000000 {
		return CPFSLT
	} else if (raw & 0b1111111000000000) == 0b0110101000000000 {
		return CLRF
	} else if (raw & 0b1111111000000000) == 0b0110111000000000 {
		return MOVWF
	} else if (raw & 0b1111111000000000) == 0b0000001000000000 {
		return MULWF
	} else if (raw & 0b1111111000000000) == 0b0110110000000000 {
		return NEGF
	} else if (raw & 0b1111111000000000) == 0b0110100000000000 {
		return SETF
	} else if (raw & 0b1111111000000000) == 0b0110011000000000 {
		return TSTFSZ
	} else if (raw & 0b1111111000000000) == 0b1110110000000000 {
		return CALL
	} else if (raw & 0b1111110000000000) == 0b0010010000000000 {
		return ADDWF
	} else if (raw & 0b1111110000000000) == 0b0010000000000000 {
		return ADDWFC
	} else if (raw & 0b1111110000000000) == 0b0001010000000000 {
		return ANDWF
	} else if (raw & 0b1111110000000000) == 0b0001110000000000 {
		return COMF
	} else if (raw & 0b1111110000000000) == 0b0000010000000000 {
		return DECF
	} else if (raw & 0b1111110000000000) == 0b0010110000000000 {
		return DECFSZ
	} else if (raw & 0b1111110000000000) == 0b0100110000000000 {
		return DCFSNZ
	} else if (raw & 0b1111110000000000) == 0b0010100000000000 {
		return INCF
	} else if (raw & 0b1111110000000000) == 0b0011110000000000 {
		return INCFSZ
	} else if (raw & 0b1111110000000000) == 0b0100100000000000 {
		return INFSNZ
	} else if (raw & 0b1111110000000000) == 0b0001000000000000 {
		return IORWF
	} else if (raw & 0b1111110000000000) == 0b0101000000000000 {
		return MOVF
	} else if (raw & 0b1111110000000000) == 0b0011010000000000 {
		return RLCF
	} else if (raw & 0b1111110000000000) == 0b0100010000000000 {
		return RLNCF
	} else if (raw & 0b1111110000000000) == 0b0011000000000000 {
		return RRCF
	} else if (raw & 0b1111110000000000) == 0b0100000000000000 {
		return RRNCF
	} else if (raw & 0b1111110000000000) == 0b0101010000000000 {
		return SUBFWB
	} else if (raw & 0b1111110000000000) == 0b0101110000000000 {
		return SUBWF
	} else if (raw & 0b1111110000000000) == 0b0101100000000000 {
		return SUBWFB
	} else if (raw & 0b1111110000000000) == 0b0011100000000000 {
		return SWAPF
	} else if (raw & 0b1111110000000000) == 0b0001100000000000 {
		return XORWF
	} else if (raw & 0b1111100000000000) == 0b1101000000000000 {
		return BRA
	} else if (raw & 0b1111100000000000) == 0b1101100000000000 {
		return RCALL
	} else if (raw & 0b1111000000000000) == 0b1100000000000000 {
		return MOVFF
	} else if (raw & 0b1111000000000000) == 0b1001000000000000 {
		return BCF
	} else if (raw & 0b1111000000000000) == 0b1000000000000000 {
		return BSF
	} else if (raw & 0b1111000000000000) == 0b1011000000000000 {
		return BTFSC
	} else if (raw & 0b1111000000000000) == 0b1010000000000000 {
		return BTFSS
	} else if (raw & 0b1111000000000000) == 0b0111000000000000 {
		return BTG
	} else if (raw & 0b1111000000000000) == 0b1111000000000000 {
		return NOP1
	} else {
		return ILLEGAL
	}
}

type (
	ByteOriented        Instruction
	ByteToByte          Instruction
	BitOriented         Instruction
	Literal             Instruction
	ControlGotoHigh     Instruction
	ControlGotoLow      Instruction
	ControlCallHigh     Instruction
	ControlCallLow      Instruction
	ControlBranch       Instruction
	ControlBranchStatus Instruction
	ControlReturn       Instruction
	LoadFsrHigh         Instruction
	LoadFsrLow          Instruction
	TableOp             Instruction
	XinstFsr            Instruction
	MovsfHighMovss      Instruction
	MovsfLow            Instruction
)

func (inst ByteOriented) F() uint8 {
	return uint8(inst & 0x00FF)
}

func (inst ByteOriented) A() bool {
	return inst&0x0100 != 0
}

func (inst ByteOriented) D() bool {
	return inst&0x0200 != 0
}

func (inst ByteToByte) F() uint16 {
	return uint16(inst & 0x0FFF)
}

func (inst BitOriented) F() uint8 {
	return uint8(inst & 0x00FF)
}

func (inst BitOriented) A() bool {
	return inst&0x0100 != 0
}

func (inst BitOriented) Bit() uint8 {
	return uint8((inst & 0x0E00) >> 11)
}

func (inst Literal) K() uint8 {
	return uint8(inst & 0x00FF)
}

func (inst ControlGotoHigh) Literal() uint8 {
	return uint8(inst & 0x00FF)
}

func (inst ControlGotoLow) Literal() uint16 {
	return uint16(inst & 0x0FFF)
}

func (inst ControlCallHigh) Literal() uint8 {
	return uint8(inst & 0x00FF)
}

func (inst ControlCallHigh) S() bool {
	return inst&0x0100 != 0
}

func (inst ControlCallLow) Literal() uint16 {
	return uint16(inst & 0x0FFF)
}

func (inst ControlBranch) Literal() uint16 {
	return uint16(inst & 0x07FF)
}

func (inst ControlBranchStatus) Literal() uint8 {
	return uint8(inst & 0x00FF)
}

func (inst ControlReturn) S() bool {
	return inst&0x0001 != 0
}

func (inst LoadFsrHigh) K() uint8 {
	return uint8(inst & 0x000F)
}

func (inst LoadFsrHigh) F() uint8 {
	return uint8(inst&0x0030) >> 4
}

func (inst LoadFsrLow) K() uint8 {
	return uint8(inst & 0x00FF)
}

func (inst TableOp) N() uint8 {
	return uint8(inst & 0x0003)
}

func (inst XinstFsr) K() uint8 {
	return uint8(inst & 0x003F)
}

func (inst XinstFsr) F() uint8 {
	return uint8(inst&0x00C0) >> 6
}

func (inst MovsfHighMovss) Z() uint8 {
	return uint8(inst & 0x007F)
}

func (inst MovsfLow) F() uint16 {
	return uint16(inst & 0x0FFF)
}
