// Package yy converts incomplete dates to time.Time
//  Incomplete dates are:
//  -  with missing year digits
//  -  with missing year and month
//  -  with missing year,month and day
//
// Conversion is done by finding nearest valid date to reference date.
// Date validity is according std time package, for example no leap seconds.
//
// Motivation for this package is large number of financial protocols and file formats, where
// abbreviated dates are transmitted in relation to transmitted date.
// (For example expiration date printed on credit cards have 2 year digits and month: YY/MM)
//
// Supported date formats are:
//
//  YYY-MM-DD    find X such that XYYY-MM-DD is nearest valid date
//  YY-MM-DD     find XX such that XXYY-MM-DD is nearest valid date
//  Y-MM-DD      find XXX such that XXXY-MM-DD is nearest valid date
//  MM-DD        find XXXX such that XXXX-MM-DD is nearest valid date
//  DD           find XXXX-ZZ such that XXXX-ZZ-DD is nearest valid date
//  YYYY-MM      return YYYY-MM-01
//  YYY-MM       find X such that XYYY-MM-01 is nearest valid date
//  YY-MM        find XX such that XXYY-MM-01 is nearest valid date
//  Y-MM         find XXX such that XXXY-MM-01 is nearest valid date
//  MM           find XXXX such that XXXX-MM-01 is nearest valid date
//  YYYY-JJJ     year+julian day JJJ=1..365/6
//  YYY-JJJ      find X such that XYYY-JJJ is nearest valid date
//  YY-JJJ       find XX such that XXYY-JJJ is nearest valid date
//  Y-JJJ        find XXX such that XXXY-JJJ is nearest valid date
//  JJJ          find XXXX such that XXXX-JJJ is nearest valid date
//  +/-RRR       RRR days after/before today
//  YYYY-MM-DD   full date
//  YYYY         return YYYY-01-01
//  YYY          find X such that XYYY-01-01 is nearest valid date
//  YY           find XX such that XXYY-01-01 is nearest valid date
//  Y            find XXX such that XXXY-01-01 is nearest valid date
//  "nothing"    return reference date
//
// Above, 'nearest valid date' means nearest to reference date.
//
// Additionally, time components can be specified, but they don't participate in finding nearest date.
// If they are missing, hour, minute, second and fraction defaults to 0, location is copied from reference time.
package yy

import (
	"errors"
	"time"
)

var errInvalidDate = errors.New("invalid date")

// Convert IDate to time.Time.
// rt is reference time.
// All missing time parts defaults to 0.
// All missing date parts(day & month), not subject to finding, defaults to 1.
// Missing location defaults to coping location from reference time.
// If no any date component present, converts to reference date.
func Convert(rt time.Time, p *IDate) (time.Time, error) {
	y, mo, dd := rt.Date()
	var h, m, s, f int

	l := rt.Location()

	if p.L.Present() {
		l = p.L.Get()
	}

	if p.F.Present() {
		f = p.F.Get()
	}

	if p.S.Present() {
		s = p.S.Get()
	}

	if p.M.Present() {
		m = p.M.Get()
	}

	if p.H.Present() {
		h = p.H.Get()
	}

	var t Tm

	// if ! have some date   {
	if !p.R.Present() && !p.Mo.Present() && !p.D.Present() && !p.J.Present() && p.Y.Digits() == 0 {

		t.FromValues(y, mo, dd, h, m, s, f, l)
		return totm(&t, isValid)
	}

	if p.R.Present() {
		y, mo, dd = rt.AddDate(0, 0, p.R.Get()).Date()
		t.FromValues(y, time.Month(mo), dd, h, m, s, f, l)
		return totm(&t, isValid)
	}

	if p.J.Present() {
		// assert dd,mm == nil
		switch p.Y.Digits() {
		case 0:
			yf := newJ(1, y, 0, 1, p.J.Get(), h, m, s, f, l)
			return nearDateFind(rt, yf)
		case 1:
			yf := newJ(10, y/10, p.Y.Get(), 1, p.J.Get(), h, m, s, f, l)
			return nearDateFind(rt, yf)
		case 2:
			yf := newJ(100, y/100, p.Y.Get(), 1, p.J.Get(), h, m, s, f, l)
			return nearDateFind(rt, yf)
		case 3:
			yf := newJ(1000, y/1000, p.Y.Get(), 1, p.J.Get(), h, m, s, f, l)
			return nearDateFind(rt, yf)
		case 4:
			t.FromValues(p.Y.Get(), 1, p.J.Get(), h, m, s, f, l)
			return totm(&t, isValidJJJ)

		}
		panic(" year digits ???")
	}

	if p.D.Present() {
		if p.Mo.Present() {

			switch p.Y.Digits() {
			case 0:
				yf := newY(1, y, 0, p.Mo.Get(), p.D.Get(), h, m, s, f, l)
				return nearDateFind(rt, yf)
			case 1:
				yf := newY(10, y/10, p.Y.Get(), p.Mo.Get(), p.D.Get(), h, m, s, f, l)
				return nearDateFind(rt, yf)
			case 2:
				yf := newY(100, y/100, p.Y.Get(), p.Mo.Get(), p.D.Get(), h, m, s, f, l)
				return nearDateFind(rt, yf)
			case 3:
				yf := newY(1000, y/1000, p.Y.Get(), p.Mo.Get(), p.D.Get(), h, m, s, f, l)
				return nearDateFind(rt, yf)
			case 4:
				t.FromValues(p.Y.Get(), time.Month(p.Mo.Get()), p.D.Get(), h, m, s, f, l)
				return totm(&t, isValid)
			}
			panic(" year digits ???")
		}
		mf := newM(y, mo, p.D.Get(), h, m, s, f, l)
		return nearDateFind(rt, mf)
	}

	dd = 1

	if p.Mo.Present() {
		switch p.Y.Digits() {
		case 0:
			yf := newY(1, y, 0, p.Mo.Get(), 1, h, m, s, f, l)
			return nearDateFind(rt, yf)
		case 1:
			yf := newY(10, y/10, p.Y.Get(), p.Mo.Get(), 1, h, m, s, f, l)
			return nearDateFind(rt, yf)
		case 2:
			yf := newY(100, y/100, p.Y.Get(), p.Mo.Get(), 1, h, m, s, f, l)
			return nearDateFind(rt, yf)
		case 3:
			yf := newY(1000, y/1000, p.Y.Get(), p.Mo.Get(), 1, h, m, s, f, l)
			return nearDateFind(rt, yf)
		case 4:
			t.FromValues(p.Y.Get(), time.Month(p.Mo.Get()), 1, h, m, s, f, l)
			return totm(&t, isValid)
		}
		panic(" year digits ???")

	}
	mo = 1
	// assert p.Y.Digits() != 0 // dp.yyyy != nil

	switch p.Y.Digits() {

	case 1:
		yf := newY(10, y/10, p.Y.Get(), 1, 1, h, m, s, f, l)
		return nearDateFind(rt, yf)
	case 2:
		yf := newY(100, y/100, p.Y.Get(), 1, 1, h, m, s, f, l)
		return nearDateFind(rt, yf)
	case 3:
		yf := newY(1000, y/1000, p.Y.Get(), 1, 1, h, m, s, f, l)
		return nearDateFind(rt, yf)
	case 4:
		t.FromValues(p.Y.Get(), 1, 1, h, m, s, f, l)
		return totm(&t, isValid)

	default:
		panic(" ??")
	}
}

