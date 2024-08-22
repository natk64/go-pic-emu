package pic18

import (
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
