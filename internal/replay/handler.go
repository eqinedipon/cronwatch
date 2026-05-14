package replay

import (
	"encoding/json"
	"net/http"
)

// replayResponse is the JSON body returned by the replay endpoint.
type replayResponse struct {
	Dispatched []string          `json:"dispatched"`
	Errors     map[string]string `json:"errors,omitempty"`
}

// Handler returns an http.HandlerFunc that triggers a replay on POST.
func Handler(r *Replayer) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		results := r.Run(req.Context())

		resp := replayResponse{
			Dispatched: []string{},
			Errors:     make(map[string]string),
		}

		for _, res := range results {
			if res.Error != nil {
				resp.Errors[res.Job] = res.Error.Error()
			} else {
				resp.Dispatched = append(resp.Dispatched, res.Job)
			}
		}

		if len(resp.Errors) == 0 {
			resp.Errors = nil
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
