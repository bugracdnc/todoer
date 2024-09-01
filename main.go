package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Todoer struct
type Todoer struct {
	Id          string `json:"id"`
	Todo        string `json:"todo"`
	UpdateDate  string `json:"updateDate"`
	CreatedDate string `json:"createdDate"`
}

func (t *Todoer) getCreatedDate() (time.Time, error) {
	if len(t.CreatedDate) > 0 {
		dateString := t.CreatedDate
		dateString = strings.Replace(dateString, "T", " ", 1)
		dateString = dateString[:strings.IndexByte(dateString, '.')]

		return time.Parse(time.DateTime, dateString)
	} else {
		return time.Time{}, errors.New("t.CreatedDate is null")
	}
}

func buildTodoer(todo string) *Todoer {
	return &Todoer{Id: "", UpdateDate: "", CreatedDate: "", Todo: todo}
}

//end Todoer struct

// HTTP PARAMETERS
const serverPort = 8080
const basePath = "/api/v1/todoer"

//END HTTP PARAMETERS

// REQUESTS
func makeListRequest() {
	requestURL := fmt.Sprintf("http://localhost:%d%s", serverPort, basePath)
	res, errGet := http.Get(requestURL)
	if errGet != nil {
		panic(errGet)
	}

	if res.StatusCode == http.StatusOK {
		body, errReadBody := io.ReadAll(res.Body)
		if errReadBody != nil {
			panic(errReadBody)
		}
		var todoers []Todoer
		errUnmarshal := json.Unmarshal(body, &todoers)
		if errUnmarshal != nil {
			panic(errUnmarshal)
		}
		for index, todoer := range todoers {
			date, errDate := todoer.getCreatedDate()
			dateString := fmt.Sprintf("%d:%d %02d/%02d", date.Hour(), date.Minute(), date.Day(), date.Month())
			if errDate != nil {
				fmt.Printf("%d) %s\n", index+1, todoer.Todo)
			} else {
				fmt.Printf("%d) %s - %s\n", index+1, todoer.Todo, dateString)
			}
		}
	}
}

func makeAddRequest(str string) (*Todoer, error) {
	todoer := buildTodoer(str)
	requestURL := fmt.Sprintf("http://localhost:%d%s", serverPort, basePath)

	body, err := json.Marshal(todoer)

	if err != nil {
		fmt.Printf("error in json.Marshal: %s", err)
		return nil, err
	}

	res, errPost := http.Post(requestURL, "application/json", bytes.NewBuffer(body))

	if errPost != nil {
		panic(errPost)
	}

	fmt.Printf("client: status code: %d\n", res.StatusCode)
	if res.StatusCode == http.StatusCreated {
		fmt.Printf("client: Location: %s\n", res.Header.Get("Location"))
	}

	return todoer, nil
}

//END REQUESTS

// COMMANDS
func list() {
	makeListRequest()
}
func add(str string) {
	_, err := makeAddRequest(str)

	if err != nil {
		return
	}
}
func move(str string) {
	fmt.Printf("move: %s", str)
}
func remove(str string) {
	fmt.Printf("remove: %s", str)
}

//END COMMANDS

// USAGE
func printUsage(pad int, cmd string, info string) {
	fmt.Printf("\t%-"+strconv.Itoa(pad)+"v %s\n", cmd, info)
}

func usage() {
	//Usage strings
	cmds := map[string]string{
		"-h":     "returns usage information",
		"--help": "returns usage information",
		"add":    "adds new todo",
		"move":   "moves given todo to new index",
		"remove": "removes given todo"}

	//Print hardcoded usage example and commands title
	fmt.Println("usage:\n\ttodoer <command> [arguments]")
	fmt.Println()
	fmt.Println("commands:")

	//Get longest command length
	lenmax := 0
	for key := range cmds {
		if len(key) > lenmax {
			lenmax = len(key)
		}
	}

	//Print commands accommodating for the longest command
	for key, info := range cmds {
		printUsage(lenmax, key, info)
	}
}

//END USAGE

func main() {
	args := os.Args[1:]

	//fmt.Printf("debug:\tcommand=%s\n\n", args)

	if len(args) > 0 {
		switch args[0] {
		case "-h", "--help":
			usage()
		case "list":
			list()
		case "add":
			if len(args) > 1 {
				add(strings.Join(args[1:], " "))
			}
		case "move":
			if len(args) > 1 {
				move(strings.Join(args[1:], " "))
			}
		case "remove":
			if len(args) > 1 {
				remove(strings.Join(args[1:], " "))
			}
		default:
			fmt.Println("error: not implemented yet")
		}
	} else {
		usage()
	}
}
