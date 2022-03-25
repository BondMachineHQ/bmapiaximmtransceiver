package bmapiaximmtransceiver

import (
	"context"
	"log"
	"os"
	"sync"
)

func AXImmTransceiver(ctx context.Context, device string, debug bool) (chan<- uint8, <-chan uint8, <-chan struct{}) {

	fi, err := os.Lstat(device)
	if err != nil {
		log.Fatal(err)
	}

	port, err := os.OpenFile(device, os.O_APPEND|os.O_RDWR, fi.Mode())
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	src := make(chan uint8)
	dst := make(chan uint8)
	ended := make(chan struct{})
	buff := make([]byte, 100)
	wg.Add(1)
	go func() {
		if debug {
			log.Println("transceiver: receiver starting")
		}
		defer port.Close()
		defer wg.Done()
		for {
			n, err := port.Read(buff)
			if err != nil {
				break
			}
			if n == 0 {
				break
			}
			for i := 0; i < n; i++ {
				select {
				case <-ctx.Done():
					if debug {
						log.Println("transceiver: receiver exiting")
					}
					return
				case dst <- buff[i]:
				}
			}
		}
		if debug {
			log.Println("transceiver: receiver exiting")
		}
	}()

	wg.Add(1)
	go func() {
		if debug {
			log.Println("transceiver: sender staring")
		}
		defer port.Close()
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				if debug {
					log.Println("transceiver: sender exiting")
				}
				return
			case data := <-src:
				_, err := port.Write([]byte{data})
				if err != nil {
					return
				}
			}
		}
	}()

	go func() {
		if debug {
			log.Println("transceiver: started")
		}
		wg.Wait()
		ended <- struct{}{}
		if debug {
			log.Println("transceiver: exiting")
		}
	}()

	return src, dst, ended
}
