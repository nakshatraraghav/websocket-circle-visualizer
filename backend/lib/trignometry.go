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

	err := t.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer t.Stop()

	for range t.C {
		sin <- math.Sin(2*math.Pi*rand.Float64()) * radius
	}
}

func CosSampleGenerator(cos, rchan chan float64, radius float64) {
	t := ticker.NewConstant(time.Second)

	err := t.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer t.Stop()

	for range t.C {
		cos <- math.Cos(2*math.Pi*rand.Float64()) * radius
	}
}
