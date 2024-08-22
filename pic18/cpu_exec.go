package pic18

import "github.com/natk64/go-pic-emu/pic18/instruction"

func (cpu *CPU) ExecuteInstruction(inst instruction.Instruction, extendedSet bool) bool {
	opcode := inst.Opcode()
	if opcode == instruction.ILLEGAL {
		return false
	}

	switch opcode {
	case instruction.ADDFSR:
		if extendedSet {
			cpu.execADDFSR(inst)
		}
	case instruction.ADDLW:
		cpu.execADDLW(inst)
	case instruction.ADDWF:
		cpu.execADDWF(inst)
	case instruction.ADDWFC:
		cpu.execADDWFC(inst)
	case instruction.ANDLW:
		cpu.execANDLW(inst)
	case instruction.ANDWF:
		cpu.execANDWF(inst)
	case instruction.BC:
		cpu.execBC(inst)
	case instruction.BCF:
		cpu.execBCF(inst)
	case instruction.BN:
		cpu.execBN(inst)
	case instruction.BNC:
		cpu.execBNC(inst)
	case instruction.BNN:
		cpu.execBNN(inst)
	case instruction.BNOV:
		cpu.execBNOV(inst)
	case instruction.BNZ:
		cpu.execBNZ(inst)
	case instruction.BRA:
		cpu.execBRA(inst)
	case instruction.BSF:
		cpu.execBSF(inst)
	case instruction.BTFSC:
		cpu.execBTFSC(inst)
	case instruction.BTFSS:
		cpu.execBTFSS(inst)
	case instruction.BTG:
		cpu.execBTG(inst)
	case instruction.BOV:
		cpu.execBOV(inst)
	case instruction.BZ:
		cpu.execBZ(inst)
	case instruction.CALL:
		cpu.execCALL(inst)
	case instruction.CALLW:
		if extendedSet {
			cpu.execCALLW(inst)
		}
	case instruction.CLRF:
		cpu.execCLRF(inst)
	case instruction.CLRWDT:
		cpu.execCLRWDT(inst)
	case instruction.COMF:
		cpu.execCOMF(inst)
	case instruction.CPFSEQ:
		cpu.execCPFSEQ(inst)
	case instruction.CPFSGT:
		cpu.execCPFSGT(inst)
	case instruction.CPFSLT:
		cpu.execCPFSLT(inst)
	case instruction.DAW:
		cpu.execDAW(inst)
	case instruction.DECF:
		cpu.execDECF(inst)
	case instruction.DECFSZ:
		cpu.execDECFSZ(inst)
	case instruction.DCFSNZ:
		cpu.execDCFSNZ(inst)
	case instruction.GOTO:
		cpu.execGOTO(inst)
	case instruction.INCF:
		cpu.execINCF(inst)
	case instruction.INCFSZ:
		cpu.execINCFSZ(inst)
	case instruction.INFSNZ:
		cpu.execINFSNZ(inst)
	case instruction.IORLW:
		cpu.execIORLW(inst)
	case instruction.IORWF:
		cpu.execIORWF(inst)
	case instruction.LFSR:
		cpu.execLFSR(inst)
	case instruction.MOVF:
		cpu.execMOVF(inst)
	case instruction.MOVFF:
		cpu.execMOVFF(inst)
	case instruction.MOVLB:
		cpu.execMOVLB(inst)
	case instruction.MOVLW:
		cpu.execMOVLW(inst)
	case instruction.MOVSF:
		if !extendedSet {
			cpu.execMOVSF(inst)
		}
	case instruction.MOVSS:
		if !extendedSet {
			cpu.execMOVSS(inst)
		}
	case instruction.MOVWF:
		cpu.execMOVWF(inst)
	case instruction.MULLW:
		cpu.execMULLW(inst)
	case instruction.MULWF:
		cpu.execMULWF(inst)
	case instruction.NEGF:
		cpu.execNEGF(inst)
	case instruction.NOP:
		break
	case instruction.NOP1:
		break
	case instruction.POP:
		cpu.execPOP(inst)
	case instruction.PUSH:
		cpu.execPUSH(inst)
	case instruction.PUSHL:
		if !extendedSet {
			cpu.execPUSHL(inst)
		}
	case instruction.RCALL:
		cpu.execRCALL(inst)
	case instruction.RESET:
		cpu.execRESET(inst)
	case instruction.RETFIE:
		cpu.execRETFIE(inst)
	case instruction.RETLW:
		cpu.execRETLW(inst)
	case instruction.RETURN:
		cpu.execRETURN(inst)
	case instruction.RLCF:
		cpu.execRLCF(inst)
	case instruction.RLNCF:
		cpu.execRLNCF(inst)
	case instruction.RRCF:
		cpu.execRRCF(inst)
	case instruction.RRNCF:
		cpu.execRRNCF(inst)
	case instruction.SETF:
		cpu.execSETF(inst)
	case instruction.SLEEP:
		cpu.execSLEEP(inst)
	case instruction.SUBFSR:
		if !extendedSet {
			cpu.execSUBFSR(inst)
		}
	case instruction.SUBFWB:
		cpu.execSUBFWB(inst)
	case instruction.SUBLW:
		cpu.execSUBLW(inst)
	case instruction.SUBWF:
		cpu.execSUBWF(inst)
	case instruction.SUBWFB:
		cpu.execSUBWFB(inst)
	case instruction.SWAPF:
		cpu.execSWAPF(inst)
	case instruction.TBLRD:
		cpu.execTBLRD(inst)
	case instruction.TBLWT:
		cpu.execTBLWT(inst)
	case instruction.TSTFSZ:
		cpu.execTSTFSZ(inst)
	case instruction.XORLW:
		cpu.execXORLW(inst)
	case instruction.XORWF:
		cpu.execXORWF(inst)
	}

	return true
}

