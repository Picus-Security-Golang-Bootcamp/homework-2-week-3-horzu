package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Books represents a slice of Book
type Books struct {
	Books []Book `json:"books"`
}

// A Book represents an example of Book.
type Book struct {
	Id        int     `json:"id"`
	Title     string  `json:"title"`
	Page      int     `json:"page"`
	Stock     int     `json:"stock"`
	Price     string  `json:"price"`
	StockCode string  `json:"stockCode"`
	ISBN      string  `json:"ISBN"`
	Author    Authors `json:"author"`
}

type Authors struct {
	Id   int
	Name string
}

func (b *Book) setStock(stock int) *Book {
	b.Stock = stock
	return b
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

	listPtr := flag.Bool("list", false, "no argument needed.")
	searchPtr := flag.String("search", "", "a string value to search in the books list")
	getPtr := flag.Int("get", 0, "an int value to get a book from its id number")
	deletePtr := flag.Int("delete", 0, "an int value of book id to set stocks of the given book as 0")
	buyPtr := flag.Int("buy", 0, "an int value to buy a book from its id number, following by an int value as order quantity.")

	var usage = `Usage of C:\Users\merts\AppData\Local\Temp\go-build849867169\b001\exe\main.exe:
	Options:
	-get Number of workers to run concurrently. Default is 10.
	-search  Timeout for each request in seconds. Default is 30.
	-get  Timeout for each request in seconds. Default is 30.
	-delete  Timeout for each request in seconds. Default is 30.
	-buy  Timeout for each request in seconds. Default is 30.
`
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage))
	}
	flag.Parse()

	switch  {
	case *listPtr :
		list(data)
	case len(*searchPtr) > 0:
		fmt.Println(search(data, *searchPtr))
	case *getPtr > 0:
		fmt.Println(get(data, *getPtr))
	case *deletePtr>0:
		fmt.Println(delete(data, *deletePtr))
	case *buyPtr>0:
		quantity, err := strconv.Atoi(flag.Args()[0])
		if err != nil{
			fmt.Println("String conversation error!")
		}
		fmt.Println(buy(data, *buyPtr, quantity))
	}
}

// buy function reduces the given book's stock quantity as given order quantity deducted from stock quantity
func buy(data Books, bookId int, quantity int) string {
		for i, book := range data.Books {
		if bookId == book.Id && book.Stock > 0 && book.Stock >= quantity {
			newQuantity := book.Stock - quantity
			data.Books[i].setStock(newQuantity)
			newData, err := json.Marshal(data)
			if err != nil {
				return "json converting error"
			}
			err = ioutil.WriteFile("books.json", newData, 0644)
			if err != nil {
				return "Error: Couldn't write to file"
			}
			bookInfo := fmt.Sprintf("%+v", data.Books[i])
			return fmt.Sprintf("%+v", strings.Join(strings.Split(bookInfo[1:len(bookInfo)-1], " "), " "))
		} else if bookId == book.Id && book.Stock < quantity && book.Stock > 0 {
			fmt.Printf("Stock: %+v\n", book.Stock)
			return "Not enough books in stock."
		} else if bookId == book.Id && book.Stock == 0 {
			return "Book is not available for sale. Please try later"
		}
	}
	return "Given id is not valid"
}

// delete function removes given book from the available items list. It sets StockCode to 0, which means book is not available. Returns commands result as string.
func delete(data Books, bookId int) string {
	for i, book := range data.Books {
		if bookId == book.Id {
			if book.Stock == 0 {
				return "Book is already unavailable"
			}
			data.Books[i].setStock(0)
			newData, err := json.Marshal(data)
			if err != nil {
				return "Error: Couldn't convert data to the json file"
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

// get function returns information of the book of given id
func get(data Books, bookId int) string {
	for _, book := range data.Books {
		if bookId == book.Id {
			if book.Stock == 0{
				fmt.Println("Entered book is in the list but it is not on stock.")
			}
			bookInfo := fmt.Sprintf("Book ID: %d | Title: %s | Author: %s | Stock: %d | Price: %s", book.Id, book.Title, book.Author.Name, book.Stock, book.Price)
			return bookInfo
		}
	}
	return "Book with the given id is not available"
}

// search function searches given string in the books and returns matched books
func search(data Books, searchTerm string) string {
	if searchTerm == "" {
		return "No search parameter entered"
	}
	var foundBooks []string
	for _, book := range data.Books {
		if strings.Contains(strings.ToLower(book.Title), strings.ToLower(searchTerm)) && book.Stock > 0 {
			bookInfo := fmt.Sprintf("Book ID: %d | Title: %s | Author: %s | Stock: %d | Price: %s", book.Id, book.Title, book.Author.Name, book.Stock, book.Price)
			foundBooks = append(foundBooks, bookInfo)
		}
	}
	if len(foundBooks) > 0 {
		return strings.Join(foundBooks, "\n")
	}
	return "No books found"
}

// list function lists books in the list
func list(data Books) {
	for i := 0; i < len(data.Books); i++ {
		if data.Books[i].Stock > 0 {
			bookInfo := fmt.Sprintf("Book ID: %d | Title: %s | Author: %s | Stock: %d | Price: %s \n", data.Books[i].Id, data.Books[i].Title, data.Books[i].Author.Name, data.Books[i].Stock, data.Books[i].Price)
			fmt.Printf(bookInfo)
		}
	}
}
