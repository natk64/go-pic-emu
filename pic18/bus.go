package pic18

import (
	"log"

	"golang.org/x/exp/constraints"
)

type AddrMask uint8

type AddrType constraints.Integer

type BusReader[T AddrType] interface {
	BusRead(addr T) (data uint8, mask AddrMask)
}

type BusWriter[T AddrType] interface {
	BusWrite(addr T, data uint8) (mask AddrMask)
}

type BusReadWriter[T AddrType] interface {
	BusReader[T]
	BusWriter[T]
}

type (
	DataBusReader        BusReader[uint16]
	DataBusWriter        BusReader[uint16]
	DataBusReadWriter    BusReadWriter[uint16]
	ProgramBusReader     BusReader[uint32]
	ProgramBusWriter     BusReader[uint32]
	ProgramBusReadWriter BusReadWriter[uint32]
)

var (
	_ BusReadWriter[uint8] = MultiBusReadWriter[uint8](nil)
)

// MultiBusReadWriter provides a simple (but inefficient) way to map devices onto a bus.
type MultiBusReadWriter[T AddrType] []BusReadWriter[T]

func (list MultiBusReadWriter[T]) BusRead(addr T) (data uint8, mask AddrMask) {
	for _, reader := range list {
		resultData, resultMask := reader.BusRead(addr)
		resultData &= uint8(resultMask)
		data |= resultData
		mask |= resultMask
	}
	return
}

func (list MultiBusReadWriter[T]) BusWrite(addr T, data uint8) (mask AddrMask) {
	for _, writer := range list {
		resultMask := writer.BusWrite(addr, data)
		mask |= resultMask
	}
	return
}

type nullReader[T AddrType] struct{}
type nullWriter[T AddrType] struct{}
type nullReadWriter[T AddrType] struct{}

func (nullReader[T]) BusRead(addr T) (uint8, AddrMask) {
	return 0, 0
}

func (n nullWriter[T]) BusWrite(addr T, data uint8) AddrMask {
	return 0
}

func (n nullReadWriter[T]) BusRead(addr T) (uint8, AddrMask) {
	return 0, 0
}

func (n nullReadWriter[T]) BusWrite(addr T, data uint8) AddrMask {
	return 0
}

func NullReader[T AddrType]() BusReader[T] {
	return nullReader[T]{}
}

func NullWriter[T AddrType]() BusWriter[T] {
	return nullWriter[T]{}
}

func NullReadWriter[T AddrType]() BusReadWriter[T] {
	return nullReadWriter[T]{}
}

type busPrinter[T AddrType] struct {
	Inner BusReadWriter[T]
}

func (bp busPrinter[T]) BusRead(addr T) (uint8, AddrMask) {
	data, mask := bp.Inner.BusRead(addr)
	if mask == 0 {
		log.Printf("READ  $00 <- $%03x UNIMPLEMENTED\n", addr)
	} else {
		log.Printf("READ  $%02x <- $%03x\n", data&uint8(mask), addr)
	}
	return data, mask
}

func (bp busPrinter[T]) BusWrite(addr T, data uint8) AddrMask {
	mask := bp.Inner.BusWrite(addr, data)
	if mask == 0 {
		log.Printf("WRITE $%02x -> $%03x UNIMPLEMENTED\n", data, addr)
	} else {
		log.Printf("WRITE $%02x -> $%03x\n", data, addr)
	}
	return mask
}

func BusPrinter[T AddrType](bus BusReadWriter[T]) BusReadWriter[T] {
	return busPrinter[T]{Inner: bus}
}
