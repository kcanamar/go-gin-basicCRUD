package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Model
type book struct{
	ID string `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Quantity int16 `json:"quantity"`
}

// inital books
var books = []book{
	{ID: "1", Title: "One", Author: "Bob", Quantity: 1},
	{ID: "2", Title: "Two", Author: "Sue", Quantity: 2},
	{ID: "3", Title: "Three", Author: "Jaime", Quantity: 3},
}


// Helpers
func getBookById(id string) (*book, error) {
	// iterate over the books slice 
	for i, b := range books {
		// check if id matches
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

func removeBook(books []book, b *book) {
	for i := len(books) - 1; i >= 0; i-- {
        if books[i] == *b {
            copy(books[i:], books[i+1:])
            books[len(books)-1] = book{}
            books = books[:len(books)-1]
        }
    }
}

// Index Handler
func getBooks(c *gin.Context){
	// Response, data
	c.IndentedJSON(http.StatusOK, books)
}

// Create Handler
func createBook(c *gin.Context){
	// variable for newBook with type book 
	var newBook book

	// Error handling the request body
	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	// reasign the books slice with the append the request body
	books = append(books, newBook)

	// response, data
	c.IndentedJSON(http.StatusCreated, newBook)

}

// Show Handler
func showBook(c *gin.Context) {
	// request params
	id := c.Param("id")
	// fetching the book
	book, err := getBookById(id)

	// error handling
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return 
	}

	// response, data
	c.IndentedJSON(http.StatusOK, book)
}

// Update Handlers
func checkoutBook(c *gin.Context) {
	// check for query parameter
	id, ok := c.GetQuery("id")

	// error handler
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	// get the book id
	book, err := getBookById(id)

	// error handler
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return 
	}

	// error handler
	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Sorry not enough books for everyone."})
		return 
	}
	// update book quanity 
	book.Quantity -= 1

	// response, data
	c.IndentedJSON(http.StatusOK, book)
}

func checkinBook(c *gin.Context) {
	// check for query parameter
	id, ok := c.GetQuery("id")

	// error handler
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	// get the book id
	book, err := getBookById(id)

	// error handler
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return 
	}

	// update book quanity 
	book.Quantity += 1

	// response, data
	c.IndentedJSON(http.StatusOK, book)
}

// Delete handler
func destroyBook(c *gin.Context) {
	// check for query parameter
	id, ok := c.GetQuery("id")

	// error handler
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	// get the book id
	book, err := getBookById(id)

	// error handler
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return 
	}

	removeBook(books, book)

	// response , data
	c.IndentedJSON(http.StatusOK, books)
}

func main() {
	// create router
	router := gin.Default()
	// Index Route
	router.GET("/books", getBooks)
	// Show Route
	router.GET("/books/:id", showBook)
	// Create Route
	router.POST("/books", createBook)
	// Update Routes - use query parameters ex. "/checkout?id=1"
	router.PUT("/checkout", checkoutBook)
	router.PUT("/return", checkinBook)
	// Delete Route - use query parameters ex. "/checkout?id=1"
	router.DELETE("/destory", destroyBook)
	// Server listener
	router.Run()
}