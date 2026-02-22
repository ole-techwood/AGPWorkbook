package help

type HelpHandler struct {
	helpService *HelpService
}

func NewHelpHandler() *HelpHandler {
	return &HelpHandler{
		helpService: NewHelpService(),
	}
}

func (h *HelpHandler) GetHelp() {
	h.helpService.GetHelp()
}
