package main

import (
	"fmt"
	"io/ioutil"

	v1 "github.com/puppetlabs/nebula-sdk/pkg/workflow/v1"
)

func main() {
	data, err := ioutil.ReadFile("v1/testData/valid.yaml")

	if err != nil {
		fmt.Println("Could not read file")
		return
	}

	validationError := v1.Validate(string(data))

	if validationError != nil {
		fmt.Println(validationError)
	}
}
