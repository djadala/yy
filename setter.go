package yy

import (
	"strconv"
	"time"
)

// Sets value from chars in []byte
type setter interface {
	Set([]byte) error
}

// Loc indicate if timezone present/absent in incomplete date
type Loc struct {
	l *time.Location
}

// Set sets timezone
// accept `[+-]\d\d:{0,1}\d\d` or `.+`
func (t *Loc) Set(v []byte) error {
	if v[0] == '+' || v[0] == '-' {
		if v[3] == ':' {
			return t.SetHHMM(v[:3], v[4:6])
		}
		return t.SetHHMM(v[:3], v[3:5])
	}
	return t.SetName(v)
}

// Present returns if timezone is present in incomplete date
func (t *Loc) Present() bool {
	return t.l != nil
}

// Get Returns Location
func (t *Loc) Get() *time.Location {
	return t.l
}

// SetHHMM sets timezone to hh:mm
func (t *Loc) SetHHMM(hs, ms []byte) error {
	h, err := strconv.Atoi(string(hs))
	if err != nil {
		return err
	}
	m, err := strconv.Atoi(string(ms))
	if err != nil {
		return err
	}
	if h < 0 {
		h = -h
		return t.SetS(-((h*60 + m) * 60))
	}
	return t.SetS((h*60 + m) * 60)
}

// SetName sets timezone from name, special names 'z' && 'l'
// set UTC && Local timezones
func (t *Loc) SetName(v []byte) error {
	// TODO trims spaces
	zs := string(v)
	if zs == "z" {
		t.l = time.UTC
		return nil
	}
	if zs == "l" {
		t.l = time.Local
		return nil
	}

	l, err := time.LoadLocation(zs)
	if err != nil {
		return err
	}
	t.l = l
	return nil
}

// SetS sets timezone from seconds offset
func (t *Loc) SetS(s int) error {
	t.l = time.FixedZone("", s)
	return nil
}

//////////////////////

// Int indicate if various parts of date (day,month,hour,min,sec)
// present/absent in incomplete date
type Int struct {
	present bool
	val     int
}

///////////////

// SetI sets Int value to v
func (i *Int) SetI(v int) {
	i.val = v
	i.present = true
}

// Set sets Int from chars in v
func (i *Int) Set(v []byte) error {
	iv, err := strconv.Atoi(string(v))
	if err == nil {
		i.present = true
		i.val = iv
	}
	return err
}

// Present returns presense of Int value
func (i *Int) Present() bool {
	return i.present
}

// Get returns Int value
func (i *Int) Get() int {
	return i.val
}

///////////////////

// Frac indicate if fractions present/absent in incomplete date
type Frac struct {
	Int
}

// Set sets Frac from chars in v
func (f *Frac) Set(v []byte) error {
	// force allocation of new backing array
	return f.Int.Set(append(v[:len(v):len(v)], []byte("000000000000")...)[:9])
}

//////////////////////

// Year indicate number of year decimal digits in incomplete date (0 digits=no year)
type Year struct {
	digits int8
	y      int
}

// Set sets year from chars in v, number of digits is set to len(v)
func (y *Year) Set(v []byte) error {
	if len(v) == 0 {
		y.digits = 0
		return nil
	}
	yv, err := strconv.Atoi(string(v))
	if err == nil {
		y.digits = int8(len(v))
		y.y = yv
	}
	return err
}

// SetDI sets digits & year to integers d, iy
func (y *Year) SetDI(d, iy int) {
	y.digits = int8(d)
	y.y = iy
}

// Get Returns (incomplete) year
func (y *Year) Get() int {
	return y.y
}

// Digits returns number of decimal digits in year
func (y *Year) Digits() int8 {
	return y.digits
}

// IDate represent components of incomplete date
type IDate struct {
	R, J, Mo, D, H, M, S Int
	F                    Frac
	Y                    Year
	L                    Loc
}
