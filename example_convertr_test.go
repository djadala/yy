package yy_test

import (
	"fmt"
	"time"

	"github.com/djadala/yy"
)

func ExampleConvert_r() {
	var d yy.IDate
	var ref = time.Date(2013, time.June, 10, 23, 1, 2, 3, time.UTC)

	err := d.R.Set([]byte("-22"))
	if err != nil {
		panic(err)
	}

	r, err := yy.Convert(ref, &d)
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
	/////
	err = d.R.Set([]byte("22"))
	if err != nil {
		panic(err)
	}

	r, err = yy.Convert(ref, &d)
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
	/////
	err = d.R.Set([]byte("+23"))
	if err != nil {
		panic(err)
	}

	r, err = yy.Convert(ref, &d)
	if err != nil {
		panic(err)
	}
	fmt.Println(r)

	// Output: 2013-05-19 00:00:00 +0000 UTC
	// 2013-07-02 00:00:00 +0000 UTC
	// 2013-07-03 00:00:00 +0000 UTC
}
