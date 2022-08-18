package main

import (
	"fmt"
	"os"
	"time"

    "IslamWalid-FuseFS/fs"
)

type structure struct {
	String       string
	Int          int
	Bool         bool
	SubStructure subStructure
}

type subStructure struct {
	Float float32
}

func Routine(input *structure) {
	time.Sleep(time.Second * 5)
	input.String = "new string"
} 

func main() {
	var err error
	input := &structure{
		String: "str",
		Int:    18,
		Bool:   true,
		SubStructure: subStructure{
			Float: 1.3,
		},
	}

	err = os.MkdirAll("mnt", 0777)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	go Routine(input)
	err = fs.Mount("mnt", input)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