// offsetPC offsets a program counter by a specifc number of words.
// The word offset is specified by a number 2's complement number stored in an unsigned int.
// Since the program counter counts bytes, it will always increase or descrease by a factor of 2.
//
// This function implements the baviour of most branch instructions.
//
// Example:
//
//	offsetPC(1000, 1)   // 1002
//	offsetPC(1000, 255) // 998
func offsetPC(pc uint32, signedWordOffset uint8) uint32 {
	return uint32(int64(pc) + int64(int8(signedWordOffset))*2)
}

// offsetPC_11bit works like offsetPC, bit the offset is an 11 bit number instead.
func offsetPC_11bit(pc uint32, signedWordOffset uint16) uint32 {
	native := int16(signedWordOffset) - int16((signedWordOffset&(1<<10))<<1)
	return uint32(int64(pc) + int64(native)*2)
}

func (cpu *CPU) execADDFSR(inst instruction.Instruction) {
	var fsrl_offset uint16
	if instruction.XinstFsr(inst).F() == 0 {
		fsrl_offset = uint16(Registers.FSR0L)
	} else if instruction.XinstFsr(inst).F() == 1 {
		fsrl_offset = uint16(Registers.FSR1L)
	} else {
		fsrl_offset = uint16(Registers.FSR2L)
	}

	low, _ := cpu.DataBus.BusRead(fsrl_offset)
	high, _ := cpu.DataBus.BusRead(fsrl_offset + 1)
	value := (uint16(high) << 8) | uint16(low)
	value += uint16(instruction.XinstFsr(inst).K())
	cpu.DataBus.BusWrite(fsrl_offset, uint8(value))
	cpu.DataBus.BusWrite(fsrl_offset+1, uint8(value>>8))

	if instruction.XinstFsr(inst).F() == 3 {
		cpu.wreg = instruction.Literal(inst).K()
		cpu.pc = cpu.Stack.Top()
		if !cpu.stackPop() {
			return
		}

		cpu.flush = true
	}
}

func (cpu *CPU) execADDLW(inst instruction.Instruction) {
	cpu.wreg = cpu.Alu.Add(cpu.wreg, instruction.Literal(inst).K())
}

func (cpu *CPU) execADDWF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.Add(cpu.wreg, val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}
}

func (cpu *CPU) execADDWFC(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.AddWithCarry(cpu.wreg, val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}
}

func (cpu *CPU) execANDLW(inst instruction.Instruction) {
	cpu.wreg = cpu.Alu.And(cpu.wreg, instruction.Literal(inst).K())
}

func (cpu *CPU) execANDWF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.And(cpu.wreg, val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}
}

func (cpu *CPU) execBC(inst instruction.Instruction) {
	if cpu.Alu.status.C() {
		cpu.pc = offsetPC(cpu.pc, instruction.ControlBranchStatus(inst).Literal())
		cpu.flush = true
	}
}

