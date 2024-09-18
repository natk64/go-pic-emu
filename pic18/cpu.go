package pic18

import "github.com/natk64/go-pic-emu/pic18/instruction"

type InterruptState int

const (
	InterruptStateNone InterruptState = iota
	InterruptStateHighPrio
	InterruptStateLowPrio
)

type CPU struct {
	pc                 uint32
	fetchedInstruction uint16
	WReg               uint8

	pcLatchHigh  uint8
	pcLatchUpper uint8

	flush          bool
	interruptState InterruptState

	shadowWreg   uint8
	shadowStatus uint8
	shadowBsr    uint8

	Alu            ALU
	Stack          Stack
	BankController BankController
	Interrupts     InterruptController

	nextAction func(instruction.Instruction)

	Table        *TableRWController
	Sleep        *SleepController
	Config       *ConfigTable
	DataBus      DataBusReadWriter
	ProgramBus   ProgramBusReadWriter
	EventHandler CpuEventHandler
}

type CpuEventHandler interface {
	IllegalInstruction()
}

func (cpu *CPU) BusRead(addr uint16) (uint8, AddrMask) {
	switch addr {
	case Registers.PCL:
		cpu.pcLatchUpper = uint8((cpu.pc >> 16) & 0xFF)
		cpu.pcLatchHigh = uint8((cpu.pc >> 8) & 0xFF)
		pcLower := uint8(cpu.pc & 0xFF)
		return pcLower, 0xFF
	case Registers.PCLATH:
		return cpu.pcLatchHigh, 0xFF
	case Registers.PCLATU:
		return cpu.pcLatchUpper, 0xFF
	case Registers.WREG:
		return cpu.WReg, 0xFF
	default:
		return 0, 0
	}
}

func (cpu *CPU) BusWrite(addr uint16, data uint8) AddrMask {
	switch addr {
	case Registers.PCL:
		cpu.pc = (uint32(cpu.pcLatchUpper) << 16) | (uint32(cpu.pcLatchHigh) << 8) | uint32(data)
		return 0xFF
	case Registers.PCLATH:
		cpu.pcLatchHigh = data
		return 0xFF
	case Registers.PCLATU:
		cpu.pcLatchUpper = data
		return 0xFF
	case Registers.WREG:
		cpu.WReg = data
		return 0xFF
	default:
		return 0
	}
}

func (cpu *CPU) PowerOnReset() {
	cpu.flush = true
	cpu.pc = 0
	cpu.fetchedInstruction = 0
	cpu.interruptState = 0
	cpu.nextAction = nil
	cpu.pcLatchHigh = 0
	cpu.pcLatchUpper = 0
	cpu.shadowBsr = 0
	cpu.shadowStatus = 0
	cpu.shadowWreg = 0
}

func (cpu *CPU) MclrReset() {
	panic("unimplemented")
}

func (cpu *CPU) GotoInterrupt(highPrio bool) {
	cpu.shadowWreg = cpu.WReg
	cpu.shadowStatus = uint8(cpu.Alu.status)
	cpu.shadowBsr = cpu.BankController.BSR
	if highPrio {
		cpu.interruptState = InterruptStateHighPrio
	} else {
		cpu.interruptState = InterruptStateLowPrio
	}

	if !cpu.Stack.Push(cpu.pc) {
		return
	}

	cpu.flush = true
	if highPrio {
		cpu.pc = 0x08
	} else {
		cpu.pc = 0x18
	}
}

// stackPush wraps stack.Push to handle the reset logic.
// It returns true if execution can continue, or false if the CPU was reset and the caller should abort.
func (cpu *CPU) stackPush(value uint32) bool {
	becameFull := cpu.Stack.Push(value)
	if becameFull && cpu.Config.ResetOnFull {
		cpu.MclrReset()
		return false
	}

	return true
}

// stackPop wraps stack.Pop, just like stackPush.
func (cpu *CPU) stackPop() bool {
	underflow := cpu.Stack.Pop()
	if underflow && cpu.Config.ResetOnUnderflow {
		cpu.MclrReset()
		return false
	}

	return true
}

func (cpu *CPU) Tick() {
	if cpu.Interrupts.CheckHighPriority() {
		cpu.GotoInterrupt(true)
		return
	} else if cpu.Interrupts.CheckLowPriority() {
		cpu.GotoInterrupt(false)
		return
	}

	if cpu.flush {
		cpu.FetchInstruction()
		cpu.flush = false
		return
	}

	decoded := instruction.Instruction(cpu.fetchedInstruction)
	cpu.pc += 2

	if cpu.nextAction != nil {
		tmp := cpu.nextAction
		cpu.nextAction = nil
		tmp(decoded)
		cpu.FetchInstruction()
		return
	}

	ok := cpu.ExecuteInstruction(decoded, true)
	if !ok {
		if cpu.EventHandler != nil {
			cpu.EventHandler.IllegalInstruction()
		}
	}

	cpu.BankController.ApplyIndirectOp()
	cpu.FetchInstruction()
}

func (cpu *CPU) FetchInstruction() {
	fetchedLow, _ := cpu.ProgramBus.BusRead(cpu.pc)
	fetchedHigh, _ := cpu.ProgramBus.BusRead(cpu.pc + 1)
	cpu.fetchedInstruction = (uint16(fetchedHigh) << 8) | uint16(fetchedLow)
}
