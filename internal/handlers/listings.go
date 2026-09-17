package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/durgasatkar/olx-api/internal/httpx"
	"github.com/durgasatkar/olx-api/internal/middleware"
)

type Listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
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

// GetListings godoc
//
//	@Summary		Get recent listings
//	@Description	Fetches the top 10 most recent classified listings ordered by creation date.
//	@Tags			listings
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		Listing
//	@Failure		500	{object}	map[string]interface{}	"Internal server error"
//	@Router			/listings [get]
func (lh ListingHandler) GetListings(w http.ResponseWriter, r *http.Request) {
	// request scoped context
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)
	rows, err := lh.db.QueryContext(ctx, `SELECT id, title, description, price, city, created_at
		FROM listings 
		ORDER BY created_at DESC 
		LIMIT 10`)
	if err != nil {
		lh.logger.Error("listing query error", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.ErrorCodeInternalError)
		return
	}
	defer rows.Close()
	listings := []Listing{}
	for rows.Next() {
		var l Listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			lh.logger.Error("rows scan error", "error", err)
			httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.ErrorCodeInternalError)
			return
		}
		listings = append(listings, l)
	}
	lh.logger.Info("listings fetched", "total", len(listings), "request_id", requestId)
	if err := rows.Err(); err != nil {
		lh.logger.Error("row error", "error", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.ErrorCodeInternalError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(listings)

}

// DeleteListing godoc
//
//	@Summary		Delete a listing
//	@Description	Removes an active item listing from the platform using its unique ID.
//	@Tags			listings
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Listing Unique Identifier"
//	@Success		204	{string}	string	"ok"
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/listings/{id} [delete]
func (lh ListingHandler) DeleteListing(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	fmt.Println("id", id)
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)
	_, err := lh.db.ExecContext(ctx, `DELETE FROM listings WHERE id = $1`, id)
	if err != nil {
		lh.logger.Error("delete failed", "listing_id", id, "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.ErrorCodeInternalError)
	}

	w.WriteHeader(http.StatusNoContent)

	w.Write([]byte("ok"))

}

// CreateListing godoc
//
//	@Summary		Create a new listing
//	@Description	Submits and registers a new item listing to the marketplace.
//	@Tags			listings
//	@Accept			json
//	@Produce		json
//	@Param			listing	body		CreateListingRequest	true	"Listing Object Details"
//	@Success		201		{object}    CreateListingResponse
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/listings [post]
func (lh ListingHandler) CreateListing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIDFromContext(ctx)
	var req CreateListingRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lh.logger.Error("failed to decode", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "invalid body", httpx.ErrorCodeMalformedJson)
		return
	}
	lh.logger.Info("craete payload", "payload", req, "request_id", requestId)
	if err := req.Validate(); err != nil {
		var verr *ValidationError
		errors.As(err, &verr)
		lh.logger.Error(err.Error(), "request_id", requestId, "err", err)
		httpx.ValidationError(w, http.StatusUnprocessableEntity, err.Error(), httpx.ErrorCodeValidationFailed, verr.Field)
	}
	row := lh.db.QueryRowContext(ctx, `INSERT INTO listings (title, description, price, city) VALUES ($1, $2, $3, $4) RETURNING id, title, created_at`, req.Title, req.Description, req.Price, req.City)
	// var id string
	var out CreateListingResponse
	if err := row.Scan(&out.ID, &out.CreatedAt); err != nil {
		lh.logger.Error("failed to insert", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.ErrorCodeInternalError)

		return
	}
	lh.logger.Info("listings created", "request_id", requestId, "listing", out)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(out)

}
