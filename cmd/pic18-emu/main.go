package main

import (
	"fmt"
	"time"

	"github.com/natk64/go-pic-emu/pic18"
)

var _ pic18.CpuEventHandler = FuncEventHandler{}

type FuncEventHandler struct {
	SleepFunc func()
}

func (f FuncEventHandler) IllegalInstruction() {
	panic("invalid instruction")
}

func (f FuncEventHandler) Sleep() {
	if f.SleepFunc != nil {
		f.SleepFunc()
	}
}

func main() {
	run := true
	cpu := &pic18.CPU{
		Config: &pic18.ConfigTable{},
		EventHandler: FuncEventHandler{
			SleepFunc: func() {
				run = false
			},
		},
	}

	dataBus := pic18.MultiBusReadWriter[uint16]{
		cpu,
		&cpu.Alu,
		&cpu.Stack,
		&cpu.BankController,
	}

	programBus := pic18.MultiBusReadWriter[uint32]{
		pic18.Memory[uint32]{
			Data: []byte{
				0x00,
				0x00,

				// ADDLW 1
				0x01,
				0x0F,

				// BNZ -2
				0xFE,
				0xE1,

				// MOVLW 69
				0x45,
				0x0E,

				// SLEEP
				0x03,
				0x00,
			},
		},
	}

	cpu.DataBus = dataBus
	cpu.ProgramBus = programBus
	cpu.BankController.Bus = dataBus

	cpu.PowerOnReset()

	ticks := 0
	start := time.Now()
	for run {
		ticks++
		cpu.Tick()
	}

	elapsed := time.Since(start)
	fmt.Printf("%d ticks in %v, %v MHz\n", ticks, elapsed, ticks/int(elapsed.Microseconds()))

	wreg, _ := dataBus.BusRead(pic18.Registers.WREG)
	fmt.Printf("WREG: %d\n", wreg)
}
