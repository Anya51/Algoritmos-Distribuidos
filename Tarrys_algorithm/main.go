package main

import (
	"fmt"
	"time"
)

//Token ...
type Token struct {
	Sender string
}

//Neighbour ...
type Neighbour struct {
	ID   string
	From chan Token
	To   chan Token
}

func redirect(in chan Token, neigh Neighbour) {
	token := <-neigh.From
	in <- token
}

func popN(m map[string]Neighbour, key string) Neighbour {
	element := m[key]
	delete(m, key)
	return element
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
		fmt.Printf("* %s é raiz.\n", ID)
		for _, neigh := range neighs {
			var tk Token
			tk.Sender = ID

			neigh.To <- tk
			time.Sleep(100 * time.Millisecond)

			tk = <-in
			fmt.Printf("From %s to %s\n", tk.Sender, ID)

		}
		fmt.Println("Fim!")
	} else {
		tk := <-in
		fmt.Printf("From %s to %s\n", tk.Sender, ID)

		nm := nmap
		parent = popN(nm, tk.Sender)
		fmt.Printf("* %s é pai de %s\n", parent.ID, ID)

		for _, n := range nm {
			tk.Sender = ID
			n.To <- tk
			tk = <-in
			fmt.Printf("From %s to %s\n", tk.Sender, ID)
		}
		tk.Sender = ID
		parent.To <- tk
	}

}

func main() {

	pQ := make(chan Token, 10)
	qP := make(chan Token, 10)

	pR := make(chan Token, 10)
	rP := make(chan Token, 10)

	qR := make(chan Token, 10)
	rQ := make(chan Token, 10)

	rT := make(chan Token, 10)
	tR := make(chan Token, 10)

	rS := make(chan Token, 10)
	sR := make(chan Token, 10)

	sT := make(chan Token, 10)
	tS := make(chan Token, 10)

	go process("T", Token{}, Neighbour{"R", rT, tR}, Neighbour{"S", sT, tS})
	go process("S", Token{}, Neighbour{"R", rS, sR}, Neighbour{"T", tS, sT})
	go process("R", Token{}, Neighbour{"Q", qR, rQ}, Neighbour{"P", pR, rP}, Neighbour{"T", tR, rT}, Neighbour{"S", sR, rS})
	go process("Q", Token{}, Neighbour{"P", pQ, qP}, Neighbour{"R", rQ, qR})
	process("P", Token{Sender: "init"}, Neighbour{"Q", qP, pQ}, Neighbour{"R", rP, pR})

}
