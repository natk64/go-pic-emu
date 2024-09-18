package pic18

type BankController struct {
	ExtendedSet bool
	WReg        *uint8
	BSR         uint8
	FSR         [3]uint16
	newFSR      [3]uint16
	Bus         BusReadWriter[uint16]
}

type FSRAction int

const (
	FSRNone FSRAction = iota
	FSRPreInc
	FSRPostInc
	FSRPostDec
	FSRPlusW
)

const (
	FSR0H uint16 = 0xFEA
	FSR0L uint16 = 0xFE9
	FSR1H uint16 = 0xFE2
	FSR1L uint16 = 0xFE1
	FSR2H uint16 = 0xFDA
	FSR2L uint16 = 0xFD9

	INDF0    uint16 = 0xFEF
	POSTINC0 uint16 = 0xFEE
	POSTDEC0 uint16 = 0xFED
	PREINC0  uint16 = 0xFEC
	PLUSW0   uint16 = 0xFEB
	INDF1    uint16 = 0xFE7
	POSTINC1 uint16 = 0xFE6
	POSTDEC1 uint16 = 0xFE5
	PREINC1  uint16 = 0xFE4
	PLUSW1   uint16 = 0xFE3
	INDF2    uint16 = 0xFDF
	POSTINC2 uint16 = 0xFDE
	POSTDEC2 uint16 = 0xFDD
	PREINC2  uint16 = 0xFDC
	PLUSW2   uint16 = 0xFDB
)

func (controller BankController) address(location uint8, useBankSelectRegister bool) uint16 {
	if useBankSelectRegister {
		return (uint16(controller.BSR) << 8) | uint16(location)
	}

	if controller.ExtendedSet && location < 0x60 {
		return (controller.FSR[2] + uint16(location))
	}

	if location >= 0x80 {
		return 0xF00 + uint16(location)
	}

	return uint16(location)
}

func (controller *BankController) ApplyIndirectOp() {
	controller.FSR = controller.newFSR
}

func (controller BankController) Read(location uint8, useBankSelectRegister bool) uint8 {
	value, _ := controller.Bus.BusRead(controller.address(location, useBankSelectRegister))
	return value
}

func (controller BankController) Write(location uint8, value uint8, useBankSelectRegister bool) {
	controller.Bus.BusWrite(controller.address(location, useBankSelectRegister), value)
}

func (controller *BankController) BusRead(addr uint16) (uint8, AddrMask) {
	switch addr {
	case Registers.BSR:
		return controller.BSR, 0xFF
	case FSR0L:
		return uint8(controller.FSR[0]), 0xFF
	case FSR1L:
		return uint8(controller.FSR[1]), 0xFF
	case FSR2L:
		return uint8(controller.FSR[2]), 0xFF
	case FSR0H:
		return uint8(controller.FSR[0] >> 8), 0x0F
	case FSR1H:
		return uint8(controller.FSR[1] >> 8), 0x0F
	case FSR2H:
		return uint8(controller.FSR[2] >> 8), 0x0F
	case INDF0:
		return controller.indirectRead(0, FSRNone)
	case INDF1:
		return controller.indirectRead(1, FSRNone)
	case INDF2:
		return controller.indirectRead(2, FSRNone)
	case PREINC0:
		return controller.indirectRead(0, FSRPreInc)
	case PREINC1:
		return controller.indirectRead(1, FSRPreInc)
	case PREINC2:
		return controller.indirectRead(2, FSRPreInc)
	case POSTINC0:
		return controller.indirectRead(0, FSRPostInc)
	case POSTINC1:
		return controller.indirectRead(1, FSRPostInc)
	case POSTINC2:
		return controller.indirectRead(2, FSRPostInc)
	case POSTDEC0:
		return controller.indirectRead(0, FSRPostDec)
	case POSTDEC1:
		return controller.indirectRead(1, FSRPostDec)
	case POSTDEC2:
		return controller.indirectRead(2, FSRPostDec)
	case PLUSW0:
		return controller.indirectRead(0, FSRPlusW)
	case PLUSW1:
		return controller.indirectRead(1, FSRPlusW)
	case PLUSW2:
		return controller.indirectRead(2, FSRPlusW)
	}

	return 0, 0
}

