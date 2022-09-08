package main

import (
	"fmt"
	"log"
	"os"
)

func hello(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Missing name: hello <name>")
	}
	fmt.Println(hello(os.Args[1]))
}
