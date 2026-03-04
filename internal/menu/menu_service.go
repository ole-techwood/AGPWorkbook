package menu

import (
	"fmt"
	"strings"
)

type MenuService struct{}

func NewMenuService() MenuService {
	return MenuService{}
}

func (s *MenuService) GetMenu() {
	var sb strings.Builder
	sb.Grow(len(Menu) * 32) // rough estimate per line

	for _, pizza := range Menu {
		fmt.Fprintf(&sb, "%d | %s | $%.2f\n", pizza.ID, pizza.Name, pizza.Price)
	}

	fmt.Print(sb.String())
}

func (s *MenuService) GetPizzaByName(name string) (*Pizza, error) {
	for i := range Menu {
		if strings.EqualFold(Menu[i].Name, name) {
			return &Menu[i], nil
		}
	}

	return nil, fmt.Errorf("pizza not found: %s", name)
}
