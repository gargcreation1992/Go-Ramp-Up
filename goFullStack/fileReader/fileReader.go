package fileReader

import "io/ioutil"

type ioClient struct {
}

func (i *ioClient) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func NewClient() *ioClient {
	return &ioClient{}
}
