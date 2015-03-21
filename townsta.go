package townsta

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Townsta struct {
	config *Config
}

func NewTownsta() *Townsta {
	return &Townsta{}
}

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		httpError, ok := err.(HTTPError)
		if ok {
			http.Error(w, httpError.Message, httpError.Code)
			return
		}
		// Default to 500 Internal Server Error
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *Townsta) Start(args []string) {

	t.config = NewConfig()
	if err := t.config.Load(args); err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", appHandler(t.indexHandler).ServeHTTP).Methods("GET")
}

func (t *Townsta) indexHandler(w http.ResponseWriter, r *http.Request) error {
	s := NewSession(t.config)
	s.AddPath("/", "Home")
	s.render(w, r, t.config.templatePath("layout.html"), t.config.templatePath("index.html"))
	return nil
}
