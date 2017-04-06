package yy

import (
	"errors"
	"time"
)

var invalidDate = errors.New("invalid date")

/*


                              incomplete   qualified
Package yy is used to determine time from convert abbreviated dates to full time.Time instances.

                              incomplete   qualified
Package yy is used to determine time from convert abbreviated dates to full time.Time instances.


Package yy converts incomplete dates to `time.Time`
	Conversion is done by finding nearest valid date to reference date



	So what date is 99-11-29 (YY-MM-DD) ?
	Unfortunately this question does not have answer, but if the question is 'what date is
	99-11-29, if today is 2000-01-01 ?' can be defined as nearest to today posible date, so
	99-11-29 mean 1999-11-29 if today is 2000-01-01.

    Motivation for this package is large number of financial protocols and file formats, where
    abbreviated dates are transmitted in relation to transmitted date.

    supported date formats are:

    YYYY-MM-DD   full normal date
    YYY-MM-DD    find X such that XYYY-MM-DD is nearest valid date
    YY-MM-DD     find XX such that XXYY-MM-DD is nearest valid date
    Y-MM-DD      find XXX such that XXXY-MM-DD is nearest valid date
    MM-DD        find XXXX such that XXXX-MM-DD is nearest valid date
    DD           find XXXX-ZZ such that XXXX-ZZ-DD is nearest valid date
    YYYY-JJJ 	 full julian date jjj=1..365/6
    YYY-JJJ      find X such that XYYY-JJJ is nearest valid date
    YY-JJJ       find XX such that XXYY-JJJ is nearest valid date
    Y-JJJ        find XXX such that XXXY-JJJ is nearest valid date
    JJJ          find XXXX such that XXXX-JJJ is nearest valid date

	+/-RRR		 RRR days after/before today

  date formats:
  "2006-01-02T15:04:05.33@-0700"
  [date][#time][@zone]
  date=[reldate][jdate][mdate]
  reldate=+/-d{1,3}

  +ddd
  -ddd
  jjj
  y-jjj
  yy-jjj
  yyy-jjj
  yyyy-jjj
  dd
  mm-dd
  y-mm-dd
  yy-mm-dd
  yyy-mm-dd
  yyyy-mm-dd
  mm
  y-mm
  yy-mm
  yyy-mm
  yyyy-mm
  yyyy




  time formats:
  hh
  hh:mm[timezone]
  hh:mm:ss[timezone]
  hh:mm:ss.d+[timezone]

  timezone formats:
  +/-hh:mm
  +/-hhmm
  @timezoneName
  z
  l

  [+-](\d\d):{0,1}(\d\d)
  =(.+)
  z
  l


  -+nnnn
  jjj

  d
  m-d
  y-m-d
  y-m
  y



*/

// Convert IDate to time.Time
// rt is reference time
// All missing time parts defaults to 0
// All missing date parts(day & month), not subject to finding, defaults to 1
// Missing location defaults to coping location from reference time
// If no any date component present, converts to reference date
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

// преобразува от символи до дата, формата на датата е зададен в format:
// където има 'Y' на съответното място в data има цифра от годината
// броят на цифрите в годината е = на броя 'Y'
// R relative, останалите полета за дата не трябва да ги има
// J julian date, не трябва да има месец и ден
// D - ден
// M - месец
// h,m,s,f час, минути, секунди
// L timezone
// rt - reference date/time
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
