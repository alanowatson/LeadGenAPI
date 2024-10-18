package handlers

import (
	"encoding/json"
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
    playlisters = make(map[int]models.Playlister)
    playlisterID = 1
    playlisterMutex sync.RWMutex
)

func GetPlaylisters(w http.ResponseWriter, r *http.Request) {
    log.Println("GetPlaylisters function called")

    paginationParams := pagination.GetPaginationParams(r)

    var playlisters []models.Playlister
    var totalItems int

    // Query to get total count
    err := db.DB.QueryRow("SELECT COUNT(*) FROM playlisters").Scan(&totalItems)
    if err != nil {
        log.Printf("Error getting total count: %v", err)
        util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving playlisters")
        return
    }
    log.Printf("Total playlisters in database: %d", totalItems)

    // Query to get paginated results
    offset := (paginationParams.Page - 1) * paginationParams.PerPage
    rows, err := db.DB.Query("SELECT playlisterid, spotifyuserid, curatorfullname, email FROM playlisters ORDER BY playlisterid LIMIT $1 OFFSET $2",
        paginationParams.PerPage, offset)
    if err != nil {
        log.Printf("Error querying playlisters: %v", err)
        util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving playlisters")
        return
    }
    defer rows.Close()

    for rows.Next() {
        var p models.Playlister
        if err := rows.Scan(&p.ID, &p.SpotifyUserID, &p.CuratorFullName, &p.Email); err != nil {
            log.Printf("Error scanning playlister row: %v", err)
            util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving playlisters")
            return
        }
        playlisters = append(playlisters, p)
    }

    if err = rows.Err(); err != nil {
        log.Printf("Error after scanning all rows: %v", err)
        util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving playlisters")
        return
    }

    totalPages := (totalItems + paginationParams.PerPage - 1) / paginationParams.PerPage

    response := map[string]interface{}{
        "data":        playlisters,
        "page":        paginationParams.Page,
        "per_page":    paginationParams.PerPage,
        "total_items": totalItems,
        "total_pages": totalPages,
    }

    log.Printf("Responding with %d playlisters", len(playlisters))
    util.RespondWithJSON(w, http.StatusOK, response)
}

func CreatePlaylister(w http.ResponseWriter, r *http.Request) {
    var playlister models.Playlister
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&playlister); err != nil {
        errors.HandleError(w, err, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    if err := validation.ValidateStruct(playlister); err != nil {
        errors.HandleError(w, err, http.StatusBadRequest, "Validation error")
        return
    }

    playlisterMutex.Lock()
    playlister.ID = playlisterID
    playlisters[playlisterID] = playlister
    playlisterID++
    playlisterMutex.Unlock()

    util.RespondWithJSON(w, http.StatusCreated, playlister)
}

func GetPlaylister(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        util.RespondWithError(w, http.StatusBadRequest, "Invalid playlister ID")
        return
    }

    playlisterMutex.RLock()
    playlister, found := playlisters[id]
    playlisterMutex.RUnlock()

    if !found {
        util.RespondWithError(w, http.StatusNotFound, "Playlister not found")
        return
    }

    util.RespondWithJSON(w, http.StatusOK, playlister)
}

func UpdatePlaylister(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        errors.HandleError(w, err, http.StatusBadRequest, "Invalid playlister ID")
        return
    }

    var playlister models.Playlister
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&playlister); err != nil {
        errors.HandleError(w, err, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    // Validate the updated playlister data
    if err := validation.ValidateStruct(playlister); err != nil {
        errors.HandleError(w, err, http.StatusBadRequest, "Validation error")
        return
    }

    playlisterMutex.Lock()
    defer playlisterMutex.Unlock()

    if _, found := playlisters[id]; !found {
        util.RespondWithError(w, http.StatusNotFound, "Playlister not found")
        return
    }

    playlister.ID = id
    playlisters[id] = playlister

    util.RespondWithJSON(w, http.StatusOK, playlister)
}
func DeletePlaylister(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        util.RespondWithError(w, http.StatusBadRequest, "Invalid playlister ID")
        return
    }

    playlisterMutex.Lock()
    defer playlisterMutex.Unlock()

    if _, found := playlisters[id]; !found {
        util.RespondWithError(w, http.StatusNotFound, "Playlister not found")
        return
    }

    delete(playlisters, id)

    util.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
