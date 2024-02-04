package sqlite_teo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type SQLitePool struct {
	*sql.DB
}

//	type User struct {
//		email    string
//		password string
//	}
type GestureItem struct {
	ID      sql.NullInt64
	gesture string
	counter int
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

func (db *SQLitePool) QueryAll(w http.ResponseWriter, r *http.Request) {
	//im going to assume it's a get request
	tx, _ := db.Begin()

	stmt, err := tx.Prepare("SELECT * FROM gestures_table where gesture !=''")
	if err != nil {
		log.Fatal("Couldnt not query table")
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()

	if err != nil {
		log.Fatal("Couldnt query table")
		return
	}
	defer rows.Close()

	var gesturesItems []GestureItem

	for rows.Next() {
		var gestItem GestureItem

		err := rows.Scan(&gestItem.ID, &gestItem.gesture, &gestItem.counter)
		if err != nil {
			log.Fatal("Could not scan row")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		gesturesItems = append(gesturesItems, gestItem)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating over rows:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	map_array := make(map[string]int, 17)

	for _, item := range gesturesItems {
		map_array[item.gesture] = item.counter
	}
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map_array)
}

func (db *SQLitePool) CreateTable() {
	tx, _ := db.Begin()

	stmt, err := tx.Prepare("CREATE TABLE IF NOT Exists gestures_table(	id INTEGER PRIMARY KEY AUTOINCREMENT,gesture TEXT NOT NULL,counter INTEGER DEFAULT 0);")

	if err != nil {
		log.Fatal("Couldnt create table")
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec()

	if err != nil {
		log.Fatal("Couldnt create table")
		return
	}

	tx.Commit()
}

func table() []GestureItem {

	gestureList := make([]GestureItem, 18)

	str := "H-0.mp4 H-2.mp4 H-4.mp4 H-6.mp4 H-8.mp4 H-DecreaseFanSpeed.mp4 H-FanOn.mp4 H-LightOff.mp4 H-SetThermo.mp4 H-1.mp4 H-3.mp4 H-5.mp4 H-7.mp4 H-9.mp4 H-FanOff.mp4 H-IncreaseFanSpeed.mp4 H-LightOn.mp4"
	list := strings.Split(str, " ")
	fmt.Println(len(list))
	for idx, item := range list {
		gestureList[idx].gesture = item
		gestureList[idx].counter = 0
	}

	return gestureList
}

func (db *SQLitePool) ReturnCounterItems() {

	list := table()
	for _, item := range list {

		tx, _ := db.Begin()

		stmt, err := tx.Prepare("INSERT INTO gestures_table(gesture,counter) values (?,?)")

		if err != nil {
			log.Fatal("Couldnt insert values ")
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(item.gesture, item.counter)

		if err != nil {
			log.Fatal("Couldnt create table")
			return
		}

		tx.Commit()
	}
	// 	stmt, err := tx.Prepare("insert into users(email,password) values (?,?)")

}

////ignore I thought we needed users/ keeping it just in case

// func (db *SQLitePool) Create_user(w http.ResponseWriter, r *http.Request) {

// 	if !PostRequest_Checker(w, r) {
// 		return
// 	}
// 	fmt.Println("We received a request ")

// 	err := r.ParseMultipartForm(10 << 20)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("error parsing form data: %v", err), http.StatusBadRequest)
// 		return
// 	}

// 	email := r.FormValue("email")
// 	password := r.FormValue("password")

// 	fmt.Printf("email is %s a", email)

// 	encrypted_password, err := encryptors.EncryptPassword(password)

// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("error encrypting password: %v", err), http.StatusBadRequest)
// 		return
// 	}

// 	tx, err := db.Begin()

// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("error connecting to db: %v", err), http.StatusBadRequest)
// 		return
// 	}
// 	defer tx.Rollback()

// 	stmt, err := tx.Prepare("insert into users(email,password) values (?,?)")

// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("error inserting to db: %v", err), http.StatusBadRequest)
// 		return
// 	}

// 	defer stmt.Close()

// 	_, err = stmt.Exec(email, encrypted_password)

// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("error bad data passed: %v", err), http.StatusBadRequest)
// 		return
// 	}

// 	// so annoying
// 	err = tx.Commit()
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("try password again: %v", err), http.StatusBadRequest)
// 		return
// 	}

// 	w.WriteHeader(200)
// 	fmt.Fprintf(w, "Welcome to the app %s \n", email)
// 	log.Println("User has been created")
// }

// func (db *SQLitePool) Login_User(w http.ResponseWriter, r *http.Request) {

// 	PostRequest_Checker(w, r)

// 	err := r.ParseMultipartForm(10 << 20)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("error parsing form data: %v", err), http.StatusBadRequest)
// 		return
// 	}

// 	email := r.FormValue("email")
// 	password := r.FormValue("password")
// 	log.Println("email is: ", email)
// 	tx, err := db.Begin()

// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Err initializing db: %v", err), http.StatusBadRequest)
// 		return
// 	}
// 	defer tx.Rollback()

// 	stmt, err := tx.Prepare("SELECT email, password FROM users WHERE email = (?)")

// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Getting user data: %v", err), http.StatusBadRequest)
// 		return
// 	}

// 	defer stmt.Close()
// 	db_user := &User{}

// 	err = stmt.QueryRow(email).Scan(&db_user.email, &db_user.password)

// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		http.Error(w, fmt.Sprintf("Err Could not find email: %v", err), http.StatusBadRequest)
// 		return
// 	}

// 	err = tx.Commit()

// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Username or password incorrect: %v", err), http.StatusBadRequest)
// 		return
// 	}
// 	canLogin := encryptors.DecryptPassword(db_user.password, password)

// 	// not safe for practice but can front end can allow to go in to app
// 	if !canLogin {
// 		w.WriteHeader(http.StatusBadRequest)
// 		fmt.Fprintln(w, "Username or password incorrect")
// 		return
// 	}
// 	w.WriteHeader(http.StatusAccepted)
// 	fmt.Fprintln(w, "user can enter app!")

// }
