package utils

import (
	"bytes"
	qrCodeEncoder "github.com/skip2/go-qrcode"
	qrCodeDecoder "github.com/tuotoo/qrcode"
)

// return a qrCode bytes data
func CreatQRCodeBytes(content string) ([]byte, error) {
	var png []byte
	png, err := qrCodeEncoder.Encode(content, qrCodeEncoder.Medium, 256)
	if nil != err {
		return nil, err
	}
	return png, nil
}

// return qrCode real content
func ParseQRCodeBytes(data []byte) (string, error) {
	buff := new(bytes.Buffer)
	buff.Write(data)
	qrmatrix, err := qrCodeDecoder.Decode(buff)
	if nil != err {
		return "", err
	}
	return qrmatrix.Content, nil
}
