package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// nextDayHandler - рассчитать следующую дату (= сдвигу + дата создания).
func (h *SrvHand) nextDayHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Sugar().Info("START /api/nextdate ", r.Method)
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var nowTime time.Time
	if now == "" {
		nowTime = time.Now()
	} else {
		var err error
		nowTime, err = time.Parse(pattern, now)
		if err != nil {
			h.Logger.Sugar().Errorf("date(<now>) conversion error: %v", err)
			//http.Error(w, "date(<now>) conversion error", http.StatusOK)
			writeJson(w, reterror{Error: "date(<now>) conversion error"})
			return
		}
	}

	nxtDate, err := NextDate(nowTime, date, repeat)
	if err != nil {
		h.Logger.Sugar().Errorf("didn't get next date: %v", err)
		//http.Error(w, err.Error(), http.StatusOK)
		writeJson(w, reterror{Error: fmt.Sprintf("didn't get next date: %v", err)})

		return
	}

	w.Write([]byte(nxtDate))

}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	fmt.Printf("now: %v\ndstar: %v\nrepeat: %v\n", now, dstart, repeat)
	if len(repeat) == 0 {
		return "", fmt.Errorf("expected to receive a rule repeat")
	}

	date, err := time.Parse(pattern, dstart)
	if err != nil {
		return "", fmt.Errorf("couldn't parse dstart, incorrect value (%v) ", dstart)
	}

	// repeat rules
	if repeat == "y" {
		result := addYear(now, date)
		return result, nil
	}

	if d, _ := regexp.MatchString(`^d \d`, repeat); d {
		rpt := strings.Split(repeat, " ")
		if len(rpt) != 2 {
			return "", fmt.Errorf("repeat isn't correct (%v)", repeat)
		}
		countDays, err := strconv.Atoi(rpt[1])
		if err != nil {
			return "", err
		}
		if countDays < 0 || countDays > 400 {
			return "", fmt.Errorf("repeat value < 0 (or > 400) (%v)", repeat)
		}

		result := addDays(now, date, countDays)

		return result, nil

	}
	// *
	// if w, _ := regexp.MatchString(`^w `, repeat); w {
	// 	rpt := strings.Split(repeat, " ")
	// 	if len(rpt) != 2 {
	// 		return "", fmt.Errorf("repeat isn't correct (%v)", repeat)
	// 	}

	// 	result, err := addWeekDays(rpt[1], now, date)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	return result, nil
	// }
	// if m, _ := regexp.MatchString(`^m `, repeat); m {

	// }

	// return nxtDate.Format(pattern), nil
	return "", fmt.Errorf("unknown repeat value: %v", repeat)
}

// afterNow - возвращает true, если первая дата больше второй
func afterNow(date time.Time, now time.Time) bool {
	if date.Year() > now.Year() {
		return true
	}
	if date.Year() >= now.Year() && date.Month() > now.Month() {
		return true
	}
	if date.Year() >= now.Year() && date.Month() >= now.Month() && date.Day() > now.Day() {
		return true
	}
	return false
}

// addYear - задача выполняется ежегодно.
func addYear(now time.Time, dstart time.Time) string {
	nxtDate := dstart.AddDate(1, 0, 0)
	for {

		if afterNow(nxtDate, now) {
			break
		}
		nxtDate = nxtDate.AddDate(1, 0, 0)
		fmt.Println(nxtDate)
		//fmt.Printf("afterNow(%v, %v) = %v\n", nxtDate, now)
	}
	return nxtDate.Format(pattern)
}

// addDays - задача переносится на указанное число дней.
func addDays(now time.Time, dstart time.Time, count int) string {
	nxtDate := dstart.AddDate(0, 0, count)
	for {

		if afterNow(nxtDate, now) {
			break
		}
		nxtDate = nxtDate.AddDate(0, 0, count)
		fmt.Println("nxtDate:", nxtDate, "now:", now, afterNow(nxtDate, now))
	}
	return nxtDate.Format(pattern)
}

// addWeekDays - задача назначается в указанные дни недели.
func addWeekDays(rpt string, now time.Time, dstart time.Time) (string, error) {
	if len(rpt) == 0 {
		return "", fmt.Errorf("repeat value isn't correct (w %v)", rpt)
	}

	// формирование слайса - дни недели для повтора
	var x []int
	if len(x) == 1 {
		w, err := strconv.Atoi(rpt)
		if err != nil || w <= 0 || w > 7 {
			return "", fmt.Errorf("repeat value isn't correct (w %v)", rpt)
		}
		x = append(x, w)
	} else {
		m := strings.Split(rpt, ",")
		for _, item := range m {
			w, err := strconv.Atoi(item)
			if err != nil || w <= 0 || w > 7 {
				return "", fmt.Errorf("repeat value isn't correct (w %v)", rpt)
			}
			x = append(x, w)
		}
	}

	countX := len(x)

	// подобие матрицы
	// p1 - сдвиг до ближайшего дня недели.
	// 		Т.е определим текущая дата > или < дней недели из repeat.
	//		Если  > заполним 0, иначе кол-во дней полученное как разность между днями.
	// Например: сегодня среда(3), repeat(1, 5) => (0, 2)
	p1 := make([]int, countX)

	for i := 0; i < countX; i++ {
		if dstart.Day() < x[i] {
			copy(p1[i:], x[i:])
			break
		}
	}
	fmt.Printf("p1: %v", p1)
	// p2 - сдвиги для циклического увеличения.
	// Составляется по заданным данным, количество дней до следующего заданного дня недели.
	p2 := make([]int, countX)
	for i := 0; i < countX; i++ {
		if i == (countX - 1) {
			p2[i] = 7 - x[i] + x[0]
		} else {
			p2[i] = x[i+1] - x[i]
		}
	}

	var nxtDate time.Time
	// сдвиг для первой недели
	nxtDate = dstart
	for _, item := range p1 {
		nxtDate = nxtDate.AddDate(0, 0, item)
		if afterNow(nxtDate, now) {
			return nxtDate.Format(pattern), nil
		}
	}
	// сдвиг для последующих
	for {
		for _, item := range p2 {
			nxtDate = nxtDate.AddDate(0, 0, item)
			if afterNow(nxtDate, now) {
				return nxtDate.Format(pattern), nil
			}
		}
	}

}
