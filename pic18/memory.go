package pic18

type Memory[T AddrType] struct {
	Offset int
	Data   []byte
}

func (memory Memory[T]) BusRead(addr T) (uint8, AddrMask) {
	index := int(addr) - memory.Offset
	if index < 0 || index >= len(memory.Data) {
		return 0, 0
	}

	return memory.Data[index], 0xFF
}

func (memory Memory[T]) BusWrite(addr T, data uint8) AddrMask {
	index := int(addr) - memory.Offset
	if index < 0 || index >= len(memory.Data) {
		return 0
	}

	memory.Data[index] = data
	return 0xFF
}
