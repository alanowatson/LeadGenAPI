package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
    playlistCampaigns     = make(map[string]models.PlaylistCampaign)
    playlistCampaignMutex sync.RWMutex
)

const MaxAllowedMessages = 3


func GetPlaylistCampaigns(w http.ResponseWriter, r *http.Request) {
    log.Println("GetPlaylistCampaigns function called")

    paginationParams := pagination.GetPaginationParams(r)
    log.Printf("Pagination params: page=%d, per_page=%d", paginationParams.Page, paginationParams.PerPage)

    var totalItems int
    err := db.DB.QueryRow("SELECT COUNT(*) FROM playlistcampaigns").Scan(&totalItems)
    if err != nil {
        log.Printf("Error getting total count: %v", err)
        util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving playlist campaigns")
        return
    }
    log.Printf("Total playlist campaigns in database: %d", totalItems)

    totalPages := (totalItems + paginationParams.PerPage - 1) / paginationParams.PerPage
    log.Printf("Total pages: %d", totalPages)

    if paginationParams.Page > totalPages && totalPages > 0 {
        log.Printf("Requested page %d exceeds total pages %d", paginationParams.Page, totalPages)
        util.RespondWithError(w, http.StatusNotFound, "Page not found")
        return
    }

    offset := (paginationParams.Page - 1) * paginationParams.PerPage

    query := `
        SELECT playlistid, campaignid, playlisterid, referenceartists, placementstatus, numberofmessages, purchased
        FROM playlistcampaigns
        ORDER BY playlistid, campaignid
        LIMIT $1 OFFSET $2
    `
    log.Printf("Executing query: %s", query)
    rows, err := db.DB.Query(query, paginationParams.PerPage, offset)
    if err != nil {
        log.Printf("Error querying playlist campaigns: %v", err)
        util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving playlist campaigns")
        return
    }
    defer rows.Close()

    var playlistCampaigns []models.PlaylistCampaign
    for rows.Next() {
        var pc models.PlaylistCampaign
        err := rows.Scan(
            &pc.PlaylistID,
            &pc.CampaignID,
            &pc.PlaylisterId,
            &pc.ReferenceArtists,
            &pc.PlacementStatus,
            &pc.NumberOfMessages,
            &pc.Purchased,
        )
        if err != nil {
            log.Printf("Error scanning playlist campaign row: %v", err)
            util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving playlist campaigns")
            return
        }
        playlistCampaigns = append(playlistCampaigns, pc)
    }

    if err = rows.Err(); err != nil {
        log.Printf("Error after scanning all rows: %v", err)
        util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving playlist campaigns")
        return
    }

    log.Printf("Number of playlist campaigns retrieved: %d", len(playlistCampaigns))

    response := map[string]interface{}{
        "data":        playlistCampaigns,
        "page":        paginationParams.Page,
        "per_page":    paginationParams.PerPage,
        "total_items": totalItems,
        "total_pages": totalPages,
    }

    log.Printf("Response prepared with %d items", len(playlistCampaigns))
    util.RespondWithJSON(w, http.StatusOK, response)
    log.Println("GetPlaylistCampaigns function completed")
}

func GetPlaylistCampaign(w http.ResponseWriter, r *http.Request) {
    log.Println("GetPlaylistCampaign function called")

    vars := mux.Vars(r)
    playlistID, err := strconv.Atoi(vars["playlistId"])
    if err != nil {
        log.Printf("Invalid playlist ID: %v", err)
        util.RespondWithError(w, http.StatusBadRequest, "Invalid playlist ID")
        return
    }

    campaignID, err := strconv.Atoi(vars["campaignId"])
    if err != nil {
        log.Printf("Invalid campaign ID: %v", err)
        util.RespondWithError(w, http.StatusBadRequest, "Invalid campaign ID")
        return
    }

    log.Printf("Looking up playlist campaign with PlaylistID: %d and CampaignID: %d", playlistID, campaignID)

    query := `
        SELECT playlistid, campaignid, playlisterid, referenceartists,
               placementstatus, numberofmessages, purchased
        FROM playlistcampaigns
        WHERE playlistid = $1 AND campaignid = $2
    `

    var pc models.PlaylistCampaign
    err = db.DB.QueryRow(query, playlistID, campaignID).Scan(
        &pc.PlaylistID,
        &pc.CampaignID,
        &pc.PlaylisterId,
        &pc.ReferenceArtists,
        &pc.PlacementStatus,
        &pc.NumberOfMessages,
        &pc.Purchased,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("PlaylistCampaign not found with PlaylistID: %d and CampaignID: %d", playlistID, campaignID)
            util.RespondWithError(w, http.StatusNotFound, "PlaylistCampaign not found")
            return
        }
        log.Printf("Error querying playlist campaign: %v", err)
        util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving playlist campaign")
        return
    }

    log.Printf("Successfully retrieved playlist campaign with PlaylistID: %d and CampaignID: %d", playlistID, campaignID)
    util.RespondWithJSON(w, http.StatusOK, pc)
    log.Println("GetPlaylistCampaign function completed")
}

