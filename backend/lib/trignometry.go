package lib

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/tkennon/ticker"
)

func SinSampleGenerator(sin, rchan chan float64, radius float64) {
	t := ticker.NewConstant(time.Second)

	r := radius

	err := t.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer t.Stop()

	for range t.C {
		select {
		case nr := <-rchan:
			r = nr
		default:
		}

		sin <- math.Sin(2*math.Pi*rand.Float64()) * r
	}
}

func CosSampleGenerator(cos, rchan chan float64, radius float64) {
	t := ticker.NewConstant(time.Second)

	r := radius

	err := t.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer t.Stop()

	for range t.C {
		select {
		case nr := <-rchan:
			r = nr
		default:
		}

		cos <- math.Cos(2*math.Pi*rand.Float64()) * r
	}
}
