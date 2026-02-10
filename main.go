package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

//go:embed frontend/dist/*
var assets embed.FS
var app *App

func main() {
	app = NewApp()
    app.startup(nil)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/deps", handleDeps)
	mux.HandleFunc("/api/tunnel", handleTunnel)
	mux.HandleFunc("/api/stream/start", handleStartStream)
	mux.HandleFunc("/api/stream/stop", handleStopStream)

	signaling := NewSignalingServer()
	mux.HandleFunc("/ws", signaling.HandleConnections)
    mux.HandleFunc("/stream-ws", handleStreamWS)

	distFS, _ := fs.Sub(assets, "frontend/dist")
	fileServer := http.FileServer(http.FS(distFS))
    
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/invitado" {
            http.ServeFile(w, r, "./guest/index.html")
            return
        }
		if strings.HasPrefix(r.URL.Path, "/api") || strings.HasPrefix(r.URL.Path, "/ws") || strings.HasPrefix(r.URL.Path, "/stream-ws") {
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	fmt.Println("\n" + `   _____ _____  _    _ _____  `)
	fmt.Println(`  / ____|  __ \| |  | |  __ \ `)
	fmt.Println(` | |  __| |__) | |__| | |__) |`)
	fmt.Println(` | | |_ |  ___/|  __  |  _  / `)
	fmt.Println(` | |__| | |    | |  | | | \ \ `)
	fmt.Println(`  \_____|_|    |_|  |_|_|  \_\`)
	fmt.Println("\n ✨ GPHR Studio (Gopher Quick Studio) está listo!")
	fmt.Println(" -------------------------------------------")
	fmt.Println(" 📺 Panel de Control: http://localhost:8080")
	fmt.Println(" -------------------------------------------")
	fmt.Println(" 🚀 Presiona Ctrl+C para apagar el estudio.\n")
	
    go func() {
        time.Sleep(800 * time.Millisecond)
        openBrowser("http://localhost:8080")
    }()

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println(" ❌ Error al iniciar el servidor:", err)
	}
}

func handleDeps(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(app.CheckDependencies())
}
func handleTunnel(w http.ResponseWriter, r *http.Request) {
    url, err := app.ToggleTunnel()
    if err != nil { http.Error(w, err.Error(), 500); return }
    json.NewEncoder(w).Encode(map[string]string{"url": url})
}
func handleStartStream(w http.ResponseWriter, r *http.Request) {
    var req StreamReq
    json.NewDecoder(r.Body).Decode(&req)
    app.StartStream(req.Url)
    w.WriteHeader(200)
}
func handleStopStream(w http.ResponseWriter, r *http.Request) {
    app.StopStream()
    w.WriteHeader(200)
}
func handleStreamWS(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil { return }
    defer conn.Close()
    for {
        _, message, err := conn.ReadMessage()
        if err != nil { break }
        app.PushVideoChunk(message)
    }
}
type StreamReq struct { Url string `json:"url"` }