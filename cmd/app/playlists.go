package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SardorMS/project/cmd/app/middleware"
	"github.com/SardorMS/project/pkg/types"
	"github.com/gorilla/mux"
)

// handlePlaylistCreate - метод для создания плейлиста.
func (s *Server) handlePlaylistCreate(writer http.ResponseWriter, request *http.Request) {

	id, err := middleware.Authentication(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	userID, ok := mux.Vars(request)["user_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item *types.UserPlaylist

	if err = json.NewDecoder(request.Body).Decode(&item); err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := s.playlistSvc.PlaylistCreate(request.Context(), item, id, userID)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	respondJSON(writer, result)

}

// handleGetUserPlaylists - метод для получения всех плейлистов текущего пользователя.
func (s *Server) handleGetUserPlaylists(writer http.ResponseWriter, request *http.Request) {

	id, err := middleware.Authentication(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	userID, ok := mux.Vars(request)["user_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := s.playlistSvc.GetAllPlaylists(request.Context(), id, userID)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	respondJSON(writer, result)
}

// handleRemoveAllUserPlaylists - метод для удаления всех плейлистов текущего пользователя.
func (s *Server) handleRemoveAllUserPlaylists(writer http.ResponseWriter, request *http.Request) {

	id, err := middleware.Authentication(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	userID, ok := mux.Vars(request)["user_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = s.playlistSvc.RemoveAllPlaylists(request.Context(), id, userID)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	item := map[string]interface{}{"code": "200 OK", "state": "All Playlists Removed"}
	respondJSON(writer, item)
}

// handleGetUserPlaylist - метод для получения публичного плейлиста по его playlist_id.
func (s *Server) handleGetUserPlaylist(writer http.ResponseWriter, request *http.Request) {

	id, err := middleware.Authentication(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	playlistID, ok := mux.Vars(request)["playlist_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := s.playlistSvc.GetPlaylist(request.Context(), id, playlistID)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	respondJSON(writer, result)
}

// handlePlaylistUploadImage - метод для загрузки изображения для конкретного плейлиста.
func (s *Server) handlePlaylistUploadImage(writer http.ResponseWriter, request *http.Request) {

	id, err := middleware.Authentication(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	playlistID, ok := mux.Vars(request)["playlist_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item *types.Image

	if err = json.NewDecoder(request.Body).Decode(&item); err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, message, err := s.playlistSvc.PlaylistUploadImage(request.Context(), item, id, playlistID)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if message == "1" {
		data := &types.Message{
			Status:  "400 Bad Request!",
			Message: "You Already Have An Image. Please Remove It First. Thank You!",
		}
		respondJSON(writer, data)
	} else {
		respondJSON(writer, result)
	}
}

// handlePlaylistImageRemove - метод для удаления изображения для конкретного плейлиста.
func (s *Server) handlePlaylistImageRemove(writer http.ResponseWriter, request *http.Request) {

	id, err := middleware.Authentication(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	playlistID, ok := mux.Vars(request)["playlist_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := s.playlistSvc.PlaylistImageRemove(request.Context(), id, playlistID)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(writer, result)
}

// handlePlaylistChange - метод для внесения изменений в детали, для конкретного плейлиста.
func (s *Server) handlePlaylistChange(writer http.ResponseWriter, request *http.Request) {

	id, err := middleware.Authentication(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	playlistID, ok := mux.Vars(request)["playlist_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item *types.UserPlaylist

	if err = json.NewDecoder(request.Body).Decode(&item); err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := s.playlistSvc.PlaylistChange(request.Context(), item, id, playlistID)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	respondJSON(writer, result)
}

// handleTrackToPlaylist - метод для загрузки трека в конкретный плейлист.
func (s *Server) handleTrackToPlaylist(writer http.ResponseWriter, request *http.Request) {
	id, err := middleware.Authentication(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	playlistID, ok := mux.Vars(request)["playlist_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item *types.TrackInfo

	if err = json.NewDecoder(request.Body).Decode(&item); err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := s.playlistSvc.TrackToPlaylist(request.Context(), item, id, playlistID)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	respondJSON(writer, result)

}

// handleRemoveTrackFromPlaylist - метод для удаления треки из конкретного плейлиста.
func (s *Server) handleRemoveTrackFromPlaylist(writer http.ResponseWriter, request *http.Request) {
	id, err := middleware.Authentication(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	playlistID, ok := mux.Vars(request)["playlist_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item *types.TrackInfo

	if err = json.NewDecoder(request.Body).Decode(&item); err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := s.playlistSvc.RemoveTrackFromPlaylist(request.Context(), item, id, playlistID)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	respondJSON(writer, result)

}

// handleRemoveAllTrackFromPlaylist - метод для удаления всех треков из конкретного плейлиста.
func (s *Server) handleRemoveAllTracksFromPlaylist(writer http.ResponseWriter, request *http.Request) {
	id, err := middleware.Authentication(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	playlistID, ok := mux.Vars(request)["playlist_id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	result, err := s.playlistSvc.RemoveAllTracksFromPlaylist(request.Context(), id, playlistID)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	respondJSON(writer, result)

}
