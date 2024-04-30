package commands

import "fmt"

func CommandHelp(_ *CommandConfig) error {
	fmt.Println("Usage:")
	for _, command := range GetCommands() {
		fmt.Printf("%s - %s\n", command.Name, command.Description)
	}
	return nil
}
