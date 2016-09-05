package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/gorilla/mux"
	"github.com/jawher/mow.cli"
	"golang.org/x/net/context"
)

func main() {
	app := cli.App("up-coco-admin", "Universal Publishing Platform - coco devops admin tool")
	etcdURL := app.String(cli.StringOpt{Name: "etcd-url", Value: "http://localhost:2379", Desc: "etcd URL", EnvVar: "ETCD_URL"})
	port := app.Int(cli.IntOpt{Name: "port", Value: 8080, Desc: "tcp port to listen on", EnvVar: "PORT"})

	app.Action = func() {
		runServer(*port, *etcdURL)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

var index = []byte(`
<body>
<a href="etcd-all">Dump of etcd config</a>
</body>
`)

func runServer(port int, etcdURL string) {
	eh := &etcdHandlers{port, etcdURL}

	router := mux.NewRouter()
	http.Handle("/", router)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(index)
	}).Methods("GET")

	router.HandleFunc("/etcd-all", eh.dump).Methods("GET")
	router.PathPrefix("/__api").Handler(http.StripPrefix("/__api", http.FileServer(http.Dir("./api/"))))

	addr := fmt.Sprintf(":%d", port)
	log.Printf("about to listen on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

type etcdHandlers struct {
	port    int
	etcdURL string
}

func (eh *etcdHandlers) dump(w http.ResponseWriter, r *http.Request) {
	cfg := client.Config{
		Endpoints:               []string{eh.etcdURL},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second * 5,
	}

	c, err := client.New(cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	kapi := client.NewKeysAPI(c)

	response, err := kapi.Get(context.Background(), "/", &client.GetOptions{Recursive: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	// enc.SetIndent("  ", "  ") // commented out because it required go1.7

	w.Header().Set("Content-Type", "application/json")

	if err := enc.Encode(tree(response.Node)); err != nil {
		log.Printf("%v\n", err)
	}
}

func tree(n *client.Node) interface{} {
	if !n.Dir {
		return n.Value
	}
	m := make(map[string]interface{})
	for _, child := range n.Nodes {
		m[filepath.Base(child.Key)] = tree(child)
	}
	return m
}
