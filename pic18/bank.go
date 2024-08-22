package pic18

type BankController struct {
	ExtendedSet bool
	BSR         uint8
	Bus         BusReadWriter[uint16]
}

func (controller BankController) address(location uint8, useBankSelectRegister bool) uint16 {
	if useBankSelectRegister {
		return (uint16(controller.BSR) << 8) | uint16(location)
	}

	if controller.ExtendedSet && location < 0x60 {
		fsr2low, _ := controller.Bus.BusRead(Registers.FSR2L)
		fsr2high, _ := controller.Bus.BusRead(Registers.FSR2H)

		return ((uint16(fsr2high) << 8) | uint16(fsr2low)) + uint16(location)
	}

	if location >= 0x80 {
		return 0xF00 + uint16(location)
	}

	return uint16(location)
}

func (controller BankController) Read(location uint8, useBankSelectRegister bool) uint8 {
	value, _ := controller.Bus.BusRead(controller.address(location, useBankSelectRegister))
	return value
}

func (controller BankController) Write(location uint8, value uint8, useBankSelectRegister bool) {
	controller.Bus.BusWrite(controller.address(location, useBankSelectRegister), value)
}

func (controller BankController) BusRead(addr uint16) (uint8, AddrMask) {
	if addr == Registers.BSR {
		return controller.BSR, 0xFF
	}
	return 0, 0
}

func (controller *BankController) BusWrite(addr uint16, data uint8) AddrMask {
	if addr == Registers.BSR {
		controller.BSR = data
		return 0xFF
	}
	return 0
}