func CreatePlaylistCampaign(w http.ResponseWriter, r *http.Request) {
    var pc models.PlaylistCampaign
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&pc); err != nil {
        util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    if err := validation.ValidateStruct(pc); err != nil {
        util.RespondWithError(w, http.StatusBadRequest, "Validation error: " + err.Error())
        return
    }

    if pc.NumberOfMessages > MaxAllowedMessages {
        util.RespondWithError(w, http.StatusBadRequest, "Exceeds maximum allowed messages")
        return
    }

    playlistExists, err := db.PlaylistExists(pc.PlaylistID)
    if err != nil {
        util.RespondWithError(w, http.StatusInternalServerError, "Error checking playlist existence")
        return
    }
    if !playlistExists {
        util.RespondWithError(w, http.StatusBadRequest, "Referenced Playlist does not exist")
        return
    }

    campaignExists, err := db.CampaignExists(pc.CampaignID)
    if err != nil {
        util.RespondWithError(w, http.StatusInternalServerError, "Error checking campaign existence")
        return
    }
    if !campaignExists {
        util.RespondWithError(w, http.StatusBadRequest, "Referenced Campaign does not exist")
        return
    }

    tx, err := db.BeginTx()
    if err != nil {
        util.RespondWithError(w, http.StatusInternalServerError, "Error starting transaction")
        return
    }
    defer tx.Rollback()

    if err := db.CreatePlaylistCampaignTx(tx, pc); err != nil {
        util.RespondWithError(w, http.StatusInternalServerError, "Error creating PlaylistCampaign")
        return
    }

    if err := tx.Commit(); err != nil {
        util.RespondWithError(w, http.StatusInternalServerError, "Error committing transaction")
        return
    }

    util.RespondWithJSON(w, http.StatusCreated, pc)
}

func UpdatePlaylistCampaign(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    playlistID, _ := strconv.Atoi(vars["playlistId"])
    campaignID, _ := strconv.Atoi(vars["campaignId"])

    key := fmt.Sprintf("%d_%d", playlistID, campaignID)

    var pc models.PlaylistCampaign
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&pc); err != nil {
        errors.HandleError(w, err, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    if err := validation.ValidateStruct(pc); err != nil {
        errors.HandleError(w, err, http.StatusBadRequest, "Validation error")
        return
    }

    playlistCampaignMutex.Lock()
    defer playlistCampaignMutex.Unlock()

    if _, found := playlistCampaigns[key]; !found {
        util.RespondWithError(w, http.StatusNotFound, "PlaylistCampaign not found")
        return
    }

    // Ensure the IDs in the URL match the IDs in the payload
    if pc.PlaylistID != playlistID || pc.CampaignID != campaignID {
        util.RespondWithError(w, http.StatusBadRequest, "Playlist ID and Campaign ID in URL must match payload")
        return
    }

    playlistCampaigns[key] = pc

    util.RespondWithJSON(w, http.StatusOK, pc)
}

func DeletePlaylistCampaign(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    playlistID, _ := strconv.Atoi(vars["playlistId"])
    campaignID, _ := strconv.Atoi(vars["campaignId"])

    key := fmt.Sprintf("%d_%d", playlistID, campaignID)

    playlistCampaignMutex.Lock()
    defer playlistCampaignMutex.Unlock()

    if _, found := playlistCampaigns[key]; !found {
        util.RespondWithError(w, http.StatusNotFound, "PlaylistCampaign not found")
        return
    }

    delete(playlistCampaigns, key)

    util.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
