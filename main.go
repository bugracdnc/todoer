package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TODOER STRUCT
type Todoer struct {
	Id          uuid.UUID `json:"id"`
	Todo        string    `json:"todo"`
	UpdateDate  string    `json:"updateDate"`
	CreatedDate string    `json:"createdDate"`
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
	return &Todoer{Id: uuid.UUID{}, UpdateDate: "", CreatedDate: "", Todo: todo}
}

var repo []Todoer

//END TODOER STRUCT

// HTTP PARAMETERS
const serverUrl = "http://192.168.1.241"
const serverPort = 8080
const basePath = "/api/v1/todoer"

//END HTTP PARAMETERS

// REQUESTS
func makeListRequest() []Todoer {
	requestURL := fmt.Sprintf("%s:%d%s", serverUrl, serverPort, basePath)
	res, errGet := http.Get(requestURL)
	if errGet != nil {
		panic(errGet)
	}

	if res.StatusCode != http.StatusOK {
		return []Todoer{}
	}

	body, errReadBody := io.ReadAll(res.Body)
	if errReadBody != nil {
		panic(errReadBody)
	}
	var todoers []Todoer
	errUnmarshal := json.Unmarshal(body, &todoers)
	if errUnmarshal != nil {
		panic(errUnmarshal)
	}
	sort.Slice(todoers, func(i, j int) bool {
		d1, err1 := todoers[i].getCreatedDate()
		d2, err2 := todoers[j].getCreatedDate()
		if err1 == nil && err2 == nil {
			return d1.Before(d2)
		}
		if err1 != nil {
			return false
		}
		return false
	})
	return todoers
}

func makeAddRequest(str string) (bool, string) {
	todoer := buildTodoer(str)
	requestURL := fmt.Sprintf("%s:%d%s", serverUrl, serverPort, basePath)

	body, err := json.Marshal(todoer)

	if err != nil {
		panic(err)
	}

	res, errPost := http.Post(requestURL, "application/json", bytes.NewBuffer(body))

	if errPost != nil {
		panic(errPost)
	}

	if res.StatusCode == http.StatusCreated {
		return true, fmt.Sprintf("Added new to Location: %s%s\n", serverUrl, res.Header.Get("Location"))
	}
	return false, strconv.Itoa(res.StatusCode)
}

func makeDeleteRequest(uuid uuid.UUID) bool {
	client := &http.Client{}
	requestURL := fmt.Sprintf("%s:%d%s/%s", serverUrl, serverPort, basePath, uuid.String())
	req, err := http.NewRequest("DELETE", requestURL, nil)
	if err != nil {
		panic(err)
	}

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNoContent {
		return true
	} else {
		fmt.Printf("unhandled error: %d\n", res.StatusCode)
		return false
	}
}

//END REQUESTS

// COMMANDS
func list() {
	repo = makeListRequest()
	if len(repo) > 0 {
		fmt.Printf("Current List: {%d items}\n", len(repo))
		for index, todoer := range repo {
			date, errDate := todoer.getCreatedDate()
			dateString := fmt.Sprintf("%d:%d %02d/%02d", date.Hour(), date.Minute(), date.Day(), date.Month())
			if errDate != nil {
				fmt.Printf("%d) %s\n", index+1, todoer.Todo)
			} else {
				fmt.Printf("%d) %s - %s\n", index+1, todoer.Todo, dateString)
			}
		}
		fmt.Println()
	} else {
		fmt.Println("Todoers is empty!")
	}
}
func add(str string) {
	isAdded, result := makeAddRequest(str)
	if isAdded {
		fmt.Println(result)
	} else {
		fmt.Printf("error: status code %s", result)
	}
}

func remove(str string) {
	if len(repo) == 0 {
		repo = makeListRequest()
	}

	if str == "all" {
		for _, todoer := range repo {
			makeDeleteRequest(todoer.Id)
		}
		fmt.Printf("Cleared the list!")
	} else {
		index, err := strconv.Atoi(str)
		if err != nil {
			panic(err)
		}
		index--
		if index <= len(repo) {
			if makeDeleteRequest(repo[index].Id) {
				fmt.Printf("Deleted todo: %s\n\n", repo[index].Todo)
				list()
			}
		}
	}
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
		"remove": "removes given todo",
		"clear":  "clears the todo list (same as 'todoer remove all')"}

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

	repo = makeListRequest()

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
		case "remove":
			if len(args) > 1 {
				remove(strings.Join(args[1:], " "))
			}
		case "clear":
			remove("all")
		default:
			fmt.Println("error: not implemented yet")
		}
	} else {
		usage()
	}
}
