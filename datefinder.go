package yy

import (
	"time"
)

func abs(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

// датата е валидна ако след нормализиране не се промени нищо
// Date is valid if after normalization, nothing changed
func isValid(t *Tm) bool {
	var r Tm

	r.From(t.Date())
	return (*t) == r
}

// предполага че t е попълнен с julian date: month = 1 && day = jjj
// Assume day = JJJ(julian date), month=1
func isValidJJJ(t *Tm) bool {
	var r Tm
	r.From(t.Date())
	r.Month = 1
	r.Day = t.Day

	return (*t) == r
}

// при търсене на месец е валидна ако след нормализирането съвпада всичко без годината и месеца
// while finding month, date is valid, if after normalization only changed parts are year and month
func isValidM(t *Tm) bool {
	var r Tm
	r.From(t.Date())
	r.Month = t.Month
	r.Year = t.Year

	return (*t) == r
}

// return nearest to reference date/time
func nearDate(ref, d1, d2, d3 time.Time) time.Time {

	a := abs(ref.Unix() - d1.Unix())
	b := abs(ref.Unix() - d2.Unix())
	c := abs(ref.Unix() - d3.Unix())

	if a < b {
		if a < c {
			return d1
		} // else {
		return d3
		// }
	} // else {
	if b < c {
		return d2
	} // else {
	return d3
	// }
	// }
}

////////////////////////////////////////////////////////////

type dateFinder interface {
	gen(int) (Tm, bool)
}

// struct for finding year
type yearFind struct {
	dt            Tm
	yearHi, scale int
}

func newY(scale, yhi, y, mo, d, h, m, s, f int, l *time.Location) *yearFind {
	return &yearFind{
		dt: Tm{
			Year:  y,
			Month: time.Month(mo),
			Day:   d,
			Hour:  h,
			Min:   m,
			Sec:   s,
			Nsec:  f,
			Loc:   l,
		},
		yearHi: yhi,
		scale:  scale,
	}
}

func (y *yearFind) gen(i int) (Tm, bool) {
	ret := y.get(i)
	return ret, isValid(&ret)
}

func (y *yearFind) get(i int) Tm {
	return Tm{
		Year:  y.dt.Year + (y.yearHi+i)*y.scale,
		Month: y.dt.Month,
		Day:   y.dt.Day,
		Hour:  y.dt.Hour,
		Min:   y.dt.Min,
		Sec:   y.dt.Sec,
		Nsec:  y.dt.Nsec,
		Loc:   y.dt.Loc,
	}
}

// struct for finding year in julian date
type yearFindJulian struct {
	yearFind // share methods
}

func newJ(scale, yhi, y, mo, d, h, m, s, f int, l *time.Location) *yearFindJulian {
	return &yearFindJulian{
		yearFind: yearFind{
			dt: Tm{
				Year:  y,
				Month: time.Month(mo),
				Day:   d,
				Hour:  h,
				Min:   m,
				Sec:   s,
				Nsec:  f,
				Loc:   l,
			},
			yearHi: yhi,
			scale:  scale,
		},
	}
}

func (y *yearFindJulian) gen(i int) (Tm, bool) {
	ret := y.get(i)
	return ret, isValidJJJ(&ret)
}

// struct for finding month
type monthFind struct {
	dt Tm
}

func newM(y int, mo time.Month, d, h, m, s, f int, l *time.Location) *monthFind {
	return &monthFind{
		dt: Tm{
			Year:  y,
			Month: time.Month(mo),
			Day:   d,
			Hour:  h,
			Min:   m,
			Sec:   s,
			Nsec:  f,
			Loc:   l,
		},
	}
}

func (y *monthFind) get(i int) Tm {
	return Tm{
		Year:  y.dt.Year,
		Month: y.dt.Month + time.Month(i),
		Day:   y.dt.Day,
		Hour:  y.dt.Hour,
		Min:   y.dt.Min,
		Sec:   y.dt.Sec,
		Nsec:  y.dt.Nsec,
		Loc:   y.dt.Loc,
	}
}

func (y *monthFind) gen(i int) (Tm, bool) {
	ret := y.get(i)
	return ret, isValidM(&ret) // валидна е когото всичко освен месеца и годината след нормализирне съвпаднат
}

var (
	_ dateFinder = &yearFind{}
	_ dateFinder = &yearFindJulian{}
	_ dateFinder = &monthFind{}
)

///////////////////////////////////////////////////////////

// find closest date to ref from generated dates 0 +1 -1 +2 -2 +3 -3 ......
func nearDateFind(ref time.Time, v dateFinder) (time.Time, error) {

	all := make([]time.Time, 0, 4)

	if t, valid := v.gen(0); valid {
		all = append(all, t.Date())
	}

	for i := 1; i < 9; i++ {
		if t, valid := v.gen(i); valid {
			all = append(all, t.Date())
		}

		if t, valid := v.gen(-i); valid {
			all = append(all, t.Date())
		}
		if len(all) >= 3 {
			return nearDate(ref, all[0], all[1], all[2]), nil
		}
	}
	return time.Time{}, errInvalidDate
}
