package model

import (
	"errors"
	"fmt"
	"strconv"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

const noteChrome = 2

type state int

const (
	stateMenu state = iota
	stateItem
	stateOptions
	statePayment
	stateDone
	stateExit
)

type MainModel struct {
	form *huh.Form
	Data Data

	state         state
	order         OrderModel
	item          int
	width, height int
}

func (m *MainModel) setForm(form *huh.Form) {
	if m.width > 0 {
		form = form.WithWidth(m.width).WithHeight(m.height)
	}

	m.form = form
}

func (m *MainModel) resetForm() {
	m.state = stateMenu

	m.setForm(huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title(lipgloss.Sprintf("Selamat datang di %v!", m.Data.Restaurant)).
				DescriptionFunc(func() string {
					return m.order.View(m.width - noteChrome)
				}, m.order.Items),
			huh.NewSelect[int]().
				Key("option").
				Options(
					huh.NewOption("Pesan", 0),
					huh.NewOption("Bayar", 1),

					huh.NewOption("Exit", -1),
				).
				Validate(func(option int) error {
					if option == 1 && !m.order.Any() {
						return errors.New("pesanan masih kosong")
					}
					return nil
				}),
		),
	))
}

func (m *MainModel) itemForm() {
	m.state = stateItem

	options := make([]huh.Option[int], len(m.Data.Menu.Main))
	for i, item := range m.Data.Menu.Main {
		options[i] = huh.NewOption(item.Name, i)
	}

	m.setForm(huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Key("item").
				Title("Silahkan pilih menu utama kami").
				Options(options...),
		),
	))
}

func (m *MainModel) optionsForm() bool {
	menu := m.Data.Menu.Main[m.item]

	var groups []*huh.Group
	for index, detail := range menu.Details {
		options, ok := m.Data.Options[detail.Name]
		if !ok {
			continue
		}

		choices := make([]huh.Option[string], len(options))
		for k, option := range options {
			choices[k] = huh.NewOption(option.Name, option.Name)
		}

		qty := detail.Qty
		groups = append(groups, huh.NewGroup(
			huh.NewMultiSelect[string]().
				Key(detailKey(index)).
				Title(fmt.Sprintf("Pilih opsi untuk %s", detail.Name)).
				Options(choices...).
				Height(len(choices)+1).
				Limit(qty).
				Validate(func(selected []string) error {
					if len(selected) != qty {
						return fmt.Errorf("pilih %d opsi", qty)
					}
					return nil
				}),
		))
	}

	if len(groups) == 0 {
		return false
	}

	m.state = stateOptions
	m.setForm(huh.NewForm(groups...))

	return true
}

func (m *MainModel) paymentForm() {
	m.state = statePayment
	total := m.order.Total()

	m.setForm(huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title(lipgloss.Sprintf(
					"Total harga %s:\n%d",
					lipgloss.NewStyle().Faint(true).Render("(sudah termasuk PPn)"),
					total,
				)),
			huh.NewInput().
				Key("payment").
				Title("Silahkan bayar").
				Validate(func(s string) error {
					paid, err := strconv.Atoi(s)
					if err != nil {
						return errors.New("masukkan angka")
					}
					if paid < total {
						return errors.New("uang anda kurang")
					}
					return nil
				}),
		),
	))
}

func (m *MainModel) doneForm(change int) {
	m.state = stateDone

	m.setForm(huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title(fmt.Sprintf("Kembalian anda:\n%d", change)).
				Description("Terimakasih telah berkunjung!").
				Next(true),
		),
	))
}

func (m *MainModel) confirmExit() {
	m.state = stateExit

	m.setForm(huh.NewForm(
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
	))
}

func (m MainModel) selectedItem() OrderItem {
	menu := m.Data.Menu.Main[m.item]

	item := OrderItem{
		Name:    menu.Name,
		Details: make([]OrderDetail, len(menu.Details)),
		Price:   menu.Price,
	}
	for j, detail := range menu.Details {
		option, _ := m.form.Get(detailKey(j)).([]string)
		item.Details[j] = OrderDetail{
			Name:   detail.Name,
			Qty:    detail.Qty,
			Option: option,
		}
	}

	return item
}

func detailKey(detail int) string {
	return fmt.Sprintf("detail%d", detail)
}

func NewMain(data Data) (m MainModel) {
	m.Data = data
	m.resetForm()

	return
}

func (m MainModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width-1, msg.Height-1
		m.setForm(m.form)
		return m, nil
	}

	var cmds []tea.Cmd

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	switch m.state {
	case stateMenu:
		if m.form.State == huh.StateCompleted {
			switch m.form.GetInt("option") {
			case 0:
				m.itemForm()
			case 1:
				m.paymentForm()
			default:
				m.confirmExit()
			}
			cmds = append(cmds, m.form.Init())
		} else if m.form.State != huh.StateNormal {
			m.confirmExit()
			cmds = append(cmds, m.form.Init())
		}

	case stateItem:
		if m.form.State == huh.StateCompleted {
			m.item = m.form.GetInt("item")
			if !m.optionsForm() {
				m.order.Add(m.selectedItem())
				m.resetForm()
			}
			cmds = append(cmds, m.form.Init())
		} else if m.form.State != huh.StateNormal {
			m.resetForm()
			cmds = append(cmds, m.form.Init())
		}

	case stateOptions:
		if m.form.State == huh.StateCompleted {
			m.order.Add(m.selectedItem())
			m.resetForm()
			cmds = append(cmds, m.form.Init())
		} else if m.form.State != huh.StateNormal {
			m.resetForm()
			cmds = append(cmds, m.form.Init())
		}

	case statePayment:
		if m.form.State == huh.StateCompleted {
			paid, _ := strconv.Atoi(m.form.GetString("payment"))
			m.doneForm(paid - m.order.Total())
			cmds = append(cmds, m.form.Init())
		} else if m.form.State != huh.StateNormal {
			m.resetForm()
			cmds = append(cmds, m.form.Init())
		}

	case stateDone:
		if m.form.State != huh.StateNormal {
			cmds = append(cmds, tea.Quit)
		}

	case stateExit:
		exit := m.form.GetBool("exit")
		if (m.form.State == huh.StateCompleted && exit) ||
			m.form.State == huh.StateAborted {
			cmds = append(cmds, tea.Quit)
		} else if m.form.State != huh.StateNormal {
			m.resetForm()
			cmds = append(cmds, m.form.Init())
		}
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() (v tea.View) {
	v.AltScreen = true
	v.Content = m.form.View()

	return
}
