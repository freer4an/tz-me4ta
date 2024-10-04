package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
)

type Object struct {
	A int `json:"a"`
	B int `json:"b"`
}

const DEFAULT_MAX_WORKERS = 3

var (
	maxWorkers int // flag -n
	sum        atomic.Int32
	wg         sync.WaitGroup
)

func main() {
	data, err := parseJsonObjects("data.json")
	if err != nil {
		log.Fatal(err)
	}

	chJobs := make(chan Object)
	go func() {
		defer close(chJobs)
		for i := 0; i < len(data); i++ {
			chJobs <- data[i]
		}
	}()

	wg.Add(maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go func() {
			worker(chJobs)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("Sum:", int(sum.Load()))
}

func init() {
	flag.IntVar(&maxWorkers, "n", DEFAULT_MAX_WORKERS, "go run . -n 10")
	flag.Parse()
}

func parseJsonObjects(path string) ([]Object, error) {
	var data []Object
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error while reading data.json: %w", err)
	}

	if err = json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("error while parsing data.json: %w", err)
	}
	return data, nil
}

func worker(jobs <-chan Object) {
	for job := range jobs {
		sumSingleObject := job.A + job.B
		sum.Add(int32(sumSingleObject))
	}
}

// Программа должна считать файл и вычислить сумму всех чисел. Для этого она должна:
// ⁃ Нужно распараллелить вычисление по горутинам
// ⁃ Количество горутин, для параллельной обработки нужно получить при запуске программы через аргумент
// ⁃ Вывести общий результат в консоль
// ⁃ Опубликовать решение в github.com и предоставить доступ
