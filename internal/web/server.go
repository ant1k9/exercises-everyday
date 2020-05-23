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
	ExercisesTypes   []string
	ThisWeekStats    map[string]int
	LastWeekStats    map[string]int
	ChangeStats      map[string]string
	EstimatedRepeats int
}

func calculateChangeStats(lastWeekStats, thisWeekStats map[string]int) map[string]string {
	changeStats := make(map[string]string)
	for k, v := range thisWeekStats {
		if lastWeekValue, ok := lastWeekStats[k]; ok && lastWeekValue > 0 {
			value := int((float64(v)/float64(lastWeekValue) - 1.0) * 100.0)
			changeStats[k] = strconv.Itoa(value)
			if value > 0 {
				changeStats[k] = "+" + changeStats[k]
			}
			continue
		}
		changeStats[k] = "+100"
	}

	for k := range lastWeekStats {
		if _, ok := lastWeekStats[k]; !ok {
			changeStats[k] = "-100"
		}
	}

	return changeStats
}

func checkSession(r *http.Request) bool {
	if cookie, err := r.Cookie(SessionCookie); err == nil {
		return db.CheckSession(cookie.Value)
	}
	return false
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("internal/web/templates/index.html"))
	lastWeekStats, currentWeekStats := db.GetStatsForTwoWeeks()
	tmpl.Execute(w, Data{
		ExercisesTypes:   db.AllExercisesTypes(),
		ThisWeekStats:    currentWeekStats,
		LastWeekStats:    lastWeekStats,
		ChangeStats:      calculateChangeStats(lastWeekStats, currentWeekStats),
		EstimatedRepeats: db.EstimatedRepeats(lastWeekStats, currentWeekStats),
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

func ServeForever() {
	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("internal/web/static"))),
	)
	router.HandleFunc("/exercise/done", exerciseDone)
	router.HandleFunc("/login", login)
	router.HandleFunc("/", indexPage)

	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(
			"%s:%s",
			config.Conf.Server.Host, config.Conf.Server.Port,
		), router,
	))
}
