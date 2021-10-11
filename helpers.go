package etherego

import (
	"math/big"
)

/*
   TODO: make code more DRY.
   		 optimize big numbers
*/

// WeiToEth returns wei amount in eth
func WeiToEth(wei *big.Int) (*big.Float, error) {
	var (
		ok     bool
		checks = make([]bool, 0)
	)

	weiLocal, ok := big.NewFloat(0).SetString(wei.String())
	checks = append(checks, ok)

	divider, ok := big.NewFloat(0).SetString(WeiEth.String())
	checks = append(checks, ok)

	ethLocal := big.NewFloat(0).Quo(weiLocal, divider)
	if any(checks) {
		return nil, ErrTranslation
	}

	return ethLocal, nil
}

// EthToWei returns eth admount in wei
func EthToWei(eth *big.Float) (*big.Int, error) {
	var (
		ok bool
	)

	multiplier, ok := big.NewFloat(0).SetString(WeiEth.String())

	wei := big.NewFloat(0).Mul(eth, multiplier)
	v, _ := wei.Uint64()

	if !ok {
		return nil, ErrTranslation
	}
	return big.NewInt(0).SetUint64(v), nil
}

// GWeiToEth returns ETH amount of gwei
func GWeiToEth(gwei *big.Int) (*big.Float, error) {
	var (
		ok     bool
		checks = make([]bool, 0)
	)

	weiLocal, ok := big.NewFloat(0).SetString(gwei.String())
	checks = append(checks, ok)

	divider, ok := big.NewFloat(0).SetString(GweiEth.String())
	checks = append(checks, ok)

	ethLocal := big.NewFloat(0).Quo(weiLocal, divider)
	if any(checks) {
		return nil, ErrTranslation
	}

	return ethLocal, nil
}

func EthToGWei(eth *big.Float) (*big.Int, error) {
	var ok bool
	multiplier, ok := big.NewFloat(0).SetString(GweiEth.String())

	wei := big.NewFloat(0).Mul(eth, multiplier)
	v, _ := wei.Uint64()

	if !ok {
		return nil, ErrTranslation
	}
	return big.NewInt(0).SetUint64(v), nil
}

// any checks whether is any element in array is false. If element is false, any returns true
func any(bools []bool) bool {
	for _, b := range bools {
		if !b {
			return true
		}
	}
	return false
}
