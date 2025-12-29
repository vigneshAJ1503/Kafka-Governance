package main

import (
	"log"
	"net/http"

	"kafka-governance/config"
	"kafka-governance/db"
	"kafka-governance/routes"
)

func main() {
	cfg := config.Load()

	client, database, err := db.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(nil)

	db.InitTopicRepo(database)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	routes.Register(mux)

	log.Println("Server running on :", cfg.AppPort)
	http.ListenAndServe(":"+cfg.AppPort, mux)
}
