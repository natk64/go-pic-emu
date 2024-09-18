package main

import (
	"fmt"
	"log"
	"time"

	"github.com/natk64/go-pic-emu/binary"
	"github.com/natk64/go-pic-emu/pic18"
	"github.com/natk64/go-pic-emu/pic18/peripherals/eusart"
)

var _ pic18.CpuEventHandler = DefaultEventHandler{}

type DefaultEventHandler struct{}

func (DefaultEventHandler) IllegalInstruction() {
	panic("invalid instruction")
}

func main() {
	run := true
	sleep := &pic18.SleepController{
		OnSleep: func() {
			run = false
		},
		OnWakeUp: func() { run = true },
	}

	cpu := &pic18.CPU{
		Config:       &pic18.ConfigTable{},
		Stack:        pic18.Stack{Data: make([]uint32, 31)},
		Sleep:        sleep,
		Interrupts:   pic18.InterruptController{Sleep: sleep},
		EventHandler: DefaultEventHandler{},
		Table:        &pic18.TableRWController{},
	}

	ram := pic18.Memory[uint16]{Data: make([]byte, 2048)}
	eusart1 := eusart.New(1, &cpu.Interrupts)
	eusart2 := eusart.New(2, &cpu.Interrupts)

	var dataBus pic18.BusReadWriter[uint16] = pic18.MultiBusReadWriter[uint16]{
		ram,
		cpu,
		eusart1,
		eusart2,
		cpu.Table,
		&cpu.Alu,
		&cpu.Stack,
		&cpu.BankController,
		&cpu.Interrupts,
	}

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
	cpu.BankController.Bus = pic18.BusPrinter(dataBus)
	cpu.BankController.WReg = &cpu.WReg
	cpu.Table.ProgramBus = programBus

	cpu.PowerOnReset()

	ticks := 0
	start := time.Now()
	for {
		if !run {
			if time.Since(start) > time.Second*5 {
				break
			}
			time.Sleep(time.Millisecond)
			continue
		}
		ticks++
		cpu.Tick()
	}

	elapsed := time.Since(start)
	fmt.Printf("%d ticks in %v, %v MHz\n", ticks, elapsed, ticks/int(elapsed.Microseconds()))

	wreg, _ := dataBus.BusRead(pic18.Registers.WREG)
	fmt.Printf("WREG: %d\n", wreg)
}
