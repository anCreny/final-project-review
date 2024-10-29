package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sater-151/todo-list/models"
)

func Config() (port string, dbFilePath string) {
	port = os.Getenv("TODO_PORT")
	dbFilePath = os.Getenv("TODO_DBFILE")
	return port, dbFilePath
}

func GetPass() (pass string) {
	pass = os.Getenv("TODO_PASSWORD")
	return pass
}

func CheckCorrectRepeat(repeat string) error {
	switch {
	case repeat == "":
		return fmt.Errorf("repeat is empty")
	case repeat == "y":
		return nil
	case string(repeat[0]) == "d" && len(repeat) > 1:
		repeatMas := strings.Split(repeat, " ")
		d, err := strconv.Atoi(string(repeatMas[1]))
		if err != nil {
			return err
		}
		if len(repeatMas) == 2 && d <= 400 {
			return nil
		}

	case string(repeat[0]) == "w" && len(repeat) > 1:
		repeatMas := strings.Split(repeat, " ")
		weekday := strings.Split(repeatMas[1], ",")
		if len(weekday) > 7 || len(weekday) == 0 {
			return fmt.Errorf("uncorrect repeat")
		}
		for i := 0; i < len(weekday); i++ {
			buf, err := strconv.Atoi(weekday[i])
			if err != nil {
				return err
			}
			if buf > 7 || buf < 1 {
				return fmt.Errorf("uncorrect repeat")
			}
		}
		return nil
	case string(repeat[0]) == "m" && len(repeat) > 1:
		repeatMas := strings.Split(repeat, " ")
		if len(repeatMas) == 2 || len(repeatMas) == 3 {
			days := strings.Split(repeatMas[1], ",")
			for _, i := range days {
				buf, err := strconv.Atoi(i)
				if err != nil {
					return err
				}
				if buf < (-2) || buf == 0 || buf > 31 {
					return fmt.Errorf("uncorrect repeat")
				}
			}
			if len(repeatMas) == 3 {
				months := strings.Split(repeatMas[2], ",")
				for _, i := range months {
					buf, err := strconv.Atoi(i)
					if buf < 1 || buf > 12 {
						return fmt.Errorf("uncorrect repeat")
					}
					if err != nil {
						return err
					}
				}

			}
		} else {
			return fmt.Errorf("uncorrect repeat")
		}
		return nil
	}
	return fmt.Errorf("uncorrect repeat")
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	err := CheckCorrectRepeat(repeat)
	if err != nil {
		return "", err
	}
	dateParse, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}
	repeatParse := strings.Split(repeat, " ")
	switch string(repeat[0]) {
	case "y":
		for {
			dateParse = dateParse.AddDate(1, 0, 0)
			if now.Compare(dateParse) == -1 {
				return dateParse.Format("20060102"), nil
			}
		}
	case "d":
		days, _ := strconv.Atoi(repeatParse[1])
		for {
			dateParse = dateParse.AddDate(0, 0, days)
			if now.Compare(dateParse) == -1 {
				return dateParse.Format("20060102"), nil
			}

		}
	case "w":
		weekday := strings.Split(repeatParse[1], ",")
		var weekdayParse []int
		for i := 0; i < len(weekday); i++ {
			buf, _ := strconv.Atoi(weekday[i])
			if buf == 7 {
				buf = 0
			}
			weekdayParse = append(weekdayParse, buf)
		}
		for {
			dateParse = dateParse.AddDate(0, 0, 1)
			weekdayNow := dateParse.Weekday()
			for _, i := range weekdayParse {
				if i == int(weekdayNow) && now.Compare(dateParse) == -1 {
					return dateParse.Format("20060102"), nil
				}
			}
		}
	case "m":
		if len(repeatParse) == 2 {
			numbers := strings.Split(repeatParse[1], ",")
			for {
				dateParse = dateParse.AddDate(0, 0, 1)
				for _, i := range numbers {
					if i == "-1" {
						buf := dateParse.AddDate(0, 0, 1)

						if dateParse.Format("01") != buf.Format("01") && now.Compare(dateParse) == -1 {
							return dateParse.Format("20060102"), nil
						}
						continue
					}
					if i == "-2" {
						buf := dateParse.AddDate(0, 0, 2)
						if dateParse.Format("01") != buf.Format("01") && now.Compare(dateParse) == -1 {
							return dateParse.Format("20060102"), nil
						}
						continue
					}
					a, _ := strconv.Atoi(i)
					b, _ := strconv.Atoi(dateParse.Format("2"))
					if a == b && now.Compare(dateParse) == -1 {
						return dateParse.Format("20060102"), nil
					}
				}
			}
		}
		if len(repeatParse) == 3 {
			numbers := strings.Split(repeatParse[1], ",")
			months := strings.Split(repeatParse[2], ",")
			for {
				dateParse = dateParse.AddDate(0, 0, 1)
				for _, i := range numbers {
					if i == "-1" {
						buf := dateParse.AddDate(0, 0, 1)

						if dateParse.Format("01") != buf.Format("01") {
							for _, j := range months {
								c, _ := strconv.Atoi(j)
								d, _ := strconv.Atoi(dateParse.Format("1"))
								if d == c && now.Compare(dateParse) == -1 {
									return dateParse.Format("20060102"), nil
								}
							}
						}
						continue
					}
					if i == "-2" {
						buf := dateParse.AddDate(0, 0, 2)
						if dateParse.Format("01") != buf.Format("01") {
							for _, j := range months {
								c, _ := strconv.Atoi(j)
								d, _ := strconv.Atoi(dateParse.Format("1"))
								if d == c && now.Compare(dateParse) == -1 {
									return dateParse.Format("20060102"), nil
								}
							}
						}
						continue
					}
					a, _ := strconv.Atoi(i)
					b, _ := strconv.Atoi(dateParse.Format("2"))
					if a == b {
						for _, j := range months {
							c, _ := strconv.Atoi(j)
							d, _ := strconv.Atoi(dateParse.Format("1"))
							if d == c && now.Compare(dateParse) == -1 {
								return dateParse.Format("20060102"), nil
							}
						}
					}
				}
			}
		}
	}
	return "", fmt.Errorf("uncorrect date")
}

func CheckTask(task models.Task) (models.Task, error) {
	if task.Title == "" {
		return task, fmt.Errorf("title is empty")
	}
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	} else {
		t, err := time.Parse("20060102", task.Date)
		if err != nil {
			return task, err
		}
		now := time.Now().Format("20060102")
		if t.Compare(time.Now()) == -1 && now != task.Date {
			if task.Repeat == "" {
				task.Date = time.Now().Format("20060102")
			} else {
				task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					return task, err
				}
			}
		}
	}
	return task, nil
}
