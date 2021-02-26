package subset

import (
	"strconv"
	"testing"
)



func TestSubset(t *testing.T) {
	type args struct {
		backends   []string
		clientID   int
		subsetSize int
	}
	var backends = make([]string, 100000)
	for i  := 0 ; i < 100000 ; i++{
		backends[i] =  strconv.Itoa(i)
	}
	tCount := 1000
	tests := make([]args, tCount)
	for i:= 1; i <= tCount; i++{
		tests[i-1] = args{
			backends:   backends,
			clientID:   i,
			subsetSize: 100,
		}
	}

	for _, tt := range tests {
		t.Run("case" + strconv.Itoa(tt.clientID), func(t *testing.T) {
			got := Subset(tt.backends, tt.clientID, tt.subsetSize)
			t.Logf("Subset() = %v", got)

		})
	}
}