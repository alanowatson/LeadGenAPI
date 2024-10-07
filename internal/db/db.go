package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/alanowatson/LeadGenAPI/internal/models"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
    dbUser := os.Getenv("DB_USER")
    dbName := os.Getenv("DB_NAME")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbSSLMode := os.Getenv("DB_SSLMODE")

    connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=%s",
        dbUser, dbName, dbPassword, dbHost, dbSSLMode)


    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("error opening database: %w", err)
    }

    if err = DB.Ping(); err != nil {
        return fmt.Errorf("error connecting to database: %w", err)
    }

    return nil
}

func BeginTx() (*sql.Tx, error) {
    return DB.Begin()
}

func PlaylistExists(id int) (bool, error) {
    var exists bool
    err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM playlists WHERE id=$1)", id).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("error checking playlist existence: %w", err)
    }
    return exists, nil
}

func CampaignExists(id int) (bool, error) {
    var exists bool
    err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM campaigns WHERE id=$1)", id).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("error checking campaign existence: %w", err)
    }
    return exists, nil
}

func CreatePlaylistCampaignTx(tx *sql.Tx, pc models.PlaylistCampaign) error {
    _, err := tx.Exec(`
        INSERT INTO playlist_campaigns (playlist_id, campaign_id, playlister_id, reference_artists, placement_status, number_of_messages, purchased)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, pc.PlaylistID, pc.CampaignID, pc.PlaylisterId, pc.ReferenceArtists, pc.PlacementStatus, pc.NumberOfMessages, pc.Purchased)

    if err != nil {
        return fmt.Errorf("error creating playlist campaign: %w", err)
    }
    return nil
}
