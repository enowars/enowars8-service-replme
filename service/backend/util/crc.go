package util

import (
	"math/big"
)

var CRC_DEGREE int64 = 251
var CRC_POLYNOMIAL = "8561cc4ee956c6503c5da0ffacb20feabb3eb142e7645e7ff1a2067fd8e1cfb"

type CRCUtil struct {
	Degree   uint
	DegreeBI *big.Int
	Table    [256]*big.Int
	Mask     *big.Int
}

func genCRCTable(degree int64, poly *big.Int) [256]*big.Int {
	var table [256]*big.Int
	mask := new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(degree), nil),
		big.NewInt(1),
	)
	for i := 0; i < 256; i++ {
		c := big.NewInt(int64(i))
		c = new(big.Int).Lsh(c, uint(degree-8))
		m := new(big.Int).Lsh(big.NewInt(1), uint(degree-1))
		for j := 0; j < 8; j++ {
			if new(big.Int).And(c, m).Cmp(big.NewInt(0)) != 0 {
				c = new(big.Int).Xor(
					new(big.Int).Lsh(c, 1),
					poly,
				)
			} else {
				c = new(big.Int).Lsh(c, 1)
			}
		}
		c = new(big.Int).And(c, mask)
		table[i] = c
	}
	return table
}

func CRC() CRCUtil {

	degree := big.NewInt(CRC_DEGREE)

	mask := new(big.Int).Sub(
		new(big.Int).Exp(
			big.NewInt(2),
			degree,
			nil,
		),
		big.NewInt(1),
	)

	poly, _ := new(big.Int).SetString(CRC_POLYNOMIAL, 16)
	var table = genCRCTable(CRC_DEGREE, poly)

	return CRCUtil{
		Degree:   uint(CRC_DEGREE),
		DegreeBI: degree,
		Table:    table,
		Mask:     mask,
	}
}

func (crc *CRCUtil) Calculate(s []byte) *big.Int {
	result := big.NewInt(0)
	shift := crc.Degree - 8
	byteMask := big.NewInt(0xFF)

	for _, b := range s {
		idx := new(big.Int).And(
			new(big.Int).Xor(
				new(big.Int).Rsh(
					result,
					shift,
				),
				big.NewInt(int64(b)),
			),
			byteMask,
		).Uint64()

		result = new(big.Int).And(
			new(big.Int).Xor(
				crc.Table[idx],
				new(big.Int).Lsh(
					result,
					8,
				),
			),
			crc.Mask,
		)
	}

	return result
}
