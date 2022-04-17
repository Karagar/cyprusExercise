package utils

import (
	"fmt"
	"io/ioutil"
	"os"
)

//MustReadFile panic if can not read file by fileName
func MustReadFile(fileName string) []byte {
	f, err := os.Open(fileName)
	panicOnErr(err)

	byteValue, err := ioutil.ReadAll(f)
	panicOnErr(err)

	err = f.Close()
	panicOnErr(err)

	return byteValue
}

//panicOnErr panic if parameter is not nil
func panicOnErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
