package sqlite_teo

import (
	"SmartHomeBackend/main/encryptors"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type SQLitePool struct {
	*sql.DB
}
type User struct {
	email    string
	password string
}

func PostRequest_Checker(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid Request")
		return false
	}
	err := fill_Form(w, r)

	if err != nil {
		return false
	}
	return true
}

func fill_Form(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return err
	}
	return nil
}

func (db *SQLitePool) Create_user(w http.ResponseWriter, r *http.Request) {

	if !PostRequest_Checker(w, r) {
		return
	}
	fmt.Println("We received a request ")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing form data: %v", err), http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	fmt.Printf("email is %s a", email)

	encrypted_password, err := encryptors.EncryptPassword(password)

	if err != nil {
		http.Error(w, fmt.Sprintf("error encrypting password: %v", err), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()

	if err != nil {
		http.Error(w, fmt.Sprintf("error connecting to db: %v", err), http.StatusBadRequest)
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("insert into users(email,password) values (?,?)")

	if err != nil {
		http.Error(w, fmt.Sprintf("error inserting to db: %v", err), http.StatusBadRequest)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(email, encrypted_password)

	if err != nil {
		http.Error(w, fmt.Sprintf("error bad data passed: %v", err), http.StatusBadRequest)
		return
	}

	// so annoying
	err = tx.Commit()
	if err != nil {
		http.Error(w, fmt.Sprintf("try password again: %v", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "Welcome to the app %s \n", email)
	log.Println("User has been created")
}

func (db *SQLitePool) Login_User(w http.ResponseWriter, r *http.Request) {

	PostRequest_Checker(w, r)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing form data: %v", err), http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	log.Println("email is: ", email)
	tx, err := db.Begin()

	if err != nil {
		http.Error(w, fmt.Sprintf("Err initializing db: %v", err), http.StatusBadRequest)
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("SELECT email, password FROM users WHERE email = (?)")

	if err != nil {
		http.Error(w, fmt.Sprintf("Getting user data: %v", err), http.StatusBadRequest)
		return
	}

	defer stmt.Close()
	db_user := &User{}

	err = stmt.QueryRow(email).Scan(&db_user.email, &db_user.password)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, fmt.Sprintf("Err Could not find email: %v", err), http.StatusBadRequest)
		return
	}

	err = tx.Commit()

	if err != nil {
		http.Error(w, fmt.Sprintf("Username or password incorrect: %v", err), http.StatusBadRequest)
		return
	}
	canLogin := encryptors.DecryptPassword(db_user.password, password)

	// not safe for practice but can front end can allow to go in to app
	if !canLogin {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Username or password incorrect")
		return
	}
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "user can enter app!")

}
