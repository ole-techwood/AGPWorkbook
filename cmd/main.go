package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ole-techwood/PizzaCLI/internal/help"
	"github.com/ole-techwood/PizzaCLI/internal/menu"
	"github.com/ole-techwood/PizzaCLI/internal/order"
)

func main() {
	helpHandler := help.NewHelpHandler()
	menuHandler := menu.NewMenuHandler()
	orderHandler := order.NewOrderHandler()

	fmt.Println("Welcome to PizzaCLI! Type 'help' for available commands.")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("pizza> ")

		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		command, rest, _ := strings.Cut(line, " ")
		var args []string
		if rest != "" {
			args = strings.Fields(rest)
		}

		switch command {
		case "menu":
			menuHandler.GetMenu()
		case "add":
			orderHandler.AddPizzaToOrder(args)
		case "checkout":
			orderHandler.Checkout()
		case "help":
			helpHandler.GetHelp()
		case "exit", "quit":
			fmt.Println("Bye!")
			return
		default:
			fmt.Printf("Unknown command: %s\n", command)
			fmt.Println()
			helpHandler.GetHelp()
		}

		fmt.Println()
	}
}
