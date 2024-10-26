package rest

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"pulsar_alice/internal/mappers"
	"pulsar_alice/internal/models/alice"
	"pulsar_alice/pkg/middleware/user"
)

func (s *service) Query(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.Ctx(ctx)
	var req alice.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error().Err(err).Msg("Failed unmarshal data")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	devices := s.deviceProvider.Meters(ctx, user.User(ctx))

	aliceDevices := alice.PayloadState{
		UserID:  user.User(ctx),
		Devices: make([]alice.PayloadStateDevice, 0, len(devices)),
	}
	for _, reqDev := range req.Devices {
		for _, dev := range devices {
			if dev.ID != reqDev.ID {
				continue
			}
			aliceDevices.Devices = append(aliceDevices.Devices, mappers.MeterState(dev))
			break
		}
	}
	if err := json.NewEncoder(w).Encode(alice.State{
		RequestID: r.Header.Get(xRequestID),
		Payload:   aliceDevices,
	}); err != nil {
		logger.Error().Err(err).Msg("Failed marshal response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
