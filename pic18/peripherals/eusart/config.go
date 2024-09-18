package eusart

import (
	"fmt"

	"github.com/natk64/go-pic-emu/pic18"
)

func New(num int, interrupts *pic18.InterruptController) (eusart *EUSART) {
	if num == 1 {
		eusart = &EUSART{
			ModeChange:  func() {},
			TxInterrupt: interrupts.CreateInterrupt(pic18.PeripheralInterrupt("TX1", 1, 4)),
			RxInterrupt: interrupts.CreateInterrupt(pic18.PeripheralInterrupt("RC1", 1, 5)),
			Registers: Registers{
				TXSTAx:   pic18.Registers.TXSTA1,
				RCSTAx:   pic18.Registers.RCSTA1,
				TXREGx:   pic18.Registers.TXREG1,
				RCREGx:   pic18.Registers.RCREG1,
				BAUDCONx: pic18.Registers.BAUDCON1,
				SPBRGHx:  pic18.Registers.SPBRGH1,
				SPBRGx:   pic18.Registers.SPBRG1,
			},
		}
	} else if num == 2 {
		eusart = &EUSART{
			ModeChange:  func() {},
			TxInterrupt: interrupts.CreateInterrupt(pic18.PeripheralInterrupt("TX2", 3, 4)),
			RxInterrupt: interrupts.CreateInterrupt(pic18.PeripheralInterrupt("RC2", 3, 5)),
			Registers: Registers{
				TXSTAx:   pic18.Registers.TXSTA2,
				RCSTAx:   pic18.Registers.RCSTA2,
				TXREGx:   pic18.Registers.TXREG2,
				RCREGx:   pic18.Registers.RCREG2,
				BAUDCONx: pic18.Registers.BAUDCON2,
				SPBRGHx:  pic18.Registers.SPBRGH2,
				SPBRGx:   pic18.Registers.SPBRG2,
			},
		}
	} else {
		return nil
	}

	eusart.Transmit = func(data uint8, bit9 bool) {
		go func() {
			fmt.Print(string(rune(data)))
			eusart.TXDone()
		}()
	}

	return eusart
}
