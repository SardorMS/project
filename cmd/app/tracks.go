package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SardorMS/project/pkg/types"
)

// handleTrackCreate - метод для создания треков.
func (s *Server) handleTracksCreate(writer http.ResponseWriter, request *http.Request) {

	var item []*types.Track

	if err := json.NewDecoder(request.Body).Decode(&item); err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := s.tracksSvc.TracksCreate(request.Context(), item); err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		
		data := &types.Message{
			Status: "400 Bad Request",
			Message: "Something Went Wrong!",
		}
		
		respondJSON(writer, data)
	}

	data := &types.Message{
		Status: "200 OK",
		Message: "Good Done!",
	}
	respondJSON(writer, data)
}

// handleGetAllTracks - метод для вывода всех треков.
func (s *Server) handleGetAllTracks(writer http.ResponseWriter, request *http.Request) {

	result, err := s.tracksSvc.GetAllTracks(request.Context()); 
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		
		data := &types.Message{
			Status: "400 Bad Request",
			Message: "Something Went Wrong!",
		}
		
		respondJSON(writer, data)
	}

	respondJSON(writer, result)
}


// handleTracksRemove - метод для удаления всех треков.
func (s *Server) handleTracksRemove(writer http.ResponseWriter, request *http.Request) {

	if err := s.tracksSvc.TracksRemove(request.Context()); err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		
		data := &types.Message{
			Status: "400 Bad Request",
			Message: "Something Went Wrong!",
		}
		
		respondJSON(writer, data)
	}

	data := &types.Message{
		Status: "200 OK",
		Message: "Good Done!",
	}

	respondJSON(writer, data)
}