package order

import (
	"fmt"
	"strings"

	"github.com/ole-techwood/PizzaCLI/internal/menu"
)

type OrderService struct {
	menuService *menu.MenuService
}

func NewOrderService() *OrderService {
	return &OrderService{
		menuService: menu.NewMenuService(),
	}
}

var currentOrder Order

func (s *OrderService) AddPizzaToOrder(pizzaName string, quantity int) error {
	pizza, err := s.menuService.GetPizzaByName(pizzaName)
	if err != nil {
		return fmt.Errorf("failed to add pizza to order: %w", err)
	}

	currentOrder.Items = append(currentOrder.Items, OrderItem{
		Pizza:    pizza,
		Quantity: quantity,
	})

	s.printCurrentOrder()

	return nil
}

func (s *OrderService) printCurrentOrder() {
	if len(currentOrder.Items) == 0 {
		return
	}

	var sb strings.Builder
	sb.Grow(len(currentOrder.Items) * 32) // rough estimate per line

	for _, item := range currentOrder.Items {
		price := float64(item.Quantity) * item.Pizza.Price
		fmt.Fprintf(&sb, "%s | %d | $%.2f\n", item.Pizza.Name, item.Quantity, price)
	}

	fmt.Println("Current order:")
	fmt.Print(sb.String())
}
