package bmapiaximmtransceiver

import (
	"context"
	"fmt"
	"testing"
)

func TestSource(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, dest := AXImmTransceiver(ctx, "/dev/bm")

	i := 0
	for n := range dest {
		fmt.Println(i, n)
		if i == 10 {
			break
		}
		i++
	}
}
