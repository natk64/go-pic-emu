package pic18

import (
	"github.com/natk64/go-pic-emu/pic18/instruction"
)

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
	case instruction.ADDULNK:
		if extendedSet {
			cpu.execADDULNK(inst)
		}
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
		if extendedSet {
			cpu.execMOVSF(inst)
		}
	case instruction.MOVSS:
		if extendedSet {
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
		if extendedSet {
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
		if extendedSet {
			cpu.execSUBFSR(inst)
		}
	case instruction.SUBFWB:
		cpu.execSUBFWB(inst)
	case instruction.SUBLW:
		cpu.execSUBLW(inst)
	case instruction.SUBULNK:
		if extendedSet {
			cpu.execSUBULNK(inst)
		}
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
	index := int(instruction.XinstFsr(inst).F())
	if index < len(cpu.BankController.FSR) {
		cpu.BankController.FSR[index] += uint16(instruction.XinstFsr(inst).K())
	}
}

func (cpu *CPU) execADDLW(inst instruction.Instruction) {
	cpu.WReg = cpu.Alu.Add(cpu.WReg, instruction.Literal(inst).K())
}

func (cpu *CPU) execADDULNK(inst instruction.Instruction) {
	cpu.BankController.FSR[2] += uint16(instruction.Literal(inst).K())
	cpu.pc = cpu.Stack.Top()
	if !cpu.stackPop() {
		return
	}
	cpu.flush = true
}

func (cpu *CPU) execADDWF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.Add(cpu.WReg, val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.WReg = result
	}
}

func (cpu *CPU) execADDWFC(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.AddWithCarry(cpu.WReg, val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.WReg = result
	}
}

func (cpu *CPU) execANDLW(inst instruction.Instruction) {
	cpu.WReg = cpu.Alu.And(cpu.WReg, instruction.Literal(inst).K())
}

func (cpu *CPU) execANDWF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.And(cpu.WReg, val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.WReg = result
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
		cpu.shadowWreg = cpu.WReg
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
	cpu.pc |= uint32(cpu.WReg)

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
		cpu.WReg = result
	}
}

func (cpu *CPU) execCPFSEQ(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	if val == cpu.WReg {
		cpu.pc += 2
		cpu.flush = true
	}
}

func (cpu *CPU) execCPFSGT(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	if val > cpu.WReg {
		cpu.pc += 2
		cpu.flush = true
	}
}

func (cpu *CPU) execCPFSLT(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	if val < cpu.WReg {
		cpu.pc += 2
		cpu.flush = true
	}
}

func (cpu *CPU) execDAW(instruction.Instruction) {
	cpu.WReg = cpu.Alu.DecimalAdjust(cpu.WReg)
}

func (cpu *CPU) execDECF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.Sub(val, 1)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.WReg = result
	}
}

func (cpu *CPU) execDECFSZ(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := val - 1

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.WReg = result
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
		cpu.WReg = result
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
		cpu.WReg = result
	}
}

func (cpu *CPU) execINCFSZ(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := val + 1

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.WReg = result
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
		cpu.WReg = result
	}

	if result != 0 {
		cpu.pc += 2
		cpu.flush = true
	}
}

func (cpu *CPU) execIORLW(inst instruction.Instruction) {
	cpu.WReg = cpu.Alu.Or(cpu.WReg, instruction.Literal(inst).K())
}

func (cpu *CPU) execIORWF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.Or(cpu.WReg, val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.WReg = result
	}
}

