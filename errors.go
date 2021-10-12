package etherego

import "errors"

var (
	ErrTranslation = errors.New("unable to translate values")
	ErrKeyNotECDSA = errors.New("error asserting type: publicKey is not of type *ecdsa.PublicKey")
)
