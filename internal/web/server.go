package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/ant1k9/exercises-everyday/internal/config"
	"github.com/ant1k9/exercises-everyday/internal/db"
)

const SessionCookie = "exercises-everyday-session"

type Data struct {
	ExercisesTypes []string
	ThisWeekStats  map[string]int
	LastWeekStats  map[string]int
}

func checkSession(r *http.Request) bool {
	if cookie, err := r.Cookie(SessionCookie); err == nil {
		return db.CheckSession(cookie.Value)
	}
	return false
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("internal/web/templates/index.html"))
	lastWeekStats, thisWeekStats := db.GetStatsForTwoWeeks()
	tmpl.Execute(w, Data{
		ExercisesTypes: db.AllExercisesTypes(),
		ThisWeekStats:  thisWeekStats,
		LastWeekStats:  lastWeekStats,
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmpl := template.Must(template.ParseFiles("internal/web/templates/login.html"))
		tmpl.Execute(w, nil)
		return
	case http.MethodPost:
		if err := r.ParseForm(); err == nil {
			if sessValue, ok := db.Authenticate(r.FormValue("login"), r.FormValue("password")); ok {
				expiration := time.Now().Add(365 * 24 * time.Hour)
				cookie := http.Cookie{Name: SessionCookie, Value: sessValue, Expires: expiration}
				http.SetCookie(w, &cookie)
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
		}
	}

	http.NotFound(w, r)
}

func exerciseDone(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if checkSession(r) {
			if err := r.ParseForm(); err == nil {
				repeatsValue := r.FormValue("repeats")
				if repeats, err := strconv.Atoi(repeatsValue); err == nil {
					db.SaveProgress(r.FormValue("type"), repeats)
				}
			}
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

func newType(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if checkSession(r) {
			if err := r.ParseForm(); err == nil {
				db.SaveProgress(r.FormValue("type"), 0)
			}
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

func ServeForever() {
	router := mux.NewRouter()
	router.HandleFunc("/exercise/done", exerciseDone)
	router.HandleFunc("/login", login)
	router.HandleFunc("/new/type", newType)
	router.HandleFunc("/", indexPage)

	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(
			"%s:%s",
			config.Conf.Server.Host, config.Conf.Server.Port,
		), router,
	))
}