func (cpu *CPU) execLFSR(inst instruction.Instruction) {
	k := instruction.LoadFsrHigh(inst).K()
	index := int(instruction.LoadFsrHigh(inst).F())
	if index < len(cpu.BankController.FSR) {
		cpu.BankController.FSR[index] &= 0x00FF
		cpu.BankController.FSR[index] |= uint16(k) << 8
	}

	cpu.nextAction = func(next instruction.Instruction) {
		k := instruction.LoadFsrLow(next).K()
		if index < len(cpu.BankController.FSR) {
			cpu.BankController.FSR[index] &= 0xFF00
			cpu.BankController.FSR[index] |= uint16(k)
		}
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
		cpu.WReg = val
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
	cpu.WReg = instruction.Literal(inst).K()
}

func (cpu *CPU) execMOVSF(inst instruction.Instruction) {
	src_value, _ := cpu.DataBus.BusRead(uint16(cpu.BankController.FSR[2]) + uint16(instruction.MovsfHighMovss(inst).Z()))
	cpu.nextAction = func(next instruction.Instruction) {
		cpu.DataBus.BusWrite(instruction.MovsfLow(next).F(), src_value)
	}
}

func (cpu *CPU) execMOVSS(inst instruction.Instruction) {
	src_value, _ := cpu.DataBus.BusRead(uint16(cpu.BankController.FSR[2]) + uint16(instruction.MovsfHighMovss(inst).Z()))
	cpu.nextAction = func(next instruction.Instruction) {
		cpu.DataBus.BusWrite(uint16(cpu.BankController.FSR[2])+uint16(instruction.MovsfHighMovss(next).Z()), src_value)
	}
}

func (cpu *CPU) execMOVWF(inst instruction.Instruction) {
	cpu.BankController.Write(instruction.ByteOriented(inst).F(), cpu.WReg, instruction.ByteOriented(inst).A())
}

func (cpu *CPU) execMULLW(inst instruction.Instruction) {
	cpu.Alu.Mul(cpu.WReg, instruction.Literal(inst).K())
}

func (cpu *CPU) execMULWF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	cpu.Alu.Mul(cpu.WReg, val)
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
	cpu.DataBus.BusWrite(cpu.BankController.FSR[2], instruction.Literal(inst).K())
	cpu.BankController.FSR[2]--
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

	cpu.Interrupts.HighPriorityEnable = true

	if instruction.ControlReturn(inst).S() {
		cpu.WReg = cpu.shadowWreg
		cpu.Alu.status = AluStatus(cpu.shadowStatus)
		cpu.BankController.BSR = cpu.shadowBsr
	}

	cpu.flush = true
}

func (cpu *CPU) execRETLW(inst instruction.Instruction) {
	cpu.WReg = instruction.Literal(inst).K()
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
		cpu.WReg = cpu.shadowWreg
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
		cpu.WReg = result
	}
}

func (cpu *CPU) execRLNCF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.RotateLeft(val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.WReg = result
	}
}

func (cpu *CPU) execRRCF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.RotateRightCarry(val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.WReg = result
	}
}

func (cpu *CPU) execRRNCF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.RotateRight(val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.WReg = result
	}
}

func (cpu *CPU) execSETF(inst instruction.Instruction) {
	cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	cpu.BankController.Write(instruction.ByteOriented(inst).F(), 0xFF, instruction.ByteOriented(inst).A())
}

func (cpu *CPU) execSLEEP(instruction.Instruction) {
	// TODO: Clear WDT
	cpu.Sleep.Sleep()
}

func (cpu *CPU) execSUBFSR(inst instruction.Instruction) {
	index := int(instruction.XinstFsr(inst).F())
	if index < len(cpu.BankController.FSR) {
		cpu.BankController.FSR[index] -= uint16(instruction.XinstFsr(inst).K())
	}
}

func (cpu *CPU) execSUBFWB(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.SubWithBorrow(cpu.WReg, val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.WReg = result
	}
}

func (cpu *CPU) execSUBLW(inst instruction.Instruction) {
	cpu.WReg = cpu.Alu.Sub(instruction.Literal(inst).K(), cpu.WReg)
}

func (cpu *CPU) execSUBULNK(inst instruction.Instruction) {
	cpu.BankController.FSR[2] -= uint16(instruction.Literal(inst).K())
	cpu.pc = cpu.Stack.Top()
	if !cpu.stackPop() {
		return
	}
	cpu.flush = true
}

func (cpu *CPU) execSUBWF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.Sub(val, cpu.WReg)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.WReg = result
	}
}

func (cpu *CPU) execSUBWFB(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.SubWithBorrow(val, cpu.WReg)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.WReg = result
	}
}

func (cpu *CPU) execSWAPF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := (val << 4) | (val >> 4)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.WReg = result
	}
}

func (cpu *CPU) execTBLRD(inst instruction.Instruction) {
	action := TableAction(instruction.TableOp(inst).N())
	cpu.Table.TableRead(action)
	cpu.flush = true
}

func (cpu *CPU) execTBLWT(inst instruction.Instruction) {
	action := TableAction(instruction.TableOp(inst).N())
	cpu.Table.TableWrite(action)
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
	cpu.WReg = cpu.Alu.Xor(instruction.Literal(inst).K(), cpu.WReg)
}

func (cpu *CPU) execXORWF(inst instruction.Instruction) {
	val := cpu.BankController.Read(instruction.ByteOriented(inst).F(), instruction.ByteOriented(inst).A())
	result := cpu.Alu.Xor(cpu.WReg, val)

	if instruction.ByteOriented(inst).D() {
		cpu.BankController.Write(instruction.ByteOriented(inst).F(), result, instruction.ByteOriented(inst).A())
	} else {
		cpu.WReg = result
	}
}
