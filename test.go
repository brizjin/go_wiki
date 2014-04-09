package main

import (
    "fmt"
    "net/http"
    "runtime"
    "io/ioutil"
    "html/template"
    "regexp"
    "time"
)
var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Hello struct{}
type Page struct {
    Title string
    Body  []byte
}



func (p *Page) save() error {
    filename := p.Title + ".txt"
    time.Sleep(5000 * time.Millisecond)
    return ioutil.WriteFile(filename, p.Body, 0600)
}
func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func (h Hello) ServeHTTP(
    w http.ResponseWriter,
    r *http.Request) {
    fmt.Fprint(w, "Hello!")
}

type String string

type Struct struct {
    Greeting string
    Punct    string
    Who      string
}

func(s String) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request){
    //fmt.Println("OUT!")
    numCPU := runtime.NumCPU()
	fmt.Fprint(w,"Hello string! Number of CPUs: " +  fmt.Sprintf("%v",numCPU))
}

func(s Struct) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request){
	fmt.Fprint(w,"Hello STRUCT")
}
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        p = &Page{Title: title}
    }
    renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    //err := p.save()

    go p.save()
    /*if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }*/
    //http.Redirect(w, r, "/view/"+title, http.StatusFound)
    renderTemplate(w, "view", p)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
    }
}
func main() {
    
   
    
    
    http.Handle("/string", String("I'm a frayed knot."))
    http.Handle("/struct", &Struct{"Hello", ":", "Gophers!"})
    //http.Handle("/", &Hello{})
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
    http.HandleFunc("/",func (w http.ResponseWriter, r *http.Request){
        fmt.Fprintf(w,"<html><body>")
        fmt.Fprintf(w, "Hi there, I love <b>%s</b>!<br>", r.URL.Path[1:])
        vs,_ := ioutil.ReadDir(".")
        for _,v := range vs {
            name := v.Name()
            name_without_ext := name[:len(name)-4]
            name_ext := name[len(name)-3:]
            if name_ext == "txt" {
                fmt.Fprintf(w,"[<a href=\"/view/%s\">%s</a>]<br>",name_without_ext,name_without_ext)    
            }
            
        }
        fmt.Fprintf(w,"</html></body>")
        
    })

    i := 1    
    go func(){
        for{
            time.Sleep(5000 * time.Millisecond)
            filename := fmt.Sprintf("autofile%v.txt",i)
            ioutil.WriteFile(filename, []byte("auto"), 0600)
            i++
        }        
    }()

    http.ListenAndServe("localhost:4000", nil)


    /*p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
    p1.save()
    p2, _ := loadPage("TestPage")
    fmt.Println(string(p2.Body))*/


    /*m := image.NewRGBA(image.Rect(0, 0, 100, 100))
    fmt.Println(m.Bounds())
    fmt.Println(m.At(0, 0).RGBA())*/
}