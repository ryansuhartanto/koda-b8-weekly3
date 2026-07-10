package model

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

type MainModel struct {
	restaurant string

	form    *huh.Form
	exiting bool
}

func (m *MainModel) resetForm() {
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title(lipgloss.Sprintf("Selamat datang di %v!", m.restaurant)),
			huh.NewSelect[int]().
				Key("option").
				Options(
					huh.NewOption("Pesan", 0),
					huh.NewOption("Bayar", 1),

					huh.NewOption("Exit", -1),
				),
		),
	)
}

func (m *MainModel) confirmExit() {
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Key("exit").
				Title("Batalkan pesanan?"),
		),

		huh.NewGroup(
			huh.NewNote().
				Title("Terimakasih telah berkunjung!").
				Next(true),
		).WithHideFunc(func() bool {
			return !m.form.GetBool("exit")
		}),
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
		option := m.form.GetInt("option")
		if m.form.State == huh.StateCompleted && option != -1 {
			m.resetForm()
			cmds = append(cmds, m.form.Init())
		} else if m.form.State != huh.StateNormal {
			m.exiting = true
			m.confirmExit()
			cmds = append(cmds, m.form.Init())
		}
	} else {
		exit := m.form.GetBool("exit")
		if (m.form.State == huh.StateCompleted && exit) ||
			m.form.State == huh.StateAborted {
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
	v.Content = m.form.View()
	v.AltScreen = true

	return
}
