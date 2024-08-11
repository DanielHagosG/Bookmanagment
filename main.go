package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		books, err := GetBooks(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, books)
	})

	http.HandleFunc("/add-cart", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form data", http.StatusBadRequest)
				return
			}

			title := r.FormValue("title")
			author := r.FormValue("author")
			publishedDateStr := r.FormValue("published_date")
			publishedDate, err := time.Parse("2006-01-02", publishedDateStr)
			if err != nil {
				http.Error(w, "Invalid date format", http.StatusBadRequest)
				return
			}

			book := Book{
				Title:         title,
				Author:        author,
				PublishedDate: publishedDate,
			}

			err = CreateBook(db, book)
			if err != nil {
				http.Error(w, "Error adding book to the database", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	})

	// Update book handler
	http.HandleFunc("/update-book", func(w http.ResponseWriter, r *http.Request) {
		bookID := r.URL.Query().Get("book_id")
		id, err := strconv.Atoi(bookID)
		if err != nil {
			http.Error(w, "Invalid book ID", http.StatusBadRequest)
			return
		}

		book, err := GetBook(db, id)
		if err != nil {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}

		tmpl, err := template.ParseFiles("templates/update.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, book)
	})

	// Perform update handler
	http.HandleFunc("/perform-update", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		bookID := r.FormValue("book_id")
		id, err := strconv.Atoi(bookID)
		if err != nil {
			http.Error(w, "Invalid book ID", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		author := r.FormValue("author")
		publishedDateStr := r.FormValue("published_date")
		publishedDate, err := time.Parse("2006-01-02", publishedDateStr)
		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}

		book := Book{
			ID:            id,
			Title:         title,
			Author:        author,
			PublishedDate: publishedDate,
		}

		err = UpdateBook(db, book)
		if err != nil {
			http.Error(w, "Error updating book in the database", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Delete book handler
	http.HandleFunc("/delete-book", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form data", http.StatusBadRequest)
				return
			}

			bookID := r.FormValue("book_id")
			id, err := strconv.Atoi(bookID)
			if err != nil {
				http.Error(w, "Invalid book ID", http.StatusBadRequest)
				return
			}

			err = DeleteBook(db, id)
			if err != nil {
				http.Error(w, "Error deleting book from the database", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	})

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
