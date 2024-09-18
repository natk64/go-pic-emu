package pic18

type TableRWController struct {
	tablePointer uint32
	tableLatch   uint8

	ProgramBus BusReadWriter[uint32]
}

type TableAction int

const (
	TableNone    TableAction = 0
	TablePostInc TableAction = 1
	TablePostDec TableAction = 2
	TablePreInc  TableAction = 3
)

const (
	TBLPTRU = 0xFF8
	TBLPTRH = 0xFF7
	TBLPTRL = 0xFF6
	TABLAT  = 0xFF5
)

func (controller *TableRWController) BusRead(addr uint16) (uint8, AddrMask) {
	switch addr {
	case TABLAT:
		return controller.tableLatch, 0xFF
	case TBLPTRL:
		return uint8(controller.tablePointer), 0xFF
	case TBLPTRH:
		return uint8(controller.tablePointer >> 8), 0xFF
	case TBLPTRU:
		return uint8(controller.tablePointer>>16) & 0x3F, 0xFF
	default:
		return 0, 0
	}
}

func (controller *TableRWController) BusWrite(addr uint16, data uint8) AddrMask {
	switch addr {
	case TABLAT:
		controller.tableLatch = data
		return 0xFF
	case TBLPTRL:
		controller.tablePointer = (controller.tablePointer & 0x3FFF00) | uint32(data)
		return 0xFF
	case TBLPTRH:
		controller.tablePointer = (controller.tablePointer & 0x3F00FF) | (uint32(data) << 8)
		return 0xFF
	case TBLPTRU:
		controller.tablePointer = (controller.tablePointer & 0x00FFFF) | (uint32(data&0x3F) << 16)
		return 0xFF
	default:
		return 0
	}
}

func (controller *TableRWController) TableRead(action TableAction) {
	if action == TablePreInc {
		controller.tablePointer = (controller.tablePointer + 1) & 0x3FFF00
	}

	controller.tableLatch, _ = controller.ProgramBus.BusRead(controller.tablePointer)

	if action == TablePostInc {
		controller.tablePointer = (controller.tablePointer + 1) & 0x3FFF00
	} else if action == TablePostDec {
		controller.tablePointer = (controller.tablePointer - 1) & 0x3FFF00
	}
}

func (controller *TableRWController) TableWrite(action TableAction) {
	if action == TablePreInc {
		controller.tablePointer = (controller.tablePointer + 1) & 0x3FFF00
	}

	controller.ProgramBus.BusWrite(controller.tablePointer, controller.tableLatch)

	if action == TablePostInc {
		controller.tablePointer = (controller.tablePointer + 1) & 0x3FFF00
	} else if action == TablePostDec {
		controller.tablePointer = (controller.tablePointer - 1) & 0x3FFF00
	}
}
