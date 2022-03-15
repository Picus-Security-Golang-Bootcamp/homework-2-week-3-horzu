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
	Id        int     `json:"id"`
	Title     string  `json:"title"`
	Page      int     `json:"page"`
	Stock     int     `json:"stock"`
	Price     float64 `json:"price"`
	StockCode string  `json:"stockCode"`
	ISBN      string  `json:"ISBN"`
	Author    string  `json:"author"`
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
	case "buy":
		bookId := os.Args[2]
		quantity := os.Args[3]
		fmt.Println(buy(data, bookId, quantity))
	default:
		fmt.Println("Entered command is not valid!")
		usage()
	}
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("list: Lists available books.")
	fmt.Println("search <bookName>: Searches given string in the available books.")
	fmt.Println("get <bookID>: Gets book information of given id.")
	fmt.Println("delete <bookID>: Sets stock of given id's book as 0. It will be not in the list but you can get the information of the book with get command.")
	fmt.Println("buy <bookID> <quantity>: Buys given quantity of the given book and returns the new state of the book.")
}

// buy function reduces the given book's stock quantity as given order quantity
func buy(data Books, bookId string, quantity string) string {
	id, err := strconv.Atoi(bookId)
	if err != nil {
		return "String conversation error on bookID"
	}
	order, err := strconv.Atoi(quantity)
	if err != nil {
		return "String conversation error on order quantity"

	}
	for i, book := range data.Books {
		if id == book.Id && book.Stock >= order && book.Stock > 0 {
			newQuantity := book.Stock - order
			book.setStock(newQuantity)
			data.Books = append(data.Books[:i], data.Books[i+1:]...)
			data.Books = append(data.Books, book)
			newData, err := json.Marshal(data)
			if err != nil {
				return "json converting error"
			}
			err = ioutil.WriteFile("books.json", newData, 0644)
			if err != nil {
				return "Error: Couldn't write to file"
			}
			bookInfo := fmt.Sprintf("%+v", book)
			return fmt.Sprintf("%+v", strings.Join(strings.Split(bookInfo[1:len(bookInfo)-1], " "), " "))
		} else if id == book.Id && book.Stock < order && book.Stock > 0 {
			fmt.Printf("Stock: %+v\n", book.Stock)
			return "Not enough books in stock. Please order less"
		} else if id == book.Id && book.Stock >= order && book.Stock == 0 {
			return "Book is not available for sale. Please try later"
		}
	}
	return "Given id is not valid"
}

// delete function removes given book from the available items list. It sets StockCode to 0, which means book is not available. Returns commands result as string.
func delete(data Books, bookId string) string {
	id, err := strconv.Atoi(bookId)
	if err != nil {
		return ("String conversation error!")
	}
	for i, book := range data.Books {
		if id == book.Id {
			book.setStock(0)
			data.Books = append(data.Books[:i], data.Books[i+1:]...)
			data.Books = append(data.Books, book)
			newData, err := json.Marshal(data)
			if err != nil {
				return "Error: Couldn't write to file"
			}
			err = ioutil.WriteFile("books.json", newData, 0644)
			if err != nil {
				return "Error: Couldn't write to file"
			}
			return fmt.Sprintf("%s has been remevod from the list", book.Title)
		}
	}
	return "Given id is not valid"
}

func get(data Books, bookId string) string {
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
		if strings.Contains(strings.ToLower(book.Title), strings.ToLower(searchTerm)) && book.Stock > 0 {
			foundBooks = append(foundBooks, book.Title)
		}
	}
	if len(foundBooks) > 0 {
		return strings.Join(foundBooks, "\n")
	}
	return "No books found"
}

// list function lists books in the Books struct
func list(data Books) {
	for i := 0; i < len(data.Books); i++ {
		if data.Books[i].Stock > 0 {
			fmt.Printf("Book ID: %d | Title: %s | Author: %s | Stock: %d \n", data.Books[i].Id, data.Books[i].Title, data.Books[i].Author, data.Books[i].Stock)
		}
	}
}

func (b *Book) setStock(stock int) *Book {
	b.Stock = stock
	return b
}
