package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

// HandleUuid used to enable people to work with mssql GUID, not just machines
func HandleUuid(id []byte) string {
	return fmt.Sprintf("%X-%X-%X-%X-%X", id[0:4], id[4:6], id[6:8], id[8:10], id[10:])
}

// ReadJsonBody used to read request body to struct s
func ReadJsonBody(rawBody io.ReadCloser, s interface{}) error {
	body, err := ioutil.ReadAll(rawBody)
	defer rawBody.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, s)
	return err
}

// ReadUserIP used to check user IP address. If we dont care that user can spoof it - set it to True
func ReadUserIP(r *http.Request, isCheckHeader bool) string {
	IPAddress := r.RemoteAddr
	if IPAddress != "" {
		return strings.Split(IPAddress, ":")[0]
	}

	if isCheckHeader {
		IPAddress = r.Header.Get("X-Real-Ip")
		if IPAddress == "" {
			IPAddress = r.Header.Get("X-Forwarded-For")
		}
		if IPAddress != "" {
			return strings.Split(IPAddress, ",")[0]
		}
	}

	return IPAddress
}
