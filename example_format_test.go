package yy_test

import (
	"fmt"
	"time"

	"github.com/djadala/yy"
)

func ExampleFromFormat() {
	var ref = time.Date(2013, time.June, 10, 23, 1, 2, 3, time.UTC)

	r, err := yy.FromFormat([]byte("99-123"), []byte("YY-JJJ"), ref)
	if err != nil {
		panic(err)
	}
	fmt.Println(r) // year 99 julian day 123 in 2013 is 1999-05-03
	// Output: 1999-05-03 00:00:00 +0000 UTC
}
