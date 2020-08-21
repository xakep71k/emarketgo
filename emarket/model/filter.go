package model

type Filter struct {
	enable *bool
	sort   bool
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) Sort(s bool) *Filter {
	f.sort = s
	return f
}

func (f *Filter) Enable(e bool) *Filter {
	f.enable = new(bool)
	*f.enable = e
	return f
}
