package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
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

func main() {
	fmt.Println("Starting mosaic server ...")
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))
	images := http.FileServer(http.Dir("images"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	mux.Handle("/images/", http.StripPrefix("/images/", images))

	mux.HandleFunc("/", upload)
	mux.HandleFunc("/mosaic", mosaic)

	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: mux,
	}
	POSTS = MakeDB()

	fmt.Println("Mosaic server started.")
	server.ListenAndServe()
}

func upload(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("upload.html")
	t.Execute(w, nil)
}

type results struct {
	TotalTani        int
	TotalSenmonOne   int
	TotalSenmonTwo   int
	TotalSenmonThree int
	TotalSenmonFour  int
	TotalSentaku     int
	TotalYujyo       int
	TotalYu          int
	TotalRyou        int
	TotalKa          int
	TotalFuka        int
	TotalMijyuken    int
	TotalGoukaku     int
	TotalSonota      int
	Result           string
}

//  Handler function for fan-out and fan-in
func mosaic(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	record, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	major := r.FormValue("goukaku")

	TotalSenmonOne := 0
	TotalSenmonTwo := 0
	TotalSenmonThree := 0
	TotalSenmonFour := 0
	TotalSentaku := 0
	TotalTani := 0
	TotalYujyo := 0
	TotalYu := 0
	TotalRyou := 0
	TotalKa := 0
	TotalFuka := 0
	TotalMijyuken := 0
	TotalGoukaku := 0
	TotalSonota := 0

	for _, item := range record {
		TaniSu, _ := strconv.Atoi(item[6])
		if item[7] == "優" || item[7] == "良" || item[7] == "可" || item[7] == "優上" || item[7] == "合格" {
			count := 0
			for i, post := range POSTS {
				if item[2] == post.Subjects {
					if post.Category == "専門科目1" {
						TotalSenmonOne = TotalSenmonOne + TaniSu
						count = count + 1
					} else if post.Category == "専門科目2" {
						TotalSenmonTwo = TotalSenmonTwo + TaniSu
						count = count + 1
					} else if post.Category == "専門科目3" {
						TotalSenmonThree = TotalSenmonThree + TaniSu
						count = count + 1
					} else if post.Category == "専門科目4" {
						TotalSenmonFour = TotalSenmonFour + TaniSu
						count = count + 1
					}
				}
				if i == len(POSTS)-1 {
					if count == 0 {
						TotalSentaku = TotalSentaku + TaniSu
					}
				}
			}
			if item[7] == "優上" {
				TotalYujyo = TotalYujyo + TaniSu
			} else if item[7] == "優" {
				TotalYu = TotalYu + TaniSu
			} else if item[7] == "良" {
				TotalRyou = TotalRyou + TaniSu
			} else if item[7] == "可" {
				TotalKa = TotalKa + TaniSu
			} else if item[7] == "合格" {
				TotalGoukaku = TotalGoukaku + TaniSu
			}
		} else {
			if item[7] == "不可" {
				TotalFuka = TotalFuka + TaniSu
			} else if item[7] == "未受験" {
				TotalMijyuken = TotalMijyuken + TaniSu
			} else if item[7] == "その他" {
				TotalSonota = TotalSonota + TaniSu
			}
		}
	}
	TotalTani = TotalYujyo + TotalYu + TotalRyou + TotalKa + TotalGoukaku
	var result string
	if TotalTani >= 80 {
		if major == "経済" {
			if TotalSenmonOne >= 20 {
				if TotalSenmonTwo >= 18 {
					result = "卒業見込み"
				} else {
					result = "現段階では卒業できません:\n　総合単位は足りています。専門科目2(" + strconv.Itoa(18-TotalSenmonOne) + "単位不足)が不足しています。"
				}
			} else {
				if TotalSenmonTwo < 18 {
					result = "現段階では卒業できません:\n　総合単位は足りています。専門科目1(" + strconv.Itoa(20-TotalSenmonOne) + "単位不足)" + ",専門科目2(" + strconv.Itoa(18-TotalSenmonTwo) + "単位不足)が不足しています。"
				} else {
					result = "現段階では卒業できません:\n　総合単位は足りています。専門科目1(" + strconv.Itoa(20-TotalSenmonOne) + "単位不足)が不足しています。"
				}
			}
		} else if major == "経営" {
			if TotalSenmonOne >= 20 {
				if TotalSenmonThree >= 18 {
					result = "卒業見込み"
				} else {
					result = "現段階では卒業できません:\n　総合単位は足りています。専門科目3(" + strconv.Itoa(18-TotalSenmonThree) + "単位不足)が不足しています。"
				}
			} else {
				if TotalSenmonThree < 18 {
					result = "現段階では卒業できません:\n　総合単位は足りています。専門科目1(" + strconv.Itoa(20-TotalSenmonOne) + "単位不足)" + ",専門科目3(" + strconv.Itoa(18-TotalSenmonThree) + "単位不足)が不足しています。"
				} else {
					result = "現段階では卒業できません:\n　総合単位は足りています。専門科目1(" + strconv.Itoa(20-TotalSenmonOne) + "単位不足)が不足しています。"
				}
			}
		} else if major == "金融" {
			if TotalSenmonOne >= 20 {
				if TotalSenmonFour >= 18 {
					result = "卒業見込み"
				} else {
					result = "現段階では卒業できません:\n　総合単位は足りています。専門科目4(" + strconv.Itoa(18-TotalSenmonFour) + "単位不足)が不足しています。"
				}
			} else {
				if TotalSenmonFour < 18 {
					result = "現段階では卒業できません:\n　総合単位は足りています。専門科目1(" + strconv.Itoa(18-TotalSenmonOne) + "単位不足)" + ",専門科目4(" + strconv.Itoa(18-TotalSenmonFour) + "単位不足)が不足しています。"
				} else {
					result = "現段階で卒業できません:\n　総合単位は足りています。専門科目1(" + strconv.Itoa(20-TotalSenmonOne) + "単位不足)が不足しています。"
				}
			}
		}
	} else {
		if major == "経済" {
			if TotalSenmonOne < 20 {
				if TotalSenmonTwo < 18 {
					result = "専門科目1(" + strconv.Itoa(20-TotalSenmonOne) + "単位不足)" + "専門科目2(" + strconv.Itoa(18-TotalSenmonTwo) + "単位不足)が不足しています。"
				} else {
					result = "専門科目1(" + strconv.Itoa(20-TotalSenmonOne) + "単位不足)が不足しています。"
				}
			} else {
				if TotalSenmonTwo < 18 {
					result = "専門科目2(" + strconv.Itoa(18-TotalSenmonTwo) + "単位不足)" + "が不足しています。"
				} else {
					result = "ただし専門科目は足りています。"
				}
			}
		} else if major == "経営" {
			if TotalSenmonOne < 20 {
				if TotalSenmonThree < 18 {
					result = "専門科目1(" + strconv.Itoa(20-TotalSenmonOne) + "単位不足)" + "専門科目3(" + strconv.Itoa(18-TotalSenmonThree) + "単位不足)が不足しています。"
				} else {
					result = "専門科目1(" + strconv.Itoa(20-TotalSenmonOne) + "単位不足)が不足しています。"
				}
			} else {
				if TotalSenmonThree < 18 {
					result = "専門科目3(" + strconv.Itoa(18-TotalSenmonThree) + "単位不足)" + "が不足しています。"
				} else {
					result = "ただし専門科目は足りています。"
				}
			}
		} else if major == "金融" {
			if TotalSenmonOne < 20 {
				if TotalSenmonFour < 18 {
					result = "専門科目1(" + strconv.Itoa(20-TotalSenmonOne) + "単位不足)" + "専門科目4(" + strconv.Itoa(18-TotalSenmonFour) + "単位不足)が不足しています。"
				} else {
					result = "専門科目1(" + strconv.Itoa(20-TotalSenmonOne) + "単位不足)が不足しています。"
				}
			} else {
				if TotalSenmonFour < 18 {
					result = "専門科目4(" + strconv.Itoa(18-TotalSenmonFour) + "単位不足)" + "が不足しています。"
				} else {
					result = "ただし専門科目は足りています。"
				}
			}
		}
		result = "現段階では卒業できません\n" + "総合単位(" + strconv.Itoa(80-TotalTani) + "単位不足)が不足しています。" + result
	}

	ResultLists := results{TotalTani: TotalTani, TotalSenmonOne: TotalSenmonOne, TotalSenmonTwo: TotalSenmonTwo, TotalSenmonThree: TotalSenmonThree, TotalSenmonFour: TotalSenmonFour,
		TotalSentaku: TotalSentaku, TotalYujyo: TotalYujyo, TotalYu: TotalYu, TotalRyou: TotalRyou, TotalKa: TotalKa,
		TotalFuka: TotalFuka, TotalMijyuken: TotalMijyuken,
		TotalGoukaku: TotalGoukaku, TotalSonota: TotalSonota, Result: result}

	t, _ := template.ParseFiles("results.html")
	t.Execute(w, ResultLists)
}
