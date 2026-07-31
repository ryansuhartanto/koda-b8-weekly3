package model

import "testing"

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
	if len(d.Menu.Main) != 1 || d.Menu.Main[0].Price != 40000 {
		t.Errorf("Menu.Main = %+v", d.Menu.Main)
	}
	if len(d.Menu.Main[0].Details) != 2 || d.Menu.Main[0].Details[0].Qty != 1 {
		t.Errorf("Details = %+v", d.Menu.Main[0].Details)
	}
	if len(d.Menu.Extra) != 1 || d.Menu.Extra[0].Name != "Rice" {
		t.Errorf("Menu.Extra = %+v", d.Menu.Extra)
	}
	if _, ok := d.Options["Rice"]; ok {
		t.Error("Options[Rice] should be absent")
	}
}

// Options are either a bare string or a ["name", surcharge] tuple.
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
	m := NewOrder()
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
