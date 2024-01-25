package sqlite_teo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type Video struct {
	email      string
	video_type string
	path       string
}

func get_text_field(field string, w http.ResponseWriter, reader *multipart.Reader) (string, error) {
	text := make([]byte, 512)
	p, err := reader.NextPart()
	// one more field to parse, EOF is considered as failure here
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	log.Printf("Text Field is %s \n", p.FormName())

	if p.FormName() != field {
		http.Error(w, fmt.Sprintf("%s is expected", field), http.StatusBadRequest)
		return "", errors.New("email Field Expected")
	}

	_, err = p.Read(text)
	if err != nil && err != io.EOF {
		http.Error(w, "could not read the email", http.StatusInternalServerError)
		return "", err
	}
	email := string(text)
	log.Printf("Printing a  %s", email)
	//
	return email, nil
}

func (pg *SQLitePool) Upload_video(w http.ResponseWriter, r *http.Request) {
	video_uploader := &Video{}
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20+1024)
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// parse text field
	video_uploader.email, err = get_text_field("email", w, reader)
	if err != nil {
		log.Fatal(err)
		return
	}

	video_uploader.video_type, err = get_text_field("type", w, reader)
	if err != nil {
		log.Fatal(err)
		return
	}
	// we have to still parse the type

	// parse file field
	p, err := reader.NextPart()
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if p.FormName() != "file_field" {
		http.Error(w, "file_field is expected", http.StatusBadRequest)
		return
	}
	buf := bufio.NewReader(p)

	sniff, _ := buf.Peek(512)
	contentType := http.DetectContentType(sniff)
	if contentType != "application/octet-stream" {
		http.Error(w, fmt.Sprintf("%s is not accepted type", contentType), http.StatusBadRequest)
		return
	}

	// f, err := os.CreateTemp("", "hello_23")

	saveDirectory := "/home/teoortega/AndroidStudioProjects/SmartHomeBackend/" + video_uploader.email
	os.MkdirAll(saveDirectory, os.ModePerm)
	f, err := os.Create(filepath.Join(saveDirectory, "Hi_.mp4"))
	if err != nil {
		http.Error(w, "Could not create dir, check file", http.StatusInternalServerError)
		log.Printf("Could not create temp dir %v", err)
		return
	}
	defer f.Close()
	var maxSize int64 = 10 << 20 // 10 mb is the max
	lmt := io.MultiReader(buf, io.LimitReader(p, maxSize-511))
	written, err := io.Copy(f, lmt)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if written > maxSize {
		os.Remove(f.Name())
		http.Error(w, "file size over limit Max is 10 mb", http.StatusBadRequest)
		return
	}

}

// This basically checks just in case a video has been deleted
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func NewVideo(email *string, video *string, path *string) (*Video, error) {

	real, err := exists(*path)

	if err != nil {
		return nil, err
	}

	if !real {
		return nil, errors.New("path doesnt exist")
	}

	vid := &Video{*email, *video, *path}

	return vid, nil
}

// func (vid *Video) Serve_Video() {
// 	Serve
// }
