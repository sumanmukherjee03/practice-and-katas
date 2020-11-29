package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	describe()
	variant1()
}

func variant1() {
	msgs := make(chan string)
	signals := make(chan bool)

	// Here’s a non-blocking receive.
	// If a value is available on messages then select will take the <-messages case with that value.
	// If not it will immediately take the default case.
	select {
	case m := <-msgs:
		fmt.Println("Message received : ", m)
	default:
		fmt.Println("Nothing to receive")
	}

	// A non-blocking send works similarly.
	// Here msg cannot be sent to the messages channel, because the channel has no buffer and there is no receiver.
	// Therefore the default case is selected.
	msg := "hello1"
	select {
	case msgs <- msg:
		fmt.Println("Message sent : ", msg)
	default:
		fmt.Println("Nothing to send")
	}

	// We can use multiple cases above the default clause to implement a multi-way non-blocking select.
	// Here we attempt non-blocking receives on both messages and signals.
	select {
	case nm := <-msgs:
		fmt.Println("Message received : ", nm)
	case sig := <-signals:
		fmt.Println("Signal received : ", sig)
	default:
		fmt.Println("No activity")
	}
}

func variant2() {
	msgs := make(chan string, 1)
	signals := make(chan bool)
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		time.Sleep(2 * time.Second)
		select {
		case m := <-msgs:
			fmt.Println("Message received in go routine 1: ", m)
		default:
			fmt.Println("Nothing to receive")
		}
		wg.Done()
	}()

	go func() {
		msg := "hello1"
		select {
		case msgs <- msg:
			fmt.Println("Message sent : ", msg)
		default:
			fmt.Println("Nothing to send")
		}
		wg.Done()
	}()

	go func() {
		time.Sleep(2 * time.Second)
		select {
		case nm := <-msgs:
			fmt.Println("Message received in goroutine 3: ", nm)
		case sig := <-signals:
			fmt.Println("Signal received in go routine 3: ", sig)
		default:
			fmt.Println("No activity")
		}
		wg.Done()
	}()

	wg.Wait()
}

func describe() {
	str := `
Basic sends and receives on channels are blocking.
However, we can use select with a default clause to implement non-blocking sends, receives, and even non-blocking multi-way selects.


_____________________
	`
	fmt.Println(str)
}
