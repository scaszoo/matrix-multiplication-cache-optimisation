/**
The intention of this code is to experiment with cpu cache behaviour and resulting performance.
A matrix multiplication demands processing of many bytes from memory.
To avoid optimizations from a specialised matrix multiplications library a manual multiplication is performed here.

First experiment is the multiplication of matrix with itself.
Due to the nature of matrix multiplications, one matrix is read per row, the other matrix is read per column while calculating the new result value

The complexity of the standard algorithm is O(i * j * k). For quadratic matrices the run time is cubic, of the order O(n^3)

The representation of a matrix in memory is per row.
Optimal use of the cpu cache is achieved when the data that is going to be used is stored in consecutive order in memory.
So typically one matrix has the optimal data representation in memory, the second doesn't.

To test the significance of that factor, this program prepares the second matrix with pivoting the data, so this second matrix can be accessed per row as well and therefore support efficient caching for both matrices.

The necessary processing time gets displayed.

This is the version employing concurrency to further investigate cache behaviour and resulting performance.
*/
package main

import (
	"fmt"
	"math/rand"
    "time"
	"runtime"
	"sync"
)

const (
	x = 1200 // column
	y = 1200 // row
)

var (
	// matrix [row][column]
	matrix [y][x]int
	matrixTurned [y][x]int
	matrixResult[y][x]int
	matrixTurnedResult[y][x]int
	waitgroup sync.WaitGroup // create waitgroup (empty struct)
	threads = runtime.NumCPU()
)

func calculateMatrix(matrixSource *[y][x]int, matrixTarget *[y][x]int, rowStart int, rowCount int){
	var i, j, k, tmp int
	for i = rowStart;  i<rowCount; i++ { // row
		for j = 0;  j<x; j++ { // column
			tmp = 0
			for k = 0;  k<x; k++ {
				tmp += matrixSource[k][j] * matrixSource[i][k]
			}
			matrixTarget[i][j] = tmp
		}
	}
	waitgroup.Done() // decrement waitgroup counter
}

func calculateMatrixTurned(matrixSource *[y][x]int, matrixSourceTurned *[y][x]int, matrixTarget *[y][x]int, rowStart int, rowCount int){
	var i, j, k, tmp int
	for i = rowStart;  i<rowCount; i++ { // row
		for j = 0;  j<x; j++ { // column
			tmp = 0
			for k = 0;  k<x; k++ {
				tmp += matrixSource[i][k] * matrixSourceTurned[j][k]
			}
			matrixTarget[i][j] = tmp
		}
	}
	waitgroup.Done()
}

func main() {
	fmt.Printf("Number of available threads: %d\n", threads)
	fmt.Printf("# of rows: %d, # of go routines thats going to get started: %d, rows per thread: %d\n\n\n", y, threads, y % threads)
	if y % threads != 0 {
		fmt.Printf("Ratio of rows and # of go routines doesn't divide properly, please adjust y! Remainder: %d", y % threads)
		return
	}

	var i, j, k, tmp int

	// initialize matrix and "turned" matrix with the same random data
	for i := 0;  i<y; i++ {
		for j := 0;  j<x; j++ {
			tmp = rand.Intn(100) // can be surprisingly large numbers (e.g. 10^10), without significantly slowing down the calculation
			matrix[i][j] = tmp
			matrixTurned[j][i] = tmp
        }
    }

	// Concurrent: calculate result matrix with regular matrix multiplication without optimization
	fmt.Println("Concurrent: calculate result matrix with regular matrix multiplication without optimization . . .")
	start := time.Now()
	for x := 1;  x<=threads; x++ {
		waitgroup.Add(1) // increment waitgroup counter
		go calculateMatrix(&matrix, &matrixResult, y/threads*(x-1), y/threads*x)
	}
	waitgroup.Wait() // blocks here
	fmt.Printf("Time elapsed: %s\n\n", time.Since(start))

	// Concurrent: calculate result matrix with regular and optimized matrix
	fmt.Println("Concurrent: calculate result matrix with regular and optimized matrix . . .")
	start = time.Now()
	for x := 1;  x<=threads; x++ {
		waitgroup.Add(1)
		go calculateMatrixTurned(&matrix, &matrixTurned, &matrixTurnedResult, y/threads*(x-1), y/threads*x)
	}
	waitgroup.Wait() // blocks here
	fmt.Printf("Time elapsed: %s\n\n", time.Since(start))

	// make sure the same results were calculated!
	for i = 0;  i<y; i++ {
		for j = 0;  j<x; j++ {
			if matrixTurnedResult[i][j] != matrixResult[i][j] {
				fmt.Printf("CRASH BURN\n %d %d\n", matrixResult[i][j], matrixTurnedResult[i][j])
				return
			}
		}
	}

	// Single threaded: calculate result matrix with regular matrix multiplication without optimization
	fmt.Println("Single threaded: calculate result matrix with regular matrix multiplication without optimization . . .")
	start = time.Now()
	for i = 0;  i<y; i++ { // row
		for j = 0;  j<x; j++ { // column
			tmp = 0
			for k = 0;  k<x; k++ {
				tmp += matrix[k][j] * matrix[i][k]
			}
			matrixResult[i][j] = tmp
        }
    }
    fmt.Printf("Time elapsed: %s\n\n", time.Since(start))

	// Single threaded: calculate result matrix with regular and optimized matrix
	fmt.Println("Single threaded: calculate result matrix with regular and optimized matrix . . .")
	start = time.Now()
	for i = 0;  i<y; i++ { // Y ZEILE
		for j = 0;  j<x; j++ { // X ZEILE
			tmp = 0
			for k = 0;  k<x; k++ {
				tmp += matrix[i][k] * matrixTurned[j][k]
			}
			matrixTurnedResult[i][j] = tmp
        }
    }
	fmt.Printf("Time elapsed: %s\n\n", time.Since(start))

	// make sure the same results were calculated!
	for i = 0;  i<y; i++ {
		for j = 0;  j<x; j++ {
			if matrixTurnedResult[i][j] != matrixResult[i][j] {
				fmt.Printf("CRASH BURN\n %d %d\n", matrixResult[i][j], matrixTurnedResult[i][j])
				return
			}
        }
    }
}
