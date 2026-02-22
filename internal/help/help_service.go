package help

type HelpService struct{}

func NewHelpService() *HelpService {
	return &HelpService{}
}

func (s *HelpService) GetHelp() {
	println("Usage: pizzacli [command]")
	println("")
	println("Commands:")
	println("  menu  - Show pizza menu")
	println("  add <name> [quantity]  - Add a pizza to the order")
}
