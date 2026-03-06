package order

import (
	"fmt"
	"strings"

	"github.com/ole-techwood/PizzaCLI/internal/menu"
)

type OrderService struct {
	sb           strings.Builder
	menuService  *menu.MenuService
	currentOrder Order
}

func NewOrderService() *OrderService {
	return &OrderService{
		menuService: menu.NewMenuService(),
	}
}

func (s *OrderService) AddPizzaToOrder(pizzaName string, quantity int) error {
	pizza, err := s.menuService.GetPizzaByName(pizzaName)
	if err != nil {
		return fmt.Errorf("failed to add pizza to order: %w", err)
	}

	s.currentOrder.Items = append(s.currentOrder.Items, OrderItem{
		Pizza:    pizza,
		Quantity: quantity,
	})

	s.printCurrentOrder()

	return nil
}

func (s *OrderService) Checkout() error {
	if len(s.currentOrder.Items) == 0 {
		return fmt.Errorf("no items in order")
	}

	s.printCurrentOrder()
	fmt.Println("Order placed! Thank you.")
	s.currentOrder = Order{}

	return nil
}

func (s *OrderService) printCurrentOrder() {
	if len(s.currentOrder.Items) == 0 {
		return
	}

	s.sb.Reset()
	s.sb.Grow(len(s.currentOrder.Items) * 32)

	var total float64
	for _, item := range s.currentOrder.Items {
		itemTotal := float64(item.Quantity) * item.Pizza.Price
		total += itemTotal
		fmt.Fprintf(&s.sb, "%d × %s @ $%.2f = $%.2f\n", item.Quantity, item.Pizza.Name, item.Pizza.Price, itemTotal)
	}

	fmt.Println("Current order:")
	fmt.Print(s.sb.String())
	fmt.Printf("**Total: $%.2f**\n", total)
}
