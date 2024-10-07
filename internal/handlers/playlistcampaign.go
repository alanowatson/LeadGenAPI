package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/alanowatson/LeadGenAPI/internal/errors"
	"github.com/alanowatson/LeadGenAPI/internal/models"
	"github.com/alanowatson/LeadGenAPI/internal/validation"
	"github.com/alanowatson/LeadGenAPI/pkg/util"
	"github.com/gorilla/mux"
    "github.com/alanowatson/LeadGenAPI/internal/pagination"
    "github.com/alanowatson/LeadGenAPI/internal/db"

)

var (
    playlistCampaigns     = make(map[string]models.PlaylistCampaign)
    playlistCampaignMutex sync.RWMutex
)

const MaxAllowedMessages = 3


func GetPlaylistCampaigns(w http.ResponseWriter, r *http.Request) {
    paginationParams := pagination.GetPaginationParams(r)

    playlistCampaignMutex.RLock()
    defer playlistCampaignMutex.RUnlock()

    playlistCampaignList := make([]models.PlaylistCampaign, 0, len(playlistCampaigns))
    for _, pc := range playlistCampaigns {
        playlistCampaignList = append(playlistCampaignList, pc)
    }

    totalItems := len(playlistCampaignList)
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

    paginatedList := playlistCampaignList[start:end]

    response := map[string]interface{}{
        "data":        paginatedList,
        "page":        paginationParams.Page,
        "per_page":    paginationParams.PerPage,
        "total_items": totalItems,
        "total_pages": totalPages,
    }

    util.RespondWithJSON(w, http.StatusOK, response)
}

func GetPlaylistCampaign(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    playlistID, _ := strconv.Atoi(vars["playlistId"])
    campaignID, _ := strconv.Atoi(vars["campaignId"])

    key := fmt.Sprintf("%d_%d", playlistID, campaignID)

    playlistCampaignMutex.RLock()
    pc, found := playlistCampaigns[key]
    playlistCampaignMutex.RUnlock()

    if !found {
        util.RespondWithError(w, http.StatusNotFound, "PlaylistCampaign not found")
        return
    }

    util.RespondWithJSON(w, http.StatusOK, pc)
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
