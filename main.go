package main

import (
	"SmartHomeBackend/main/sqlite_teo"
	"database/sql"
	"log"
	"net/http"
)

func main() {

	log.Println("Starting to listen on http:localhost:8080")
	db, err := sql.Open("sqlite3", "videos.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlite := &sqlite_teo.SQLitePool{DB: db}

	handler := http.StripPrefix("/get_expert/video/", http.FileServer(http.Dir("video_storage/expert_videos")))
	http.Handle("/get_expert/video/", handler)

	http.HandleFunc("/create_user", sqlite.Create_user)
	http.HandleFunc("/login_user", sqlite.Login_User)
	http.HandleFunc("/upload_video", sqlite.Upload_video)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

/*

  {id:0,emailL: mateos@gmail.com, password:HelloWorld,video_list: [
    {
      type:"OpenCurtain"
      path:"/video/mateo@gmail.com"
    }
  ]}


*/
