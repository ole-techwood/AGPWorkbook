package menu

import (
	"fmt"
	"strings"
	"time"
)

type MenuService struct {
	sb         strings.Builder
	pizzaCache map[string]cacheEntry
}

func NewMenuService() *MenuService {
	return &MenuService{
		pizzaCache: make(map[string]cacheEntry),
	}
}

func (s *MenuService) GetMenu() {
	s.sb.Reset()
	s.sb.Grow(len(Menu) * 32) // rough estimate per line

	for _, pizza := range Menu {
		fmt.Fprintf(&s.sb, "%d | %s | $%.2f\n", pizza.ID, pizza.Name, pizza.Price)
	}

	fmt.Print(s.sb.String())
}

func (s *MenuService) GetPizzaByName(name string) (*Pizza, error) {
	key := strings.ToLower(strings.TrimSpace(name))

	if entry, ok := s.pizzaCache[key]; ok && time.Since(entry.cachedAt) < cacheTTL {
		fmt.Printf("cache hit: %s\n", entry.pizza.Name)
		return entry.pizza, nil
	}

	for i := range Menu {
		if strings.EqualFold(Menu[i].Name, name) {
			s.pizzaCache[key] = cacheEntry{pizza: &Menu[i], cachedAt: time.Now()}
			return &Menu[i], nil
		}
	}

	return nil, fmt.Errorf("pizza not found: %s", name)
}
