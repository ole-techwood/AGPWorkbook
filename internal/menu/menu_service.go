package menu

import (
	"fmt"
	"strings"
)

type MenuService struct{}

func NewMenuService() *MenuService {
	return &MenuService{}
}

func (s *MenuService) GetMenu() {
	var sb strings.Builder
	sb.Grow(len(Menu) * 32) // rough estimate per line

	for _, pizza := range Menu {
		fmt.Fprintf(&sb, "%d | %s | $%.2f\n", pizza.ID, pizza.Name, pizza.Price)
	}

	fmt.Print(sb.String())
}

func (s *MenuService) GetPizzaByName(name string) (Pizza, error) {
	for _, pizza := range Menu {
		if strings.EqualFold(pizza.Name, name) {
			return pizza, nil
		}
	}

	return Pizza{}, fmt.Errorf("pizza not found: %s", name)
}
