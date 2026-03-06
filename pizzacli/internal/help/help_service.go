package help

type HelpService struct{}

func NewHelpService() HelpService {
	return HelpService{}
}

func (s *HelpService) GetHelp() {
	println("Usage: [command]")
	println("")
	println("Commands:")
	println("  menu                   - Show pizza menu")
	println("  add <name> [quantity]  - Add a pizza to the order")
	println("  checkout               - Place order and clear the cart")
	println("  help                   - Show this help message")
	println("  exit                   - Exit PizzaCLI")
}