//////////////////////////////////////////////////////////////////

func getFormatData(res, data, format []byte, mask byte) []byte {
	res = res[:0]
	for i := range format {
		if format[i] == mask {
			res = append(res, data[i])
		}
	}
	return res
}

func getFormatNum(s setter, date, format []byte, mask byte) error {
	res := make([]byte, 0, len(format))

	res = getFormatData(res, date, format, mask)
	if len(res) == 0 {
		return nil
	}
	return s.Set(res) //strconv.Atoi(string(res))
}

// FromFormat converts date according to format to time.Time
// date & format are treated as strings.
//
// format define date,
// at positions of chars 'Y,M,D,J,h,m,s,f,L,R' in format,
// are expected symbols
// of 'year,month,day,julian day,hour,minute,second,fraction,timezone,relative days'
// in date. All other chars in format are ignored, corresponding positions in date also are ignored.
//
// Accepted patterns are:
//  Y      `\d{1,4}`              year, number of 'Y's is equal to number of year digits
//  M      `\d{2}`                month
//  D      `\d{2}`                day
//  J      `\d{3}`                julian day
//  h      `\d{2}`                hour
//  m      `\d{2}`                minute
//  s      `\d{2}`                seconds
//  f      `\d{1,9}`              fraction
//  L      `[+-](\d\d):?(\d\d)`   timezone offset or
//         `.+`                   timezone name
//                                Special names 'l' & 'z' are Local & UTC zones
//  R      `[+-]?\d+`             relative days
//
// rt are reference time.
func FromFormat(date, format []byte, rt time.Time) (time.Time, error) {

	//fmt.Printf("%s %s\n", date, format)
	var p IDate
	err := getFormatNum(&p.R, date, format, 'R')
	if err != nil {
		return time.Time{}, err
	}
	err = getFormatNum(&p.Y, date, format, 'Y')
	if err != nil {
		return time.Time{}, err
	}
	err = getFormatNum(&p.Mo, date, format, 'M')
	if err != nil {
		return time.Time{}, err
	}
	err = getFormatNum(&p.J, date, format, 'J')
	if err != nil {
		return time.Time{}, err
	}
	err = getFormatNum(&p.D, date, format, 'D')
	if err != nil {
		return time.Time{}, err
	}

	err = getFormatNum(&p.H, date, format, 'h')
	if err != nil {
		return time.Time{}, err
	}
	err = getFormatNum(&p.M, date, format, 'm')
	if err != nil {
		return time.Time{}, err
	}
	err = getFormatNum(&p.S, date, format, 's')
	if err != nil {
		return time.Time{}, err
	}
	err = getFormatNum(&p.F, date, format, 'f')
	if err != nil {
		return time.Time{}, err
	}

	err = getFormatNum(&p.L, date, format, 'L')
	if err != nil {
		return time.Time{}, err
	}

	r, e := Convert(rt, &p)
	if e != nil {
		return time.Time{}, e
	}
	return r, nil
}