func (cpu *CPU) execBCF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.BitOriented(inst).F(), instruction.BitOriented(inst).A())
	val &= ^(1 << instruction.BitOriented(inst).Bit())
	cpu.BankController.Write(instruction.BitOriented(inst).F(), val, instruction.BitOriented(inst).A())
}

func (cpu *CPU) execBN(inst instruction.Instruction) {
	if cpu.Alu.status.N() {
		cpu.pc = offsetPC(cpu.pc, instruction.ControlBranchStatus(inst).Literal())
		cpu.flush = true
	}
}

func (cpu *CPU) execBNC(inst instruction.Instruction) {
	if !(cpu.Alu.status.C()) {
		cpu.pc = offsetPC(cpu.pc, instruction.ControlBranchStatus(inst).Literal())
		cpu.flush = true
	}
}

func (cpu *CPU) execBNN(inst instruction.Instruction) {
	if !(cpu.Alu.status.N()) {
		cpu.pc = offsetPC(cpu.pc, instruction.ControlBranchStatus(inst).Literal())
		cpu.flush = true
	}
}

func (cpu *CPU) execBNOV(inst instruction.Instruction) {
	if !(cpu.Alu.status.OV()) {
		cpu.pc = offsetPC(cpu.pc, instruction.ControlBranchStatus(inst).Literal())
		cpu.flush = true
	}
}

func (cpu *CPU) execBNZ(inst instruction.Instruction) {
	if !(cpu.Alu.status.Z()) {
		cpu.pc = offsetPC(cpu.pc, instruction.ControlBranchStatus(inst).Literal())
		cpu.flush = true
	}
}

func (cpu *CPU) execBRA(inst instruction.Instruction) {
	cpu.pc = offsetPC_11bit(cpu.pc, instruction.ControlBranch(inst).Literal())
	cpu.flush = true
}

func (cpu *CPU) execBSF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.BitOriented(inst).F(), instruction.BitOriented(inst).A())
	val |= (1 << instruction.BitOriented(inst).Bit())
	cpu.BankController.Write(instruction.BitOriented(inst).F(), val, instruction.BitOriented(inst).A())
}

func (cpu *CPU) execBTFSC(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.BitOriented(inst).F(), instruction.BitOriented(inst).A())
	val &= (1 << instruction.BitOriented(inst).Bit())
	if val == 0 {
		cpu.pc += 2
		cpu.flush = true
	}
}

func (cpu *CPU) execBTFSS(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.BitOriented(inst).F(), instruction.BitOriented(inst).A())
	val &= (1 << instruction.BitOriented(inst).Bit())
	if val != 0 {
		cpu.pc += 2
		cpu.flush = true
	}
}

func (cpu *CPU) execBTG(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.BitOriented(inst).F(), instruction.BitOriented(inst).A())
	val ^= (1 << instruction.BitOriented(inst).Bit())
	cpu.BankController.Write(instruction.BitOriented(inst).F(), val, instruction.BitOriented(inst).A())
}

func (cpu *CPU) execBOV(inst instruction.Instruction) {
	if cpu.Alu.status.OV() {
		cpu.pc = offsetPC(cpu.pc, instruction.ControlBranchStatus(inst).Literal())
		cpu.flush = true
	}
}

func (cpu *CPU) execBZ(inst instruction.Instruction) {
	if cpu.Alu.status.Z() {
		cpu.pc = offsetPC(cpu.pc, instruction.ControlBranchStatus(inst).Literal())
		cpu.flush = true
	}
}

func (cpu *CPU) execCALL(inst instruction.Instruction) bool {
	if !cpu.stackPush(cpu.pc) {
		return false
	}

	fetched_low, _ := cpu.ProgramBus.BusRead(cpu.pc)
	fetched_high, _ := cpu.ProgramBus.BusRead(cpu.pc + 1)
	next_instruction := (uint16(fetched_high) << 8) | uint16(fetched_low)
	pcHigh := instruction.ControlCallLow(next_instruction).Literal() << 8
	pcLow := instruction.ControlCallHigh(inst).Literal()

	cpu.pc = (uint32(pcHigh) | uint32(pcLow)) << 1

	if instruction.ControlCallHigh(inst).S() {
		cpu.shadowWreg = cpu.wreg
		cpu.shadowStatus = uint8(cpu.Alu.status)
		cpu.shadowBsr = cpu.BankController.BSR
	}

	cpu.flush = true
	return true
}

