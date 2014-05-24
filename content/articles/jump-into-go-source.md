---
layout: article
title: Jump into Go's Source
---

People ask all the time "How does this work in Go?" or "How is that
implemented in Go?". If you have a similar question, maybe it is time
you look at the source code. Understandably, this is not people's first
instinct. However, Go's source code is very approachable, in part due to
the simplicity of the language.

If you haven't already, installing Go from
[source](http://golang.org/doc/install/source) is simple and tends to be
the preferred installation method for many Go users. Or you can browse
the source tree at <http://golang.org/src>.

So, why not get started? Here is a selection of great literature on Go's internals, as well as a simplified and annotated tree of Go's
source with an emphasis on the runtime:

* [Design Documents](https://code.google.com/p/go-wiki/wiki/DesignDocuments)
on go-wiki
* [Go Data Structures: Interfaces](http://research.swtch.com/interfaces)
(2009) by Russ Cox
* [The Go scheduler](http://morsmachine.dk/go-scheduler), and
[The Go netpoller](http://morsmachine.dk/go-netpoller) (2013) by
Daniel Morsing

<pre>
./<a href="http://golang.org/src">src</a>
├── <a href="http://golang.org/src/cmd">cmd</a> 
│   ├── <a href="http://golang.org/src/cmd/cgo">cgo</a>: cgo preprocessor
│   ├── <a href="http://golang.org/src/cmd/gc">gc</a>: Go compiler
│   ├── <a href="http://golang.org/src/cmd/go">go</a>: main CLI for Go users
│   └── ...rest of toolchain
└── <a href="http://golang.org/src/pkg">pkg</a>
    ├── <a href="http://golang.org/src/pkg/go">go</a>: tooling of Go code
    ├── <a href="http://golang.org/src/pkg/runtime">runtime</a>
    │   ├── <a href="http://golang.org/src/pkg/runtime/runtime.h">runtime.h</a>: data structures, scheduler states
    │   ├── <a href="http://golang.org/src/pkg/runtime/proc.c">proc.c</a>: scheduler
    │   ├── <a href="http://golang.org/src/pkg/runtime/mgc0.c">mgc0.c</a>: garbage collector
    │   ├── <a href="http://golang.org/src/pkg/runtime/chan.c">chan.c</a>: channels
    │   ├── <a href="http://golang.org/src/pkg/runtime/slice.c">slice.c</a>: slices
    │   ├── <a href="http://golang.org/src/pkg/runtime/hashmap.c">hashmap.c</a>: maps
    │   ├── <a href="http://golang.org/src/pkg/runtime/cgocall.c">cgocall.c</a>: cgo
    │   └── ...rest of runtime
    └── ...rest of standard library
</pre>

For the Go tools, such as godoc, vet, and cover, you will have to visit
[go.tools](https://code.google.com/p/go.tools) on Google Code. If you
want to keep up with development, the
[golang-dev](http://groups.google.com/group/golang-dev) and
[golang-nuts](http://groups.google.com/group/golang-nuts) mailing lists
are worth following, as well as the [Go development
dashboard](https://go-dev.appspot.com/).

Whew. If you got that far, consider
[contributing](http://golang.org/doc/contribute.html)!
