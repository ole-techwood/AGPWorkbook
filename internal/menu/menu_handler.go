package menu

type MenuHandler struct {
	menuService *MenuService
}

func NewMenuHandler() *MenuHandler {
	return &MenuHandler{
		menuService: NewMenuService(),
	}
}

func (h *MenuHandler) GetMenu() {
	h.menuService.GetMenu()
}
