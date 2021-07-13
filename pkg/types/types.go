package types

import "time"

//=============================================================
//                  Структуры для Application

// DataApp - ...
type DataApp struct {
	ClientID      string   `json:"client_id"`
	RedirectURI   string   `json:"redirect_uri"`
	CodeChallenge string   `json:"code_challenge"`
	CodeVerifier  string   `json:"code_verifier"`
	State         string   `json:"state"`
	Scope         []string `json:"scope"`
	Code          string   `json:"code"`
}

// Application - ...
type Application struct {
	ID           int64     `json:"id"`
	ClientID     string    `json:"client_id"`
	RedirectURI  string    `json:"redirect_uri"`
	State        string    `json:"state"`
	Scope        []string  `json:"scope"`
	CodeVerifier string    `json:"code_verifier"`
	Code         string    `json:"code"`
	Created      time.Time `json:"created"`
}

// Token - ...
type Token struct {
	ClientID     string `json:"client_id"`
	RedirectURI  string `json:"redirect_uri"`
	Code         string `json:"code"`
	CodeVerifier string `json:"code_verifier"`
}

// Tokens - ...
type Tokens struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	Scope        []string  `json:"scope"`
	ExpiresIn    int64     `json:"expires_in"`
	RefreshToken string    `json:"refresh_token"`
	Created      time.Time `json:"created"`
}

// TokenRefresh - ...
type TokenRefresh struct {
	RefreshToken string `json:"refresh_token"`
}

// Token - ...
type TokenUpdate struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   int64     `json:"expires_in"`
	Scope       []string  `json:"scope"`
	Created     time.Time `json:"created"`
}

//========================================================================
//                          Структуры для USER

// Registration - ...
type Registration struct {
	AccessToken string `json:"access_token"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Product     string `json:"product"`
	Birthdate   string `json:"birthdate"`
}

// UserPrivate - ...
type UserPrivate struct {
	ID          int64     `json:"id"`
	Country     string    `json:"country"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	Endpoint    string    `json:"href"`
	UserID      string    `json:"user_id"`
	Images      []Image   `json:"images"`
	Product     string    `json:"product"`
	Birthdate   string    `json:"birthdate"`
	URI         URI       `json:"uri"`
	Created     time.Time `json:"created"`
}

// Image - ...
type Image struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// URI - ...
type URI string

// Auth - ...
type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login - ...
type Login struct {
	ID          int64  `json:"id"`
	Country     string `json:"country"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	UserID      string `json:"user_id"`
	Product     string `json:"product"`
}

// AnyUser - ...
type AnyUser struct {
	DisplayName string  `json:"display_name"`
	Endpoint    string  `json:"href"`
	UserID      string  `json:"user_id"`
	Images      []Image `json:"images"`
	URI         URI     `json:"uri"`
}

// ImageUpdate - ...
type ImageUpdate struct {
	ID          int64     `json:"id"`
	Country     string    `json:"country"`
	DisplayName string    `json:"display_name"`
	Endpoint    string    `json:"href"`
	UserID      string    `json:"user_id"`
	Images      []Image   `json:"images"`
	URI         URI       `json:"uri"`
	Product     string    `json:"product"`
	Created     time.Time `json:"created"`
}

//==================================================================================
//                             Плейлитсы и треки

// UserPlaylist - ...
type UserPlaylist struct {
	Name         string  `json:"name"`
	Descriptions string  `json:"descriptions"`
	IsPublic     bool    `json:"is_public"`
	Images       []Image `json:"images"`
}

// Playlists - ...
type Playlists struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	OwnerName    string    `json:"owner_name"`
	Descriptions string    `json:"descriptions"`
	IsPublic     bool      `json:"is_public"`
	Images       []Image   `json:"images"`
	Tracks       []Track   `json:"tracks"`
	PlaylistID   string    `json:"playlist_id"`
	Endpoint     string    `json:"href"`
	URI          URI       `json:"uri"`
	Created      time.Time `json:"created"`
}

// PlaylistImage - для проверки фотографии, есть она или нет.
type ImageChecker struct {
	Images []Image `json:"images"`
}

// Message - структура для ответов во время проферки, когда выдало ошибку.
type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// TrackInfo - ...
type TrackInfo struct {
	Name       string `json:"name"`
	ArtistName string `json:"artist_name"`
	URI        URI    `json:"uri"`
}

// Track - ...
type Track struct {
	ID          int64
	TrackID     string    `json:"track_id"`
	Name        string    `json:"name"`
	ArtistName  string    `json:"artist_name"`
	AlbumName   string    `json:"album_name"`
	TrackNumber int64     `json:"track_number"`
	Genres      []string  `json:"genres"`
	Duration    int64     `json:"duration"`
	Endpoint    string    `json:"href"`
	URI         URI       `json:"uri"`
	Created     time.Time `json:"created"`
}
