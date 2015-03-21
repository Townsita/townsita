package townsita

import (
	"html/template"
	"net/http"
)

type PathLink struct {
	URL   string
	Label string
}

type Session struct {
	td   TemplateData
	path []*PathLink
}

type TemplateData map[string]interface{}

func NewSession(c *Config) *Session {
	return &Session{
		td:   NewTemplateData(c),
		path: []*PathLink{},
	}
}

func NewTemplateData(c *Config) TemplateData {
	td := make(TemplateData)
	td.Set("Title", "Townsita")
	return td
}

func (s *Session) getHelpers() template.FuncMap {
	return template.FuncMap{
		"slug": slug,
	}
}

func (s *Session) render(w http.ResponseWriter, r *http.Request, filenames ...string) error {
	t := template.New("layout.html")
	// Add helper functions
	t.Funcs(s.getHelpers())
	// Add path
	s.td.Set("Path", s.path)
	return template.Must(t.ParseFiles(filenames...)).Execute(w, s.td)
}

func (td TemplateData) Set(name string, value interface{}) {
	td[name] = value
}

func (s *Session) Set(name string, value interface{}) {
	s.td.Set(name, value)
}

func (s *Session) AddPath(url, label string) {
	s.path = append(s.path, &PathLink{url, label})
}
