package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

//Post Structure to store db
type Post struct {
	Subjects string
	Tani     string
	Category string
}

//POSTS Structure to store db
var POSTS []Post

//MakeDB for calculating tani
func MakeDB() []Post {
	fmt.Println("Start populating tiles db ...")
	// reading a CSV file
	file, err := os.Open("posts.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	record, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	var posts []Post
	for _, item := range record {
		post := Post{Subjects: item[0], Tani: item[1], Category: item[2]}
		posts = append(posts, post)
	}
	return posts
}
