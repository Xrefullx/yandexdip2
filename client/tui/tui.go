package tui

import (
	"github.com/Xrefullx/YanDip/client/model"
	"github.com/Xrefullx/YanDip/client/services"
	"github.com/rivo/tview"
)

type TUI struct {
	app           *tview.Application
	secretService services.SecretService
}

// NewTUI creates a new TUI instance
func NewTUI(app *tview.Application, secretService services.SecretService) *TUI {
	return &TUI{
		app:           app,
		secretService: secretService,
	}
}

// SetQ sets up the TUI layout and starts the event loop
func (t *TUI) SetQ() error {
	hello := tview.NewTextView().SetText("Hello, world!")

	grid := tview.NewGrid().SetRows(3).SetColumns(30).AddItem(hello, 0, 0, 1, 1, 0, 0, true)

	t.app.SetRoot(grid, true)

	return t.app.Run()
}

// AddAuth adds a new auth secret using the SecretService
func (t *TUI) AddAuth(auth model.Auth) (int64, error) {
	return t.secretService.AddAuth(auth)
}

// AddCard adds a new card secret using the SecretService
func (t *TUI) AddCard(card model.Card) (int64, error) {
	return t.secretService.AddCard(card)
}

// AddBinary adds a new binary secret using the SecretService
func (t *TUI) AddBinary(filePath string, title string, description string) (int64, error) {
	return t.secretService.AddBinary(filePath, title, description)
}

// UpdateSecret updates a secret using the SecretService
func (t *TUI) UpdateSecret(secret model.Secret) error {
	return t.secretService.UpdateSecret(secret)
}

// DeleteSoftSecret soft deletes a secret using the SecretService
func (t *TUI) DeleteSoftSecret(id int64) error {
	return t.secretService.DeleteSoftSecret(id)
}

// GetSecret gets a secret using the SecretService
func (t *TUI) GetSecret(id int64) (model.Secret, error) {
	return t.secretService.GetSecret(id)
}
