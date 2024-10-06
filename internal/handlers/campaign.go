package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/alanowatson/LeadGenAPI/internal/errors"
	"github.com/alanowatson/LeadGenAPI/internal/models"
	"github.com/alanowatson/LeadGenAPI/internal/validation"
    "github.com/alanowatson/LeadGenAPI/internal/pagination"
	"github.com/alanowatson/LeadGenAPI/pkg/util"
	"github.com/gorilla/mux"
)

var (
    campaigns     = make(map[int]models.Campaign)
    campaignID    = 1
    campaignMutex sync.RWMutex
)

func GetCampaigns(w http.ResponseWriter, r *http.Request) {
    paginationParams := pagination.GetPaginationParams(r)

    campaignMutex.RLock()
    defer campaignMutex.RUnlock()

    campaignList := make([]models.Campaign, 0, len(campaigns))
    for _, campaign := range campaigns {
        campaignList = append(campaignList, campaign)
    }

    totalItems := len(campaignList)
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

    paginatedList := campaignList[start:end]

    response := map[string]interface{}{
        "data":        paginatedList,
        "page":        paginationParams.Page,
        "per_page":    paginationParams.PerPage,
        "total_items": totalItems,
        "total_pages": totalPages,
    }

    util.RespondWithJSON(w, http.StatusOK, response)
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
        errors.HandleError(w, err, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    if err := validation.ValidateStruct(campaign); err != nil {
        errors.HandleError(w, err, http.StatusBadRequest, "Validation error")
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
        errors.HandleError(w, err, http.StatusBadRequest, "Invalid campaign ID")
        return
    }

    var campaign models.Campaign
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&campaign); err != nil {
        errors.HandleError(w, err, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    if err := validation.ValidateStruct(campaign); err != nil {
        errors.HandleError(w, err, http.StatusBadRequest, "Validation error")
        return
    }

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
