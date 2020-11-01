package magazine

type Controller struct {
	m *Model
	v *View
}

func NewController(m *Model, v *View) (*Controller, error) {
	if err := v.PrepareStaticContent(m.GetSortedRecords()); err != nil {
		return nil, err
	}
	return &Controller{m: m, v: v}, nil
}

func (m *Controller) PageIterator() []struct{} {
	return make([]struct{}, len(m.m.GetSortedRecords()))
}

func (m *Controller) Page(n int) ([]byte, error) {
	return m.v.RenderNthPage(n)
}
