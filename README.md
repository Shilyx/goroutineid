# goroutineid

get current goroutine id

example:

```golang
package main

import (
	"fmt"
	"time"

	"github.com/shilyx/goroutineid"
)

func main() {
	fmt.Println("main", goroutineid.Get())

	go func() {
		fmt.Println("func", goroutineid.Get())
		time.Sleep(time.Hour)
	}()

	go func() {
		fmt.Println("func2", goroutineid.Get())
		time.Sleep(time.Hour)
	}()

	fmt.Println("main", goroutineid.Get())

	<-make(chan interface{})
}
```