package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const URI_R = "https://api.chucknorris.io/jokes/random"
const URI_C = "https://api.chucknorris.io/jokes/random?category="
const URI_ALL_C = "https://api.chucknorris.io/jokes/categories"

type Data struct {
	Value string `json:"value"`
}

type States struct {
	Categories []string `json:"categories"`
	Value      string   `json:"value"`
}

func GetJokeCate(wg *sync.WaitGroup) error {
	defer wg.Done()
	rand.Seed(time.Now().UnixNano())
	var data States
	arr, err := GetCategories()
	if err != nil {
		return err
	}
	nameCate := arr[rand.Intn(len(arr))]
	resp, err := http.Get(fmt.Sprintf(URI_C+"%s", nameCate))
	if err != nil {
		return err
	}
	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(buffer, &data); err != nil {
		return err
	}
	file, err := os.OpenFile(nameCate, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)
	readString, _ := reader.ReadString('\n')
	if strings.Contains(readString, data.Value) {
		return nil
	}
	file.Write([]byte(data.Value + "\n"))
	file.Close()
	return nil
}

func GetCategories() ([]string, error) {
	var arrCate []string
	resp, err := http.Get(URI_ALL_C)
	if err != nil {
		return nil, err
	}
	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buffer, &arrCate); err != nil {
		return nil, err
	}
	return arrCate, nil
}

func GetRandomJoke() {
	var data Data
	resp, err := http.Get(URI_R)
	if err != nil {
		log.Println(err)
	}
	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	if err := json.Unmarshal(buffer, &data); err != nil {
		log.Println(err)
	}
	fmt.Println(data.Value)
}

func main() {
	wg := &sync.WaitGroup{}
	fs := flag.NewFlagSet("dump", flag.ExitOnError)
	count := fs.Int("n", 5, "Write in file categories joke")
	fs.String("random", "", "View random joke")

	if len(os.Args) < 2 {
		fmt.Println("expected one more arguments")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "dump":
		if err := fs.Parse(os.Args[2:]); err != nil {
			log.Println(err)
		}
		for i := 0; i < *count; i++ {
			wg.Add(1)
			go GetJokeCate(wg)
		}
		fmt.Println("Wait... write in files")
	case "random":
		GetRandomJoke()
	default:
		fmt.Println("Use command: \tdump -n 5 ; for give state and write in files\n\t\trandom ; for view random joke")
	}
	wg.Wait()
}
