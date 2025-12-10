package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	numDays        = 7
	minDaysInMonth = -2
	maxDaysInMonth = 31
	minMonth       = 1
	maxMonth       = 12
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
			writeJson(w, reterror{Error: "date(<now>) conversion error"})
			return
		}
	}

	h.Logger.Sugar().Infof("nowTime: %v, date = %v, repeat = %v", nowTime, date, repeat)
	nxtDate, err := NextDate(nowTime, date, repeat)
	if err != nil {
		h.Logger.Sugar().Errorf("didn't get next date: %v", err)
		writeJson(w, reterror{Error: fmt.Sprintf("didn't get next date: %v", err)})
		return
	}

	w.Write([]byte(nxtDate))

}

// NextDate - находит ближайшую дату срабатывания.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	//	fmt.Printf("now: %v\ndstar: %v\nrepeat: %v\n", now, dstart, repeat)
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
	// w ...
	if w, _ := regexp.MatchString(`^w `, repeat); w {
		result, err := nextWeekDay(now, date, repeat)
		if err != nil {
			return "", err
		}
		return result, nil
	}
	// m ...
	if m, _ := regexp.MatchString(`^m `, repeat); m {
		result, err := nextMonthDay(now, date, repeat)
		if err != nil {
			return "", err
		}
		return result, nil
	}

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
		// fmt.Println(nxtDate)
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
		//fmt.Println("nxtDate:", nxtDate, "now:", now, afterNow(nxtDate, now))
	}
	return nxtDate.Format(pattern)
}

// nextWeekDay - рассчитать следующую дату.
func nextWeekDay(now time.Time, dstart time.Time, repeat string) (string, error) {
	// формирование дней недели из repeat
	checkDays, err := getRepeatValues(repeat)
	if err != nil {
		return "", err
	}
	// создадим одномерную матрицу по кол-ву дней недели.
	// в те дни, в которые должна быть назначена задача поставим 1.
	//
	matrixWeek := make([]int, numDays)

	for _, item := range checkDays {
		matrixWeek[item-1] = 1
	}

	var nxtDate time.Time
	// //date := time.Now()
	// date := now
	// if afterNow(date, dstart) {
	// 	nxtDate = date
	// } else {
	// 	nxtDate = dstart.AddDate(0, 0, 2)
	// }
	nxtDate = dstart

	for {
		nxtDate = nxtDate.AddDate(0, 0, 1)
		i := int(nxtDate.Weekday())

		if int(nxtDate.Weekday()) == 0 {
			// Sunday
			if matrixWeek[6] == 1 {
				if afterNow(nxtDate, now) {
					break
				}
			}
		} else {
			// Monday-Saturday
			if matrixWeek[i-1] == 1 {
				if afterNow(nxtDate, now) {
					break
				}
			}
		}
	}

	return nxtDate.Format(pattern), nil
}
func getRepeatValues(repeat string) ([]int, error) {
	if len(repeat) == 0 {
		return []int{}, fmt.Errorf("expected to receive a rule repeat")
	}
	rpt := strings.Split(repeat, " ")
	if len(rpt) != 2 {
		return []int{}, fmt.Errorf("repeat isn't correct (%v)", repeat)
	}
	rptValues := strings.Split(rpt[1], ",")
	if len(rptValues) == 0 {
		return []int{}, fmt.Errorf("repeat isn't correct (%v)", repeat)
	}

	var checkDays []int
	for _, i := range rptValues {
		r, err := strconv.Atoi(i)
		if err != nil {
			return []int{}, err
		}
		if r < 0 || r > 7 {
			return []int{}, fmt.Errorf("repeat value should be > 0 and < 7: %v", repeat)
		}
		checkDays = append(checkDays, r)
	}
	return checkDays, nil
}

