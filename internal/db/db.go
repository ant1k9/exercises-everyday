package db

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/ant1k9/exercises-everyday/internal/config"
)

var Db *sqlx.DB

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

func AllExercisesTypes() []string {
	var types []string
	Db.Select(&types, "SELECT DISTINCT type FROM exercises ORDER BY type")
	return types
}

func Authenticate(login, password string) (string, bool) {
	var id int
	var dbPass string
	Db.QueryRow("SELECT id, password FROM users WHERE login = $1", login).
		Scan(&id, &dbPass)

	if err := bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(password)); err == nil {
		return newSession(id), true
	}
	return "", false
}

func CheckSession(sessValue string) bool {
	var exist int
	Db.QueryRow("SELECT COUNT(1) FROM sessions WHERE value = $1", sessValue).Scan(&exist)
	return exist > 0
}

func extractCurrentWeek() int {
	var currentWeek int
	Db.QueryRow("SELECT EXTRACT(week FROM date) FROM exercises ORDER BY date DESC LIMIT 1").
		Scan(&currentWeek)
	return currentWeek
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

func newSession(id int) string {
	value := sha256.New()
	for i := 0; i < 5; i++ {
		value.Write([]byte(strconv.Itoa(time.Now().Nanosecond())))
	}
	sessValue := hex.EncodeToString(value.Sum(nil))
	Db.MustExec("INSERT INTO sessions (user_id, value) VALUES($1, $2)", id, sessValue)
	return sessValue
}

func SaveProgress(t string, repeats int) {
	Db.MustExec("INSERT INTO exercises (type, repeats) VALUES ($1, $2)", t, repeats)
}
