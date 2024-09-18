package eusart

import (
	"github.com/natk64/go-pic-emu/pic18"
)

type EUSART struct {
	synchronous          bool
	tx_en                bool
	tx_9bit_en           bool
	sync_clk_master_mode bool
	async_send_break     bool
	async_high_speed     bool

	sp_en              bool
	rx_9bit_en         bool
	single_receive     bool
	continuous_receive bool
	address_detect     bool

	framing_err bool
	overrun_err bool

	rx_active           bool
	async_invert_rx     bool
	wakeup_en           bool
	clock_data_polarity bool
	baud_rate_16bit_en  bool

	baud_rate uint16

	tx_reg  uint8
	rx_reg  uint8
	tx_bit9 bool
	rx_bit9 bool

	tsr_loaded   bool
	txreg_loaded bool

	ModeChange  func()
	Transmit    func(data uint8, bit9 bool)
	RxInterrupt pic18.Interrupt
	TxInterrupt pic18.Interrupt

	Registers Registers
}

type Registers struct {
	TXSTAx   uint16
	RCSTAx   uint16
	TXREGx   uint16
	RCREGx   uint16
	BAUDCONx uint16
	SPBRGHx  uint16
	SPBRGx   uint16
}

func (eusart *EUSART) loadTSR() {
	eusart.txreg_loaded = false
	eusart.tsr_loaded = true
	eusart.Transmit(eusart.tx_reg, eusart.tx_bit9 && eusart.tx_9bit_en)
	eusart.TxInterrupt.Raise()
}

func (eusart *EUSART) BusRead(addr uint16) (uint8, pic18.AddrMask) {
	read_bits := func(b7, b6, b5, b4, b3, b2, b1, b0 bool) uint8 {
		return (bInt(b7) << 7) | (bInt(b6) << 6) | (bInt(b5) << 5) | (bInt(b4) << 4) |
			(bInt(b3) << 3) | (bInt(b2) << 2) | (bInt(b1) << 1) | (bInt(b0) << 0)
	}

	if addr == eusart.Registers.TXSTAx {
		return read_bits(
			eusart.sync_clk_master_mode,
			eusart.tx_9bit_en,
			eusart.tx_en,
			eusart.synchronous,
			eusart.async_send_break,
			eusart.async_high_speed,
			!eusart.tsr_loaded,
			eusart.tx_bit9,
		), 0xFF
	} else if addr == eusart.Registers.RCSTAx {
		return read_bits(
			eusart.sp_en,
			eusart.rx_9bit_en,
			eusart.single_receive,
			eusart.continuous_receive,
			eusart.address_detect,
			eusart.framing_err,
			eusart.overrun_err,
			eusart.rx_bit9,
		), 0xFF
	} else if addr == eusart.Registers.TXREGx {
		return eusart.tx_reg, 0xFF
	} else if addr == eusart.Registers.RCREGx {
		return eusart.rx_reg, 0xFF
	} else if addr == eusart.Registers.BAUDCONx {
		return read_bits(
			false,
			!eusart.rx_active,
			eusart.async_invert_rx,
			eusart.clock_data_polarity,
			eusart.baud_rate_16bit_en,
			false,
			eusart.wakeup_en,
			false,
		), 0xFB
	} else if addr == eusart.Registers.SPBRGx {
		return uint8(eusart.baud_rate & 0xFF), 0xFF
	} else if addr == eusart.Registers.SPBRGHx {
		return uint8((eusart.baud_rate >> 8) & 0xFF), 0xFF
	}

	return 0, 0
}

func (eusart *EUSART) BusWrite(addr uint16, value uint8) pic18.AddrMask {
	if addr == eusart.Registers.TXSTAx {
		eusart.sync_clk_master_mode = bitTest(value, 7)
		eusart.tx_9bit_en = bitTest(value, 6)
		eusart.tx_en = bitTest(value, 5)
		eusart.synchronous = bitTest(value, 4)
		eusart.async_send_break = bitTest(value, 3)
		eusart.async_high_speed = bitTest(value, 2)
		// Bit 1 is read only.
		eusart.tx_bit9 = bitTest(value, 0)

		eusart.ModeChange()
		if eusart.tx_en {
			eusart.TxInterrupt.Raise()
		}

		return 0xFF
	} else if addr == eusart.Registers.RCSTAx {
		eusart.sp_en = bitTest(value, 7)
		eusart.rx_9bit_en = bitTest(value, 6)
		eusart.single_receive = bitTest(value, 5)
		eusart.continuous_receive = bitTest(value, 4)
		eusart.address_detect = bitTest(value, 3)
		eusart.rx_bit9 = bitTest(value, 0)

		if !eusart.continuous_receive {
			eusart.overrun_err = false
		}

		eusart.ModeChange()
		return 0xFF
	} else if addr == eusart.Registers.TXREGx {
		eusart.tx_reg = value
		eusart.txreg_loaded = true
		eusart.TxInterrupt.Clear()
		if !eusart.tsr_loaded {
			eusart.loadTSR()
		}
		return 0xFF
	} else if addr == eusart.Registers.RCREGx {
		return 0xFF
	} else if addr == eusart.Registers.BAUDCONx {
		eusart.async_invert_rx = bitTest(value, 5)
		eusart.clock_data_polarity = bitTest(value, 4)
		eusart.baud_rate_16bit_en = bitTest(value, 3)
		eusart.wakeup_en = bitTest(value, 1)
		eusart.ModeChange()
		return 0xFB
	} else if addr == eusart.Registers.SPBRGx {
		eusart.baud_rate = (eusart.baud_rate & 0xFF00) | uint16(value)
		eusart.ModeChange()
		return 0xFF
	} else if addr == eusart.Registers.SPBRGHx {
		eusart.baud_rate = (eusart.baud_rate & 0x00FF) | (uint16(value) << 8)
		eusart.ModeChange()
		return 0xFF
	}

	return 0x00
}

func (eusart *EUSART) ImportRX(data uint8, bit9 bool) {
	if !eusart.sp_en {
		return
	}

	if !eusart.continuous_receive && !eusart.single_receive {
		return
	}

	eusart.single_receive = false
	eusart.rx_reg = data
	if eusart.rx_9bit_en {
		eusart.rx_bit9 = bit9
	}

	eusart.RxInterrupt.Raise()
}

func (eusart *EUSART) TXDone() {
	if eusart.txreg_loaded {
		eusart.loadTSR()
	} else {
		eusart.tsr_loaded = false
	}
}

func bInt(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func bitTest(val uint8, bit int) bool {
	return val&(1<<bit) != 0
}
