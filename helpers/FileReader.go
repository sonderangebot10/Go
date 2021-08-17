package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func GenericReadFile(fileLoc string, el interface{}) {

	jsonFile, err := os.Open(fileLoc)

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &el)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	jsonFile.Close()
}
