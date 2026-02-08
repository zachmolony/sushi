package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// LocalFileServer serves .glb/.gltf files on its own HTTP port.
// This avoids issues with Wails' asset server and Vite proxying.
type LocalFileServer struct {
	Port int
}

func StartLocalFileServer() *LocalFileServer {
	mux := http.NewServeMux()
	mux.HandleFunc("/localfile/", serveLocalFile)

	// Find a free port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Printf("[fileserver] failed to listen: %v\n", err)
		return &LocalFileServer{Port: 0}
	}
	port := listener.Addr().(*net.TCPAddr).Port
	fmt.Printf("[fileserver] serving on http://127.0.0.1:%d\n", port)

	go func() {
		if err := http.Serve(listener, mux); err != nil {
			fmt.Printf("[fileserver] server error: %v\n", err)
		}
	}()

	return &LocalFileServer{Port: port}
}

func serveLocalFile(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	fmt.Printf("[fileserver] request: %s -> path=%s\n", r.URL.String(), filePath)

	if filePath == "" {
		http.Error(w, "missing path parameter", http.StatusBadRequest)
		return
	}

	// Security: only serve .glb / .gltf files
	lower := strings.ToLower(filePath)
	if !strings.HasSuffix(lower, ".glb") && !strings.HasSuffix(lower, ".gltf") {
		http.Error(w, "forbidden file type", http.StatusForbidden)
		return
	}

	// Check file exists
	info, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	// Set appropriate content type
	contentType := "model/gltf-binary"
	if strings.HasSuffix(lower, ".gltf") {
		contentType = "model/gltf+json"
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	// Open and serve the file
	f, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "cannot open file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	http.ServeContent(w, r, filepath.Base(filePath), info.ModTime(), f)
}
