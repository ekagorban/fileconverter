package main

import (
	"fileconverter/internal/converter"
	"log"
)

func main() {
	if err := converter.Run(); err != nil {
		log.Fatalf("converter run error: %v", err)
		return
	}
}
