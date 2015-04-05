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

func New(c *Config, da DataAdapter) *Townsita {
	return &Townsita{
		config: c,
		da:     da,
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

func (t *Townsita) GetHTTPHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", appHandler(t.indexHandler).ServeHTTP).Methods("GET")
	r.HandleFunc("/auth/login.html", appHandler(t.logInHandler).ServeHTTP).Methods("GET", "POST")
	r.HandleFunc("/auth/register.html", appHandler(t.registerHandler).ServeHTTP).Methods("GET", "POST")
	r.HandleFunc("/auth/logout.html", appHandler(t.logOutHandler).ServeHTTP).Methods("GET", "POST")
	r.HandleFunc("/user/profile.html", appHandler(t.myProfileHandler).ServeHTTP).Methods("GET")
	r.HandleFunc("/user/profile/{id}/{slug}", appHandler(t.userProfileHandler).ServeHTTP).Methods("GET")
	r.HandleFunc("/message/new/{id}/{slug}", appHandler(t.newMessageHandler).ServeHTTP).Methods("GET", "POST")
	r.HandleFunc("/message/view/{id}/{slug}", appHandler(t.viewMessageHandler).ServeHTTP).Methods("GET")
	r.HandleFunc("/message/address/{id}/{slug}", appHandler(t.addressMessageHandler).ServeHTTP).Methods("GET")
	return r
}

func (t *Townsita) indexHandler(w http.ResponseWriter, r *http.Request) error {
	s := NewSession(t.config, r)
	s.Set("MessageTypes", t.da.MustGetMessageTypes())
	s.AddPath("/", "Home")
	return s.render(w, r, t.config.templatePath("layout.html"), t.config.templatePath("index.html"))
}

func (t *Townsita) newMessageHandler(w http.ResponseWriter, r *http.Request) error {
	s := NewSession(t.config, r)
	if !s.Logged() {
		http.Redirect(w, r, "/auth/login.html", http.StatusTemporaryRedirect)
	}
	vars := mux.Vars(r)
	if vars["id"] == "" {
		return HTTPError{
			nil,
			"Bad Request.",
			http.StatusBadRequest,
		}
	}
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
	s.Set("TypeId", vars["id"])
	s.AddPath("/", "Home")
	s.AddPath("/", "New Message")
	return s.render(w, r, t.config.templatePath("layout.html"), t.config.templatePath("new.html"))
}

func (t *Townsita) validateMessage(r *http.Request) (*Message, ValidationErrors) {
	var ve ValidationErrors
	var message Message
	message.TypeID = r.FormValue("type_id")
	message.Headline = r.FormValue("headline")
	return &message, ve
}

func (t *Townsita) viewMessageHandler(w http.ResponseWriter, r *http.Request) error {
	s := NewSession(t.config, r)
	vars := mux.Vars(r)
	if vars["id"] == "" {
		return HTTPError{
			nil,
			"Bad Request.",
			http.StatusBadRequest,
		}
	}
	message, err := t.da.GetMessageById(vars["id"])
	if err != nil {
		return err
	}
	s.Set("Message", message)
	return s.render(w, r, t.config.templatePath("layout.html"), t.config.templatePath("view.html"))
}

func (t *Townsita) addressMessageHandler(w http.ResponseWriter, r *http.Request) error {
	s := NewSession(t.config, r)
	vars := mux.Vars(r)
	if vars["id"] == "" {
		return HTTPError{
			nil,
			"Bad Request.",
			http.StatusBadRequest,
		}
	}
	message, err := t.da.GetMessageById(vars["id"])
	if err != nil {
		return err
	}
	s.Set("Message", message)
	return s.render(w, r, t.config.templatePath("layout.html"), t.config.templatePath("address.html"))
}

func (t *Townsita) logInHandler(w http.ResponseWriter, r *http.Request) error {
	s := NewSession(t.config, r)
	if s.Logged() {
		http.Redirect(w, r, "/user/profile.html", http.StatusTemporaryRedirect)
	}
	if r.Method == "POST" {
		user, ve := t.validateUserLogin(r)
		if len(ve) == 0 {
			userID, err := t.da.LoginUser(user)
			if err != nil {
				return err
			}
			user.ID = userID
			s.logInUser(user, w)
			// Redirect to the message page
			http.Redirect(w, r, "/user/profile.html", http.StatusFound)
			return nil
		}
	}
	return s.render(w, r, t.config.templatePath("layout.html"), t.config.templatePath("auth/login.html"))
}

func (t *Townsita) registerHandler(w http.ResponseWriter, r *http.Request) error {
	s := NewSession(t.config, r)
	if s.Logged() {
		http.Redirect(w, r, "/user/profile.html", http.StatusTemporaryRedirect)
	}
	if r.Method == "POST" {
		user, ve := t.validateUserRegister(r)
		if len(ve) == 0 {
			userID, err := t.da.RegisterUser(user)
			if err != nil {
				return err
			}
			user.ID = userID
			s.logInUser(user, w)
			// Redirect to the message page
			http.Redirect(w, r, "/user/profile.html", http.StatusFound)
			return nil
		}
	}
	return s.render(w, r, t.config.templatePath("layout.html"), t.config.templatePath("auth/register.html"))
}

func (t *Townsita) logOutHandler(w http.ResponseWriter, r *http.Request) error {
	s := NewSession(t.config, r)
	s.logOutUser(w)
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

func (t *Townsita) validateUserLogin(r *http.Request) (*User, ValidationErrors) {
	var ve ValidationErrors
	user := NewUser()
	return user, ve
}

func (t *Townsita) validateUserRegister(r *http.Request) (*User, ValidationErrors) {
	var ve ValidationErrors
	user := NewUser()
	return user, ve
}

func (t *Townsita) myProfileHandler(w http.ResponseWriter, r *http.Request) error {
	s := NewSession(t.config, r)
	if !s.Logged() {
		http.Redirect(w, r, "/auth/login.html", http.StatusTemporaryRedirect)
	}
	user, err := t.da.LoadUserByID(s.userId)
	if err != nil {
		return err
	}
	s.Set("User", user)
	s.Set("OwnMessages", t.da.GetOwnMessages(user.ID, 10, 0))
	s.Set("ReceivedMessages", t.da.GetReceivedMessages(user.ID, 10, 0))
	s.AddPath("/", "Home")
	return s.render(w, r, t.config.templatePath("layout.html"), t.config.templatePath("user/my_profile.html"))
}

func (t *Townsita) userProfileHandler(w http.ResponseWriter, r *http.Request) error {
	s := NewSession(t.config, r)
	vars := mux.Vars(r)
	if vars["id"] == "" {
		return HTTPError{
			nil,
			"Bad Request.",
			http.StatusBadRequest,
		}
	}
	user, err := t.da.LoadUserByID(vars["id"])
	if err != nil {
		return err
	}
	s.Set("User", user)
	s.AddPath("/", "Home")
	return s.render(w, r, t.config.templatePath("layout.html"), t.config.templatePath("user/user_profile.html"))
}
