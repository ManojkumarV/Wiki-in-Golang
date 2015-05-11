/**
* @author Manojkumar V
*/
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string, extension string) (*Page, error) {
	filename := title + extension
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	style, _ := loadPage("css/style", ".css")
	fmt.Fprintf(w, "<style>%s</style>", style.Body)
	fmt.Fprintf(w, "<div class='editarea'><h1 style='text-align:center'>New Topic</h1>"+
		"<form action=\"/add/\" method=\"POST\">"+
		"<input type='text' name='newtopic' placeholder='Topic' class='topicbox'><br /><br /><textarea name=\"newcontent\" class='editbox'></textarea><br>"+
		"<input type=\"submit\" value=\"Save\" class='submitbutton'>"+
		"</form>")
	fmt.Fprintf(w, "<a href='../wiki/'>Back</a><br /><br /></div>")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	var content string
	style, _ := loadPage("css/style", ".css")
	fmt.Fprintf(w, "<!doctype html><html><style>%s</style>", style.Body)
	title := r.URL.Path[len("/wiki/"):]
	if len(title) > 1 {
		fmt.Fprintf(w, "<title>%s</title>", title)
		fmt.Fprintf(w, "<body><div class='title'><br />%s</div>", title)
	} else {
		fmt.Fprintf(w, "<title>Online Wiki</title>")
		fmt.Fprintf(w, "<body><div class='title'><br />Welcome to online wiki. Developed using Go lang</div>")
	}
	fmt.Fprintf(w, "<div class='main'>")
	fmt.Fprintf(w, "<div class='left'>")
	fmt.Fprintf(w, "<h3 style='text-align:center'>Other Topics</h3>")
	files, _ := ioutil.ReadDir("data")
	for _, f := range files {
		fmt.Fprintf(w, "<div style='padding-top:5px'><a href='"+f.Name()[:len(f.Name())-4]+"' class='pagelink'>"+f.Name()[:len(f.Name())-4]+"</a></div>")
	}
	fmt.Fprintf(w, "</div>")
	if len(title) < 1 {
		p, _ := loadPage("about", ".txt")
		content = string(p.Body)
		content = strings.Replace(content, "<", "&lt", -1)
		content = strings.Replace(content, ">", "&gt", -1)
		content = strings.Replace(content, "\t", "&nbsp;&nbsp;&nbsp;&nbsp;", -1)
		content = strings.Replace(content, "\n", "<br />", -1)
		fmt.Fprintf(w, "<div class='content'>%s</div>", content)
	} else if _, err := os.Stat("data/" + title + ".txt"); os.IsNotExist(err) {
		fmt.Fprintf(w, "<div class='content'><h2>Content Unavailable</h2></div>")
	} else {
		p, _ := loadPage("data/"+title, ".txt")
		content = string(p.Body)
		content = strings.Replace(content, "<", "&lt", -1)
		content = strings.Replace(content, ">", "&gt", -1)
		content = strings.Replace(content, "\t", "&nbsp;&nbsp;&nbsp;&nbsp;", -1)
		content = strings.Replace(content, "\n", "<br />", -1)
		fmt.Fprintf(w, "<div class='content'>%s</div>", content)
		fmt.Fprintf(w, "<a href='../edit/"+title+"'>Edit</a>")
	}
	fmt.Fprintf(w, "</div><div class='footer'>")
	fmt.Fprintf(w, "<br /><a href='../new/' style='padding-left:20px'>Add New Content</a>")
	fmt.Fprintf(w, "</div>")
	fmt.Fprintf(w, "</body></html>")
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	style, _ := loadPage("css/style", ".css")
	fmt.Fprintf(w, "<style>%s</style>", style.Body)
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage("data/"+title, ".txt")
	if err != nil {
		p = &Page{Title: title}
	}
	fmt.Fprintf(w, "<div class='editarea'><h1 style='text-align:center'>Editing %s</h1>"+
		"<form action=\"/save/%s\" method=\"POST\">"+
		"<textarea name=\"body\" class='editbox'>%s</textarea><br>"+
		"<input type=\"submit\" value=\"Save\" class='submitbutton'>"+
		"</form>",
		title, p.Title, p.Body)
	fmt.Fprintf(w, "<a href='../wiki/"+title+"'>Back</a><br /><br /></div>")
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	title = title + ".txt"
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/wiki/"+title[5:len(title)-4], http.StatusFound)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("newtopic")
	title = "data/" + title + ".txt"
	body := r.FormValue("newcontent")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/wiki/"+title[5:len(title)-4], http.StatusFound)
}

func main() {
	http.HandleFunc("/wiki/", viewHandler)
	http.HandleFunc("/add/", addHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/new/", newHandler)
	http.HandleFunc("/", addHandler)
	http.ListenAndServe(":6060", nil)
}
