package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
	"github.com/marcoguerri/go-tpm-tcti/abrmd"
)

var certPath = flag.String("path", "", "path where to write EK certificate")

func main() {
	flag.Parse()
	if *certPath == "" {
		log.Fatalf("please provide path of certificate via `path` flag")
	}

	log.Printf("creating new broker connection...")
	broker, err := abrmd.NewBroker()
	if err != nil {
		log.Panicf("could not create new broker: %v", err)
	}
	buff, err := tpm2.NVRead(broker, tpmutil.Handle(0x1c0000a))
	if err != nil {
		log.Fatalf("could not read nvram index: %v", err)
	}
	err = ioutil.WriteFile(*certPath, buff, 0644)
	if err != nil {
		log.Panicf("could not write certificate: %w", err)
	}
	log.Printf("certificate written in %s", *certPath)
}
