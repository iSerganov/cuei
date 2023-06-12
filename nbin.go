package cuei

import (
	"fmt"
	"math/big"
)

type Nbin struct {
	Bites big.Int
}

// Add bytes to NBin.Bites for encoding
func (nb *Nbin) AddBytes(str string, nbits uint) {
	t := new(big.Int)
	t.SetBytes([]byte(str))
	o := nb.Bites.Lsh(&nb.Bites, nbits)
	nb.Bites = *nb.Bites.Add(o, t)
}

// Add64 left shift NBin.Bites by nbits and add uint64 val
func (nb *Nbin) Add64(val uint64, nbits uint) {
	t := new(big.Int)
	t.SetUint64(val)
	o := nb.Bites.Lsh(&nb.Bites, nbits)
	nb.Bites = *nb.Bites.Add(o, t)
}

// Add32 left shift NBin.Bites by nbits and add uint32 val
func (nb *Nbin) Add32(val uint32, nbits uint) {
	u := uint64(val)
	nb.Add64(u, nbits)
}

// Add16 left shift NBin.Bites by nbits and add uint16 val
func (nb *Nbin) Add16(val uint16, nbits uint) {
	u := uint64(val)
	nb.Add64(u, nbits)
}

// Add8 left shift NBin.Bites by nbits and add uint8 val
func (nb *Nbin) Add8(val uint8, nbits uint) {
	u := uint64(val)
	nb.Add64(u, nbits)
}

// AddFlag left shift NBin.Bites by 1 and add bool val
func (nb *Nbin) AddFlag(val bool) {
	if val == true {
		nb.Add64(1, 1)
	} else {
		nb.Add64(0, 1)
	}
}

// Add90k left shift NBin.Bites by nbits and add val as ticks
func (nb *Nbin) Add90k(val float64, nbits uint) {
	u := uint64(val * float64(90000.0))
	nb.Add64(u, nbits)
}

// AddHex64 left shift NBin.Bites by nbits and add hex string as uint64
func (nb *Nbin) AddHex64(val string, nbits uint) {
	u := new(big.Int)
	_, err := fmt.Sscan(val, u)
	if err != nil {
		fmt.Println("error scanning value:", err)
	} else {
		fmt.Println(u.Uint64())
		nb.Add64(u.Uint64(), nbits)
	}
}

// Reserve num bits by setting them to 1
func (nb *Nbin) Reserve(num int) {

	for i := 0; i < num; i++ {
		nb.Add64(1, 1)
	}
}
