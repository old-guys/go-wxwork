package lib_wxwork

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

var (
	Utils UtilsConfig
)

type UtilsConfig struct {
	PKCS7Padding pKCS7PaddingConfig
	NetworkOrder networkOrderConfig
}

type pKCS7PaddingConfig struct {

}

type networkOrderConfig struct {

}

const BLOCK_SIZE = 32

func (c *pKCS7PaddingConfig) AddPadding(data []byte) []byte {
	length := len(data)
	padding_amount := BLOCK_SIZE - (length % BLOCK_SIZE)

	if padding_amount == 0 {
		padding_amount = BLOCK_SIZE
	}

	char := rune(padding_amount)
	for index := 0; index < padding_amount; index++ {
		data = bytes.Join([][]byte{data, []byte(string(char))}, nil)
	}

	return data
}

func (c *pKCS7PaddingConfig) RemovePadding(data []byte) []byte {
	length := len(data)
	padding_amount := int(data[length - 1])

	if padding_amount < 1 || padding_amount > BLOCK_SIZE {
		padding_amount = 0
	}

	return data[:(length - padding_amount)]
}

func (c *networkOrderConfig) BytesToInt(data []byte) int32 {
	var num int32

	bytesBuffer := bytes.NewBuffer(data)
	binary.Read(bytesBuffer, binary.BigEndian, &num)

	return num
}

func (c *networkOrderConfig) IntToBytes(number int) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int32(number))
	if err != nil {
		fmt.Println("IntToBytes Binary write err:", err)
	}

	return buf.Bytes()
}