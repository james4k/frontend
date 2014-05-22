// +build ignore

package main

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/russross/blackfriday"
	"j4k.co/fmatter"
	"j4k.co/layouts"
)

type blogData struct {
	baseData
	*blogPost
	Title       string
	Description string
}

type blogIndexData struct {
	baseData
	Posts     []*blogPost
	Recent    []*blogPost
	Remaining []*blogPost
}

type blogPost struct {
	Layout      string
	Title       string
	Description string
	Content     template.HTML
	Path        string
	PathTitle   string
	Date        string
	ShortDate   string
	datetime    time.Time
}

type blogPostSlice []*blogPost

func (b blogPostSlice) Len() int {
	return len(b)
}

func (b blogPostSlice) Less(i, j int) bool {
	return b[i].datetime.After(b[j].datetime)
}

func (b blogPostSlice) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func init() {
	routePrecached("/blog/{yyyy}/{mm}/{title}", NewBlogPrecache)
	routePrecached("/blog", NewBlogIndexPrecache)
}

var postNameRE = regexp.MustCompile(`^([0-9]{4})-([0-9]{2})-([0-9]{2})-([0-9A-z\-]+)$`)

var posts struct {
	sync.RWMutex
	m map[string]*blogPost
}

func init() {
	posts.m = blogPosts()
}

type BlogPrecache struct {
}

func NewBlogPrecache() renderer {
	return &BlogPrecache{}
}

func (p *BlogPrecache) Map(emit emitFunc) {
	for _, v := range posts.m {
		emit("yyyy", v.datetime.Format("2006"),
			"mm", v.datetime.Format("01"),
			"title", v.PathTitle)
	}
}

func (p *BlogPrecache) Render(w io.Writer, vars map[string]string) {
	key := fmt.Sprintf("/blog/%s/%s/%s", vars["yyyy"], vars["mm"], vars["title"])
	posts.RLock()
	post, ok := posts.m[key]
	posts.RUnlock()
	if !ok {
		errorlog.Printf("blog post not found: %s\n", key)
		return
	}

	data := new(blogData)
	data.MenderMap = menderMap
	data.blogPost = post
	data.Title = post.Title
	data.Description = post.Description

	/*
		tmpl := template.New("blog")
		tmpl.Funcs(funcs)
		_, err := tmpl.Parse(post.Content)
		if err != nil {
			errorlog.Println(err)
			return
		}
	*/
	err := layouts.ExecuteHTML(w, post.Layout, post.Content, data)
	if err != nil {
		errorlog.Println(err)
		return
	}
}

type BlogIndexPrecache struct {
}

func NewBlogIndexPrecache() renderer {
	return &BlogIndexPrecache{}
}

func (p *BlogIndexPrecache) Render(w io.Writer, vars map[string]string) {
	const recentCount = 3

	data := new(blogIndexData)
	data.MenderMap = menderMap
	data.NoIndex = true

	posts.RLock()
	for _, post := range posts.m {
		data.Posts = append(data.Posts, post)
	}
	posts.RUnlock()

	sort.Sort(blogPostSlice(data.Posts))
	if len(data.Posts) > recentCount+1 {
		data.Recent = data.Posts[:recentCount]
		data.Remaining = data.Posts[recentCount:]
	} else {
		data.Recent = data.Posts
		data.Remaining = nil
	}

	tmpl := template.New("blog")
	tmpl.Funcs(funcs)
	bytes, err := fmatter.ReadFile("content/templates/blog.html", &data.baseData)
	if err != nil {
		errorlog.Println(err)
		return
	}
	_, err = tmpl.Parse(string(bytes))
	if err != nil {
		errorlog.Println(err)
		return
	}
	err = layouts.Execute(w, "base", tmpl, data)
	if err != nil {
		errorlog.Println(err)
		return
	}
}

func blogPosts() map[string]*blogPost {
	matches, err := filepath.Glob(filepath.Join("content/blog", "*.md"))
	if err != nil {
		errorlog.Println(err)
		return nil
	}

	p := make(map[string]*blogPost, len(matches))
	for _, file := range matches {
		base := filepath.Base(file)
		ext := filepath.Ext(base)
		withoutExt := base[:len(base)-len(ext)]
		rm := postNameRE.FindStringSubmatch(withoutExt)
		if len(rm) < 5 {
			errorlog.Println(fmt.Errorf("bad blog post name: %s", base))
			return nil
		}

		post := new(blogPost)
		post.datetime, err = time.Parse("2006 01 02",
			fmt.Sprintf("%s %s %s", rm[1], rm[2], rm[3]))
		bytes, err := fmatter.ReadFile(file, &post)
		if err != nil {
			errorlog.Println(err)
			return nil
		}
		post.Content = template.HTML(blackfriday.MarkdownBasic(bytes))
		if post.Layout == "" {
			post.Layout = "blog"
		}
		post.Date = post.datetime.Format("January 2, 2006")
		post.ShortDate = post.datetime.Format("2006/01/02")
		post.Path = fmt.Sprintf("/blog/%s/%s/%s", rm[1], rm[2], rm[4])
		post.PathTitle = rm[4]
		p[post.Path] = post
	}
	return p
}
