package townsita

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Townsita struct {
	config *Config
	da     DataAdapter
}

type ValidationErrors []string

func New(da DataAdapter) *Townsita {
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
	r.HandleFunc("/message/new/{id}/{slug}", appHandler(t.newMessageHandler).ServeHTTP).Methods("GET", "POST")
	return r
}

func (t *Townsita) indexHandler(w http.ResponseWriter, r *http.Request) error {
	s := NewSession(t.config, r)
	s.Set("MessageTypes", t.da.MustGetMessageTypes())
	s.AddPath("/", "Home")
	return s.render(w, r, t.config.templatePath("layout.html"), t.config.templatePath("index.html"))
}

func (t *Townsita) newMessageHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	if vars["id"] == "" {
		return HTTPError{
			nil,
			"Bad Request.",
			http.StatusBadRequest,
		}
	}
	s := NewSession(t.config, r)
	var ve ValidationErrors
	var message *Message
	// Handle message post
	if r.Method == "POST" {
		message, ve = t.validateMessage(r)
		if len(ve) == 0 {
			messageId, err := t.da.SaveMessage(message, s.getUser())
			if err != nil {
				return err
			}
			// Redirect to the message page
			url := "/message/view/" + messageId + "/" + slug(message.Headline)
			http.Redirect(w, r, url, http.StatusFound)
			return nil
		}
	}
	s.Set("MessageTypes", t.da.MustGetMessageSubTypes(vars["id"]))
	s.Set("Message", message)
	s.Set("ValidationErrors", ve)
	s.AddPath("/", "Home")
	s.AddPath("/", "New Message")
	return s.render(w, r, t.config.templatePath("layout.html"), t.config.templatePath("new.html"))
}

func (t *Townsita) validateMessage(r *http.Request) (*Message, ValidationErrors) {
	var ve ValidationErrors
	return nil, ve
}
