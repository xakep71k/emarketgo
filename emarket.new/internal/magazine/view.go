package magazine

import (
	"bytes"
	"emarket/internal/records"
	"text/template"
)

type View struct {
	webRoot string
	appPage *template.Template
	static  *staticPages
}

type staticPages struct {
	contact  []byte
	delivery []byte
	notfound []byte
	home     []byte
	history  []byte
	neworder []byte
	list     [][]byte
}

func NewView(webRoot string) (*View, error) {
	var err error
	m := &View{webRoot: webRoot}
	m.appPage, err = template.New("app page").Funcs(defaultTemplateFuncs()).Parse(appHTML)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *View) PrepareStaticContent(magazs []*records.Magazine) error {
	m.static = new(staticPages)
	contact, err := newContactPage()
	if err != nil {
		return err
	}
	m.static.contact, err = m.renderPage(contact)
	if err != nil {
		return err
	}

	magazinesPagesHTML, err := prepareHTMLMagazinesPages(arrangeMagazinesPerPage(magazs))
	if err != nil {
		return nil
	}

	m.static.home, err = m.renderPage(newHomePage(magazinesPagesHTML[0]))
	if err != nil {
		return nil
	}
	return nil
}

func (m *View) RenderNthPage(n int) ([]byte, error) {
	return nil, nil
}

func (m *View) renderPage(p *page) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	if err := m.appPage.Execute(buf, p); err != nil {
		return nil, err
	}

	return doMinify(buf.Bytes(), "text/html")
}
