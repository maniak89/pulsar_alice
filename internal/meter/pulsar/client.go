package pulsar

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sigurn/crc16"

	"pulsar_alice/internal/meter"
)

const (
	payloadPosition        = 6
	crc16Len               = 2
	reqIDLen               = 2
	responseTailLen        = crc16Len + reqIDLen
	minResponseLen         = payloadPosition + responseTailLen
	responsePayloadReadLen = payloadPosition + responseTailLen + 4
)

var (
	readCommandRequest = []byte{0x01, 0x0e, 0x01, 0x00, 0x00, 0x00}
	crcTable           = crc16.MakeTable(crc16.CRC16_MODBUS)
)

type addressOrError struct {
	address []byte
	err     error
}

type client struct {
	config Config

	meterAddress addressOrError
}

func New(config Config) *client {
	return &client{
		config: config,
	}
}

func (c *client) Init(ctx context.Context) error {
	logger := log.Ctx(ctx).With().Str("address", c.config.Address).Str("meter", c.config.Meter).Logger()

	addr, err := parseAddress(c.config.Meter)
	if err != nil {
		logger.Error().Err(err).Str("meter", c.config.Meter).Msg("Failed parse meter addr")

		c.meterAddress = addressOrError{err: err}

		return nil
	}

	c.meterAddress = addressOrError{address: addr}

	return nil
}

func parseAddress(in string) ([]byte, error) {
	const expectedStringLen = 8
	if l := len(in); l < expectedStringLen {
		in = strings.Repeat("0", expectedStringLen-l) + in
	}

	return hex.DecodeString(in)
}

func (c *client) Value(ctx context.Context) meter.Value {
	logger := log.Ctx(ctx).With().Str("address", c.config.Address).Str("meter", c.config.Meter).Logger()

	if c.meterAddress.err != nil {
		return meter.Value{Address: c.config.Meter, Error: c.meterAddress.err}
	}

	dealer := net.Dialer{
		Timeout: c.config.Timeout,
	}

	conn, err := dealer.Dial("tcp4", c.config.Address)
	if err != nil {
		logger.Error().Err(err).Msg("Failed connect to address")

		return meter.Value{Address: c.config.Meter, Error: fmt.Errorf("connect: %w", err)}
	}

	reqR, err := rand.Int(rand.Reader, big.NewInt(math.MaxUint16))
	if err != nil {
		logger.Error().Err(err).Msg("Failed generate request id")

		return meter.Value{Address: c.config.Meter, Error: fmt.Errorf("generate request id: %w", err)}
	}

	reqID := uint16ToByte(uint16(reqR.Uint64()))

	if err := conn.SetWriteDeadline(time.Now().Add(c.config.Timeout)); err != nil {
		logger.Error().Err(err).Msg("Failed set write deadline")
	}

	req := makeReadRequest(c.meterAddress.address, reqID)
	for {
		n, err := conn.Write(req)
		if err != nil {
			logger.Error().Err(err).Msg("Failed write deadline")

			return meter.Value{Address: c.config.Meter, Error: fmt.Errorf("write: %w", err)}
		}

		req = req[n:]
		if len(req) == 0 {
			break
		}
	}

	if err := conn.SetReadDeadline(time.Now().Add(c.config.Timeout)); err != nil {
		logger.Error().Err(err).Msg("Failed set read deadline")
	}

	buffer := make([]byte, responsePayloadReadLen)
	var readed int
	for readed < len(buffer) {
		n, err := conn.Read(buffer[readed:])
		if err != nil {
			logger.Error().Err(err).Msg("Failed read")

			return meter.Value{Address: c.config.Meter, Error: fmt.Errorf("read: %w", err)}
		}

		readed += n
	}

	buffer = buffer[:readed]

	if err := validateResponse(buffer, c.meterAddress.address, reqID); err != nil {
		logger.Error().Err(err).Msg("Failed validate response")

		return meter.Value{Address: c.config.Meter, Error: fmt.Errorf("validate: %w", err)}
	}

	return meter.Value{Address: c.config.Meter, Value: extractMeterValue(buffer)}
}

func crc16Calc(data []byte) []byte {
	return uint16ToByte(crc16.Checksum(data, crcTable))
}

func uint16ToByte(in uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, in)

	return b
}

func makeReadRequest(devID, reqID []byte) []byte {
	b := make([]byte, 0, 14)
	b = append(b, devID...)
	b = append(b, readCommandRequest...)
	b = append(b, reqID...)
	b = append(b, crc16Calc(b)...)
	return b
}

func validateResponse(response, devID, reqID []byte) error {
	if len(response) < minResponseLen {
		return errors.New("too small response len")
	}

	for i := range devID {
		if response[i] != devID[i] {
			return errors.New("address mismatched")
		}
	}

	tail := response[len(response)-responseTailLen:]
	for i := range reqID {
		if tail[i] != reqID[i] {
			return errors.New("request id mismatched")
		}
	}

	tail = tail[reqIDLen:]

	crc16 := crc16Calc(response[:len(response)-crc16Len])
	for i := range crc16 {
		if tail[i] != crc16[i] {
			return errors.New("checksum mismatched")
		}
	}

	return nil
}

func extractMeterValue(in []byte) float64 {
	value := binary.LittleEndian.Uint32(in[payloadPosition : len(in)-responseTailLen])
	return float64(value) / 1000
}
