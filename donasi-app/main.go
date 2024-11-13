package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Struktur untuk mendekode data donasi
type DonationMessage struct {
	Amount  int    `json:"amount"`
	Message string `json:"message"`
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	tcpListener net.Listener
	clients     = make(map[*websocket.Conn]bool)
	mu          sync.Mutex
)

func main() {
	// Menjalankan WebSocket Server untuk Pendonor dan Penerima
	http.HandleFunc("/ws", handleWebSocket)

	// Menjalankan server TCP untuk pemrosesan donasi
	var err error
	tcpListener, err = net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("Error memulai TCP listener:", err)
	}
	defer tcpListener.Close()

	go handleTCP()

	// Menjalankan server UDP untuk notifikasi donasi
	go handleUDP()

	// Menjalankan server HTTP untuk menampilkan halaman HTML
	http.Handle("/", http.FileServer(http.Dir("./static"))) // Menyajikan file statis (halaman HTML)
	fmt.Println("Server dimulai di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler WebSocket untuk mengelola koneksi
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	fmt.Println("Koneksi WebSocket baru terhubung")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			return
		}

		// Parsing pesan JSON dari pendonor
		var donation DonationMessage
		err = json.Unmarshal(msg, &donation)
		if err != nil {
			log.Println("Error parsing JSON:", err)
			continue
		}

		// Meneruskan pesan donasi ke penerima
		sendToReceiver(donation)
	}
}

// Fungsi untuk mengirim pesan ke semua klien yang terhubung
func sendToReceiver(donation DonationMessage) {
	mu.Lock()
	defer mu.Unlock()

	// Mengubah objek donasi menjadi JSON
	msg, err := json.Marshal(donation)
	if err != nil {
		log.Println("Error marshaling donation:", err)
		return
	}

	// Mengirim pesan donasi ke semua klien WebSocket (penerima)
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println(err)
		}
	}
}

// Server TCP untuk menerima donasi dari klien
func handleTCP() {
	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			log.Println("Error menerima koneksi TCP:", err)
			continue
		}

		go handleTCPConnection(conn)
	}
}

// Menangani koneksi TCP individual untuk pemrosesan donasi
func handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("Error membaca dari koneksi TCP:", err)
			break
		}
		// Mengirim pesan donasi ke penerima via UDP
		sendUDPNotification(buf[:n])
	}
}

// Server UDP untuk menerima notifikasi donasi
func handleUDP() {
	// Membuat koneksi UDP
	udpAddr, err := net.ResolveUDPAddr("udp", ":8082")
	if err != nil {
		log.Fatal("Error menyelesaikan alamat UDP:", err)
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("Error mendengarkan di UDP:", err)
	}
	defer udpConn.Close()

	fmt.Println("Server UDP mendengarkan di :8082")

	buf := make([]byte, 1024)
	for {
		n, addr, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Error menerima paket UDP:", err)
			continue
		}
		fmt.Printf("Menerima %s dari %s\n", string(buf[:n]), addr)
	}
}

// Mengirim notifikasi UDP (untuk donasi baru)
func sendUDPNotification(msg []byte) {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:8082")
	if err != nil {
		log.Println("Error menyelesaikan alamat UDP:", err)
		return
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Println("Error menghubungi UDP:", err)
		return
	}
	defer udpConn.Close()

	_, err = udpConn.Write(msg)
	if err != nil {
		log.Println("Error mengirim pesan UDP:", err)
	}
}
