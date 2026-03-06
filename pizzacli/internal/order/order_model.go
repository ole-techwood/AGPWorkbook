package order

import "github.com/ole-techwood/PizzaCLI/internal/menu"

type OrderItem struct {
	Pizza    *menu.Pizza
	Quantity int
}

type Order struct {
	Items []OrderItem
}
