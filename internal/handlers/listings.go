package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListingHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{
		db:     db,
		logger: logger,
	}
}

func (lh ListingHandler) GetListings(w http.ResponseWriter, r *http.Request) {
	// request scoped context
	ctx := r.Context()
	rows, err := lh.db.QueryContext(ctx, `SELECT id, title, description, price, city, created_at
		FROM listings 
		ORDER BY created_at DESC 
		LIMIT 10`)
	if err != nil {
		lh.logger.Error("listing query error", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	listings := []listing{}
	for rows.Next() {
		var l listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			lh.logger.Error("rows scan error", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		listings = append(listings, l)
	}
	lh.logger.Info("listings fetched", "total", len(listings))
	if err := rows.Err(); err != nil {
		lh.logger.Error("row error", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(listings)

}

func (lh ListingHandler) DeleteListing(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	fmt.Println("id", id)
	ctx := r.Context()
	_, err := lh.db.ExecContext(ctx, `DELETE FROM listings WHERE id = $1`, id)
	if err != nil {
		lh.logger.Error("delete failed", "listing_id", id, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)

	w.Write([]byte("ok"))

}
