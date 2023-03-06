package utils

import (
	"encoding/json"
	"github.com/labstack/gommon/random"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

// IP finds public IP address.
func IP() string {
	ip := "127.0.0.1"

	client := http.Client{Timeout: 5 * time.Second}
	response, err := client.Get("https://api.ipify.org?format=json")
	if err != nil || response.StatusCode != http.StatusOK {
		return ip
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return ip
	}

	var s struct {
		IP string `json:"ip"`
	}
	if err = json.Unmarshal(body, &s); err != nil {
		return ip
	}

	return s.IP
}

// Token generates a random token
func Token() string {
	return random.String(32)
}

// FreePort finds a free port.
func FreePort() (int, error) {
	address, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}

	defer func() {
		err = listener.Close()
	}()

	return listener.Addr().(*net.TCPAddr).Port, err
}

// DirectoryExist checks if the given directory path exists or not.
func DirectoryExist(path string) bool {
	if stat, err := os.Stat(path); os.IsNotExist(err) || !stat.IsDir() {
		return false
	}
	return true
}
