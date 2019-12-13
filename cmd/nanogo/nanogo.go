package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/kkty/nanogo"
)

func main() {
	flag.Parse()
	b, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	program := nanogo.Parse(string(b))
	program.Run(os.Stdout)
}
