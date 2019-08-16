package utils

import (
	"github.com/skip2/go-qrcode"
)

// return a qrCode bytes data
func CreatQRCodeBytes(content string) ([]byte, error) {
	var png []byte
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if nil != err {
		return nil, err
	}
	return png, nil
}
