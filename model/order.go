package model

import (
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
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

func (detail OrderDetail) View() string {
	s := detail.Name + " " + strconv.Itoa(detail.Qty)
	if len(detail.Option) > 0 {
		s += " " + lipgloss.NewStyle().Faint(true).
			Render("("+strings.Join(detail.Option, ", ")+")")
	}

	return s
}

func (item OrderItem) View() string {
	details := make([]string, len(item.Details))
	for i, detail := range item.Details {
		details[i] = detail.View()
	}

	return strings.Join(details, " + ")
}

type OrderModel struct {
	Items []OrderItem
}

func (m *OrderModel) Add(item OrderItem) {
	m.Items = append(m.Items, item)
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

func (m OrderModel) View(width int) string {
	if !m.Any() {
		return lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Render("Pesanan anda masih kosong!"),
			lipgloss.NewStyle().Faint(true).Render("Silahkan pesan di menu utama."),
		)
	}

	t := table.New().Headers("Menu", "Detail", "Harga")
	if width > 0 {
		t.Width(width)
	}
	for _, item := range m.Items {
		t.Row(item.Name, item.View(), strconv.Itoa(item.Price))
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render("Pesanan:"),
		t.Render(),
	)
}
