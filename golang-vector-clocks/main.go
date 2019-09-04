package main

import (
	"fmt"
	"time"
)

//Message ...
type Message struct {
	Body      string
	Timestamp []int
}

func event(pid int, clock []int) []int {
	clock[pid-1]++
	fmt.Printf("Event in process pid=%v. Counter=%v\n", pid, clock)
	return clock
}

func check(x, y []int) []int {
	for index, elementY := range y {
		if x[index] > elementY {
			y[index] = x[index]

		}
	}

	return y
}

func calcTimestamp(recvTimestamp, clock []int, pid int) []int {
	check(recvTimestamp, clock)
	clock[pid-1]++
	return clock

}

func sendMessage(ch chan Message, pid int, clock []int) []int {
	clock[pid-1]++
	ch <- Message{"Test msg!!!", clock}
	fmt.Printf("Message sent from pid=%v. Counter=%v\n", pid, clock)
	return clock

}

func receiveMessage(ch chan Message, pid int, clock []int) []int {
	message := <-ch
	clock = calcTimestamp(message.Timestamp, clock, pid)
	fmt.Printf("Message received at pid=%v. Counter=%v\n", pid, clock)
	return clock
}

func processOne(ch12, ch21 chan Message) {
	pid := 1
	clock := []int{0, 0, 0}
	clock = event(pid, clock)
	clock = sendMessage(ch12, pid, clock)
	clock = event(pid, clock)
	clock = receiveMessage(ch21, pid, clock)
	clock = event(pid, clock)

}

func processTwo(ch12, ch21, ch23, ch32 chan Message) {
	pid := 2
	clock := []int{0, 0, 0}
	clock = receiveMessage(ch12, pid, clock)
	clock = sendMessage(ch21, pid, clock)
	clock = sendMessage(ch23, pid, clock)
	clock = receiveMessage(ch32, pid, clock)

}

func processThree(ch23, ch32 chan Message) {
	pid := 3
	clock := []int{0, 0, 0}
	clock = receiveMessage(ch23, pid, clock)
	clock = sendMessage(ch32, pid, clock)

}

func main() {
	oneTwo := make(chan Message, 100)
	twoOne := make(chan Message, 100)
	twoThree := make(chan Message, 100)
	threeTwo := make(chan Message, 100)

	go processOne(oneTwo, twoOne)
	go processTwo(oneTwo, twoOne, twoThree, threeTwo)
	go processThree(twoThree, threeTwo)

	time.Sleep(5 * time.Second)
}
