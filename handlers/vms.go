package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fuddata/anyvm/models"
	"github.com/fuddata/anyvm/providers"
)

func ListVMsHandler(cm *providers.CloudManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		provider := r.URL.Query().Get("provider")
		var vms []models.VM

		if provider != "" {
			p := cm.GetProvider(strings.ToLower(provider))
			if p == nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(models.APIResponse{
					Success: false,
					Error:   "Invalid provider specified",
				})
				return
			}
			result, err := p.ListVMs()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(models.APIResponse{
					Success: false,
					Error:   err.Error(),
				})
				return
			}
			vms = result
		} else {
			for _, p := range cm.GetAllProviders() {
				result, err := p.ListVMs()
				if err != nil {
					continue // In production, log and handle errors
				}
				vms = append(vms, result...)
			}
		}

		json.NewEncoder(w).Encode(models.APIResponse{
			Success: true,
			Data:    vms,
		})
	}
}
