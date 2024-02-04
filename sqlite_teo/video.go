package sqlite_teo

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Video struct {
	last_name  string
	video_type string
	path       string
}

func get_text_field(field string, w http.ResponseWriter, reader *multipart.Reader) (string, error) {
	p, err := reader.NextPart()
	// one more field to parse, EOF is considered as failure here
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	log.Printf("Text Field is %s \n", p.FormName())

	if p.FormName() != field {
		log.Printf("Could not find expected %v", err)
		http.Error(w, fmt.Sprintf("%s is expected", field), http.StatusBadRequest)
		return "", errors.New("last name Field Expected")
	}
	var buf bytes.Buffer
	_, err = io.CopyN(&buf, p, 512) // Copy at most 512 bytes
	if err != nil && err != io.EOF {
		http.Error(w, "could not read the last name", http.StatusInternalServerError)
		return "", err
	}
	last_name := buf.String()

	if !utf8.ValidString(last_name) {
		log.Printf("Invalid UTF-8 string: %s", last_name)
		// Handle the error or return an error if necessary
		return "", errors.New("invalid UTF-8 string")
	}
	log.Printf("Printing a  %s", last_name)
	//
	return last_name, nil
}

func (pg *SQLitePool) Upload_video(w http.ResponseWriter, r *http.Request) {
	video_uploader := &Video{}
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20+1024)
	reader, err := r.MultipartReader()
	if err != nil {
		log.Printf("Did not send a multipart%v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// parse text field
	last_name, err := get_text_field("last_name", w, reader)
	video_uploader.last_name = last_name
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	video_uploader.video_type, err = get_text_field("video_type", w, reader)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// we have to still parse the type

	// parse file field
	p, err := reader.NextPart()
	if err != nil && err != io.EOF {
		log.Printf("Could not read part properly %v", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if p.FormName() != "file_field" {
		log.Printf("File_field was expected %v", err)

		http.Error(w, "file_field is expected", http.StatusBadRequest)
		return
	}
	buf := bufio.NewReader(p)

	sniff, _ := buf.Peek(512)
	contentType := http.DetectContentType(sniff)
	if contentType != "video/mp4" {
		log.Printf("the content type is not accepted %s %v", contentType, err)

		http.Error(w, fmt.Sprintf("%s is not accepted type", contentType), http.StatusBadRequest)
		return
	}

	savedir := strings.Clone(video_uploader.last_name)

	//we are going to query our counter
	trial_number := strconv.Itoa(pg.QueryCounter(video_uploader.video_type))
	savefilename := video_uploader.video_type + "_" + "TRIAL_" + trial_number + "_" + savedir + ".mp4"

	check_utf8(savedir)

	savedir = strings.TrimSpace(savedir)
	saveDirectory, err := CreateDir(&savedir)

	if err != nil {
		http.Error(w, "Could not create dir", http.StatusInternalServerError)
		log.Printf("Something is really wrong")
		return
	}

	log.Printf("Created directory Succesfully at %s \n", saveDirectory)
	log.Printf("File name is going to be %s", savefilename)

	f, err := os.Create(filepath.Join(saveDirectory, savefilename))
	if err != nil {
		http.Error(w, "Could not create dir, check file", http.StatusInternalServerError)
		log.Printf("Could not create file %v", err)
		return
	}
	defer f.Close()

	var maxSize int64 = 20 << 20 // 10 mb is the max
	lmt := io.MultiReader(buf, io.LimitReader(p, maxSize-511))
	written, err := io.Copy(f, lmt)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Created file succesfuly at %s not writting to it", savefilename)

	if written > maxSize {
		log.Printf("File size too big")
		os.Remove(f.Name())
		http.Error(w, "file size over limit Max is 10 mb", http.StatusBadRequest)
		return
	}

	log.Printf("Success added file to %s \n", filepath.Join(saveDirectory, savefilename))
	w.WriteHeader(200)
	fmt.Fprintf(w, "Hello, Succesful upload %s\n", filepath.Join(saveDirectory, savefilename))

}
func check_utf8(str string) {
	c, _ := utf8.DecodeRuneInString(str)
	if c != '.' && c != ',' && c != '?' && c != '“' && c != '”' {
		fmt.Println("Ok")
	} else {
		fmt.Println("Not ok:", c)
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
