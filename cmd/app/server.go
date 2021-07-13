package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SardorMS/project/cmd/app/middleware"
	"github.com/SardorMS/project/pkg/playlists"
	"github.com/SardorMS/project/pkg/spotify"
	"github.com/SardorMS/project/pkg/tracks"
	"github.com/SardorMS/project/pkg/users"
	"github.com/gorilla/mux"
)

var (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
)

// Server - представляет собой логический сервер приложения.
type Server struct {
	mux         *mux.Router
	spotSvc     *spotify.Service
	userSvc     *users.Service
	playlistSvc *playlists.Service
	tracksSvc   *tracks.Service
}

// NewServer - функция констркутор для создание сервера.
func NewServer(mux *mux.Router, spotSvc *spotify.Service, userSvc *users.Service,
	playlistSvc *playlists.Service, tracksSvc *tracks.Service) *Server {
	return &Server{
		mux:         mux,
		spotSvc:     spotSvc,
		userSvc:     userSvc,
		playlistSvc: playlistSvc,
		tracksSvc:   tracksSvc,
	}
}

// ServeHTTP - метод для запуска сервера.
func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

// Init - инициализация сервера (регистрация и обработка запросов).
func (s *Server) Init() {

	s.mux.Use(middleware.Logger)

	// Запрос на подтверждение
	// accounts.spotify.com/authorize
	s.mux.HandleFunc("/authorize", s.handleVerify).Methods(GET)

	// Перенаправление на регистрацию, адрес который вы регистрируете в Dashborde.
	// localhost:8888/callback
	s.mux.HandleFunc("/callback", s.handleUserCheck).Methods(GET)

	// accounts.spotify.com/api/token
	s.mux.HandleFunc("/api/token", s.handleTokenExchange).Methods(POST)

	// accounts.spotify.com/api/token/refresh
	s.mux.HandleFunc("/api/token/{refresh}", s.handleTokenRefresh).Methods(POST)

	// accounts.spotify.com/{region}/signup
	s.mux.HandleFunc("/{region}/signup", s.handleUserRegistration).Methods(POST)

	// accounts.spotify.com/ru/login/
	s.mux.HandleFunc("/{region}/login", s.handleUserLogin).Methods(POST)

	userMd := middleware.Authenticate(s.userSvc.IDByToken)

	//                              Юзеры.
	// api.spotify.com/
	usersSubrouter := s.mux.PathPrefix("/v1/").Subrouter()
	usersSubrouter.Use(userMd)
	usersSubrouter.HandleFunc("/me", s.handleGetCurrentUser).Methods(GET)
	usersSubrouter.HandleFunc("/users/{user_id}", s.handleGetUserByID).Methods(GET)
	usersSubrouter.HandleFunc("/me/image", s.handleUserImageUpload).Methods(POST)
	usersSubrouter.HandleFunc("/me/image", s.handleUserImageRemove).Methods(GET)

	//                               Плейлисты.
	// api.spotify.com/
	playSubrouter := s.mux.PathPrefix("/v1/").Subrouter()
	playSubrouter.Use(userMd)
	playSubrouter.HandleFunc("/users/{user_id}/playlists", s.handlePlaylistCreate).Methods(POST)
	playSubrouter.HandleFunc("/users/{user_id}/playlists", s.handleGetUserPlaylists).Methods(GET)
	playSubrouter.HandleFunc("/users/{user_id}/playlists", s.handleRemoveAllUserPlaylists).Methods(DELETE)
	playSubrouter.HandleFunc("/playlists/{playlist_id}", s.handleGetUserPlaylist).Methods(GET)
	playSubrouter.HandleFunc("/playlists/{playlist_id}/images", s.handlePlaylistUploadImage).Methods(PUT)
	playSubrouter.HandleFunc("/playlists/{playlist_id}/images", s.handlePlaylistImageRemove).Methods(DELETE)
	playSubrouter.HandleFunc("/playlists/{playlist_id}", s.handlePlaylistChange).Methods(PUT)
	playSubrouter.HandleFunc("/playlists/{playlist_id}/tracks", s.handleTrackToPlaylist).Methods(POST)
	playSubrouter.HandleFunc("/playlists/{playlist_id}/tracks", s.handleRemoveTrackFromPlaylist).Methods(PUT)
	playSubrouter.HandleFunc("/playlists/{playlist_id}/tracks", s.handleRemoveAllTracksFromPlaylist).Methods(DELETE)

	//                                Треки.
	s.mux.HandleFunc("/v1/tracks", s.handleTracksCreate).Methods(POST)
	s.mux.HandleFunc("/v1/tracks", s.handleGetAllTracks).Methods(GET)
	s.mux.HandleFunc("/v1/tracks", s.handleTracksRemove).Methods(DELETE)
}

// RequestURL - возвращает занчение по ключу query.
func RequestURL(r *http.Request, param string) string {
	return r.URL.Query().Get(param)
}

// respondJSON - Ответ в виде JSON.
func respondJSON(w http.ResponseWriter, item interface{}) {

	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
	}
}