func (cpu *CPU) execCALLW(instruction.Instruction) bool {
	if !cpu.stackPush(cpu.pc) {
		return false
	}

	cpu.pc &= 0xFFFFFF00
	cpu.pc |= uint32(cpu.wreg)

	cpu.flush = true
	return true
}

func (cpu *CPU) execCLRF(inst instruction.Instruction) {
	cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	cpu.BankController.Write(instruction.ByteOriented(inst).F(), 0, instruction.ByteOriented(inst).A())
	cpu.Alu.status |= statusZ
}

func (cpu *CPU) execCLRWDT(instruction.Instruction) {
	// TODO:
}

func (cpu *CPU) execCOMF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.Complement(val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}
}

func (cpu *CPU) execCPFSEQ(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	if val == cpu.wreg {
		cpu.pc += 2
		cpu.flush = true
	}
}

func (cpu *CPU) execCPFSGT(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	if val > cpu.wreg {
		cpu.pc += 2
		cpu.flush = true
	}
}

func (cpu *CPU) execCPFSLT(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	if val < cpu.wreg {
		cpu.pc += 2
		cpu.flush = true
	}
}

func (cpu *CPU) execDAW(instruction.Instruction) {
	cpu.wreg = cpu.Alu.DecimalAdjust(cpu.wreg)
}

func (cpu *CPU) execDECF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.Sub(val, 1)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}
}

func (cpu *CPU) execDECFSZ(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := val - 1

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}

	if result == 0 {
		cpu.pc += 2
		cpu.flush = true
	}
}

func (cpu *CPU) execDCFSNZ(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := val - 1

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}

	if result != 0 {
		cpu.pc += 2
		cpu.flush = true
	}
}

func (cpu *CPU) execGOTO(inst instruction.Instruction) {
	fetched_low, _ := cpu.ProgramBus.BusRead(cpu.pc)
	fetched_high, _ := cpu.ProgramBus.BusRead(cpu.pc + 1)
	low_word := uint16(fetched_high)<<8 | uint16(fetched_low)
	next_instruction := instruction.ControlGotoLow(low_word)

	address := ((uint32(next_instruction.Literal()) << 8) | uint32(instruction.ControlGotoHigh(inst).Literal())) << 1
	cpu.pc = address
	cpu.flush = true
}

func (cpu *CPU) execINCF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.Add(val, 1)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}
}

func (cpu *CPU) execINCFSZ(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := val + 1

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}

	if result == 0 {
		cpu.pc += 2
		cpu.flush = true
	}
}

func (cpu *CPU) execINFSNZ(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := val + 1

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}

	if result != 0 {
		cpu.pc += 2
		cpu.flush = true
	}
}

func (cpu *CPU) execIORLW(inst instruction.Instruction) {
	cpu.wreg = cpu.Alu.Or(cpu.wreg, instruction.Literal(inst).K())
}

func (cpu *CPU) execIORWF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.Or(cpu.wreg, val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}
}

func (cpu *CPU) execLFSR(inst instruction.Instruction) {
	var fsrl_offset uint16
	if instruction.LoadFsrHigh(inst).F() == 0 {
		fsrl_offset = uint16(Registers.FSR0L)
	} else if instruction.LoadFsrHigh(inst).F() == 1 {
		fsrl_offset = uint16(Registers.FSR1L)
	} else {
		fsrl_offset = uint16(Registers.FSR2L)
	}

	cpu.DataBus.BusWrite(fsrl_offset+1, instruction.LoadFsrHigh(inst).K())
	cpu.nextAction = func(next instruction.Instruction) {
		cpu.DataBus.BusWrite(fsrl_offset, instruction.LoadFsrLow(next).K())
	}
}

func (cpu *CPU) execMOVF(inst instruction.Instruction) {
	byteInst := instruction.ByteOriented(inst)
	val := cpu.BankController.Read(byteInst.F(), byteInst.A())
	cpu.Alu.status &= ^statusN
	cpu.Alu.status &= ^statusZ

	if val == 0 {
		cpu.Alu.status |= statusZ
	}
	if val&0b10000000 != 0 {
		cpu.Alu.status |= statusN
	}

	if byteInst.D() {
		cpu.BankController.Write(byteInst.F(), val, byteInst.A())
	} else {
		cpu.wreg = val
	}
}

