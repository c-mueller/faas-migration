package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// -----------------------------------------------------------------------------
// Input Configuration
// -----------------------------------------------------------------------------

var (
	endpoint = flag.String("endpoint", "undefined", "The endpoint to the ToDo API implementation")
)

func init() {
	flag.Parse()
}

// -----------------------------------------------------------------------------
// Model Definitions
// -----------------------------------------------------------------------------

// This type defines the input object used to create a ToDo Item
type InsertRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
type ToDoItem struct {
	ID                 string `json:"ID"`
	Title              string `json:"title"`
	Description        string `json:"description"`
	InsertionTimestamp int64  `json:"insertion_timestamp"`
	DoneTimestamp      int64  `json:"done_timestamp"`
}

// -----------------------------------------------------------------------------
// API Interactions
// -----------------------------------------------------------------------------

type api string

// Request all items from the api
func (endpoint api) ListItems() ([]ToDoItem, error) {
	res, err := http.Get(fmt.Sprintf("%slst", endpoint))
	if err != nil {
		return nil, err
	}
	// Best Practice: Close body when Function is done...
	defer res.Body.Close()

	return nil, nil
}

func (endpoint api) InsertItem(i InsertRequest) (*ToDoItem, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	res, err := http.Post(fmt.Sprintf("%sput", endpoint), "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	// Best Practice: Close body when Function is done...
	defer res.Body.Close()

	resdata, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var item ToDoItem
	err = json.Unmarshal(resdata, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// -----------------------------------------------------------------------------
// Main Method
// -----------------------------------------------------------------------------

func main() {
	appUrl := *endpoint
	if appUrl == "undefined" {
		fmt.Println("The -endpoint Flag is required!")
		os.Exit(1)
	}
	// Append Trailing / if it is missing
	if !strings.HasSuffix(appUrl, "/") {
		appUrl += "/"
	}
	api := api(appUrl)
	api.CheckAvailability()
	items := api.InsertItems()
	api.CheckListItems(items)
}

// -----------------------------------------------------------------------------
// Tests
// -----------------------------------------------------------------------------

func (endpoint api) CheckAvailability() {
	// Call list to check if API is up
	fmt.Printf("Checking if API is Available....")
	_, err := endpoint.ListItems()
	if err != nil {
		fmt.Printf("Fail!\n")
		fmt.Printf("Ping API by calling List has Failed!\nError: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Ok!\n")
}

func (endpoint api) InsertItems() []ToDoItem {
	fmt.Printf("Inserting 50 ToDo Items...")
	insertionRequests := make([]InsertRequest, 0)
	items := make([]ToDoItem, 0)
	for i := 0; i < 50; i++ {
		ireq := InsertRequest{
			Title:       fmt.Sprintf("Todo-Item-#%d", i),
			Description: GenerateRandomString(300),
		}
		insertionRequests = append(insertionRequests, ireq)

		item, err := endpoint.InsertItem(ireq)
		if err != nil {
			fmt.Printf("Fail!\nError: %s\n", err.Error())
			os.Exit(1)
		}
		items = append(items, *item)
	}
	fmt.Printf("Done!\nValidating Responses....")
	if len(items) != len(insertionRequests) {
		fmt.Printf("Fail!\nResponse Count is not equal to Request Count\n")
		os.Exit(1)
	}
	for k, req := range insertionRequests {
		res := items[k]
		if res.Title != req.Title {
			if len(items) != len(insertionRequests) {
				fmt.Printf("Fail!\nA Title Did not match\n")
				os.Exit(1)
			}
		}
		if res.Description != req.Description {
			if len(items) != len(insertionRequests) {
				fmt.Printf("Fail!\nA Description Did not match\n")
				os.Exit(1)
			}
		}
	}
	fmt.Printf("Done!\n")
	return items
}

func (endpoint api) CheckListItems(items []ToDoItem) {
	
}

// -----------------------------------------------------------------------------
// Utility Methods
// -----------------------------------------------------------------------------

const randomStringCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

var charsetRunes []rune

func init() {
	charsetRunes = bytes.Runes([]byte(randomStringCharset))
	rand.Seed(time.Now().Unix())
}

func GenerateRandomString(l int) string {
	id := ""
	for i := 0; i < l; i++ {
		idx := rand.Intn(len(charsetRunes))
		id += string(charsetRunes[idx])
	}
	return id
}