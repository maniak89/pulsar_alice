package alice

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"pulsar_alice/internal/mappers"
	"pulsar_alice/internal/models/alice"
	"pulsar_alice/internal/models/common"
	"time"

	"github.com/rs/zerolog/log"
)

type client struct {
	config          Config
	callbackAddress string
	client          *http.Client
}

func New(config Config) *client {
	return &client{
		config: config,
		client: &http.Client{
			Timeout: config.RequestTimeout,
		},
		callbackAddress: config.Address + "/api/v1/skills/" + config.SkillID + "/callback/state",
	}
}

func (c *client) NotifyMetersChanged(ctx context.Context, meters []*common.Meter) error {
	if len(meters) == 0 {
		return nil
	}

	logger := log.Ctx(ctx)
	blob, err := json.Marshal(alice.State{
		TS: time.Now().Unix(),
		Payload: alice.PayloadState{
			UserID:  meters[0].UserID,
			Devices: mappers.MetersState(meters),
		},
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed marshal body")
		return err
	}
	logger.Trace().Str("url", c.callbackAddress).Bytes("blob", blob).Msg("Prepare request body")
	req, err := http.NewRequest(http.MethodPost, c.callbackAddress, bytes.NewReader(blob))
	if err != nil {
		logger.Error().Err(err).Msg("Failed create request object")
		return err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "OAuth "+c.config.OAuth2Token)
	resp, err := c.client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed make request")
		return err
	}
	if resp.StatusCode != http.StatusOK {
		blob, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error().Err(err).Msg("Failed read body")
		}
		logger.Error().Str("status", resp.Status).Bytes("response", blob).Msg("status")
		return nil
	}
	logger.Debug().Str("status", resp.Status).Msg("status")

	return nil
}
