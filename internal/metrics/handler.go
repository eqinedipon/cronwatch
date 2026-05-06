package metrics

import (
	"encoding/json"
	"net/http"
	"sort"
)

// Handler returns an http.HandlerFunc that serialises all collected
// job metrics as a JSON array, sorted by job name.
func Handler(c *Collector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		all := c.All()
		sort.Slice(all, func(i, j int) bool {
			return all[i].JobName < all[j].JobName
		})

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(all); err != nil {
			http.Error(w, "failed to encode metrics", http.StatusInternalServerError)
		}
	}
}
