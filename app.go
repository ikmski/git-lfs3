package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ikmski/git-lfs3/adapter"
)

const (
	contentMediaType = "application/vnd.git-lfs"
	metaMediaType    = "application/vnd.git-lfs+json"
)

type app struct {
	config serverConfig
	router *mux.Router
}

func newApp(
	conf serverConfig,
	batchController adapter.BatchController,
	transferController adapter.TransferController,
	lockController adapter.LockController) *app {

	a := &app{
		config: conf,
	}

	r := mux.NewRouter()

	// Batch
	r.Methods("POST").Path("/{user}/{repo}/objects/batch").MatcherFunc(MetaMatcher).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) { batchController.Batch(newContext(w, r)) })

	// Transfer
	r.Methods("GET").Path("/{user}/{repo}/objects/{oid}").MatcherFunc(ContentMatcher).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) { transferController.Download(newContext(w, r)) })
	r.Methods("PUT").Path("/{user}/{repo}/objects/{oid}").MatcherFunc(ContentMatcher).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) { transferController.Upload(newContext(w, r)) })

	// Lock
	r.Methods("GET").Path("/{user}/{repo}/locks").MatcherFunc(MetaMatcher).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) { lockController.List(newContext(w, r)) })
	r.Methods("POST").Path("/{user}/{repo}/locks/verify").MatcherFunc(MetaMatcher).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) { lockController.Verify(newContext(w, r)) })
	r.Methods("POST").Path("/{user}/{repo}/locks").MatcherFunc(MetaMatcher).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) { lockController.Lock(newContext(w, r)) })
	r.Methods("POST").Path("/{user}/{repo}/locks/{id}/unlock").MatcherFunc(MetaMatcher).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) { lockController.Unlock(newContext(w, r)) })

	a.router = r

	return a
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	a.router.ServeHTTP(w, r)
}

func (a *app) serve() error {

	s := &http.Server{
		Handler: a.router,
		Addr:    fmt.Sprintf(":%d", a.config.Port),
	}

	if a.config.Tls {
		return s.ListenAndServeTLS(a.config.CertFile, a.config.KeyFile)
	} else {
		return s.ListenAndServe()
	}
}

func ContentMatcher(r *http.Request, m *mux.RouteMatch) bool {
	mediaParts := strings.Split(r.Header.Get("Accept"), ";")
	mt := mediaParts[0]
	return mt == contentMediaType
}

func MetaMatcher(r *http.Request, m *mux.RouteMatch) bool {
	mediaParts := strings.Split(r.Header.Get("Accept"), ";")
	mt := mediaParts[0]
	return mt == metaMediaType
}
