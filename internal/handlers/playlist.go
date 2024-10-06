package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/alanowatson/LeadGenAPI/internal/errors"
	"github.com/alanowatson/LeadGenAPI/internal/models"
	"github.com/alanowatson/LeadGenAPI/internal/validation"
	"github.com/alanowatson/LeadGenAPI/pkg/util"
	"github.com/gorilla/mux"
    "github.com/alanowatson/LeadGenAPI/internal/pagination"

)

var (
    playlists     = make(map[int]models.Playlist)
    playlistID    = 1
    playlistMutex sync.RWMutex
)

func GetPlaylists(w http.ResponseWriter, r *http.Request) {
    paginationParams := pagination.GetPaginationParams(r)

    playlistMutex.RLock()
    defer playlistMutex.RUnlock()

    playlistList := make([]models.Playlist, 0, len(playlists))
    for _, playlist := range playlists {
        playlistList = append(playlistList, playlist)
    }

    totalItems := len(playlistList)
    totalPages := (totalItems + paginationParams.PerPage - 1) / paginationParams.PerPage

    if paginationParams.Page > totalPages {
        util.RespondWithError(w, http.StatusNotFound, "Page not found")
        return
    }

    start := (paginationParams.Page - 1) * paginationParams.PerPage
    end := start + paginationParams.PerPage
    if end > totalItems {
        end = totalItems
    }

    paginatedList := playlistList[start:end]

    response := map[string]interface{}{
        "data":        paginatedList,
        "page":        paginationParams.Page,
        "per_page":    paginationParams.PerPage,
        "total_items": totalItems,
        "total_pages": totalPages,
    }

    util.RespondWithJSON(w, http.StatusOK, response)
}

func GetPlaylist(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        util.RespondWithError(w, http.StatusBadRequest, "Invalid playlist ID")
        return
    }

    playlistMutex.RLock()
    playlist, found := playlists[id]
    playlistMutex.RUnlock()

    if !found {
        util.RespondWithError(w, http.StatusNotFound, "Playlist not found")
        return
    }

    util.RespondWithJSON(w, http.StatusOK, playlist)
}

func CreatePlaylist(w http.ResponseWriter, r *http.Request) {
    var playlist models.Playlist
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&playlist); err != nil {
        errors.HandleError(w, err, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    if err := validation.ValidateStruct(playlist); err != nil {
        errors.HandleError(w, err, http.StatusBadRequest, "Validation error")
        return
    }

    playlistMutex.Lock()
    playlist.ID = playlistID
    playlists[playlistID] = playlist
    playlistID++
    playlistMutex.Unlock()

    util.RespondWithJSON(w, http.StatusCreated, playlist)
}

func UpdatePlaylist(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        errors.HandleError(w, err, http.StatusBadRequest, "Invalid playlist ID")
        return
    }

    var playlist models.Playlist
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&playlist); err != nil {
        errors.HandleError(w, err, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    if err := validation.ValidateStruct(playlist); err != nil {
        errors.HandleError(w, err, http.StatusBadRequest, "Validation error")
        return
    }

    playlistMutex.Lock()
    defer playlistMutex.Unlock()

    if _, found := playlists[id]; !found {
        util.RespondWithError(w, http.StatusNotFound, "Playlist not found")
        return
    }

    playlist.ID = id
    playlists[id] = playlist

    util.RespondWithJSON(w, http.StatusOK, playlist)
}

func DeletePlaylist(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        util.RespondWithError(w, http.StatusBadRequest, "Invalid playlist ID")
        return
    }

    playlistMutex.Lock()
    defer playlistMutex.Unlock()

    if _, found := playlists[id]; !found {
        util.RespondWithError(w, http.StatusNotFound, "Playlist not found")
        return
    }

    delete(playlists, id)

    util.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
