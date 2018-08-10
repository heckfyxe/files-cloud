package main

import (
	"database/sql"
	"files-cloud/handlers"
	"files-cloud/models"
	"files-cloud/utils/transbytes"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	temp, err := template.New("index.html").ParseFiles("templates/index.html")
	if err != nil {
		panic(err)
	}

	database, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	handler := &handlers.Handler{Database: database}

	http.Handle("/css-styles/",
		http.StripPrefix("/css-styles/",
			http.FileServer(http.Dir("css-styles"))))

	http.HandleFunc("/files/", handler.Load)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.Upload(w, r)
		}

		rows, err := database.Query("SELECT _id, filename, size FROM files")
		if err != nil {
			log.Println(err)
		}

		var files models.Files
		for rows.Next() {
			var file models.File
			var size float64

			rows.Scan(&file.Id, &file.Name, &size)

			file.Size = transbytes.SizeToString(size)

			files = append(files, &file)
		}
		rows.Close()

		temp.Execute(w, files)
	})

	port := os.Getenv("PORT")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
