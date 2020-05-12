package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func reader(resultChan chan int, goRoutinesPool chan int, done *bool, goRoutinesCount *int) {

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		goRoutinesPool <- 1
		*goRoutinesCount++

		go parser(scanner.Text(), resultChan, goRoutinesPool, goRoutinesCount)

	}

	*done = true

}

func parser(url string, resultChan chan int, goRoutinesPool chan int, goRoutinesCount *int) {

	resp, err := http.Get(url)

	result := 0

	if err == nil {

		body, err := ioutil.ReadAll(resp.Body)

		if err == nil {

			result = strings.Count(strings.ToUpper(string(body)), "GO")

			println("Count for "+url+": ", result)

		}

		_ = resp.Body.Close()

	}

	*goRoutinesCount--

	println("test")

	resultChan <- result

	_ = <-goRoutinesPool
}


func main() {

	resultChan := make(chan int)
	defer close(resultChan)

	goRoutinesPool := make(chan int, 5)
	defer close(goRoutinesPool)

	done := false

	goRoutinesCount := 0

	go reader(resultChan, goRoutinesPool, &done, &goRoutinesCount)

	func() {

		var total int
		for {

			total += <-resultChan

			if done && goRoutinesCount <= 0 {
				break
			}

		}

		log.Printf("Total: %v", total)

	}()

}
