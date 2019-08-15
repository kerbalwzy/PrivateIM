package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
)

// return the hash value of data. user MD5
func BytesDataHash(data []byte) string {
	h := md5.New()
	h.Write(data)
	hashValue := hex.EncodeToString(h.Sum(nil))
	return hashValue
}

// return the hash value of an file which recv in Gin
func GinFormFileHash(file *multipart.FileHeader) (string, []byte, error) {
	// read data from file
	buff := new(bytes.Buffer)
	fp, err := file.Open()
	defer fp.Close()
	if nil != err {
		return "", nil, err
	}
	n, err := io.Copy(buff, fp)
	if nil != err || n == 0 {
		msg := fmt.Sprintf("read data from file error:%v", err)
		return "", nil, errors.New(msg)
	}
	return BytesDataHash(buff.Bytes()), buff.Bytes(), nil

}

// upload the file to local dir
func UploadFileToLocal(data []byte, path string) error {
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}

// todo upload file data to oss or other static file server
func UploadFileToCloud() error {
	return nil
}
