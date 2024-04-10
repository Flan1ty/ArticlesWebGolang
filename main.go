package main

import (
	"fmt"

	"net/http"

	"html/template"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

//type User struct {
//	Name                  string
//	Age                   uint16
//	Money                 int16
//	Avg_grades, Happiness float64
//	Hobbies               []string
//}
//
//func (user User) getAllInfo() string {
//	return fmt.Sprintf("User name: %s. He is %d and he has money equal: %d",
//		user.Name, user.Age, user.Money)
//}
//
//func (user *User) setNewName(newName string) {
//	user.Name = newName
//}
//
//func home_page(w http.ResponseWriter, r *http.Request) {
//	bob := User{"Bob", 25, -50, 4.2, 0.8, []string{"Computers", "science", "Music"}}
//	//fmt.Fprintln(w, `<h1>Main Text</h1>
//	//<b>Main Text</b>`)
//
//	tmpl, _ := template.ParseFiles("templates/index.html")
//	tmpl.Execute(w, bob)
//}
//
//func contacts_page(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintf(w, "Contacts page")
//}
//
//func handleRequest() {
//	http.HandleFunc("/", home_page)
//	http.HandleFunc("/contacts/", contacts_page)
//	http.ListenAndServe(":8080", nil)
//}

//type User struct {
//	name string `json:"name"`
//	age  uint16 `json:"age"`
//}
//
//func main() {
//
//	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/users?charset=utf8")
//	if err != nil {
//		panic(err)
//	}
//
//	defer db.Close()

// Установка данных
//insert, err := db.Query("INSERT INTO `users` (`name`, `age`) VALUES ('Bob', 35)")
//if err != nil {
//	panic(err)
//}
//defer insert.Close()

// Выборка данных

//res, err := db.Query("SELECT `name`, `age` FROM `users`")
//if err != nil {
//	panic(err)
//}
//
//for res.Next() {
//	var user User
//	err = res.Scan(&user.name, &user.age)
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println(fmt.Sprintf("User: %s with age: %d", user.name, user.age))
//}

//handleRequest()
//}

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

var posts = []Article{}

var showPost = Article{}

func index(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/users?charset=utf8")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Выборка данных
	res, err := db.Query("SELECT * FROM `article`")
	if err != nil {
		panic(err)
	}

	posts = []Article{}

	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}

		posts = append(posts, post)
	}

	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, req *http.Request) {
	title := req.FormValue("title")
	anons := req.FormValue("anons")
	full_text := req.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/users?charset=utf8")
		if err != nil {
			panic(err)
		}

		defer db.Close()

		//  Установка данных
		insert, err := db.Query(fmt.Sprintf("INSERT INTO `article` (`title`, `anons`, `full_text`)"+
			"VALUES ('%s', '%s', '%s')", title, anons, full_text))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

}

func show_post(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/users?charset=utf8")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Выборка данных
	res, err := db.Query(fmt.Sprintf("SELECT * FROM `article` where `id` = '%s'",
		vars["id"]))

	if err != nil {
		panic(err)
	}

	showPost = Article{}

	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}

		showPost = post

	}

	t.ExecuteTemplate(w, "show", showPost)
}

func handleFunc() {
	router := mux.NewRouter()
	router.HandleFunc("/", index)
	router.HandleFunc("/create", create).Methods("GET")
	router.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")
	router.HandleFunc("/save_article", save_article).Methods("POST")

	http.Handle("/", router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
