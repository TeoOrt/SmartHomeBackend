package main

import (
	"SmartHomeBackend/main/sqlite_teo"
	"database/sql"
	"log"
	"net/http"
)

func main() {
	/*
		To make this easier this code servers the expert videos
		You can upload the video
		and you can retrieve the amount of times the video has been posted
	*/

	log.Println("Starting to listen on http:localhost:8080")
	db, err := sql.Open("sqlite3", "videos.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlite := &sqlite_teo.SQLitePool{DB: db} // gets client driver for SQLITE
	sqlite.CreateTable()                     // creates sql table
	// sqlite.ReturnCounterItems() //only when creating new table
	// sqlite.ReturnCounterItems() already filled the table out

	handler := http.StripPrefix("/get_expert/video/", http.FileServer(http.Dir("video_storage/expert_videos")))
	http.HandleFunc("/upload_video", sqlite.Upload_video) //post request for user to upload vid
	http.HandleFunc("/get_counter", sqlite.QueryAll)      // gets the counter values to see the rial
	http.Handle("/get_expert/video/", handler)            // gets the video from the video_storage/expert_videos directory
	// users can get expert video

	log.Fatal(http.ListenAndServe(":8080", nil)) // listens in port 8080

}

/*

  {id:0,emailL: mateos@gmail.com, password:HelloWorld,video_list: [
    {
      type:"OpenCurtain"
      path:"/video/mateo@gmail.com"
    }
  ]}


*/