func (controller *BankController) BusWrite(addr uint16, data uint8) AddrMask {
	switch addr {
	case Registers.BSR:
		controller.BSR = data
		return 0xFF
	case FSR0L:
		return controller.setFSRL(0, data)
	case FSR1L:
		return controller.setFSRL(1, data)
	case FSR2L:
		return controller.setFSRL(2, data)
	case FSR0H:
		return controller.setFSRH(0, data)
	case FSR1H:
		return controller.setFSRH(1, data)
	case FSR2H:
		return controller.setFSRH(2, data)
	case INDF0:
		return controller.indirectWrite(0, data, FSRNone)
	case INDF1:
		return controller.indirectWrite(1, data, FSRNone)
	case INDF2:
		return controller.indirectWrite(2, data, FSRNone)
	case PREINC0:
		return controller.indirectWrite(0, data, FSRPreInc)
	case PREINC1:
		return controller.indirectWrite(1, data, FSRPreInc)
	case PREINC2:
		return controller.indirectWrite(2, data, FSRPreInc)
	case POSTINC0:
		return controller.indirectWrite(0, data, FSRPostInc)
	case POSTINC1:
		return controller.indirectWrite(1, data, FSRPostInc)
	case POSTINC2:
		return controller.indirectWrite(2, data, FSRPostInc)
	case POSTDEC0:
		return controller.indirectWrite(0, data, FSRPostDec)
	case POSTDEC1:
		return controller.indirectWrite(1, data, FSRPostDec)
	case POSTDEC2:
		return controller.indirectWrite(2, data, FSRPostDec)
	case PLUSW0:
		return controller.indirectWrite(0, data, FSRPlusW)
	case PLUSW1:
		return controller.indirectWrite(1, data, FSRPlusW)
	case PLUSW2:
		return controller.indirectWrite(2, data, FSRPlusW)
	}
	return 0
}

func (BankController) isIndfRegister(addr uint16) bool {
	return addr == INDF0 || addr == PREINC0 || addr == POSTINC0 || addr == POSTDEC0 || addr == PLUSW0 ||
		addr == INDF1 || addr == PREINC1 || addr == POSTINC1 || addr == POSTDEC1 || addr == PLUSW1 ||
		addr == INDF2 || addr == PREINC2 || addr == POSTINC2 || addr == POSTDEC2 || addr == PLUSW2
}

func (controller *BankController) indirectAddress(fsr int, action FSRAction) (addr uint16) {
	if action == FSRNone {
		return controller.FSR[fsr]
	}

	if action == FSRPlusW {
		return (controller.FSR[fsr] + uint16(*controller.WReg)) & 0xFFF
	}

	if action == FSRPreInc {
		controller.newFSR[fsr] = controller.FSR[fsr] + 1
		controller.newFSR[fsr] &= 0xFFF
		return controller.FSR[fsr]
	}

	if action == FSRPostInc {
		addr = controller.FSR[fsr]
		controller.newFSR[fsr] = controller.FSR[fsr] + 1
		controller.newFSR[fsr] &= 0xFFF
	}

	if action == FSRPostDec {
		addr = controller.FSR[fsr]
		controller.newFSR[fsr] = controller.FSR[fsr] - 1
		controller.newFSR[fsr] &= 0xFFF
	}

	return addr
}

func (controller *BankController) indirectRead(fsr int, action FSRAction) (uint8, AddrMask) {
	address := controller.indirectAddress(fsr, action)
	if controller.isIndfRegister(address) {
		return 0, 0
	}
	return controller.Bus.BusRead(address)
}

func (controller *BankController) indirectWrite(fsr int, data uint8, action FSRAction) AddrMask {
	address := controller.indirectAddress(fsr, action)
	if controller.isIndfRegister(address) {
		return 0
	}
	return controller.Bus.BusWrite(address, data)
}

func (controller *BankController) setFSRH(num int, value uint8) AddrMask {
	controller.FSR[num] &= 0x00FF
	controller.FSR[num] |= uint16(value) << 8
	controller.newFSR[num] = controller.FSR[num]
	return 0x0F
}

func (controller *BankController) setFSRL(num int, value uint8) AddrMask {
	controller.FSR[num] &= 0xFF00
	controller.FSR[num] |= uint16(value)
	controller.newFSR[num] = controller.FSR[num]
	return 0xFF
}
