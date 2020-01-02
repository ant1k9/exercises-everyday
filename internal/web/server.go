package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/ant1k9/exercises-everyday/internal/config"
	"github.com/ant1k9/exercises-everyday/internal/db"
)

type Data struct {
	ExercisesTypes []string
	ThisWeekStats  map[string]int
	LastWeekStats  map[string]int
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

func exerciseDone(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err == nil {
			repeatsValue := r.FormValue("repeats")
			if repeats, err := strconv.Atoi(repeatsValue); err == nil {
				db.SaveProgress(r.FormValue("type"), repeats)
			}
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func newType(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err == nil {
			db.SaveProgress(r.FormValue("type"), 0)
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func ServeForever() {
	router := mux.NewRouter()
	router.HandleFunc("/exercise/done", exerciseDone)
	router.HandleFunc("/new/type", newType)
	router.HandleFunc("/", indexPage)

	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(
			"%s:%s",
			config.Conf.Server.Host, config.Conf.Server.Port,
		), router,
	))
}
