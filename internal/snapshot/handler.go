package snapshot

import (
	"encoding/json"
	"net/http"
)

// Handler returns an HTTP handler that serves the latest snapshot as JSON.
// If no snapshot has been captured yet it returns 204 No Content.
func Handler(s *Snapshotter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		snap := s.Latest()
		if snap == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(snap); err != nil {
			http.Error(w, "failed to encode snapshot", http.StatusInternalServerError)
		}
	}
}
