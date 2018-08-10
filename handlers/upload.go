package handlers

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
)

func (h Handler) Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("upload_file")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	var id string
	h.Database.QueryRow(
		"INSERT INTO files(filename, size) VALUES ($1, $2) RETURNING _id",
		handler.Filename,
		handler.Size,
	).Scan(&id)

	var path bytes.Buffer
	path.WriteString("cloud-files/")
	path.WriteString(id)

	localFile, err := os.Create(path.String())
	if err != nil {
		log.Println(err)
		return
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, file)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
