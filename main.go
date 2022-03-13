package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	}

	// for i := 0; i < len(data.Books); i++ {
	// 	fmt.Println("Book ID: " + strconv.Itoa(data.Books[i].Id))
	// 	fmt.Println("Book Title: " + data.Books[i].Title)
	// 	fmt.Println("Book Page: " + strconv.Itoa(data.Books[i].Page))
	// 	fmt.Println("Book Stock: " + strconv.Itoa(data.Books[i].Stock))
	// 	fmt.Println("Book Price: " + data.Books[i].Price)
	// 	fmt.Println("Book StockCode: " + strconv.Itoa(data.Books[i].StockCode))
	// 	fmt.Println("Book Isbn: " + strconv.Itoa(data.Books[i].Isbn))
	// 	fmt.Println("Book Author: " + data.Books[i].Author)
	// 	fmt.Print("\n")
	// }
}

func list(data Books) {
	for i := 0; i < len(data.Books); i++ {
		fmt.Println("Book #", i+1, "Title: "+data.Books[i].Title)
	}
}
