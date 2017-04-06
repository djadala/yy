# yy


Package yy converts incomplete dates to `time.Time`

Conversion is done by finding nearest valid date to reference date.

Incomplete dates are:
* with missing year digits
* with missing year and month
* with missing year,month and day

Date validity is according std `time` package, for example no leap seconds


examples:

```go
package main

import (
	"fmt"
	"github.com/djadala/yy"
	"time"
)

func main() {
	var d yy.IDate
	var ref = time.Date(2013, time.June, 10, 23, 1, 2, 3, time.UTC)

	d.Mo.SetI(2)
	d.D.SetI(29)

	r, err := yy.Convert(ref, &d)
	if err != nil {
		panic(err)
	}
	fmt.Println(r) // 29 Feb in  2013 is 2012-02-29
}

```



```go
package main

import (
	"fmt"
	"github.com/djadala/yy"
	"time"
)

func main() {
	var ref = time.Date(2013, time.June, 10, 23, 1, 2, 3, time.UTC)

	r, err := yy.FromFormat([]byte("99-123"), []byte("YY-JJJ"), ref)
	if err != nil {
		panic(err)
	}
	fmt.Println(r) // year 99 julian day 123 in 2013 is 1999-05-03
}
```

[godoc](http://godoc.org/github.com/djadala/yy)
