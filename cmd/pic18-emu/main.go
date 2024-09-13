package main

import (
	"fmt"
	"log"
	"time"

	"github.com/natk64/go-pic-emu/binary"
	"github.com/natk64/go-pic-emu/pic18"
)

var _ pic18.CpuEventHandler = DefaultEventHandler{}

type DefaultEventHandler struct{}

func (DefaultEventHandler) IllegalInstruction() {
	panic("invalid instruction")
}

func main() {
	run := true
	sleep := &pic18.SleepController{
		OnSleep:  func() { run = false },
		OnWakeUp: func() { run = true },
	}

	cpu := &pic18.CPU{
		Config:       &pic18.ConfigTable{},
		Stack:        pic18.Stack{Data: make([]uint32, 31)},
		Sleep:        sleep,
		Interrupts:   pic18.InterruptController{Sleep: sleep},
		EventHandler: DefaultEventHandler{},
	}

	dataBus := pic18.MultiBusReadWriter[uint16]{
		cpu,
		&cpu.Alu,
		&cpu.Stack,
		&cpu.BankController,
		&cpu.Interrupts,
	}

	tmr0 := cpu.Interrupts.CreateInterrupt(pic18.InterruptConfig{
		DebugLabel: "TMR0",
		Peripheral: true,
		Request:    pic18.InterruptFlag{Register: pic18.Registers.INTCON, Bit: 2},
		Priority:   pic18.InterruptFlag{Register: pic18.Registers.INTCON2, Bit: 2},
		Enable:     pic18.InterruptFlag{Register: pic18.Registers.INTCON, Bit: 5},
	})
	program, err := binary.ReadIHexFile("output/program.hex")
	if err != nil {
		log.Fatalln(err)
	}

	programBus := pic18.MultiBusReadWriter[uint32]{
		pic18.Memory[uint32]{
			Data: program,
		},
	}

	cpu.DataBus = dataBus
	cpu.ProgramBus = programBus
	cpu.BankController.Bus = dataBus

	cpu.PowerOnReset()

	ticks := 0
	start := time.Now()
	for {
		if !run {
			if time.Since(start) > time.Second*20 {
				break
			}
			time.Sleep(time.Millisecond)
			continue
		}
		ticks++
		cpu.Tick()
	}

	tmr0.Raise()
	elapsed := time.Since(start)
	fmt.Printf("%d ticks in %v, %v MHz\n", ticks, elapsed, ticks/int(elapsed.Microseconds()))

	wreg, _ := dataBus.BusRead(pic18.Registers.WREG)
	fmt.Printf("WREG: %d\n", wreg)
}
