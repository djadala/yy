package yy

import (
	"fmt"
	"testing"
	"time"
)

////////////////

var ref = time.Date(2013, time.June, 10, 23, 1, 2, 3, time.UTC)

const tf = `2006-01-02 15:04:05.999999999 -0700 MST`

type testData struct {
	in, fmt  string
	ref, out string
}

var ta = []testData{
	{
		in:  "2017-11-03",
		fmt: "YYYY-MM-DD",
		out: `2017-11-03 00:00:00.000000000 +0000 UTC`,
	},
	{
		in:  "917-11-03",
		fmt: "YYY-MM-DD",
		out: `1917-11-03 00:00:00.000000000 +0000 UTC`,
	},
	{
		in:  "01-11-03",
		fmt: "YY-MM-DD",
		out: `2001-11-03 00:00:00.000000000 +0000 UTC`,
	},
	{
		in:  "1-11-03",
		fmt: "Y-MM-DD",
		out: `2011-11-03 00:00:00.000000000 +0000 UTC`,
	},
	{
		in:  "10-03",
		fmt: "MM-DD",
		out: `2013-10-03 00:00:00.000000000 +0000 UTC`,
	},

	{
		in:  "03",
		fmt: "DD",
		out: `2013-06-03 00:00:00.000000000 +0000 UTC`,
	},

	{
		in:  "-20",
		fmt: "RRR",
		out: `2013-05-21 00:00:00.000000000 +0000 UTC`,
	},
	{
		in:  "+40",
		fmt: "RRR",
		out: `2013-07-20 00:00:00.000000000 +0000 UTC`,
	},

	{
		in:  "123",
		fmt: "JJJ",
		out: `2013-05-03 00:00:00.000000000 +0000 UTC`,
	},
	{
		in:  "9-123",
		fmt: "Y-JJJ",
		out: `2009-05-03 00:00:00.000000000 +0000 UTC`,
	},
	{
		in:  "99-123",
		fmt: "YY-JJJ",
		out: `1999-05-03 00:00:00.000000000 +0000 UTC`,
	},
	{
		in:  "199-123",
		fmt: "YYY-JJJ",
		out: `2199-05-03 00:00:00.000000000 +0000 UTC`,
	},
	{
		in:  "1964-123",
		fmt: "YYYY-JJJ",
		out: `1964-05-02 00:00:00.000000000 +0000 UTC`,
	},

	{
		in:  "1234",
		fmt: "YYYY",
		out: `1234-01-01 00:00:00.000000000 +0000 UTC`,
	},
	{
		in:  "123",
		fmt: "YYY",
		out: `2123-01-01 00:00:00.000000000 +0000 UTC`,
	},
	{
		in:  "98",
		fmt: "YY",
		out: `1998-01-01 00:00:00.000000000 +0000 UTC`,
	},
	{
		in:  "9",
		fmt: "Y",
		out: `2009-01-01 00:00:00.000000000 +0000 UTC`,
	},

	{
		in:  "5678-02",
		fmt: "YYYY-MM",
		out: `5678-02-01 00:00:00.000000000 +0000 UTC`,
	},
	{
		in:  "804-04",
		fmt: "YYY-MM",
		out: `1804-04-01 00:00:00.000000000 +0000 UTC`,
	},
	{
		in:  "44-03",
		fmt: "YY-MM",
		out: `2044-03-01 00:00:00.000000000 +0000 UTC`,
	},
	{
		in:  "7-11",
		fmt: "Y-MM",
		out: `2017-11-01 00:00:00.000000000 +0000 UTC`,
	},

	{
		in:  "10",
		fmt: "MM",
		out: `2013-10-01 00:00:00.000000000 +0000 UTC`,
	},
	{
		in:  "02-29",
		fmt: "MM-DD",
		out: `2012-02-29 00:00:00.000000000 +0000 UTC`,
	},
	{
		in:  "31",
		fmt: "DD",
		out: `2013-01-31 00:00:00.000000000 +0000 UTC`,
		ref: `2013-02-01 00:00:00.000000000 +0000 UTC`,
	},

	///////////////////////////////////////////////////////////////////////////
	// time

	{
		in:  "11:22:33.1234",
		fmt: "hh:mm:ss.ffff",
		out: `2013-06-10 11:22:33.1234 +0000 UTC`,
	},
	{
		in:  "11:22:33",
		fmt: "hh:mm:ss",
		out: `2013-06-10 11:22:33.000000000 +0000 UTC`,
	},
	{
		in:  "11:22",
		fmt: "hh:mm",
		out: `2013-06-10 11:22:00.000000000 +0000 UTC`,
	},
	{
		in:  "11",
		fmt: "hh",
		out: `2013-06-10 11:00:00.000000000 +0000 UTC`,
	},

	{
		in:  "EET",
		fmt: "LLL",
		out: `2013-06-10 00:00:00.000000000 +0300 EST`,
	},
	{
		in:  "z",
		fmt: "L",
		out: `2013-06-10 00:00:00.000000000 +0000 UTC`,
	},
	{
		in:  "+0300",
		fmt: "LLLLL",
		out: `2013-06-10 00:00:00.000000000 +0300 EET`,
	},
	{
		in:  "+03:00",
		fmt: "LLLLLL",
		out: `2013-06-10 00:00:00.000000000 +0300 EET`,
	},
	{
		in:  "l",
		fmt: "L",
		ref: time.Date(2013, time.June, 10, 23, 1, 2, 3, time.UTC).String(),
		out: time.Date(2013, time.June, 10, 0, 0, 0, 0, time.Local).String(),
	},
}

func Test_01(t *testing.T) {

	var (
		n time.Time
	)

	for i := range ta {

		if ta[i].ref == "" {
			n = ref
		} else {
			var err error
			n, err = time.Parse(tf, ta[i].ref)
			if err != nil {
				t.Fatal(err)
			}
		}

		dt, err := FromFormat([]byte(ta[i].in), []byte(ta[i].fmt), n)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		o, err := time.Parse(tf, ta[i].out)
		if err != nil {
			t.Fatal(err)
		}
		// fmt.Println(dt)
		if !dt.Equal(o) {
			fmt.Println(dt, o)
			fmt.Printf("%s %s\n", ta[i].in, ta[i].fmt)
			t.Error("times dont match")
			t.Fail()

		}

	}

}
