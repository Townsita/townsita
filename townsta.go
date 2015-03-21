package townsita

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Townsita struct {
	config *Config
	da     *DataAdapter
}

func New(da *DataAdapter) *Townsita {
	return &Townsita{
		da: da,
	}
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

func (t *Townsita) GetHTTPHandler(args []string) http.Handler {

	t.config = NewConfig()
	if err := t.config.Load(args); err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", appHandler(t.indexHandler).ServeHTTP).Methods("GET")
	return r
}

func (t *Townsita) indexHandler(w http.ResponseWriter, r *http.Request) error {
	s := NewSession(t.config)
	s.Set("MessageTypes", t.da.MustGetMessageTypes())
	s.AddPath("/", "Home")
	return s.render(w, r, t.config.templatePath("layout.html"), t.config.templatePath("index.html"))
}
