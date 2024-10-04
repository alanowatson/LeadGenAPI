package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "sync"

    "github.com/gorilla/mux"
    "github.com/alanowatson/LeadGenAPI/internal/models"
    "github.com/alanowatson/LeadGenAPI/pkg/util"
    "github.com/alanowatson/LeadGenAPI/internal/validation"

)

var (
    campaigns     = make(map[int]models.Campaign)
    campaignID    = 1
    campaignMutex sync.RWMutex
)

func GetCampaigns(w http.ResponseWriter, r *http.Request) {
    campaignMutex.RLock()
    defer campaignMutex.RUnlock()

    campaignList := make([]models.Campaign, 0, len(campaigns))
    for _, campaign := range campaigns {
        campaignList = append(campaignList, campaign)
    }

    util.RespondWithJSON(w, http.StatusOK, campaignList)
}

func GetCampaign(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        util.RespondWithError(w, http.StatusBadRequest, "Invalid campaign ID")
        return
    }

    campaignMutex.RLock()
    campaign, found := campaigns[id]
    campaignMutex.RUnlock()

    if !found {
        util.RespondWithError(w, http.StatusNotFound, "Campaign not found")
        return
    }

    util.RespondWithJSON(w, http.StatusOK, campaign)
}

func CreateCampaign(w http.ResponseWriter, r *http.Request) {
    var campaign models.Campaign
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&campaign); err != nil {
        util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    if err := validation.ValidateStruct(campaign); err != nil {
        util.RespondWithError(w, http.StatusBadRequest, err.Error())
        return
    }

    campaignMutex.Lock()
    campaign.ID = campaignID
    campaigns[campaignID] = campaign
    campaignID++
    campaignMutex.Unlock()

    util.RespondWithJSON(w, http.StatusCreated, campaign)
}

func UpdateCampaign(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        util.RespondWithError(w, http.StatusBadRequest, "Invalid campaign ID")
        return
    }

    var campaign models.Campaign
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&campaign); err != nil {
        util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    campaignMutex.Lock()
    defer campaignMutex.Unlock()

    if _, found := campaigns[id]; !found {
        util.RespondWithError(w, http.StatusNotFound, "Campaign not found")
        return
    }

    campaign.ID = id
    campaigns[id] = campaign

    util.RespondWithJSON(w, http.StatusOK, campaign)
}

func DeleteCampaign(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        util.RespondWithError(w, http.StatusBadRequest, "Invalid campaign ID")
        return
    }

    campaignMutex.Lock()
    defer campaignMutex.Unlock()

    if _, found := campaigns[id]; !found {
        util.RespondWithError(w, http.StatusNotFound, "Campaign not found")
        return
    }

    delete(campaigns, id)

    util.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
