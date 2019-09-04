package main

import (
	"fmt"
	"math"
	"time"
)

//Msg ...
type Msg struct {
	Dist   float64
	Sender string
}

//Neighbour ...
type Neighbour struct {
	ID    string
	From  chan Msg
	To    chan Msg
	Value float64
}

func redirect(in chan Msg, neigh Neighbour) {
	for {
		msg := <-neigh.From
		in <- msg
	}
}

func process1(ID string, msg Msg, neighs ...Neighbour) {
	var parent Neighbour
	in := make(chan Msg, 100)
	nmap := make(map[string]Neighbour)
	for _, neigh := range neighs {
		nmap[neigh.ID] = neigh
		go redirect(in, neigh)
	}

	if msg.Sender == "init" {
		fmt.Printf("* %s é raiz.\n", ID)

		for _, neigh := range neighs {
			msg.Sender = ID
			neigh.To <- msg
		}

	} else {
		for {
			ms := <-in
			res := ms.Dist + nmap[ms.Sender].Value

			if res < msg.Dist {
				msg.Dist = res
				fmt.Printf("%s:\n msg.Dist: %f\n ms.Sender: %s\n\n", ID, msg.Dist, ms.Sender)

				parent = nmap[ms.Sender]

				msg.Sender = ID

				for _, neigh := range neighs {

					if neigh.ID != parent.ID {
						neigh.To <- msg
					}
				}
			}
		}

	}
}

func main() {

	pQ := make(chan Msg, 1)
	qP := make(chan Msg, 1)

	ValueQP := 1.0

	pR := make(chan Msg, 1)
	rP := make(chan Msg, 1)

	ValuePR := 3.0

	qR := make(chan Msg, 1)
	rQ := make(chan Msg, 1)

	ValueQR := 1.0

	rT := make(chan Msg, 1)
	tR := make(chan Msg, 1)

	ValueRT := 4.0

	rS := make(chan Msg, 1)
	sR := make(chan Msg, 1)

	ValueRS := 1.0

	sT := make(chan Msg, 1)
	tS := make(chan Msg, 1)

	ValueST := 1.0

	go process1("T", Msg{Dist: math.Inf(1)}, Neighbour{"R", rT, tR, ValueRT}, Neighbour{"S", sT, tS, ValueST})
	go process1("S", Msg{Dist: math.Inf(1)}, Neighbour{"R", rS, sR, ValueRS}, Neighbour{"T", tS, sT, ValueST})
	go process1("R", Msg{Dist: math.Inf(1)}, Neighbour{"Q", qR, rQ, ValueQR}, Neighbour{"P", pR, rP, ValuePR}, Neighbour{"T", tR, rT, ValueRT}, Neighbour{"S", sR, rS, ValueRS})
	go process1("Q", Msg{Dist: math.Inf(1)}, Neighbour{"P", pQ, qP, ValueQP}, Neighbour{"R", rQ, qR, ValueQR})
	process1("P", Msg{Sender: "init", Dist: 0.0}, Neighbour{"Q", qP, pQ, ValueQP}, Neighbour{"R", rP, pR, ValuePR})
	time.Sleep(3 * time.Second)
}
