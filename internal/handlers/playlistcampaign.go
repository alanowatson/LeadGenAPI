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
)

var (
    playlistCampaigns     = make(map[string]models.PlaylistCampaign)
    playlistCampaignMutex sync.RWMutex
)

func GetPlaylistCampaigns(w http.ResponseWriter, r *http.Request) {
    playlistCampaignMutex.RLock()
    defer playlistCampaignMutex.RUnlock()

    playlistCampaignList := make([]models.PlaylistCampaign, 0, len(playlistCampaigns))
    for _, pc := range playlistCampaigns {
        playlistCampaignList = append(playlistCampaignList, pc)
    }

    util.RespondWithJSON(w, http.StatusOK, playlistCampaignList)
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
        errors.HandleError(w, err, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    if err := validation.ValidateStruct(pc); err != nil {
        errors.HandleError(w, err, http.StatusBadRequest, "Validation error")
        return
    }

    // Additional validation to check if referenced Playlist and Campaign exist
    if _, exists := playlists[pc.PlaylistID]; !exists {
        util.RespondWithError(w, http.StatusBadRequest, "Referenced Playlist does not exist")
        return
    }
    if _, exists := campaigns[pc.CampaignID]; !exists {
        util.RespondWithError(w, http.StatusBadRequest, "Referenced Campaign does not exist")
        return
    }

    key := fmt.Sprintf("%d_%d", pc.PlaylistID, pc.CampaignID)

    playlistCampaignMutex.Lock()
    playlistCampaigns[key] = pc
    playlistCampaignMutex.Unlock()

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
