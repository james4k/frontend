package main

import (
	"log"
	"net/http"
	"path/filepath"

	"j4k.co/pages"
)

type article struct {
	Title string
	Path  string
}

type articleHandler struct {
	dir      string
	notfound http.Handler
	m        map[string]http.Handler
	slice    []article
}

func collectArticles(dir string, notfound http.Handler) *articleHandler {
	h := &articleHandler{
		dir:      dir,
		notfound: notfound,
	}
	if notfound == nil {
		h.notfound = http.HandlerFunc(http.NotFound)
	}
	err := h.read()
	if err != nil {
		log.Fatalln(err)
	}
	return h
}

func (h *articleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path
	if key[0] == '/' {
		key = key[1:]
	}
	a, ok := h.m[key]
	if !ok {
		h.notfound.ServeHTTP(w, r)
		return
	}
	a.ServeHTTP(w, r)
}

func (h *articleHandler) read() error {
	glob := filepath.Join(h.dir, "*.md")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return err
	}
	h.m = make(map[string]http.Handler, len(matches))
	h.slice = make([]article, len(matches))
	for i, file := range matches {
		base := filepath.Base(file)
		ext := filepath.Ext(base)
		withoutExt := base[:len(base)-len(ext)]
		rel, err := filepath.Rel("content", file)
		if err != nil {
			return err
		}
		p := pages.Static(rel, &h.slice[i])
		h.slice[i].Path = "/" + withoutExt
		h.m[withoutExt] = p
	}
	return nil
}

type articles []article

func (a articles) Len() int {
	return len(a)
}

func (a articles) Index(i int) interface{} {
	return a[i]
}

func (a articles) Slice(i, j int) collection {
	return a[i:j]
}

func (h *articleHandler) Collection() interface{} {
	return h.slice
}
