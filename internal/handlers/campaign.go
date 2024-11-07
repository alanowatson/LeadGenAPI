package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
    "database/sql"

	"github.com/alanowatson/LeadGenAPI/internal/db"
	"github.com/alanowatson/LeadGenAPI/internal/errors"
	"github.com/alanowatson/LeadGenAPI/internal/models"
	"github.com/alanowatson/LeadGenAPI/internal/pagination"
	"github.com/alanowatson/LeadGenAPI/internal/validation"

	"github.com/alanowatson/LeadGenAPI/pkg/util"
	"github.com/gorilla/mux"
)

var (
    campaigns     = make(map[int]models.Campaign)
    campaignID    = 1
    campaignMutex sync.RWMutex
)

func GetCampaigns(w http.ResponseWriter, r *http.Request) {
    log.Println("GetCampaigns function called")

    paginationParams := pagination.GetPaginationParams(r)
    log.Printf("Pagination params: page=%d, per_page=%d", paginationParams.Page, paginationParams.PerPage)

    var totalItems int
    err := db.DB.QueryRow("SELECT COUNT(*) FROM campaigns").Scan(&totalItems)
    if err != nil {
        log.Printf("Error getting total count: %v", err)
        util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving campaigns")
        return
    }
    log.Printf("Total campaigns in database: %d", totalItems)

    totalPages := (totalItems + paginationParams.PerPage - 1) / paginationParams.PerPage
    log.Printf("Total pages: %d", totalPages)

    if paginationParams.Page > totalPages && totalPages > 0 {
        log.Printf("Requested page %d exceeds total pages %d", paginationParams.Page, totalPages)
        util.RespondWithError(w, http.StatusNotFound, "Page not found")
        return
    }

    offset := (paginationParams.Page - 1) * paginationParams.PerPage

    query := `
        SELECT campaignid, campaignname, referenceartists, trello_link, spotify_link, launchdate, promoted_artist
        FROM campaigns
        ORDER BY campaignid
        LIMIT $1 OFFSET $2
    `
    log.Printf("Executing query: %s", query)
    rows, err := db.DB.Query(query, paginationParams.PerPage, offset)
    if err != nil {
        log.Printf("Error querying campaigns: %v", err)
        util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving campaigns")
        return
    }
    defer rows.Close()

    var campaigns []models.Campaign
    for rows.Next() {
        var c models.Campaign
        err := rows.Scan(
            &c.ID,
            &c.CampaignName,
            &c.ReferenceArtists,
            &c.TrelloLink,
            &c.SpotifyLink,
            &c.LaunchDate,
            &c.PromotedArtist,
        )
        if err != nil {
            log.Printf("Error scanning campaign row: %v", err)
            util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving campaigns")
            return
        }
        campaigns = append(campaigns, c)
    }

    if err = rows.Err(); err != nil {
        log.Printf("Error after scanning all rows: %v", err)
        util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving campaigns")
        return
    }

    log.Printf("Number of campaigns retrieved: %d", len(campaigns))

    response := map[string]interface{}{
        "data":        campaigns,
        "page":        paginationParams.Page,
        "per_page":    paginationParams.PerPage,
        "total_items": totalItems,
        "total_pages": totalPages,
    }

    log.Printf("Response prepared: %+v", response)
    util.RespondWithJSON(w, http.StatusOK, response)
    log.Println("GetCampaigns function completed")
}

func GetCampaign(w http.ResponseWriter, r *http.Request) {
    log.Println("GetCampaign function called")

    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        log.Printf("Invalid campaign ID: %v", err)
        util.RespondWithError(w, http.StatusBadRequest, "Invalid campaign ID")
        return
    }

    log.Printf("Looking up campaign with ID: %d", id)

    query := `
        SELECT campaignid, campaignname, referenceartists, trello_link, spotify_link,
               launchdate, promoted_artist
        FROM campaigns
        WHERE campaignid = $1
    `

    var c models.Campaign
    err = db.DB.QueryRow(query, id).Scan(
        &c.ID,
        &c.CampaignName,
        &c.ReferenceArtists,
        &c.TrelloLink,
        &c.SpotifyLink,
        &c.LaunchDate,
        &c.PromotedArtist,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("Campaign not found with ID: %d", id)
            util.RespondWithError(w, http.StatusNotFound, "Campaign not found")
            return
        }
        log.Printf("Error querying campaign: %v", err)
        util.RespondWithError(w, http.StatusInternalServerError, "Error retrieving campaign")
        return
    }

    log.Printf("Successfully retrieved campaign with ID: %d", id)
    util.RespondWithJSON(w, http.StatusOK, c)
    log.Println("GetCampaign function completed")
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
