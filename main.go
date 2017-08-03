package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/pedrocelso/go-rest-service/lib/config"
)

func main() {
	fmt.Println("\n\n================config.MustInit()====================")
	spew.Dump(config.MustInit())
	fmt.Println("====================================")
}
