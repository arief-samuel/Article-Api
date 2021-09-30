package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qkgo/yin"
)

func main() {
	db, _ := sql.Open("sqlite3", "C:\\Users\\ephra-samuel\\Documents\\c#\\Linode-Ubuntu\\article.db")
	newArticleDb := NewArticleDb(db)

	r := chi.NewRouter()
	r.Use(yin.SimpleLogger)

	r.Get("/api/articles", func(w http.ResponseWriter, r *http.Request) {
		res, _ := yin.Event(w, r)
		articles := newArticleDb.Get()
		res.SendJSON(articles)
	})

	r.Post("/api/articles", func(w http.ResponseWriter, r *http.Request) {
		res, req := yin.Event(w, r)
		body := map[string]string{}
		req.BindBody(&body)
		article := Article{
			Title:      body["Title"],
			Category:   body["Category"],
			Url:        body["Url"],
			Created_At: time.Now().Format(time.RFC3339),
		}
		newArticleDb.Add(article)
		res.SendStatus(204)
	})

	http.ListenAndServe(":3000", r)
}

func (articleDb *ArticleDb) Add(article Article) {
	stmt, _ := articleDb.DB.Prepare(`
	INSERT INTO Articles
(Title, Category, Url, Created_At)
VALUES(?, ?, ?, ?);
	`)
	stmt.Exec(article.Title,
		article.Category,
		article.Url,
		article.Created_At)
}

func (articleDb *ArticleDb) Get() []Article {
	articles := []Article{}
	rows, _ := articleDb.DB.Query(`
	SELECT Id, Title, Category, Url, Created_At
	FROM Articles;
	`)
	var id int
	var title, category, url, created_At string

	for rows.Next() {
		rows.Scan(&id, &title, &category, &url, &created_At)
		article := Article{
			Id:         id,
			Title:      title,
			Category:   category,
			Url:        url,
			Created_At: created_At,
		}
		articles = append(articles, article)
	}
	return articles

}

func NewArticleDb(db *sql.DB) *ArticleDb {
	return &ArticleDb{
		DB: db,
	}
}

type ArticleDb struct {
	DB *sql.DB
}

type Article struct {
	Id         int
	Title      string
	Category   string
	Url        string
	Created_At string
}
