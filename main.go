package main

import ("net/http"
        "html/template"
        "database/sql"
        _ "github.com/mattn/go-sqlite3"
        "fmt"
        "time")

type Page struct {
  Title string
  Content string
}

type Article struct {
  Id int
  Title string
  Content template.HTML
  Date string
}

type Articles struct {
  Article []Article
}

type ArticlePage struct {
  Title string
  Articles []Article
}

var tpl *template.Template
var db *sql.DB

func init() {
  db, _ = sql.Open("sqlite3", "./data.db")
  statement, _ := db.Prepare(`
  CREATE TABLE IF NOT EXISTS ARTICLES (
  ID INTEGER PRIMARY KEY AUTOINCREMENT,
  TITLE TEXT, 
  CONTENT TEXT,
  DATE DATE DEFAULT (datetime('now','localtime'))
  );
  `)
  statement.Exec()
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func display(w http.ResponseWriter, tmpl string, data interface{}) {
  tpl.ExecuteTemplate(w, tmpl, data)
}

func dbArticleInsert(title string, content string) {
  statement, _ := db.Prepare("INSERT INTO ARTICLES (TITLE, CONTENT) VALUES (?, ?)")
  //tdate := time.Now()
  //fmt.Println(tdate)
  statement.Exec(title,content)
}

func dbRead() []Article{
    rows, _ := db.Query(`
    SELECT ID, TITLE, CONTENT, DATE FROM ARTICLES`)
    articles := make([]Article, 0)
    var id int
    var title string
    var content string
    var date time.Time
    var sdate string
    var html template.HTML
    for rows.Next() {
        rows.Scan(&id,&title,&content,&date)
        html = template.HTML(content)
        sdate = date.Format("02/01/2006 15:04")
        articles = append(articles, Article{id,title,html,sdate})
    }
    return articles
}

func indxHdl(w http.ResponseWriter, r *http.Request){
  articles := dbRead()
	display(w, "main", &ArticlePage{Title: "Home", Articles: articles })
}

func artclHdl(w http.ResponseWriter, r *http.Request){
  artid := r.FormValue("id")
  row := db.QueryRow("SELECT * FROM ARTICLES WHERE ID = $1", artid)
  var id int
  var title string
  var content string
  var date time.Time
  var sdate string
  var html template.HTML
  //var article Article
  row.Scan(&id,&title,&content,&date)
  html = template.HTML(content)
  sdate = date.Format("02/01/2006 15:04")
  article := Article{Id: id, Title: title, Content: html, Date: sdate}
  fmt.Println(article)
  display(w, "article-perma", article)
}

func admn_editHdl(w http.ResponseWriter, r *http.Request){
  title := r.FormValue("Title")
  content := r.FormValue("trumbowyg-edit")
  if title != "" {
    dbArticleInsert(title, content)
  }
	display(w, "admin_edit", &Page{Title: "Admin", Content:"Yolo"})
}

func admnHdl(w http.ResponseWriter, r *http.Request){
  articles := dbRead()
	display(w, "admin", &ArticlePage{Title: "Admin", Articles: articles})
}

/*
func getrss() RssFeed{
  var rss RssFeed
  return rss
}

func rssHdl(w http.ResponseWriter, r *http.Request){
  rss := getrss()
  display(w, "feed", &Feed{RSS: rss}
}*/

func main() {
  http.HandleFunc("/", indxHdl)
  http.HandleFunc("/article", artclHdl)
  http.HandleFunc("/admin", admnHdl)
  http.HandleFunc("/edit", admn_editHdl)
 /* http.HandleFunc("/feed", rssHdl)*/
  http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
  http.ListenAndServe(":8000", nil)
  fmt.Println("Serving")
}
