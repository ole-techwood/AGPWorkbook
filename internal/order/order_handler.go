package order

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ole-techwood/PizzaCLI/internal/help"
)

type OrderHandler struct {
	orderService OrderService
	helpService  help.HelpService
}

// Fields
var (
	errPizzaNameNotSpecified = errors.New("pizza name not specified")
	errQuantityTooLow        = errors.New("quantity must be at least 1")
)

// Constructor
func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		orderService: OrderService{},
		helpService:  help.NewHelpService(),
	}
}

// Methods
func (h *OrderHandler) AddPizzaToOrder(args []string) {
	pizzaName, quantity, err := h.parseAddArgs(args)

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("")
		h.helpService.GetHelp()
		return
	}

	err = h.orderService.AddPizzaToOrder(pizzaName, quantity)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("")
		h.helpService.GetHelp()
	}
}

func (h *OrderHandler) Checkout() {
	err := h.orderService.Checkout()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("")
		h.helpService.GetHelp()
	}
}

func (h *OrderHandler) parseAddArgs(args []string) (string, int, error) {
	if len(args) == 0 {
		return "", 0, errPizzaNameNotSpecified
	}

	// By default quantity is always 1
	quantity := 1

	lastArg := args[len(args)-1]
	if parsedQuantity, err := strconv.Atoi(lastArg); err == nil {
		quantity = parsedQuantity
		args = args[:len(args)-1]
	}

	name := strings.Join(args, " ")
	if name == "" {
		return "", 0, errPizzaNameNotSpecified
	}

	if quantity < 1 {
		return "", 0, errQuantityTooLow
	}

	return name, quantity, nil
}
