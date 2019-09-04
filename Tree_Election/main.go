package main

import (
	"fmt"
	"time"
)

//Token ...
type Token struct {
	Sender   string
	Priority int
}

//Neighbour ...
type Neighbour struct {
	ID   string
	From chan Token
	To   chan Token
}

func redirect(in chan Token, neigh Neighbour) {
	for {
		token := <-neigh.From
		in <- token
	}
}

func process(ID string, token Token, neighs ...Neighbour) {
	var parent Neighbour
	in := make(chan Token, 1)
	nmap := make(map[string]Neighbour)
	for _, neigh := range neighs {
		nmap[neigh.ID] = neigh
		go redirect(in, neigh)
	}

	if token.Sender == "init" {
		x := token
		for _, neigh := range neighs {
			token.Sender = ID
			neigh.To <- token
			tk := <-in
			if x.Priority < tk.Priority {
				x = tk
			}
		}

		for _, neigh := range neighs {
			neigh.To <- x
		}

	} else {
		tk := <-in
		parent = nmap[tk.Sender]
		x := token

		if x.Priority < tk.Priority {
			x = tk
		}

		token.Sender = ID
		neighs[0].To <- x

		for index, neigh := range neighs {

			if index == 0 {
				continue
			}
			token.Sender = ID
			neigh.To <- token

			tk := <-in

			if x.Priority < tk.Priority {
				x = tk
			}
		}
		token.Sender = ID
		parent.To <- x

		tk1 := <-in
		d := tk1
		fmt.Println(ID)
		fmt.Println(d)
		for _, neigh := range neighs {
			neigh.To <- tk1
		}

	}

}

func main() {

	qT := make(chan Token, 1)
	tQ := make(chan Token, 1)
	pS := make(chan Token, 1)
	sP := make(chan Token, 1)
	pR := make(chan Token, 1)
	rP := make(chan Token, 1)
	pQ := make(chan Token, 1)
	qP := make(chan Token, 1)

	t := 40
	s := 6
	r := 10
	q := 70
	p := 9

	go process("T", Token{Priority: t}, Neighbour{"Q", qT, tQ})
	go process("S", Token{Priority: s}, Neighbour{"P", pS, sP})
	go process("R", Token{Priority: r}, Neighbour{"P", pR, rP})
	go process("Q", Token{Priority: q}, Neighbour{"T", tQ, qT}, Neighbour{"P", pQ, qP})
	process("P", Token{"init", p}, Neighbour{"Q", qP, pQ}, Neighbour{"R", rP, pR}, Neighbour{"S", sP, pS})
	time.Sleep(3 * time.Second)
}
