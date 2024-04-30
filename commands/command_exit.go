package commands

import (
	"fmt"
	"os"
)

func CommandExit(_ *CommandConfig) error {
	fmt.Println("Bye!")
	os.Exit(0)
	return nil
}
