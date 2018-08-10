package handlers

import (
	"bytes"
	"database/sql"
	"files-cloud/models"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func (h Handler) Load(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	fileID := path[len(path)-1]

	// if fileID isn't numeric
	if _, err := strconv.ParseInt(fileID, 10, 0); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var file models.File
	row := h.Database.QueryRow("SELECT * FROM files WHERE _id=$1", fileID)
	err := row.Scan(&file.Id, &file.Name, &file.Size)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		log.Println(err)
		return
	}

	var filePath bytes.Buffer
	filePath.WriteString("cloud-files/")
	filePath.WriteString(fileID)
	f, err := os.Open(filePath.String())
	if err != nil {
		if err == os.ErrNotExist {
			http.NotFound(w, r)
			return
		}
		log.Println(err)
		return
	}
	defer f.Close()

	contentType, err := GetFileContentType(f)
	if err != nil {
		contentType = "multipart/form-data"
	}
	f.Seek(0, 0)

	w.Header().Set("Content-Disposition", "attachment; filename="+file.Name)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", file.Size)

	io.Copy(w, f)
}

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	log.Println("load: content-type " + contentType)

	return contentType, nil
}
