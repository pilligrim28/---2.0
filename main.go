package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

type Radio struct {
  ID      string `json:"id"`
  Name    string `json:"name"`
  Address string `json:"address"`
}

var (
  dbFile = "db.json"
  mu     sync.Mutex
  upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
  }
)

func main() {
  // Инициализация Gin
  r := gin.Default()

  // Статические файлы фронтенда
  r.StaticFS("/", http.FS(frontendFS))

  // API
  api := r.Group("/api")
  {
    api.GET("/radios", getRadios)
    api.POST("/radios", addRadio)
  }

  // WebSocket
  r.GET("/ws", func(c *gin.Context) {
    conn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
    defer conn.Close()
    for {
      // Чтение сообщений (опционально)
      _, msg, _ := conn.ReadMessage()
      println("WebSocket:", string(msg))
    }
  })

  // TCP-сервер
  go startTCPServer()

  // UDP-сервер
  go startUDPServer()

  // Запуск
  r.Run(":3000")
}

func startTCPServer() {
  ln, _ := net.Listen("tcp", ":5000")
  defer ln.Close()
  for {
    conn, _ := ln.Accept()
    go handleTCPConnection(conn)
  }
}

func handleTCPConnection(conn net.Conn) {
  defer conn.Close()
  buf := make([]byte, 1024)
  n, _ := conn.Read(buf)
  println("TCP:", string(buf[:n]))
}

func startUDPServer() {
  pc, _ := net.ListenPacket("udp", ":5001")
  defer pc.Close()
  buf := make([]byte, 1024)
  for {
    n, addr, _ := pc.ReadFrom(buf)
    println("UDP:", string(buf[:n]), "from", addr)
  }
}

func getRadios(c *gin.Context) {
  mu.Lock()
  defer mu.Unlock()

  file, _ := os.ReadFile(dbFile)
  var radios []Radio
  json.Unmarshal(file, &radios)
  c.JSON(http.StatusOK, radios)
}

func addRadio(c *gin.Context) {
  mu.Lock()
  defer mu.Unlock()

  var radio Radio
  c.BindJSON(&radio)

  file, _ := os.ReadFile(dbFile)
  var radios []Radio
  json.Unmarshal(file, &radios)

  radio.ID = fmt.Sprintf("%d", len(radios)+1)
  radios = append(radios, radio)

  data, _ := json.Marshal(radios)
  os.WriteFile(dbFile, data, 0644)

  c.JSON(http.StatusCreated, radio)
}