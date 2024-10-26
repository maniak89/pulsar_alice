package rest

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"pulsar_alice/internal/mappers"
	"pulsar_alice/internal/models/alice"
	"pulsar_alice/pkg/middleware/user"
)

func (s *service) Devices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.Ctx(ctx)

	devices := s.deviceProvider.Meters(ctx, user.User(ctx))

	aliceDevices := alice.Devices{
		UserID:  user.User(ctx),
		Devices: make([]alice.Device, 0, len(devices)),
	}

	for _, d := range devices {
		aliceDevices.Devices = append(aliceDevices.Devices, mappers.MeterDevice(d))
	}

	response := alice.Response{
		RequestID: r.Header.Get(xRequestID),
		Payload:   aliceDevices,
	}
	blob, err := json.Marshal(&response)
	if err != nil {
		logger.Error().Err(err).Msg("Failed marshal response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	logger.Trace().Bytes("response", blob).Msg("Write response")
	if _, err := w.Write(blob); err != nil {
		logger.Error().Err(err).Msg("Failed write response")
	}
}
