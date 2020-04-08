package main

import (
	"fmt"
	"time"
	"sync"
)

func sumArr (arr []int, wg *sync.WaitGroup, i int) int{

	defer wg.Done() //Add
	sum:=0
	for _,value := range arr {	
		sum += value
	}

	fmt.Printf("Sum%v : %v \n",i, sum) 
	return sum
}

func main() {

	// arr1 := []int{1,2,3,4,5}
	// arr2 := []int{6,7,8,9,10}
	// sumGo1 := 0
	// sumGo2 := 0

	// var wg sync.WaitGroup //Add
	// wg.Add(2) //Add

	// // Goroutine 1
	// go func() {
	// 	sumGo1 = sumArr(arr1, &wg, 1)	
	// }()
	// // Goroutine 2
	// go func() {
	// 	sumGo2 = sumArr(arr2, &wg, 2)
	// }()

	// wg.Wait() //Add
	// fmt.Printf("sumGo1 + sumGo2: %v \n", sumGo1 + sumGo2)
	// time.Sleep(2 * time.Second)
	number := 589
    sqrch := make(chan int)
    cubech := make(chan int)
    go calcSquares(number, sqrch)
    go calcCubes(number, cubech)
	squares, cubes := <-sqrch, <-cubech
	fmt.Println("Final output", squares )
	time.Sleep(2 * time.Second)
    fmt.Println("Final output", squares + cubes)
}

func calcSquares(number int, squareop chan int) {  
    sum := 0
    for number != 0 {
        digit := number % 10
        sum += digit * digit
        number /= 10
    }
    squareop <- sum
}

func calcCubes(number int, cubeop chan int) {  
    sum := 0 
    for number != 0 {
        digit := number % 10
        sum += digit * digit * digit
        number /= 10
    }
    cubeop <- sum
} 
