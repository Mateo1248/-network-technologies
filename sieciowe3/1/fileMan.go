package main

import (
	"io/ioutil"
	"log"
	"os"
)

type FileManager struct {
	filename string
}

func (fm FileManager) read() []byte {
	content, err := ioutil.ReadFile(fm.filename)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

func (fm FileManager) write(s string) {
	f, err := os.Create(fm.filename)
	if err != nil {
		log.Fatal(err)
	} else {
		f.WriteString(s)
	}
}
