package scaler

import (
	"strconv"
	"strings"
	"time"
)

type Expression string

func (e Expression) Match(t time.Time) bool {
	a := strings.Split(string(e), " ")

	if len(a) != 6 {
		return false
	}

	minute := pattern(a[0])
	hour := pattern(a[1])
	day := pattern(a[2])
	month := pattern(a[3])
	year := pattern(a[4])
	weekday := pattern(a[5])

	return minute.Match(t.Minute(), convert) &&
		hour.Match(t.Hour(), convert) &&
		day.Match(t.Day(), convert) &&
		month.Match(int(t.Month()), convertMonth) &&
		year.Match(t.Year(), convert) &&
		weekday.Match(int(t.Weekday()), convertWeekday)
}

type pattern string

type convertFunc func(s string) (int, error)

func convert(s string) (int, error) {
	return strconv.Atoi(s)
}

var months = map[string]int{
	"Jan": 1,
	"Feb": 2,
	"Mar": 3,
	"Apr": 4,
	"May": 5,
	"Jun": 6,
	"Jul": 7,
	"Aug": 8,
	"Sep": 9,
	"Oct": 10,
	"Nov": 11,
	"Dec": 12,
}

func convertMonth(s string) (int, error) {
	n, ok := months[s]
	if ok {
		return n, nil
	}
	return strconv.Atoi(s)
}

var weekdays = map[string]int{
	"Sun": 0,
	"Mon": 1,
	"Tue": 2,
	"Wed": 3,
	"Thu": 4,
	"Fri": 5,
	"Sat": 6,
}

func convertWeekday(s string) (int, error) {
	n, ok := weekdays[s]
	if ok {
		return n, nil
	}
	return strconv.Atoi(s)
}

func (p pattern) Match(n int, cf convertFunc) bool {
	if p == "*" {
		return true
	}

	a := strings.Split(string(p), ",")

	for _, b := range a {
		c := strings.Split(b, "-")

		if len(c) == 0 || len(c) >= 3 {
			continue
		}

		if len(c) == 2 {
			d, e := c[0], c[1]

			if len(d) >= 1 && len(e) >= 1 {
				f, err := cf(d)
				if err != nil {
					continue
				}

				g, err := cf(e)
				if err != nil {
					continue
				}

				if n >= f && n <= g {
					return true
				}

				continue
			}

			if len(d) >= 1 && len(e) == 0 {
				f, err := cf(d)
				if err != nil {
					continue
				}

				if n >= f {
					return true
				}

				continue
			}

			if len(d) == 0 && len(e) >= 1 {
				f, err := cf(e)
				if err != nil {
					return false
				}

				if n <= f {
					return true
				}

				continue
			}

			continue
		}

		d, err := cf(c[0])
		if err != nil {
			continue
		}

		if n == d {
			return true
		}
	}

	return false
}
