package pic18

type AluStatus uint8

const (
	statusC AluStatus = 1 << iota
	statusDC
	statusZ
	statusOV
	statusN
)

func (status AluStatus) C() bool  { return status&statusC != 0 }
func (status AluStatus) DC() bool { return status&statusDC != 0 }
func (status AluStatus) Z() bool  { return status&statusZ != 0 }
func (status AluStatus) OV() bool { return status&statusOV != 0 }
func (status AluStatus) N() bool  { return status&statusN != 0 }

type ALU struct {
	status      AluStatus
	productHigh uint8
	productLow  uint8
}

func (alu *ALU) Add(a, b uint8) uint8 {
	result := uint16(a) + uint16(b)
	result8bit := uint8(result)
	alu.status = updateStatusAll(alu.status, result)
	return result8bit
}

func (alu *ALU) Sub(a, b uint8) uint8 {
	result := uint16(a) - uint16(b)
	result8bit := uint8(result)
	alu.status = updateStatusAll(alu.status, result)
	return result8bit
}

func (alu *ALU) AddWithCarry(a, b uint8) uint8 {
	carryBit := alu.status & statusC
	result := uint16(a) + uint16(b) + uint16(carryBit)
	alu.status = updateStatusAll(alu.status, result)
	return uint8(result)
}

func (alu *ALU) SubWithBorrow(a, b uint8) uint8 {
	var borrow uint16
	if (alu.status & statusC) == 0 {
		borrow = 1
	}

	result := uint16(a) + uint16(b) + uint16(borrow)
	alu.status = updateStatusAll(alu.status, result)
	return uint8(result)
}

func (alu *ALU) And(a, b uint8) uint8 {
	result := a & b
	alu.status = updateStatusN_Z(alu.status, result)
	return result
}

func (alu *ALU) Or(a, b uint8) uint8 {
	result := a | b
	alu.status = updateStatusN_Z(alu.status, result)
	return result
}

func (alu *ALU) Xor(a, b uint8) uint8 {
	result := a ^ b
	alu.status = updateStatusN_Z(alu.status, result)
	return result
}

func (alu *ALU) Negate(a uint8) uint8 {
	result := uint16(^a + 1)
	alu.status = updateStatusAll(alu.status, result)
	return uint8(result)
}

func (alu *ALU) RotateLeft(a uint8) uint8 {
	value := uint16(a)
	value <<= 1
	if value&(1<<8) != 0 {
		value |= 1
	}

	alu.status = updateStatusN_Z(alu.status, uint8(value))
	return uint8(value)
}

func (alu *ALU) RotateRight(a uint8) uint8 {
	alu.status &= ^statusC
	result := a >> 1
	if a&1 != 0 {
		result |= (1 << 7)
	}

	alu.status = updateStatusN_Z(alu.status, result)
	return result
}

func (alu *ALU) RotateLeftCarry(a uint8) uint8 {
	value := uint16(a)
	value <<= 1
	if alu.status&statusC != 0 {
		value |= 1
	}

	alu.status &= ^statusC

	if value&0b100000000 != 0 {
		alu.status |= statusC
	}

	alu.status = updateStatusN_Z(alu.status, uint8(value))
	return uint8(value)
}

func (alu *ALU) RotateRightCarry(a uint8) uint8 {
	alu.status &= ^statusC
	result := a >> 1
	if a&1 != 0 {
		alu.status |= statusC
	}

	alu.status = updateStatusN_Z(alu.status, result)
	return result
}

func (alu *ALU) DecimalAdjust(a uint8) uint8 {
	low := a & 0x0F
	high := a & 0xF0
	var result uint16

	if low > 9 || (alu.status&statusDC != 0) {
		result = uint16(low) + 6
	} else {
		result = uint16(low)
	}

	if high > 9 || (alu.status&statusC != 0) {
		result |= uint16((high + 6) << 4)
	} else {
		result |= uint16(high << 4)
	}

	alu.status &= ^statusC
	if result&(1<<8) != 0 {
		alu.status |= statusC
	}

	return uint8(result)
}

func (alu *ALU) Complement(a uint8) uint8 {
	result := ^a
	alu.status = updateStatusN_Z(alu.status, result)
	return result
}

func (alu *ALU) Mul(a, b uint8) {
	result := uint16(a) * uint16(b)
	alu.productHigh = uint8(result >> 8)
	alu.productLow = uint8(result)
}

// TODO: DC Flag
func updateStatusAll(status AluStatus, result uint16) AluStatus {
	var newStatus AluStatus
	if (result & 0b10000000) != 0 {
		newStatus |= statusN
	}
	if (status & statusN) != (newStatus & statusN) {
		newStatus |= statusOV
	}
	if (result & 0xFF) == 0 {
		newStatus |= statusZ
	}
	if (result & 0b100000000) != 0 {
		newStatus |= statusC
	}

	return newStatus
}

func updateStatusN_Z(status AluStatus, result uint8) AluStatus {
	status &= ^statusZ
	status &= ^statusN
	if result&0b10000000 != 0 {
		status = statusN
	}
	if result == 0 {
		status |= statusZ
	}
	return status
}

func (alu ALU) BusRead(addr uint16) (uint8, AddrMask) {
	if addr == Registers.STATUS {
		return uint8(alu.status), 0x1F
	}
	return 0, 0
}

func (alu *ALU) BusWrite(addr uint16, data uint8) AddrMask {
	if addr == Registers.STATUS {
		alu.status = AluStatus(data & 0x1F)
		return 0x1F
	}
	return 0
}
