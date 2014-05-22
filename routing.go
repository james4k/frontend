package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"j4k.co/goimport"
	"j4k.co/goimport/github"
	"j4k.co/mender"
	"j4k.co/pages"
)

type notFoundHandler struct {
	pages.Template
}

func (h *notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	h.Render(w, nil)
}

// this is all kind of dumb. j4k.co/pages is an old idea and just isn't
// working. need to take a deep look at similar packages that could
// replace it.
func assemble() http.Handler {
	assetMap := map[string]string{}
	pages.NewDefault("content/pages", "content/layouts")
	pages.Funcs(template.FuncMap{
		"cdn": func(s string) string {
			return "/" + s
		},
		"asset": func(s string) string {
			m, ok := assetMap[s]
			if !ok {
				return s
			}
			return m
		},
	})
	if production {
		pages.SetPrecache(true)
	}

	j4kco := mux.NewRouter().StrictSlash(true)
	j4kco.Handle("/", http.RedirectHandler("http://james4k.com", 302))

	notfound := pages.Dynamic("404.html", &notFoundHandler{})
	//articles := pages.Dir("articles", nil, notfound)
	articles := collectArticles("content/articles", notfound)
	root := mux.NewRouter().StrictSlash(true)
	if production {
		j4kcoAndPkgs := goimport.Handle(j4kco, github.Packages{
			User:             "james4k",
			FilterByHomepage: "^(http://)?j4k.co",
		})
		root.Host("j4k.co").Handler(j4kcoAndPkgs)
	}
	root.Handle("/", pages.Static("index.html", map[string]interface{}{
		"articles": articles.Collection(),
	}))
	root.Handle("/about", pages.Static("about.html", nil))
	//root.Handle("/games", pages.Static("games.html", nil))
	//root.Handle("/software", pages.Static("software.html", nil))
	assets := root.PathPrefix("/assets/")
	root.Handle("/{any}", articles)
	root.NotFoundHandler = notfound

	if production {
		assetMap = mender.VersionMap("content/mend-versions.json")
		assets.Handler(http.FileServer(http.Dir("content")))
	} else {
		assetServer := mender.Watch("content/mend.json", "content",
			http.FileServer(http.Dir("content")), os.Stderr)
		assetServer.OnChange = func() {
			assetMap = assetServer.VersionMap()
			log.Println("--- mend ---")
			for _, v := range assetMap {
				log.Println(v)
			}
			log.Println("------------")
		}
		assets.Handler(assetServer)
	}
	return root
}
