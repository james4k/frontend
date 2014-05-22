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
	h.m = map[string]http.Handler{}
	glob := filepath.Join(h.dir, "*.md")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return err
	}
	for _, file := range matches {
		base := filepath.Base(file)
		ext := filepath.Ext(base)
		withoutExt := base[:len(base)-len(ext)]
		rel, err := filepath.Rel("content/pages", file)
		if err != nil {
			return err
		}
		h.m[withoutExt] = pages.Static(rel, nil)
	}
	return nil
}

func (h *articleHandler) Collection() interface{} {
	return h.slice
}
