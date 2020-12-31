package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const matrixDimension int = 250

var (
	matrixA [matrixDimension][matrixDimension]int
	matrixB [matrixDimension][matrixDimension]int
	result  [matrixDimension][matrixDimension]int
	rwMutex = sync.RWMutex{}                  // Create a reader writer lock. We need the reader lock to get read access for reading the rows/cols of the input matrix
	cond    = sync.NewCond(rwMutex.RLocker()) // The condition needs to be able to unlock the reader lock of the rwMutex. So pass that to the condition
	wg      = sync.WaitGroup{}                // This is gonna be used by the main function to know when the work for computing a row is done
)

func main() {
	describe()
	rand.Seed(time.Now().UnixNano())

	start := time.Now()
	fmt.Println("Starting now")

	wg.Add(matrixDimension) // Use `matrixDimension` number of workers to compute the rows for each set of matrix multiplication
	for r := 0; r < matrixDimension; r++ {
		go computeRow(r) // Create a goroutine for each row in the `matrixDimension x matrixDimension` resulting matrix for calculating a row for various different matrices
	}

	// Multiply 100 sets of matrices
	for i := 0; i < 100; i++ {
		wg.Wait()      // Wait until the worker goroutines that are supposed to compute the rows are ready for each matrix or finished the multiplication for a previous set of matrices
		rwMutex.Lock() // Acquire the writers lock so that the matrix can be generated
		generateRandMatrix(&matrixA)
		generateRandMatrix(&matrixB)
		wg.Add(matrixDimension) // Add the number of goroutines back to the wait group so that the next matrix multiplication can be done
		rwMutex.Unlock()        // Unlock the writers lock so that the readers lock can be acquired by the worker threadpool
		// Broadcast to the goroutines to wake up and start computing the rows for matrixA x matrixB
		// Make sure to unlock the mutex before you signal/broadcast the child threads to continue
		// because unlike a cond.Wait() a cond.Signal() or cond.Broadcast() does not unlock the mutex.
		cond.Broadcast()
	}

	elapsed := time.Since(start)
	fmt.Println("Processing took : ", elapsed)
}

func generateRandMatrix(m *[matrixDimension][matrixDimension]int) {
	for i := 0; i < matrixDimension; i++ {
		for j := 0; j < matrixDimension; j++ {
			m[i][j] = rand.Intn(10) - 5
		}
	}
}

// This function will be run as a goroutine
// So start an infinite loop so that it can continuously keep computing rows for matrixA x matrixB at a specific row number
// matrixA and matrixB will keep changing over time
func computeRow(r int) {
	rwMutex.RLock() // Acquire the reader lock for comnputing the row for the resulting matrix
	for {
		// This wait group is marked as done to signify that row computation of the resulting matrix for the previous set of matrices is complete
		// and that the worker is ready to proceed with the next set of matrices
		wg.Done()
		cond.Wait() // Wait to receive a broadcast signal from the main proc to know that matrices are generated and it is ok to compute rows
		for c := 0; c < matrixDimension; c++ {
			for x := 0; x < matrixDimension; x++ {
				result[r][c] = result[r][c] + matrixA[r][x]*matrixB[x][c]
			}
		}
	}
}

func describe() {
	str := `
This example is for nxn matrix multiplication with n parallel goroutines to reduce the runtime.
The normal simple process is O(n^3).
Each goroutine here computes a row for a specific row number.


_____________________
	`
	fmt.Println(str)
}
