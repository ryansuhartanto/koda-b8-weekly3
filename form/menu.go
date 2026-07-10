package form

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

type MainModel struct {
	restaurant string
	form       *huh.Form

	exiting bool
}

func (m *MainModel) resetForm() {
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title(lipgloss.Sprintf("Selamat datang di %v!", m.restaurant)),
		),
	)
}

func (m *MainModel) confirmExit() {
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("exit").
				Title("Exit?"),
		),
	)
}

func NewMain(restaurant string) (m MainModel) {
	m.restaurant = restaurant
	m.resetForm()

	return
}

func (m MainModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if !m.exiting {
		if m.form.State != huh.StateNormal {
			m.exiting = true
			m.confirmExit()
			cmds = append(cmds, m.form.Init())
		}
	} else {
		exit := m.form.GetBool("exit")
		if exit || m.form.State == huh.StateAborted {
			cmds = append(cmds, tea.Quit)
		} else if m.form.State != huh.StateNormal {
			m.exiting = false
			m.resetForm()
			cmds = append(cmds, m.form.Init())
		}
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() (v tea.View) {
	v.SetContent(m.form.View())
	return
}
