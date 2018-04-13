package main

import (
	_ "wxwork/workers"
	"github.com/benmanns/goworker"
	"fmt"
)

func main() {
	if err := goworker.Work(); err != nil {
		fmt.Println("Error:", err)
	}
}
