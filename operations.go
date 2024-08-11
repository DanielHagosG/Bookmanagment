package main

import (
	"database/sql"
	"log"
	"time"
)

type Book struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	PublishedDate time.Time `json:"publishedDate"`
}

func CreateBook(db *sql.DB, book Book) error {
	_, err := db.Exec("INSERT INTO books (title, author, published_date) VALUES (?, ?, ?)", book.Title, book.Author, book.PublishedDate)
	return err
}

func GetBooks(db *sql.DB) ([]Book, error) {
	rows, err := db.Query("SELECT id, title, author, published_date FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		var publishedDateStr string
		err = rows.Scan(&b.ID, &b.Title, &b.Author, &publishedDateStr)
		if err != nil {
			log.Println(err)
			continue
		}
		b.PublishedDate, err = time.Parse("2006-01-02", publishedDateStr)
		if err != nil {
			log.Println(err)
			continue
		}
		books = append(books, b)
	}

	return books, nil
}

func GetBook(db *sql.DB, bookID int) (*Book, error) {
	row := db.QueryRow("SELECT id, title, author, published_date FROM books WHERE id = ?", bookID)

	var b Book
	var publishedDateStr string
	err := row.Scan(&b.ID, &b.Title, &b.Author, &publishedDateStr)
	if err != nil {
		return nil, err
	}
	b.PublishedDate, err = time.Parse("2006-01-02", publishedDateStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &b, nil
}

func UpdateBook(db *sql.DB, book Book) error {
	_, err := db.Exec("UPDATE books SET title = ?, author = ?, published_date = ? WHERE id = ?", book.Title, book.Author, book.PublishedDate, book.ID)
	return err
}

func DeleteBook(db *sql.DB, bookID int) error {
	_, err := db.Exec("DELETE FROM books WHERE id = ?", bookID)
	return err
}
