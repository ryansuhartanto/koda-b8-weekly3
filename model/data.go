package model

import "encoding/json"

type Data struct {
	Restaurant `json:"name"`
	Menu       `json:"menu"`
	Options    `json:"options"`
}

func NewData(data []byte) (d Data) {
	if err := json.Unmarshal(data, &d); err != nil {
		panic(err)
	}

	return
}

type Restaurant string

type Menu struct {
	Main  []MenuMain  `json:"main"`
	Extra []MenuExtra `json:"extra"`
}

type MenuMain struct {
	Name    string       `json:"name"`
	Details []MainDetail `json:"details"`
	Price   int          `json:"price"`
}

type MainDetail struct {
	Name string `json:"name"`
	Qty  int    `json:"qty"`
}

type MenuExtra struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Options map[string][]Option

type Option struct {
	Name          string
	ExtraPrice    int
	HasExtraPrice bool
}

func (o *Option) UnmarshalJSON(data []byte) error {
	var s string
	if json.Unmarshal(data, &s) == nil {
		o.Name, o.HasExtraPrice = s, false
		return nil
	}

	var tuple []json.RawMessage
	if err := json.Unmarshal(data, &tuple); err != nil {
		return err
	}

	if err := json.Unmarshal(tuple[0], &o.Name); err != nil {
		return err
	}
	if err := json.Unmarshal(tuple[1], &o.ExtraPrice); err != nil {
		return err
	}
	o.HasExtraPrice = true
	return nil
}

func (o Option) MarshalJSON() ([]byte, error) {
	if o.HasExtraPrice {
		return json.Marshal([]any{o.Name, o.ExtraPrice})
	}
	return json.Marshal(o.Name)
}
