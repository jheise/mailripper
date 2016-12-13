package main

import (
	// standard
	"fmt"
	"io/ioutil"
	"os"

	// external
	"github.com/sg3des/eml"
)

func usage() {
	fmt.Println("usage: mailripper filename")
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	srcFile := os.Args[1]

	srcData, err := ioutil.ReadFile(srcFile)
	if err != nil {
		panic(err)
	}

	rawMsg, err := eml.ParseRaw(srcData)
	if err != nil {
		panic(err)
	}

	msg, err := eml.Process(rawMsg)
	if err != nil {
		panic(err)
	}

	for _, attachment := range msg.Attachments {
		fmt.Println(attachment.Filename)
		ioutil.WriteFile(attachment.Filename, attachment.Data, 0644)
	}
}
