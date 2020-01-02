package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/ant1k9/exercises-everyday/internal/config"
)

var (
	Db *sqlx.DB
)

type exercise struct {
	Type    string `sql:"type"`
	Repeats int    `sql:"repeats"`
	Week    int    `sql:"week"`
}

func init() {
	var err error
	connectionParams := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Conf.Database.Host,
		config.Conf.Database.Port,
		config.Conf.Database.Username,
		config.Conf.Database.Password,
		config.Conf.Database.Name,
	)

	Db, err = sqlx.Connect("postgres", connectionParams)
	if err != nil {
		log.Fatal(err)
	}
}

func InitialMigrate() {
	Db.MustExec(InitialMigration)
}

func SaveProgress(t string, repeats int) {
	Db.MustExec("INSERT INTO exercises (type, repeats) VALUES ($1, $2)", t, repeats)
}

func AllExercisesTypes() []string {
	var types []string
	Db.Select(&types, "SELECT DISTINCT type FROM exercises ORDER BY type")
	return types
}

func GetStatsForTwoWeeks() (map[string]int, map[string]int) {
	lastWeekStats := make(map[string]int)
	currentWeekStats := make(map[string]int)

	currentWeek := extractCurrentWeek()

	var exercises []exercise
	Db.Select(
		&exercises,
		`SELECT type, SUM(repeats) repeats, EXTRACT(week FROM date) week
			FROM exercises WHERE EXTRACT(year FROM date) = EXTRACT(year FROM now())
			AND EXTRACT(week FROM date) <= $1 GROUP BY type, week `, currentWeek,
	)

	for _, exercise := range exercises {
		if exercise.Week == currentWeek {
			currentWeekStats[exercise.Type] = exercise.Repeats
		} else {
			lastWeekStats[exercise.Type] = exercise.Repeats
		}
	}
	return lastWeekStats, currentWeekStats
}

func extractCurrentWeek() int {
	var currentWeek int
	Db.QueryRow("SELECT EXTRACT(week FROM date) FROM exercises ORDER BY date DESC LIMIT 1").
		Scan(&currentWeek)
	return currentWeek
}
