package handlers

import (
	"database/sql"
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
   log.Printf("Pagination params: page=%d, per_page=%d", paginationParams.Page, paginationParams.PerPage)

   var totalItems int
   err := db.DB.QueryRow("SELECT COUNT(*) FROM playlisters").Scan(&totalItems)
   if err != nil {
       log.Printf("Error getting total count: %v", err)
       util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving playlisters")
       return
   }
   log.Printf("Total playlisters in database: %d", totalItems)

   totalPages := (totalItems + paginationParams.PerPage - 1) / paginationParams.PerPage
   log.Printf("Total pages: %d", totalPages)

   if paginationParams.Page > totalPages && totalPages > 0 {
       log.Printf("Requested page %d exceeds total pages %d", paginationParams.Page, totalPages)
       util.RespondWithError(w, http.StatusNotFound, "Page not found")
       return
   }

   offset := (paginationParams.Page - 1) * paginationParams.PerPage

   query := `
       SELECT playlisterid, spotifyuserid, curatorfullname, email,
              instagram, facebook, whatsapp, lastcontacted,
              preferredlanguage, followupstatus
       FROM playlisters
       ORDER BY playlisterid
       LIMIT $1 OFFSET $2
   `
   log.Printf("Executing query: %s", query)
   rows, err := db.DB.Query(query, paginationParams.PerPage, offset)
   if err != nil {
       log.Printf("Error querying playlisters: %v", err)
       util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving playlisters")
       return
   }
   defer rows.Close()

   var playlisters []models.Playlister
   for rows.Next() {
       var p models.Playlister
       err := rows.Scan(
           &p.ID,
           &p.SpotifyUserID,
           &p.CuratorFullName,
           &p.Email,
           &p.Instagram,
           &p.Facebook,
           &p.Whatsapp,
           &p.LastContacted,
           &p.PreferredLanguage,
           &p.FollowupStatus,
       )
       if err != nil {
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

   log.Printf("Number of playlisters retrieved: %d", len(playlisters))

   response := map[string]interface{}{
       "data":        playlisters,
       "page":        paginationParams.Page,
       "per_page":    paginationParams.PerPage,
       "total_items": totalItems,
       "total_pages": totalPages,
   }

   log.Printf("Response prepared with %d items", len(playlisters))
   util.RespondWithJSON(w, http.StatusOK, response)
   log.Println("GetPlaylisters function completed")
}

func GetPlaylister(w http.ResponseWriter, r *http.Request) {
    log.Println("GetPlaylister function called")

    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        log.Printf("Invalid playlister ID: %v", err)
        util.RespondWithError(w, http.StatusBadRequest, "Invalid playlister ID")
        return
    }
    log.Printf("Looking up playlister with ID: %d", id)

    query := `
        SELECT playlisterid, spotifyuserid, curatorfullname, email,
               instagram, facebook, whatsapp, lastcontacted,
               preferredlanguage, followupstatus
        FROM playlisters
        WHERE playlisterid = $1
    `

    var p models.Playlister
    err = db.DB.QueryRow(query, id).Scan(
        &p.ID,
        &p.SpotifyUserID,
        &p.CuratorFullName,
        &p.Email,
        &p.Instagram,
        &p.Facebook,
        &p.Whatsapp,
        &p.LastContacted,
        &p.PreferredLanguage,
        &p.FollowupStatus,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("Playlister not found with ID: %d", id)
            util.RespondWithError(w, http.StatusNotFound, "Playlister not found")
            return
        }
        log.Printf("Error querying playlister: %v", err)
        util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving playlister")
        return
    }

    log.Printf("Successfully retrieved playlister with ID: %d", id)
    util.RespondWithJSON(w, http.StatusOK, p)
    log.Println("GetPlaylister function completed")
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
