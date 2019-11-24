package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Starting mosaic server ...")
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

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
	TotalGoukau      int
	TotalSonota      int
	Result           string
}

//  Handler function for fan-out and fan-in
func mosaic(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("uploaded")
	r.ParseForm()
	major := r.FormValue("goukaku")
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	record, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

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
	TotalGoukau := 0
	TotalSonota := 0

	for _, item := range record {
		TaniSu, _ := strconv.Atoi(item[6])
		if item[7] == "優" || item[7] == "良" || item[7] == "可" || item[7] == "優上" || item[7] == "合格" {
			for _, post := range POSTS {
				if item[2] == post.Subjects {
					if post.Category == "専門科目1" {
						TotalSenmonOne = TotalSenmonOne + TaniSu
					} else if post.Category == "専門科目2" {
						TotalSenmonTwo = TotalSenmonTwo + TaniSu
					} else if post.Category == "専門科目3" {
						TotalSenmonThree = TotalSenmonThree + TaniSu
					} else if post.Category == "専門科目4" {
						TotalSenmonFour = TotalSenmonFour + TaniSu
					} else {
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
				TotalGoukau = TotalGoukau + TaniSu
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
	TotalTani = TotalYujyo + TotalYu + TotalRyou + TotalKa + TotalGoukau
	var result string
	if TotalTani >= 80 {
		if major == "経済" {
			if TotalSenmonOne >= 20 {
				if TotalSenmonTwo >= 18 {
					result = "卒業見込み"
				} else {
					result = "専門科目2不足"
				}
			} else {
				if TotalSenmonTwo < 18 {
					result = "専門科目1,2不足"
				} else {
					result = "専門科目1不足"
				}
			}
		} else if major == "経営" {
			if TotalSenmonOne >= 20 {
				if TotalSenmonThree >= 18 {
					result = "卒業見込み"
				} else {
					result = "専門科目3不足"
				}
			} else {
				if TotalSenmonThree < 18 {
					result = "専門科目1,3不足"
				} else {
					result = "専門科目1不足"
				}
			}
		} else if major == "金融" {
			if TotalSenmonOne >= 20 {
				if TotalSenmonFour >= 18 {
					result = "卒業見込み"
				} else {
					result = "専門科目4不足"
				}
			} else {
				if TotalSenmonFour < 18 {
					result = "専門科目1,4不足"
				} else {
					result = "専門科目1不足"
				}
			}
		}
	} else {
		result = "総単位数不足"
	}

	ResultLists := results{TotalTani: TotalTani, TotalSenmonOne: TotalSenmonOne, TotalSenmonTwo: TotalSenmonTwo, TotalSenmonThree: TotalSenmonThree, TotalSenmonFour: TotalSenmonFour,
		TotalSentaku: TotalSentaku, TotalYujyo: TotalYujyo, TotalYu: TotalYu, TotalRyou: TotalRyou, TotalKa: TotalKa,
		TotalFuka: TotalFuka, TotalMijyuken: TotalMijyuken,
		TotalGoukau: TotalGoukau, TotalSonota: TotalSonota, Result: result}

	t, _ := template.ParseFiles("results.html")
	t.Execute(w, ResultLists)
}
