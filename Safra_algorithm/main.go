package main

import (
	"fmt"
	"math"
	"sync"
)

//Token ...
type Token struct {
	Sender string
	sum    int
}

//Msg ...
type Msg struct {
	Dist   float64
	Sender string
}

//Neighbour ...
type Neighbour struct {
	ID    string
	From  chan interface{}
	To    chan interface{}
	Value float64
}

func redirect(in1 chan Token, in2 chan Msg, neigh Neighbour) {
	for v := range neigh.From {
		switch t := v.(type) {
		case Token:
			in1 <- t
		case Msg:
			in2 <- t
		default:
			fmt.Printf("Other: %v \n", t)
		}
	}
}

func printToken(token Token, ID string, count *int) {
	fmt.Printf("From %s to %s. sum: %d, count: %d\n", token.Sender, ID, token.sum, *count)
}

func printMsg(msg Msg, ID string, count int) {
	fmt.Printf("** %s: msg.Dist: %f msg.Sender: %s count: %d\n\n", ID, msg.Dist, msg.Sender, count)
}

func sh(ID string, w *sync.WaitGroup, wg *sync.WaitGroup, count *int, in chan Token, msg Msg, parent Neighbour, nmap map[string]Neighbour, neighs []Neighbour) {

	if msg.Sender == "init" {
		cicles := 0
		size := len(neighs)
		token := Token{}
		for {
			token.sum = 0

			fmt.Printf("========Start=========\n")
			fmt.Printf("cicle: %d\n", cicles)
			cicles = cicles + 1
			token.Sender = ID
			token.sum = token.sum + *count
			wg.Wait()

			neighs[0].To <- token

			for i := 1; i < size; i++ {
				tk := <-in
				printToken(tk, ID, count)
				tk.Sender = ID
				wg.Wait()
				neighs[i].To <- tk
			}

			token = <-in
			printToken(token, ID, count)
			fmt.Printf("=========End==========\n\n")
			if token.sum == 0 {
				fmt.Println("Fim!")
				break
			}
		}

	} else {

		for tk := range in {
			tk.sum = tk.sum + *count

			printToken(tk, ID, count)
			for _, neigh := range neighs {
				if parent.ID == "" {
					parent = nmap[tk.Sender]
					//	fmt.Printf("* %s é parent de %s\n", parent.ID, ID)
				}

				if parent.ID != neigh.ID {
					tk.Sender = ID
					//	wg.Wait()

					neigh.To <- tk
					tk = <-in
					printToken(tk, ID, count)
				}
			}
			tk.Sender = ID
			parent.To <- tk
		}

	}
	w.Done()

}

func process1(ID string, w *sync.WaitGroup, msg Msg, neighs ...Neighbour) {
	var parent Neighbour
	var wg sync.WaitGroup

	in1 := make(chan Token, 100)
	in2 := make(chan Msg, 100)

	nmap := make(map[string]Neighbour)
	for _, neigh := range neighs {
		nmap[neigh.ID] = neigh
		go redirect(in1, in2, neigh)
	}

	count := 0
	wg.Add(1)
	if msg.Sender == "init" {
		message := msg
		printMsg(msg, ID, count)
		msg.Sender = ID
		for _, neigh := range neighs {
			neigh.To <- msg
			count++
		}

		go sh(ID, w, &wg, &count, in1, message, parent, nmap, neighs)
		wg.Done()

		for range in2 {
			count--
		}

	} else {

		go sh(ID, w, &wg, &count, in1, msg, parent, nmap, neighs)

		for ms := range in2 {
			count--
			neigh := nmap[ms.Sender]
			res := ms.Dist + neigh.Value

			if res < msg.Dist {

				msg.Dist = res
				parent = neigh

				for _, neigh := range neighs {

					if neigh.ID != parent.ID {
						msg.Sender = ID
						neigh.To <- msg
						count++
					}
				}
			}
			printMsg(msg, ID, count)
			wg.Done()
			wg.Add(1)
		}

	}
}

func main() {
	var w sync.WaitGroup

	pQ := make(chan interface{}, 10)
	qP := make(chan interface{}, 10)

	ValueQP := 1.0

	pR := make(chan interface{}, 10)
	rP := make(chan interface{}, 10)

	ValuePR := 3.0

	qR := make(chan interface{}, 10)
	rQ := make(chan interface{}, 10)

	ValueQR := 1.0

	rT := make(chan interface{}, 10)
	tR := make(chan interface{}, 10)

	ValueRT := 4.0

	rS := make(chan interface{}, 10)
	sR := make(chan interface{}, 10)

	ValueRS := 1.0

	sT := make(chan interface{}, 10)
	tS := make(chan interface{}, 10)

	ValueST := 1.0
	w.Add(1)
	go process1("T", &w, Msg{Dist: math.Inf(1)}, Neighbour{"R", rT, tR, ValueRT}, Neighbour{"S", sT, tS, ValueST})
	go process1("S", &w, Msg{Dist: math.Inf(1)}, Neighbour{"R", rS, sR, ValueRS}, Neighbour{"T", tS, sT, ValueST})
	go process1("R", &w, Msg{Dist: math.Inf(1)}, Neighbour{"Q", qR, rQ, ValueQR}, Neighbour{"P", pR, rP, ValuePR}, Neighbour{"T", tR, rT, ValueRT}, Neighbour{"S", sR, rS, ValueRS})
	go process1("Q", &w, Msg{Dist: math.Inf(1)}, Neighbour{"P", pQ, qP, ValueQP}, Neighbour{"R", rQ, qR, ValueQR})
	go process1("P", &w, Msg{Sender: "init", Dist: 0.0}, Neighbour{"Q", qP, pQ, ValueQP}, Neighbour{"R", rP, pR, ValuePR})
	w.Wait()
}
