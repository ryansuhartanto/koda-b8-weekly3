package form

import (
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

type Main struct {
	Exit bool
}

var _ MenuForm = (*Main)(nil)

func (m *Main) Form() *huh.Form {
	restaurant := lipgloss.NewStyle().
		Foreground(lipgloss.BrightGreen).
		Background(lipgloss.Black).
		Render("Wingstop")

	lipgloss.Printf("Selamat datang di %v!\n", restaurant)

	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Exit?").
				Value(&m.Exit),
		),
	)
}
