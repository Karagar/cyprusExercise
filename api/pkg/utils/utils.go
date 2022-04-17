package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

//MustReadFile panic if can not read file by fileName
func MustReadFile(fileName string) []byte {
	f, err := os.Open(fileName)
	PanicOnErr(err)

	byteValue, err := ioutil.ReadAll(f)
	PanicOnErr(err)

	err = f.Close()
	PanicOnErr(err)

	return byteValue
}

//PanicOnErr panic if parameter is not nil
func PanicOnErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func HandleUuid(id []byte) string {
	return fmt.Sprintf("%X-%X-%X-%X-%X", id[0:4], id[4:6], id[6:8], id[8:10], id[10:])
}

func ReadJsonBody(rawBody io.ReadCloser, s interface{}) error {
	body, err := ioutil.ReadAll(rawBody)
	defer rawBody.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, s)
	return err
}
