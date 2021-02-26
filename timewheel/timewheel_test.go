package timewheel

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeWheel_AfterFunc(t *testing.T) {

	tests := []struct {
		name   string
		interval time.Duration
	}{
		{name:"2s", interval: 1*time.Second},
		{name:"4s", interval: 4*time.Second},
		{name:"10s", interval: 10*time.Second},
		{name:"10s", interval: 10*time.Second},
		{name:"10s", interval: 10*time.Second},
		{name:"20s", interval: 20*time.Second},
		{name:"30s", interval: 30*time.Second},

		{name:"100s", interval: 100*time.Second},

	}
	allDone := make(chan bool, len(tests))
	tw, _ := New(1*time.Second, 10)
	tw.Start()
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			startTs := time.Now()
			interval := tt.interval
			tw.AfterFunc(interval, func() {
				s := startTs
				e := time.Now()
				diff := e.Sub(s).Seconds()
				fmt.Printf("s is: %v end is: %v sub is: %v after %v s exec\n",s, e, diff , interval.Seconds())
				allDone <- true
			})

		})
	}

	for k := range  allDone{
		fmt.Println(k)
	}
}