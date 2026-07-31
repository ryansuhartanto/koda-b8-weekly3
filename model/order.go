package model

import (
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2/table"
)

type OrderDetail struct {
	Name   string
	Qty    int
	Option []string
}

type OrderItem struct {
	Name    string
	Details []OrderDetail
	Price   int
}

func (item OrderItem) Note() *huh.Note {
	t := table.New().Headers("Detail", "Qty", "Opsi")
	for _, detail := range item.Details {
		t.Row(detail.Name, strconv.Itoa(detail.Qty), strings.Join(detail.Option, ", "))
	}

	return huh.NewNote().
		Title(item.Name).
		Description(t.Render() + "\nHarga: " + strconv.Itoa(item.Price))
}

type OrderModel struct {
	form *huh.Form

	Items []OrderItem
}

func (m *OrderModel) resetForm() {
	if !m.Any() {
		m.form = huh.NewForm(
			huh.NewGroup(
				huh.NewNote().
					Title("Pesanan anda masih kosong!").
					Description("Silahkan pesan di menu utama."),
			).WithShowHelp(false),
		)

		return
	}

	fields := make([]huh.Field, len(m.Items))
	for i, item := range m.Items {
		fields[i] = item.Note()
	}

	m.form = huh.NewForm(
		huh.NewGroup(fields...).
			Title("Pesanan:").
			WithShowHelp(false),
	)
}

func NewOrder() (m OrderModel) {
	m.resetForm()

	return
}

func (m *OrderModel) Add(item OrderItem) {
	m.Items = append(m.Items, item)
	m.resetForm()
}

func (m OrderModel) Any() bool {
	return len(m.Items) > 0
}

func (m OrderModel) Total() (total int) {
	for _, item := range m.Items {
		total += item.Price
	}

	return
}

func (m OrderModel) Init() tea.Cmd {
	return nil
}

func (m OrderModel) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m OrderModel) View() (v tea.View) {
	v.Content = m.form.View()

	return
}
