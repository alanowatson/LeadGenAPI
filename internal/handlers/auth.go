package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/alanowatson/LeadGenAPI/internal/middleware"
    "github.com/alanowatson/LeadGenAPI/pkg/util"
)

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()

    // In a real application, you would validate the username and password against a database
    if req.Username != "admin" || req.Password != "password" {
        util.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
        return
    }

    token, err := middleware.GenerateToken(req.Username)
    if err != nil {
        util.RespondWithError(w, http.StatusInternalServerError, "Could not generate token")
        return
    }

    util.RespondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}
