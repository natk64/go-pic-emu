package pic18

type Stack struct {
	data    *[31]uint32
	pointer uint8

	full      bool
	underflow bool
}

func (stack *Stack) Push(value uint32) bool {
	if int(stack.pointer) == len(stack.data) {
		// Stack was already full, this doesn't cause a reset so just return true
		stack.full = true
		return true
	}

	stack.data[stack.pointer] = value
	stack.pointer++

	if int(stack.pointer) == len(stack.data) {
		// Stack became full
		stack.full = true
	}

	return !stack.full
}

func (stack *Stack) Pop() bool {
	if stack.pointer == 0 {
		stack.underflow = true
		return false
	}

	stack.pointer--
	return true
}

func (stack *Stack) Top() uint32 {
	if stack.pointer == 0 {
		return 0
	}
	return stack.data[stack.pointer-1]
}

func (stack *Stack) SetTop(value uint32) {
	if stack.pointer == 0 {
		return
	}
	stack.data[stack.pointer-1] = value
}

func (stack *Stack) BusRead(addr uint16) (uint8, AddrMask) {
	switch addr {
	case Registers.STKPTR:
		stkptr := stack.pointer & 0x1F
		if stack.full {
			stkptr |= 1 << 7
		}
		if stack.underflow {
			stkptr |= 1 << 6
		}
		return stkptr, 0xDF
	case Registers.TOSL:
		return uint8(stack.Top() & 0xFF), 0xFF
	case Registers.TOSH:
		return uint8((stack.Top() >> 8) & 0xFF), 0xFF
	case Registers.TOSU:
		return uint8((stack.Top() >> 16) & 0xFF), 0xFF
	}

	return 0, 0
}

func (stack *Stack) BusWrite(addr uint16, data uint8) AddrMask {
	switch addr {
	case Registers.STKPTR:
		if stack.full {
			stack.full = (data & (1 << 7)) != 0
		}
		if stack.underflow {
			stack.underflow = (data & (1 << 6)) != 0
		}
		stack.pointer = data & 0x1F
		return 0xDF
	case Registers.TOSL:
		stack.SetTop((stack.Top() & 0xFFFF00) | uint32(data))
		return 0xFF
	case Registers.TOSH:
		stack.SetTop((stack.Top() & 0xFF00FF) | (uint32(data) << 8))
		return 0xFF
	case Registers.TOSU:
		stack.SetTop((stack.Top() & 0x00FFFF) | (uint32(data) << 16))
		return 0xFF
	}

	return 0x00
}
