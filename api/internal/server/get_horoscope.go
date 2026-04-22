package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/joswayski/creditcardhoroscope/api/internal/horoscopes"
)

func (s *Server) GetHoroscope(w http.ResponseWriter, r *http.Request) {
	horoscopeId := r.PathValue("id")
	if horoscopeId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid horoscope ID",
		})
		return
	}

	// Check if its in the cache first
	horoscope, ok := s.HoroscopeCache.Get(horoscopeId)
	if ok {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(horoscope)
		return
	}

	var publicHoroscope horoscopes.PublicHoroscope

	err := s.DB.QueryRow(r.Context(), `
	SELECT external_id, horoscope, created_at
	FROM generations
	WHERE external_id = $1 AND is_public = true`,
		horoscopeId,
	).Scan(&publicHoroscope.ID, &publicHoroscope.Horoscope, &publicHoroscope.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Horoscope not found",
		})
		return
	}

	s.HoroscopeCache.Set(publicHoroscope)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Failed to fetch the horoscope :(",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(publicHoroscope)

}
