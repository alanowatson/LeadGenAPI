package errors

import (
    "fmt"
    "log"
    "net/http"

    "github.com/alanowatson/LeadGenAPI/pkg/util"
)

func HandleError(w http.ResponseWriter, err error, status int, message string) {
    log.Printf("Error: %v", err)
    util.RespondWithError(w, status, fmt.Sprintf("%s: %v", message, err))
}
