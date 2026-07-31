package model

import (
	"strings"
	"testing"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

const stubJSON = `{
	"name": "Test Restaurant",
	"menu": {
		"main": [
			{
				"name": "Combo A",
				"details": [{"name": "Chicken", "qty": 1}, {"name": "Drink", "qty": 1}],
				"price": 40000
			}
		],
		"extra": [{"name": "Rice", "price": 13000}]
	},
	"options": {
		"Chicken": ["Lemon Pepper", ["Garlic Parmesan", 2000]],
		"Drink": ["Sprite", "Coca Cola"]
	}
}`

func TestNewData(t *testing.T) {
	d := NewData([]byte(stubJSON))

	if d.Restaurant != "Test Restaurant" {
		t.Errorf("Restaurant = %q, want %q", d.Restaurant, "Test Restaurant")
	}
	if len(d.Main) != 1 || d.Main[0].Price != 40000 {
		t.Errorf("Menu.Main = %+v", d.Main)
	}
	if len(d.Main[0].Details) != 2 || d.Main[0].Details[0].Qty != 1 {
		t.Errorf("Details = %+v", d.Main[0].Details)
	}
	if len(d.Extra) != 1 || d.Extra[0].Name != "Rice" {
		t.Errorf("Menu.Extra = %+v", d.Extra)
	}
	if _, ok := d.Options["Rice"]; ok {
		t.Error("Options[Rice] should be absent")
	}
}

func TestOptionUnmarshal(t *testing.T) {
	opts := NewData([]byte(stubJSON)).Options["Chicken"]

	if len(opts) != 2 {
		t.Fatalf("len(Chicken) = %d, want 2", len(opts))
	}
	if opts[0] != (Option{Name: "Lemon Pepper"}) {
		t.Errorf("opts[0] = %+v", opts[0])
	}
	want := Option{Name: "Garlic Parmesan", ExtraPrice: 2000, HasExtraPrice: true}
	if opts[1] != want {
		t.Errorf("opts[1] = %+v, want %+v", opts[1], want)
	}
}

func TestOrder(t *testing.T) {
	var m OrderModel
	if m.Any() {
		t.Error("new order should be empty")
	}
	if total := m.Total(); total != 0 {
		t.Errorf("empty Total() = %d, want 0", total)
	}

	a := OrderItem{Name: "Combo A", Price: 43000, Details: []OrderDetail{
		{Name: "Chicken", Qty: 1, Option: []string{"Lemon Pepper"}},
	}}
	b := OrderItem{Name: "Combo B", Price: 70000}

	m.Add(a)
	if !m.Any() {
		t.Error("Any() = false after Add")
	}
	if total := m.Total(); total != 43000 {
		t.Errorf("Total() = %d, want 43000", total)
	}

	m.Add(b)
	m.Add(a)
	if total := m.Total(); total != 156000 {
		t.Errorf("Total() = %d, want 156000", total)
	}
}

func TestOptionsFormShowsEveryOption(t *testing.T) {
	m := NewMain(NewData([]byte(stubJSON)))
	m.item = 0
	if !m.optionsForm() {
		t.Fatal("optionsForm() = false, want true")
	}
	m.form.Init()
	m.form.Update(nil)

	view := m.form.View()
	for _, want := range []string{"Lemon Pepper", "Garlic Parmesan"} {
		if !strings.Contains(view, want) {
			t.Errorf("option %q not visible in:\n%s", want, view)
		}
	}
}

func TestOrderViewRendersItems(t *testing.T) {
	var m OrderModel
	if got := m.View(0); !strings.Contains(got, "Pesanan anda masih kosong!") {
		t.Errorf("empty order view missing placeholder:\n%s", got)
	}

	m.Add(OrderItem{Name: "Combo A", Price: 43000, Details: []OrderDetail{
		{Name: "Chicken", Qty: 2, Option: []string{"Lemon Pepper", "Soga"}},
		{Name: "Rice", Qty: 1, Option: []string{"Plain Rice"}},
	}})
	m.Add(OrderItem{Name: "Combo B", Price: 70000, Details: []OrderDetail{
		{Name: "Fries", Qty: 1},
	}})

	view := m.View(0)
	for _, want := range []string{
		"Combo A", "Chicken", "Lemon Pepper", "43000",
		"Combo B", "Fries", "70000",
	} {
		if !strings.Contains(view, want) {
			t.Errorf("order view missing %q:\n%s", want, view)
		}
	}

	if got := strings.Count(view, "┌"); got != 1 {
		t.Errorf("order view has %d tables, want 1 combined:\n%s", got, view)
	}
	if !strings.Contains(view, "Chicken 2 ") || !strings.Contains(view, " + ") {
		t.Errorf("order view does not join details into one cell:\n%s", view)
	}
}

func TestOrderViewRespectsWidth(t *testing.T) {
	var m OrderModel
	m.Add(OrderItem{Name: "Chicken Combo Rice 1", Price: 43000, Details: []OrderDetail{
		{Name: "Chicken", Qty: 1, Option: []string{"Golden Glitch"}},
		{Name: "Rice", Qty: 1, Option: []string{"Plain Rice"}},
		{Name: "Drink", Qty: 1, Option: []string{"Coca Cola"}},
	}})

	const width = 60
	for _, line := range strings.Split(m.View(width), "\n") {
		if got := lipgloss.Width(line); got > width {
			t.Errorf("line is %d cols, want <= %d:\n%s", got, width, line)
		}
	}
}

func TestNoteChromeMatchesHuh(t *testing.T) {
	const formWidth = 79

	render := func(s string) string {
		f := huh.NewForm(huh.NewGroup(
			huh.NewNote().Title("T").Description(s),
			huh.NewSelect[int]().Key("o").Options(huh.NewOption("Pesan", 0)),
		)).WithWidth(formWidth).WithHeight(23)
		f.Init()
		f.Update(nil)
		return f.View()
	}

	fits := strings.Repeat("X", formWidth-noteChrome)
	if !strings.Contains(render(fits), fits) {
		t.Errorf("%d cols wrapped inside a note; noteChrome is too small", len(fits))
	}

	over := strings.Repeat("X", formWidth-noteChrome+1)
	if strings.Contains(render(over), over) {
		t.Errorf("%d cols fit inside a note; noteChrome is too large", len(over))
	}
}
