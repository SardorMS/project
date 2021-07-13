package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {

	// Уникальные данные для вашего приложения.
	// Их нужно установить в переменных окружения.
	client_id := "0bbb4ddc8b544c1393d5595175974de9"
	client_secret := "775eb523365c4614aa03c36802a3cbad"

	// Адрес перенаправления
	// адрес который Spotify будет делать перенавравлять для вашего приложения.
	redirectURI := url.PathEscape("http://localhost:8888/callback")

	// Верификационный код
	codeVerifier := RandString(43)

	// Код вызова.
	hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.StdEncoding.EncodeToString(hash[:])

	// Ещё один State код - для дополнительной защиты.
	state := RandString(16)

	item := map[string]interface{}{
		"ClientID":      client_id,
		"ClientSecret":  client_secret,
		"RedirectURI":   redirectURI,
		"CodeVerifier":  codeVerifier,
		"CodeChallenge": codeChallenge,
		"State":         state,
	}
	data, err := json.Marshal(item)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Вывод в json файл, для более удобного просмотра.
	err = os.WriteFile("pkce.json", data, 0777)
	if err != nil {
		log.Println(err)
		return
	}
}

// RandString - генерирует случайную строку в зависимости от кол-во символов.
func RandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789-~_.")
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
