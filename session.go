package townsita

import (
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
)

const sessionName = "townsita"

type PathLink struct {
	URL   string
	Label string
}

type Session struct {
	td     TemplateData
	r      *http.Request
	path   []*PathLink
	store  sessions.Store
	userId string
}

type TemplateData map[string]interface{}

func NewSession(c *Config, r *http.Request) *Session {
	return &Session{
		td:    NewTemplateData(c),
		r:     r,
		path:  []*PathLink{},
		store: sessions.NewCookieStore([]byte("something-very-secret")),
	}
}

func NewTemplateData(c *Config) TemplateData {
	td := make(TemplateData)
	td.Set("Title", "Townsita")
	return td
}

func (s *Session) Logged() bool {
	session, _ := s.store.Get(s.r, sessionName)
	userId, found := session.Values["userId"]
	if found {
		s.userId, _ = userId.(string)
	}
	return found
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

func (s *Session) getUser() *User {
	// TODO: Do me
	return &User{"1"}
}

func (s *Session) loginUser(user *User, w http.ResponseWriter) {
	session, _ := s.store.Get(s.r, sessionName)
	session.Values["userId"] = user.ID
	session.Save(s.r, w)
}
