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
	Price     float64 `json:"price"`
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
	case "delete":
		bookId := os.Args[2]
		fmt.Println(delete(data, bookId))
	}
}

func (b *Book) setStockCode(stockCode int) *Book {
	b.StockCode = stockCode
	return b
}

// delete function removes given book from the available items list. It sets StockCode to 0, which means book is not available.
func delete(data Books, bookId string) string{
	id, err := strconv.Atoi(bookId)
	if err != nil {
		return ("String conversation error!")
	}
	for i, book := range data.Books {
		if id == book.Id {
			book.setStockCode(0)
			data.Books = append(data.Books[:i], data.Books[i+1:]...)
			data.Books = append(data.Books, book)
			newData, err := json.Marshal(data)
			if err != nil{
				fmt.Println(err)
			}
			err = ioutil.WriteFile("books.json", newData, 0644)
			if err != nil {
				fmt.Println("Error: Couldn't write to file")
			}
			return fmt.Sprintf("%s has been remevod from the list", book.Title)
		}
	}
	return "Given id is not valid"
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

// search function searches given string in the books and returns matched books
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

// list function lists books in the Books struct
func list(data Books) {
	for i := 0; i < len(data.Books); i++ {
		if data.Books[i].StockCode == 1 {
			fmt.Printf("#%d Book: %s | Author: %s | Stock: %d \n", i+1,data.Books[i].Title, data.Books[i].Author, data.Books[i].Stock)
		}
	}
}
