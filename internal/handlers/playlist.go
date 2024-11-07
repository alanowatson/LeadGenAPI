package handlers

import (
	"encoding/json"
    "database/sql"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/alanowatson/LeadGenAPI/internal/db"
	"github.com/alanowatson/LeadGenAPI/internal/errors"
	"github.com/alanowatson/LeadGenAPI/internal/models"
	"github.com/alanowatson/LeadGenAPI/internal/pagination"
	"github.com/alanowatson/LeadGenAPI/internal/validation"
	"github.com/alanowatson/LeadGenAPI/pkg/util"
	"github.com/gorilla/mux"
)

var (
    playlists     = make(map[int]models.Playlist)
    playlistID    = 1
    playlistMutex sync.RWMutex
)

func GetPlaylists(w http.ResponseWriter, r *http.Request) {
    log.Println("GetPlaylists function called")

    paginationParams := pagination.GetPaginationParams(r)
    log.Printf("Pagination params: page=%d, per_page=%d", paginationParams.Page, paginationParams.PerPage)

    var totalItems int
    err := db.DB.QueryRow("SELECT COUNT(*) FROM playlists").Scan(&totalItems)
    if err != nil {
        log.Printf("Error getting total count: %v", err)
        util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving playlists")
        return
    }
    log.Printf("Total playlists in database: %d", totalItems)

    totalPages := (totalItems + paginationParams.PerPage - 1) / paginationParams.PerPage
    log.Printf("Total pages: %d", totalPages)

    if paginationParams.Page > totalPages && totalPages > 0 {
        log.Printf("Requested page %d exceeds total pages %d", paginationParams.Page, totalPages)
        util.RespondWithError(w, http.StatusNotFound, "Page not found")
        return
    }

    offset := (paginationParams.Page - 1) * paginationParams.PerPage

    query := `
        SELECT playlistid, playlisterid, playlistspotifyid, numberoffollowers, current_playlist_name, lastfollowercountdate, last_exposed
        FROM playlists
        ORDER BY playlistid
        LIMIT $1 OFFSET $2
    `
    log.Printf("Executing query: %s", query)
    rows, err := db.DB.Query(query, paginationParams.PerPage, offset)
    if err != nil {
        log.Printf("Error querying playlists: %v", err)
        util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving playlists")
        return
    }
    defer rows.Close()

       var playlists []models.Playlist
    for rows.Next() {
        var p models.Playlist
        err := rows.Scan(
            &p.ID,
            &p.PlaylisterId,
            &p.PlaylistSpotifyId,
            &p.NumberOfFollowers,
            &p.CurrentPlaylistName,
            &p.LastFollowerCountDate,
            &p.LastExposed,
        )
        if err != nil {
            log.Printf("Error scanning playlist row: %v", err)
            util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving playlists")
            return
        }
        playlists = append(playlists, p)
    }

    if err = rows.Err(); err != nil {
        log.Printf("Error after scanning all rows: %v", err)
        util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving playlists")
        return
    }

    log.Printf("Number of playlists retrieved: %d", len(playlists))

    response := map[string]interface{}{
        "data":        playlists,
        "page":        paginationParams.Page,
        "per_page":    paginationParams.PerPage,
        "total_items": totalItems,
        "total_pages": totalPages,
    }

    util.RespondWithJSON(w, http.StatusOK, response)
    log.Println("GetPlaylists function completed")
}

func GetPlaylist(w http.ResponseWriter, r *http.Request) {
    log.Println("GetPlaylist function called")

    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        log.Printf("Invalid playlist ID: %v", err)
        util.RespondWithError(w, http.StatusBadRequest, "Invalid playlist ID")
        return
    }

    log.Printf("Looking up playlist with ID: %d", id)

    query := `
        SELECT playlistid, playlisterid, playlistspotifyid, numberoffollowers,
               current_playlist_name, lastfollowercountdate, last_exposed
        FROM playlists
        WHERE playlistid = $1
    `

    var p models.Playlist
    err = db.DB.QueryRow(query, id).Scan(
        &p.ID,
        &p.PlaylisterId,
        &p.PlaylistSpotifyId,
        &p.NumberOfFollowers,
        &p.CurrentPlaylistName,
        &p.LastFollowerCountDate,
        &p.LastExposed,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("Playlist not found with ID: %d", id)
            util.RespondWithError(w, http.StatusNotFound, "Playlist not found")
            return
        }
        log.Printf("Error querying playlist: %v", err)
        util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving playlist")
        return
    }

    log.Printf("Successfully retrieved playlist with ID: %d", id)
    util.RespondWithJSON(w, http.StatusOK, p)
    log.Println("GetPlaylist function completed")
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
