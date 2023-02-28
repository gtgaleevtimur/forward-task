// Package repository - реализует логику БД.
package repository

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocarina/gocsv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Model - структура миграции БД.
type Model struct {
	INN         string `csv:"G1"`
	Name        string `csv:"G2"`
	Surname     string `csv:"G3"`
	LastSurName string `csv:"G4"`
	Count       string `csv:"G5"`
}

type Repository struct {
	db *gorm.DB
}

// NewRepository - создает новое хранилище.
func NewRepository() *Repository {
	r := &Repository{}
	return r
}

// UpdateRegistry - обновляет базу данных, читая из csv файла скаченного в ../pkg
func (r *Repository) UpdateRegistry() (err error) {
	log.Println("Update registry started.")
	wd, err := newPath()
	if err != nil {
		return
	}
	file, err := os.Open(wd + "/pkg/temp.csv")
	if err != nil {
		return
	}
	defer file.Close()
	var entries []Model
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = ';'
		return r
	})
	err = gocsv.UnmarshalFile(file, &entries)
	if err != nil {
		return
	}
	r.db, err = gorm.Open(sqlite.Open(wd + "/pkg/gorm.db"))
	if err != nil {
		return
	}
	if ok := r.db.Migrator().HasTable(&Model{}); ok {
		err = r.db.Migrator().DropTable(&Model{})
		if err != nil {
			return
		}
	}
	err = r.db.Migrator().CreateTable(&Model{})
	if err != nil {
		return
	}
	result := r.db.CreateInBatches(entries, 100)
	if result.Error != nil {
		return result.Error
	}
	log.Println("Update registry ended.")
	return nil
}

// Get - возвращает количество зарегистрированных ЮЛ по ИНН.
func (r *Repository) Get(inn string) (count string, err error) {
	wd, err := newPath()
	if err != nil {
		return "", err
	}
	r.db, err = gorm.Open(sqlite.Open(wd + "/pkg/gorm.db"))
	if err != nil {
		return
	}
	r.db.Table("models").Select("count").Where("inn = ?", inn).Scan(&count)
	return
}

// init - скачивает csv файл с сайта налоговой службы.
func init() {
	wd, err := newPath()
	if err != nil {
		log.Fatal(err)
	}
	commonResponse, err := http.Get("https://www.nalog.gov.ru/opendata/7707329152-massfounders/")
	if err != nil {
		log.Fatal(err)
	}
	defer commonResponse.Body.Close()
	if commonResponse.StatusCode != 200 {
		fmt.Println("Received non 200 response code")
	}
	doc, err := goquery.NewDocumentFromReader(commonResponse.Body)
	if err != nil {
		log.Fatal(err)
	}

	response, err := http.Get(doc.Find(`.border_table > tbody:nth-child(1) > tr:nth-child(9) > td:nth-child(3) > a:nth-child(1)`).Text())
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatal(err)
	}
	file, err := os.OpenFile(wd+"/pkg/temp.csv", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func newPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for !strings.HasSuffix(wd, "forward-task") {
		wd = filepath.Dir(wd)
	}
	return wd, nil
}
