package yy_test

import (
	"fmt"
	"time"

	"github.com/djadala/yy"
)

func ExampleConvert() {
	var d yy.IDate
	var ref = time.Date(2013, time.June, 10, 23, 1, 2, 3, time.UTC)

	d.Mo.SetI(2)
	d.D.SetI(29)

	r, err := yy.Convert(ref, &d)
	if err != nil {
		panic(err)
	}
	fmt.Println(r) // 29 Feb in  2013 is 2012-02-29
	// Output: 2012-02-29 00:00:00 +0000 UTC
}
