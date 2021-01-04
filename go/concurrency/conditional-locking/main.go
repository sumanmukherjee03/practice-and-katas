package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	balance        int = 0
	mutex              = sync.Mutex{}
	moneyDeposited     = sync.NewCond(&mutex) // A conditional needs the address of the mutex/lock to be passed in
)

func main() {
	describe()
	go stingy()
	go spendy()
	time.Sleep(5 * time.Second) // We arent using any waitgroups, so give it enough time to finish
	// You will notice that Stingy will run a 100 times and complete
	// But spendy will run 50 times and not complete because the bank balance will reach 0 and no more
	// money will get deposited, so spendy will keep waiting
	// Eventually, the main proc will exit killing the spendy goroutine as well.
	fmt.Println("Final balance : ", balance)
}

func stingy() {
	for i := 0; i < 100; i++ {
		mutex.Lock() // Acquire the lock to mutate the state variable balance
		balance += 10
		str := fmt.Sprintf("STINGY COUNT : %d - Deposited 20 dollars. Current balance is : %d", i, balance)
		// A signal is sent to the conditional - remember this can only handle a single goroutine.
		// Use broadcast when multiple goroutines being blocked
		// When the other goroutine receives this signal, it will wake up and try to acquire the lock again.
		moneyDeposited.Signal()
		mutex.Unlock() // Release lock after signalling that money has been deposited
		fmt.Println(str)
		time.Sleep(3 * time.Millisecond)
	}
	fmt.Println(">>>>>>> Stingy done")
}

func spendy() {
	for i := 0; i < 100; i++ {
		mutex.Lock()         // Acquire lock to mutate the state variable - balance
		for balance-20 < 0 { // If there is not enough balance, wait
			// Wait operation will block until a "signal" is received on the conditional to continue further from the other goroutine.
			// The mutex needs to be acquired first and then the conditional is checked. If the condition is not met, the lock is released.
			// This gurantees that the other goroutine picks up the mutex and continues with it's work.
			// If you acquire the "mutex.Lock()" after this conditonal wait, then you will see an error like this - fatal error: sync: unlock of unlocked moneyDeposited.Wait()
			moneyDeposited.Wait()
		}
		balance -= 20
		str := fmt.Sprintf("SPENDY COUNT : %d - Withdrew 20 dollars. Current balance is : %d", i, balance)
		mutex.Unlock() // Release lock when the state variable mutation has been done
		fmt.Println(str)
		time.Sleep(3 * time.Millisecond)
	}
	fmt.Println(">>>>>> Spendy done")
}

func describe() {
	str := `
This example demonstrates simple conditonal locking.
Stingy saves 10 dollars a day.
Spendy spends 20 dollars a day.
They both operate on the same bank account.
But Spendy cant spend if there is less than 20 dollars in the bank account.

_____________________
	`
	fmt.Println(str)
}
