package main

import (
	"fmt"
	"os"

	"github.com/ole-techwood/PizzaCLI/internal/help"
	"github.com/ole-techwood/PizzaCLI/internal/menu"
	"github.com/ole-techwood/PizzaCLI/internal/order"
)

func main() {
	helpHandler := help.NewHelpHandler()

	if len(os.Args) < 2 {
		fmt.Println("Command not specified")
		fmt.Println("")
		helpHandler.GetHelp()

		return
	}

	command := os.Args[1]

	switch command {
	case "menu":
		menuHandler := menu.NewMenuHandler()

		menuHandler.GetMenu()
	case "add":
		orderHandler := order.NewOrderHandler()

		orderHandler.AddPizzaToOrder(os.Args[2:])
	case "help":
		helpHandler.GetHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("")
		helpHandler.GetHelp()
	}
}
