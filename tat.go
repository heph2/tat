package main

import (
	"io/ioutil"
	"log"
	"os/exec"
	"fmt"
	"os"
	"html/template"
	"flag"
	"strings"
	"net/http"
	gem "tildegit.org/nihilazo/go-gemtext"
)

const (
	usage = `usage: %s
Generate Static Websites with tat

Options:
`
)

var (
	pathPtr = flag.String("dir", ".", "Directory to serve")
	servePtr = flag.Bool("serve", false, "Serve static files")
	buildPtr = flag.Bool("build", false, "Build blog")
	initPtr = flag.Bool("init", false, "Generate files")
)

// Post struct
type Post struct {
	Title   string
	Content template.HTML
}

// HomePage struct
type HomePage struct {
	Post    Post
	Posts   []string
}

func generateHTML(path string) (template.HTML, error) {
	gemfile, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("0.0: Cannot read gemfile: %w", err)
	}
	gemtext, err := gem.ParsePage(string(gemfile))
	if err != nil {
		return "", fmt.Errorf("0.0: Cannot parse gemfile: %w", err)
	}
	renderedHTML, err := (gem.RenderHTML(&gemtext))
	if err != nil {
		return "", fmt.Errorf("0.0: Cannot convert to HTML: %w", err)
	}

	return template.HTML(renderedHTML), nil	
}

func sanitizeName(gemtextName string) string {
	return strings.ReplaceAll(gemtextName, ".gmi", ".html")
}

func generateRoot() {
	var postsName []string
	tmpRoot, _ := template.ParseFiles("layouts/index.tmpl.html")
	posts, err := ioutil.ReadDir("posts/")
	if err != nil {
		log.Fatal(err)
	}
	
	// Generate HTML for posts
	for _, post := range posts {
		postsName = append(postsName, sanitizeName(post.Name()))
		tmpPost, _ := template.ParseFiles("layouts/post.tmpl.html")
		postHTML, err := generateHTML("posts/" + post.Name())
		postStr := Post{Title: post.Name(), Content: postHTML}
		if err != nil {
			log.Println(err)
		}
		f, err := os.Create("out/posts/" + sanitizeName(post.Name()))
		defer f.Close()
		if err != nil {
			log.Println(err)
		}
		if err = tmpPost.Execute(f, postStr); err != nil {
			log.Println(err)
		}
	}
		
	postHTML, err := generateHTML("pages/home.gmi")
	if err != nil {
		log.Println(err)
		return
	}
		
	home := HomePage{
		Post: Post{
			Title: "Home",
			Content: postHTML,
		},
		Posts: postsName,
	}
	
	f, err := os.Create("out/index.html")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	
	if err = tmpRoot.Execute(f, home); err != nil {
		log.Fatal(err)
	}
}

func createDirs() {
	if err := os.MkdirAll("out/pages", os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll("out/posts", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// Copy assets directory to output
	cmd := exec.Command("cp", "-R", "assets", "out/")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func initBlog() {
	if err := os.MkdirAll("pages", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll("layouts", os.ModePerm); err != nil {
		log.Fatal(err)
	}	

	if err := os.MkdirAll("assets", os.ModePerm); err != nil {
		log.Fatal(err)
	}	

	if err := os.MkdirAll("posts", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("layouts/index.tmpl.html",
		[]byte(
`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
  <title>{{ .Post.Title }}</title>
</head>
<body>
  {{ .Post.Content }}
  <ul>
    {{ range .Posts }}
    <li><a href="/posts/{{ . }}">{{ . }}</a></li>
    {{ end }}
  </ul>
</body>
</html>
`), 0666); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("layouts/post.tmpl.html",
		[]byte(
			`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
  <title>{{.Title}}</title>
</head>
<body>
  {{.Content}}
</body>
</html>
`), 0666); err != nil {
		log.Fatal(err)
	}
}

func runServer() {
	fs := http.FileServer(http.Dir("./out"))
	http.Handle("/", fs)

	log.Println("Listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	if *initPtr {
		initBlog()
		return
	}
	
	if *buildPtr {
		createDirs()
		generateRoot()
		return
	}
	
	if *servePtr {
		createDirs()
		generateRoot()
		runServer()
		return
	}

	fmt.Fprintf(flag.CommandLine.Output(), usage, os.Args[0])	
	flag.PrintDefaults()
}