func nextMonthDay(now time.Time, dstart time.Time, repeat string) (string, error) {
	// формирование дней недели из repeat
	// currentDate := time.Now()
	// var date time.Time
	// if afterNow(currentDate, dstart) {
	// 	date = currentDate
	// } else {
	// 	date = dstart
	// }
	date := dstart
	// разбираем repeat и выбираем месяцы и дни
	checkDays, checkMonth, err := getRepeatValuesMoth(repeat)
	if err != nil {
		return "", err
	}

	// формирование вспомогательных матриц: месяцы(13 элементов), 12х32 - дни.
	matrixMonth := make([]int, maxMonth+1)
	var matrixDays [][]int

	// заполнение матрицы месяцы: если в правиле не указан месяц, то всем месяцам [1-12] = 1; иначе руководствоваться правилом
	if len(checkMonth) == 0 {
		for k := 1; k < len(matrixMonth); k++ {
			matrixMonth[k] = 1
		}
	} else {
		for _, item := range checkMonth {
			matrixMonth[item] = 1
		}
	}
	currentMonth := int(date.Month())
	// обход матрицы месяцев, если параметр == 0, то в соответсующей матрице дней установим []int{}.
	// иначе оценим праметры правила, и соответствующим дням правила установим 1.
	for monthNumber, isMonth := range matrixMonth {

		if isMonth == 1 {
			// определяем год, берём на основе значения now. Год нужен для нахождения высокосного года.
			var year int
			if monthNumber <= currentMonth {
				year = now.Year()
			} else {
				year = now.Year() + 1
			}
			length := getDaysInMonth(year, time.Month(monthNumber)) // кол-во дней в месяце

			//выделим 32 элемента на матрицу дней.
			// -1, -2 - установим 1, для остальных назначим 1 соотв-ому дню правила
			days := make([]int, maxDaysInMonth+1)

			for _, daysNumber := range checkDays {
				switch daysNumber {
				case -1:
					days[length] = 1
				case -2:
					days[length-1] = 1
				default:
					days[daysNumber] = 1
				}
			}
			matrixDays = append(matrixDays, days)
		} else {
			days := []int{}
			matrixDays = append(matrixDays, days)
		}
	}

	for {
		date = date.AddDate(0, 0, 1)
		m := int(date.Month())
		day := date.Day()

		if matrixMonth[m] == 1 {
			if matrixDays[m][day] == 1 {
				if afterNow(date, now) {
					break
				}
			}
		}
	}

	return date.Format(pattern), nil
}
func getDaysInMonth(year int, month time.Month) int {
	// Находим первый день текущего месяца
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)

	// Переходим к первому дню следующего месяца
	// AddDate(years, months, days)
	firstDayOfNextMonth := firstDay.AddDate(0, 1, 0)

	// Вычитаем один день, чтобы получить последний день текущего месяца
	lastDayOfCurrentMonth := firstDayOfNextMonth.AddDate(0, 0, -1)

	// Возвращаем день (количество дней)
	return lastDayOfCurrentMonth.Day()
}

// m <через запятую от 1 до 31, -1, -2> [через запятую от 1 до 12]
func getRepeatValuesMoth(repeat string) ([]int, []int, error) {
	if len(repeat) == 0 {
		return []int{}, []int{}, fmt.Errorf("expected to receive a rule repeat")
	}
	rpt := strings.Split(repeat, " ")
	if len(rpt) < 2 || len(rpt) > 3 {
		return []int{}, []int{}, fmt.Errorf("repeat isn't correct 2<(%v)<4", repeat)
	}
	// в правиле только дни
	if len(rpt) == 2 {
		days, err := getDaysMnth(rpt[1])
		if err != nil {
			return []int{}, []int{}, err
		}
		return days, []int{}, nil
	}
	// в правиле дни и месяцы
	days, months, err := getDaysAndMonths(rpt[1], rpt[2])
	if err != nil {
		return []int{}, []int{}, err
	}
	return days, months, nil
}

// getDaysMnth - выбрать дни правила.
func getDaysMnth(rptValues string) ([]int, error) {
	days := strings.Split(rptValues, ",")
	if len(days) == 0 {
		return []int{}, fmt.Errorf("repeat isn't correct (%v)", rptValues)
	}

	var checkDays []int
	for _, i := range days {
		d, err := strconv.Atoi(i)
		if err != nil {
			return []int{}, err
		}
		if d < minDaysInMonth || d > maxDaysInMonth {
			return []int{}, fmt.Errorf("repeat value should be -1, -2, 1-31: %v", rptValues)
		}
		checkDays = append(checkDays, d)
	}
	return checkDays, nil
}

// getDaysAndMonths - выбрать дни и месяцы для правила <m ..>
func getDaysAndMonths(rptdays string, rptmon string) ([]int, []int, error) {
	days, err := getDaysMnth(rptdays)
	if err != nil {
		return []int{}, []int{}, fmt.Errorf("repeat isn't correct rptdays=(%v)", rptdays)
	}
	months, err := getMonths(rptmon)
	if err != nil {
		return []int{}, []int{}, fmt.Errorf("repeat isn't correct rptmon=(%v)", rptmon)
	}
	return days, months, nil

}

// getMonths - выбрать месяцы для правила <m ..>
func getMonths(rptValues string) ([]int, error) {
	months := strings.Split(rptValues, ",")
	if len(months) == 0 {
		return []int{}, fmt.Errorf("repeat isn't correct (%v)", rptValues)
	}

	var checkMons []int
	for _, i := range months {
		m, err := strconv.Atoi(i)
		if err != nil {
			return []int{}, err
		}
		if m < minMonth || m > maxMonth {
			return []int{}, fmt.Errorf("repeat value should be 1-12: %v", rptValues)
		}
		checkMons = append(checkMons, m)
	}
	return checkMons, nil
}