func (cpu *CPU) execMOVFF(inst instruction.Instruction) {
	val, _ := cpu.DataBus.BusRead(instruction.ByteToByte(inst).F())
	cpu.nextAction = func(next instruction.Instruction) {
		dest := instruction.ByteToByte(next).F()
		if dest == Registers.PCL || dest == Registers.TOSU || dest == Registers.TOSH || dest == Registers.TOSL {
			return
		}

		cpu.DataBus.BusWrite(dest, val)
	}
}

func (cpu *CPU) execMOVLB(inst instruction.Instruction) {
	cpu.BankController.BSR = instruction.Literal(inst).K() & 0xF
}

func (cpu *CPU) execMOVLW(inst instruction.Instruction) {
	cpu.wreg = instruction.Literal(inst).K()
}

func (cpu *CPU) execMOVSF(inst instruction.Instruction) {

	fsr2l, _ := cpu.DataBus.BusRead(Registers.FSR2L)
	fsr2h, _ := cpu.DataBus.BusRead(Registers.FSR2L + 1)
	fsr2_val := (uint16(fsr2h) << 8) | uint16(fsr2l)
	src_value, _ := cpu.DataBus.BusRead(fsr2_val + uint16(instruction.MovsfHighMovss(inst).Z()))
	cpu.nextAction = func(next instruction.Instruction) {
		cpu.DataBus.BusWrite(instruction.MovsfLow(next).F(), src_value)
	}
}

func (cpu *CPU) execMOVSS(inst instruction.Instruction) {
	fsr2l, _ := cpu.DataBus.BusRead(Registers.FSR2L)
	fsr2h, _ := cpu.DataBus.BusRead(Registers.FSR2L + 1)
	fsr2_val := (uint16(fsr2h) << 8) | uint16(fsr2l)
	src_value, _ := cpu.DataBus.BusRead(fsr2_val + uint16(instruction.MovsfHighMovss(inst).Z()))
	cpu.nextAction = func(next instruction.Instruction) {
		fsr2l, _ := cpu.DataBus.BusRead(Registers.FSR2L)
		fsr2h, _ := cpu.DataBus.BusRead(Registers.FSR2L + 1)
		fsr2_val := (uint16(fsr2h) << 8) | uint16(fsr2l)
		cpu.DataBus.BusWrite(fsr2_val+uint16(instruction.MovsfHighMovss(next).Z()), src_value)
	}
}

func (cpu *CPU) execMOVWF(inst instruction.Instruction) {
	cpu.BankController.Write(instruction.ByteOriented(inst).F(), cpu.wreg, instruction.ByteOriented(inst).A())
}

func (cpu *CPU) execMULLW(inst instruction.Instruction) {
	cpu.Alu.Mul(cpu.wreg, instruction.Literal(inst).K())
}

func (cpu *CPU) execMULWF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	cpu.Alu.Mul(cpu.wreg, val)
}

func (cpu *CPU) execNEGF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.Negate(val)
	cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
}

func (cpu *CPU) execPOP(instruction.Instruction) {
	cpu.stackPop()
}

func (cpu *CPU) execPUSH(instruction.Instruction) {
	cpu.stackPush(cpu.pc)
}

func (cpu *CPU) execPUSHL(inst instruction.Instruction) {
	fsr2l, _ := cpu.DataBus.BusRead(Registers.FSR2L)
	fsr2h, _ := cpu.DataBus.BusRead(Registers.FSR2L + 1)
	fsr2_val := (uint16(fsr2h) << 8) | uint16(fsr2l)
	cpu.DataBus.BusWrite(fsr2_val, instruction.Literal(inst).K())
	fsr2_val--
	cpu.DataBus.BusWrite(Registers.FSR2L, uint8(fsr2_val&0xFF))
	cpu.DataBus.BusWrite(Registers.FSR2L+1, uint8((fsr2_val>>8)&0xFF))
}

func (cpu *CPU) execRCALL(inst instruction.Instruction) {
	if !cpu.stackPush(cpu.pc) {
		return
	}
	cpu.pc = offsetPC_11bit(cpu.pc, instruction.ControlBranch(inst).Literal())
	cpu.flush = true
}

func (cpu *CPU) execRESET(instruction.Instruction) {
	cpu.MclrReset()
}

func (cpu *CPU) execRETFIE(inst instruction.Instruction) {
	cpu.pc = cpu.Stack.Top()
	if !cpu.stackPop() {
		return
	}

	if cpu.SetGlobalInterruptEnable != nil {
		cpu.SetGlobalInterruptEnable(cpu.highPriorityInterrupt, true)
	}

	if instruction.ControlReturn(inst).S() {
		cpu.wreg = cpu.shadowWreg
		cpu.Alu.status = AluStatus(cpu.shadowStatus)
		cpu.BankController.BSR = cpu.shadowBsr
	}

	cpu.flush = true
}

