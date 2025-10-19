package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jawahars16/jebi/cmd"
)

func main() {
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
