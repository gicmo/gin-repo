package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/docopt/docopt-go"

	"github.com/gorilla/mux"

	. "github.com/G-Node/gin-repo/common"
	"github.com/G-Node/gin-repo/ssh"
	"github.com/G-Node/gin-repo/wire"
	"path/filepath"
	"strings"
	"gopkg.in/libgit2/git2go.v23"
)

func getRepoDir() string {
	dir := os.Getenv("GIN_REPO_KEYDIR")

	if dir == "" {
		dir = "."
	}

	return dir
}

func translatePath(vpath string, uid string) string {
	dir := os.Getenv("GIN_REPO_DIR")

	if dir == "" {
		dir = "."
	}

	if strings.HasPrefix(vpath, "'") && strings.HasSuffix(vpath, "'") {
		vpath = vpath[1 : len(vpath)-1]
	}

	path := filepath.Join(dir, uid, vpath)

	if !strings.HasSuffix(path, ".git") {
		path += ".git"
	}

	path, err := filepath.Abs(path)

	//TODO: propagate the error
	if err != nil {
		return path
	}

	fmt.Fprintf(os.Stderr, "[D] tp: %s@%s -> %s\n", uid, vpath, path)

	return path
}

func repoAccess(w http.ResponseWriter, r *http.Request) {
	log.Printf("repoAccess: %s @ %v", r.Method, r.URL.String())

	decoder := json.NewDecoder(r.Body)
	var query wire.RepoAccessQuery
	err := decoder.Decode(&query)

	if err != nil || query.Path == "" || query.User == "" || query.Method == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Error precessing request: %v", err)
		return
	}

	//TODO: check access here

	path := translatePath(query.Path, query.User)
	w.Write([]byte(path))
}

func lookupUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("lookupUser: %s @ %s", r.Method, r.URL.String())

	query := r.URL.Query()

	val, ok := query["key"]

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	keys := ssh.ReadKeysInDir(getRepoDir())

	if key, ok := keys[val[0]]; ok {
		user := User{Uid: key.Comment, Keys: []ssh.Key{key}}

		data, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(data)
		if err != nil {
			log.Panicf("[W] %v", err)
		}

	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}

func createRepo(w http.ResponseWriter, r *http.Request) {
	log.Printf("createRepo: %s @ %s", r.Method, r.URL.String())

	decoder := json.NewDecoder(r.Body)
	var creat wire.CreateRepo
	err := decoder.Decode(&creat)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Error precessing request: %v", err)
		return
	} else if creat.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Error precessing request: name missing")
		return
	}

	//TODO: check name for sanity

	vars := mux.Vars(r)
	user := vars["user"]

	path := translatePath(creat.Name, user)

	_, err = git.InitRepository(path, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
}


func main() {
	usage := `gin repo daemon.

Usage:
  gin-repod
  gin-repod -h | --help
  gin-repod --version

Options:
  -h --help     Show this screen.
  `

	args, _ := docopt.Parse(usage, nil, true, "gin repod 0.1a", false)
	fmt.Println(args)

	r := mux.NewRouter()
	r.HandleFunc("/intern/user/lookup", lookupUser).Methods("GET")
	r.HandleFunc("/intern/repos/access", repoAccess).Methods("POST")

	r.HandleFunc("/users/{user}/repos", createRepo).Methods("POST")

	http.Handle("/", r)

	hostport := ":8888"
	log.Printf("Listening on " + hostport)
	log.Fatal(http.ListenAndServe(hostport, nil))
}
