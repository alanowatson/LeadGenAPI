package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "sync"

    "github.com/gorilla/mux"
    "github.com/alanowatson/LeadGenAPI/internal/models"
    "github.com/alanowatson/LeadGenAPI/pkg/util"
)

var (
    playlisters = make(map[int]models.Playlister)
    playlisterID = 1
    playlisterMutex sync.RWMutex
)

func GetPlaylisters(w http.ResponseWriter, r *http.Request) {
    playlisterMutex.RLock()
    defer playlisterMutex.RUnlock()

    playlisterList := make([]models.Playlister, 0, len(playlisters))
    for _, playlister := range playlisters {
        playlisterList = append(playlisterList, playlister)
    }

    util.RespondWithJSON(w, http.StatusOK, playlisterList)
}

func CreatePlaylister(w http.ResponseWriter, r *http.Request) {
    var playlister models.Playlister
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&playlister); err != nil {
        util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    playlisterMutex.Lock()
    defer playlisterMutex.Unlock()

    playlister.ID = playlisterID
    playlisters[playlisterID] = playlister
    playlisterID++

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
        util.RespondWithError(w, http.StatusBadRequest, "Invalid playlister ID")
        return
    }

    var playlister models.Playlister
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&playlister); err != nil {
        util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

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
