package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func readCsv(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Can't read file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Can't read csv "+filePath, err)
	}

	return records
}

func shuffleQuiz(array [][]string) [][]string {
	dest := make([][]string, len(array))
	perm := rand.Perm(len(array))
	for i, v := range perm {
		dest[v] = array[i]
	}
	return dest
}

func main() {
	file := flag.String("file", "problem.csv", "csv file path")
	limit := flag.Int("limit", 30, "duration of the quiz")
	shuffle := flag.Bool("shuffle", false, "shuffle the quiz")
	flag.Parse()

	records := readCsv(*file)
	if *shuffle {
		records = shuffleQuiz(records)
	}

	correct := 0
	reader := bufio.NewReader(os.Stdin)
	size := len(records)

	fmt.Println("Press Enter to start...")
	reader.ReadBytes('\n')

	timer := time.NewTimer(time.Duration(*limit) * time.Second)
	for _, record := range records {
		fmt.Print(record[0] + " ")
		answerCh := make(chan string)
		go func() {
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			answerCh <- strings.ToLower(text)
		}()

		select {
		case <-timer.C:
			fmt.Printf("\n\n%d out of %d", correct, size)
			return
		case answer := <-answerCh:
			if answer == strings.ToLower(record[1]) {
				correct++
			}
		}
	}

	fmt.Printf("\n%d out of %d", correct, size)
}
