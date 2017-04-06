package yy

import (
	"time"
)

// Tm helps to record date components && check date validity
type Tm struct {
	Year                      int
	Month                     time.Month
	Day, Hour, Min, Sec, Nsec int
	Loc                       *time.Location
}

// convert Tm to time.Time, validating first Tm via f
func totm(t *Tm, f func(*Tm) bool) (time.Time, error) {
	if f(t) {
		return t.Date(), nil
	}
	return time.Time{}, invalidDate

}

// Converts Tm to time.Time
func (t *Tm) Date() time.Time {
	return time.Date(t.Year, t.Month, t.Day, t.Hour, t.Min, t.Sec, t.Nsec, t.Loc)
}

// Sets fields in Tm from time.Time
func (t *Tm) From(d time.Time) {
	t.Year, t.Month, t.Day = d.Date()
	t.Hour, t.Min, t.Sec = d.Clock()
	t.Nsec = d.Nanosecond()
	t.Loc = d.Location()
}

// Sets fields in Tm from separate values
func (t *Tm) FromValues(y int, mo time.Month, d, h, m, s, f int, l *time.Location) {
	t.Year, t.Month, t.Day = y, mo, d
	t.Hour, t.Min, t.Sec = h, m, s
	t.Nsec = f
	t.Loc = l
}

// returns if Tm represent valid date/time
// Validity is according to std time package,
// for example leap seconds are invalid
func (t *Tm) IsValid() bool {
	return isValid(t)

}
