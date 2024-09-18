package pic18

type InterruptController struct {
	InterruptPriorityEnable bool

	HighPriorityEnable bool
	LowPriorityEnable  bool

	DoGotoHighPriority bool
	DoGotoLowPriority  bool

	Sleep *SleepController

	sources   []*interruptSource
	sourceBus map[uint16]MultiBusReadWriter[uint16]
}

type interruptSource struct {
	Enable       bool
	Flag         bool
	HighPriority bool

	index      int
	controller *InterruptController
	config     InterruptConfig
}

func (src interruptSource) BusRead(addr uint16) (data uint8, mask AddrMask) {
	if addr == src.config.Request.Register {
		if src.Flag {
			data |= 1 << src.config.Request.Bit
		}
		mask |= 1 << src.config.Request.Bit
	}
	if addr == src.config.Enable.Register {
		if src.Enable {
			data |= 1 << src.config.Enable.Bit
		}
		mask |= 1 << src.config.Enable.Bit
	}
	if addr == src.config.Priority.Register {
		if src.HighPriority {
			data |= 1 << src.config.Priority.Bit
		}
		mask |= 1 << src.config.Priority.Bit
	}
	return
}

func (src *interruptSource) BusWrite(addr uint16, data uint8) (mask AddrMask) {
	if addr == src.config.Request.Register {
		src.Flag = addr&uint16(1<<src.config.Request.Bit) != 0
		mask |= 1 << src.config.Request.Bit
	}
	if addr == src.config.Enable.Register {
		src.Flag = addr&uint16(1<<src.config.Enable.Bit) != 0
		mask |= 1 << src.config.Enable.Bit
	}
	if addr == src.config.Priority.Register {
		src.Flag = addr&uint16(1<<src.config.Priority.Bit) != 0
		mask |= 1 << src.config.Priority.Bit
	}
	return
}

func (src *interruptSource) Raise() {
	if src.controller == nil || src.index >= len(src.controller.sources) {
		panic("illegal interrupt")
	}
	src.Flag = true
	src.controller.raiseInterrupt(src.controller.sources[src.index])
}

func (src *interruptSource) Clear() {
	src.Flag = false
}

// raiseInterrupt should be called when the interrupt request flag of a source is set.
// This function itself neither reads nor writes that flag.
func (controller *InterruptController) raiseInterrupt(src *interruptSource) {
	if controller.InterruptPriorityEnable {
		controller.raiseInterruptDefault(src)
	} else {
		controller.raiseInterruptCompatibilityMode(src)
	}
}

func (controller *InterruptController) raiseInterruptDefault(src *interruptSource) {
	if !src.Enable {
		// Interrupt source is disabled
		return
	}

	if src.HighPriority && !controller.HighPriorityEnable {
		// High priority interrupts are disabled
		return
	}

	if !src.HighPriority && !controller.LowPriorityEnable {
		// Low priority interrupts are disabled
		return
	}

	if src.HighPriority {
		controller.HighPriorityEnable = false
		controller.DoGotoHighPriority = true
	} else {
		controller.LowPriorityEnable = false
		controller.DoGotoLowPriority = true
	}

	controller.Sleep.WakeUp()
}

func (controller *InterruptController) raiseInterruptCompatibilityMode(src *interruptSource) {
	globalInterruptEnable := controller.HighPriorityEnable
	peripheralInterruptEnable := controller.LowPriorityEnable

	if !src.Enable {
		// Interrupt source is disabled
		return
	}

	if !globalInterruptEnable {
		// All interrupts are disabled
		return
	}

	if src.config.Peripheral && !peripheralInterruptEnable {
		// Peripheral interrupts are disabled
		return
	}

	controller.HighPriorityEnable = false
	controller.DoGotoHighPriority = true
	controller.Sleep.WakeUp()
}

func (controller *InterruptController) CheckHighPriority() bool {
	tmp := controller.DoGotoHighPriority
	controller.DoGotoHighPriority = false
	return tmp
}

func (controller *InterruptController) CheckLowPriority() bool {
	tmp := controller.DoGotoLowPriority
	controller.DoGotoLowPriority = false
	return tmp
}

func (controller *InterruptController) BusRead(addr uint16) (data uint8, mask AddrMask) {
	bus, ok := controller.sourceBus[addr]
	if ok {
		data, mask = bus.BusRead(addr)
	}

	if addr == Registers.INTCON {
		if controller.HighPriorityEnable {
			data |= 0x80
		}
		if controller.LowPriorityEnable {
			data |= 0x40
		}
		mask |= 0x80 | 0x40
	} else if addr == Registers.RCON {
		if controller.InterruptPriorityEnable {
			data |= 0x80
		}
		mask |= 0x80
	}

	return
}

func (controller *InterruptController) BusWrite(addr uint16, data uint8) (mask AddrMask) {
	bus, ok := controller.sourceBus[addr]
	if ok {
		mask |= bus.BusWrite(addr, data)
	}

	if addr == Registers.INTCON {
		controller.HighPriorityEnable = data&0x80 != 0
		controller.LowPriorityEnable = data&0x40 != 0
		mask |= 0x80 | 0x40
	} else if addr == Registers.RCON {
		controller.InterruptPriorityEnable = data&0x80 != 0
		mask |= 0x80
	}

	return
}

type Interrupt interface {
	Raise()
	Clear()
}

type InterruptFlag struct {
	Register uint16
	Bit      uint8
}

type InterruptConfig struct {
	DebugLabel string
	Peripheral bool

	Request  InterruptFlag
	Priority InterruptFlag
	Enable   InterruptFlag
}

func (controller *InterruptController) CreateInterrupt(config InterruptConfig) Interrupt {
	src := &interruptSource{
		index:      len(controller.sources),
		controller: controller,
		config:     config,
		Enable:     true,
	}

	if controller.sourceBus == nil {
		controller.sourceBus = make(map[uint16]MultiBusReadWriter[uint16])
	}

	bus := controller.sourceBus
	bus[config.Request.Register] = append(bus[config.Request.Register], src)
	if config.Enable.Register != config.Request.Register {
		bus[config.Enable.Register] = append(bus[config.Enable.Register], src)
	}
	if config.Priority.Register != config.Request.Register && config.Priority.Register != config.Enable.Register {
		bus[config.Priority.Register] = append(bus[config.Priority.Register], src)
	}

	controller.sources = append(controller.sources, src)
	return src
}

type interruptSourceReg struct {
	enable   uint16
	request  uint16
	priority uint16
}

var peripheralInterruptRegisters = [5]interruptSourceReg{
	{enable: 0xF9D, request: 0xF9E, priority: 0xF9F},
	{enable: 0xFA0, request: 0xFA1, priority: 0xFA2},
	{enable: 0xFA3, request: 0xFA4, priority: 0xFA5},
	{enable: 0xFB6, request: 0xFB7, priority: 0xFB8},
	{enable: 0xF76, request: 0xF77, priority: 0xF78},
}

// PeripheralInterrupt creates a new peripheral interrupt config.
//
// registerNum starts at 1.
func PeripheralInterrupt(debugLabel string, registerNum, bitNum int) InterruptConfig {
	register := peripheralInterruptRegisters[registerNum-1]

	return InterruptConfig{
		DebugLabel: debugLabel,
		Peripheral: true,
		Enable:     InterruptFlag{Register: register.enable, Bit: uint8(bitNum)},
		Request:    InterruptFlag{Register: register.request, Bit: uint8(bitNum)},
		Priority:   InterruptFlag{Register: register.priority, Bit: uint8(bitNum)},
	}
}