func (cpu *CPU) execRETLW(inst instruction.Instruction) {
	cpu.wreg = instruction.Literal(inst).K()
	cpu.pc = cpu.Stack.Top()
	if !cpu.stackPop() {
		return
	}

	cpu.flush = true
}

func (cpu *CPU) execRETURN(inst instruction.Instruction) {
	cpu.pc = cpu.Stack.Top()
	if !cpu.stackPop() {
		return
	}

	if instruction.ControlReturn(inst).S() {
		cpu.wreg = cpu.shadowWreg
		cpu.Alu.status = AluStatus(cpu.shadowStatus)
		cpu.BankController.BSR = cpu.shadowBsr
	}

	cpu.flush = true
}

func (cpu *CPU) execRLCF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.RotateLeftCarry(val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}
}

func (cpu *CPU) execRLNCF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.RotateLeft(val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}
}

func (cpu *CPU) execRRCF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.RotateRightCarry(val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}
}

func (cpu *CPU) execRRNCF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.RotateRight(val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}
}

func (cpu *CPU) execSETF(inst instruction.Instruction) {
	cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	cpu.BankController.Write(instruction.ByteOriented(inst).F(), 0xFF, instruction.ByteOriented(inst).A())
}

func (cpu *CPU) execSLEEP(instruction.Instruction) {
	// TODO: Clear WDT
	if cpu.EventHandler != nil {
		cpu.EventHandler.Sleep()
	}
}

func (cpu *CPU) execSUBFSR(inst instruction.Instruction) {
	var fsrl_offset uint16
	if instruction.XinstFsr(inst).K() == 0 {
		fsrl_offset = uint16(Registers.FSR0L)
	} else if instruction.XinstFsr(inst).K() == 1 {
		fsrl_offset = uint16(Registers.FSR1L)
	} else {
		fsrl_offset = uint16(Registers.FSR2L)
	}

	low, _ := cpu.DataBus.BusRead(fsrl_offset)
	high, _ := cpu.DataBus.BusRead(fsrl_offset + 1)
	value := (uint16(high) << 8) | uint16(low)
	value -= uint16(instruction.XinstFsr(inst).K())
	cpu.DataBus.BusWrite(fsrl_offset, uint8(value))
	cpu.DataBus.BusWrite(fsrl_offset+1, uint8(value>>8))

	if instruction.XinstFsr(inst).F() == 3 {
		cpu.wreg = instruction.Literal(inst).K()
		cpu.pc = cpu.Stack.Top()
		if !cpu.stackPop() {
			return
		}

		cpu.flush = true
	}
}

func (cpu *CPU) execSUBFWB(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.SubWithBorrow(cpu.wreg, val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}
}

func (cpu *CPU) execSUBLW(inst instruction.Instruction) {
	cpu.wreg = cpu.Alu.Sub(instruction.Literal(inst).K(), cpu.wreg)
}

func (cpu *CPU) execSUBWF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.Sub(val, cpu.wreg)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}
}

func (cpu *CPU) execSUBWFB(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.SubWithBorrow(val, cpu.wreg)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}
}

func (cpu *CPU) execSWAPF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := (val << 4) | (val >> 4)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}
}

func (cpu *CPU) execTBLRD(inst instruction.Instruction) {
	action := TblPtrAction(instruction.TableOp(inst).N())
	if cpu.TableRead != nil {
		cpu.TableRead(action)
	}
	cpu.flush = true
}

func (cpu *CPU) execTBLWT(inst instruction.Instruction) {
	action := TblPtrAction(instruction.TableOp(inst).N())
	if cpu.TableWrite != nil {
		cpu.TableWrite(action)
	}
	cpu.flush = true
}

func (cpu *CPU) execTSTFSZ(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	if val == 0 {
		cpu.pc += 2
		cpu.flush = true
	}
}

func (cpu *CPU) execXORLW(inst instruction.Instruction) {
	cpu.wreg = cpu.Alu.Xor(instruction.Literal(inst).K(), cpu.wreg)
}

func (cpu *CPU) execXORWF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.Xor(cpu.wreg, val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.wreg = result
	}
}
