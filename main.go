package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/tamayika/gaq/pkg/gaq"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Cannot read data from stdin. %v", err)
	}

	node := gaq.MustParse(string(data))
	nodeJSON, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		log.Fatalf("Cannot marshal node. %v", err)
	}
	fmt.Println(string(nodeJSON))
}
