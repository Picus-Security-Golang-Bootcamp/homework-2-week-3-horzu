package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Books struct {
	Books []Book `json:"books"`
}

type Book struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Page      int    `json:"page"`
	Stock     int    `json:"stock"`
	Price     string `json:"price"`
	StockCode int    `json:"stockCode"`
	Isbn      int    `json:"ISBN"`
	Author    string `json:"author"`
}

func main() {
	jsonFile, err := os.Open("books.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	values, _ := ioutil.ReadAll(jsonFile)
	data := Books{}
	json.Unmarshal(values, &data)

	command := os.Args[1]
	switch command {
	case "list":
		list(data)
	case "search":
		searchTerm := strings.Join(os.Args[2:], " ")
		fmt.Println(search(data, searchTerm))
	case "get":
		searchedBookId := os.Args[2]
		fmt.Println(get(data, searchedBookId))
	}
}
func get(data Books, bookId string) string{
	id, err := strconv.Atoi(bookId)
	if err != nil {
		return ("String conversation error!")
	}
	for _, book := range data.Books {
		if id == book.Id {
			return book.Title
		}
	}
	return "Given id is not valid"
}

func search(data Books, searchTerm string) string {
	var foundBooks []string
	for _, book := range data.Books {
		if strings.Contains(strings.ToLower(book.Title), strings.ToLower(searchTerm)) && book.StockCode == 1 {
			foundBooks = append(foundBooks, book.Title)
		}
	}
	if len(foundBooks) > 0 {
		for _, found := range foundBooks {
			return found
		}
	}
	return "No books found"
}

func list(data Books) {
	for i := 0; i < len(data.Books); i++ {
		fmt.Println("Book #", i+1, "Title: "+data.Books[i].Title)
	}
}
