package main

import (
	"log"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
	"github.com/marcoguerri/tpm2-tcti/abrmd"
)

func main() {
	log.Printf("creating new broker connection...")
	broker, err := abrmd.NewBroker()
	if err != nil {
		log.Panicf("could not create new broker: %v", err)
	}
	buff, err := tpm2.NVRead(broker, tpmutil.Handle(0x1c0000a))
	if err != nil {
		log.Fatalf("could not read nvram index: %v", err)
	}
	log.Printf("received buffer %x", buff)
}
