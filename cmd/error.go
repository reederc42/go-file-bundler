package cmd

import (
	"fmt"
	"os"
)

func hdlerr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
