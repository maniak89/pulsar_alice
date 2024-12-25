package forward

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	cookie_SESSION = "PHPSESSID"
	coldName       = "Холодное водоснабжение"
	hotName        = "Горячее водоснабжение"
)

var (
	accountRegexp      = regexp.MustCompile(`\?accountId=(\d+)`)
	webAccountIDRegexp = regexp.MustCompile(`webAccountId\s?=\s?(\d+)`)
	accountIDRegexp    = regexp.MustCompile(`accountId\s?=\s?(\d+)`)
	accountPIDRegexp   = regexp.MustCompile(`accountPid\s?=\s?(\d+)`)
	dbIDRegexp         = regexp.MustCompile(`dbId\s?=\s?(\d+)`)
)

type Client struct {
	config Config
}

type Session struct {
	config Config
	client *http.Client
}

func New(config Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) newClient() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Timeout: c.config.Timout,
		Jar:     jar,
	}, nil
}

func (c *Client) StartSession(ctx context.Context) (*Session, error) {
	logger := log.Ctx(ctx).With().Str("method", "startSession").Logger()
	req, err := http.NewRequest(http.MethodGet, c.config.Address, nil)
	if err != nil {
		logger.Error().Err(err).Msg("create request")

		return nil, err
	}

	req = req.WithContext(ctx)

	cl, err := c.newClient()
	if err != nil {
		logger.Error().Err(err).Msg("create client")

		return nil, err
	}

	resp, err := cl.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("make request")

		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error().Int("status", resp.StatusCode).Msg("status")

		return nil, errors.New("invalid status")
	}

	result := &Session{
		config: c.config,
		client: cl,
	}

	return result, result.auth(ctx)
}

func (c *Session) auth(ctx context.Context) error {
	logger := log.Ctx(ctx).With().Str("method", "auth").Logger()

	data := url.Values{
		"AUTH_FORM":     []string{"Y"},
		"TYPE":          []string{"AUTH"},
		"backurl":       []string{"/auth/"},
		"USER_LOGIN":    []string{c.config.Login},
		"USER_PASSWORD": []string{c.config.Password},
		"Login":         []string{"Войти"},
	}

	req, err := http.NewRequest(http.MethodPost, c.config.Address+"/auth/?login=yes", strings.NewReader(data.Encode()))
	if err != nil {
		logger.Error().Err(err).Msg("create request")

		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("make request")

		return err
	}

	defer resp.Body.Close()

	if len(resp.Cookies()) == 0 {
		return errors.New("auth failed")
	}

	return nil
}

func (c *Session) getAccounts(ctx context.Context) ([]string, error) {
	logger := log.Ctx(ctx).With().Str("method", "getAccounts").Logger()

	req, err := http.NewRequest(http.MethodGet, c.config.Address+"/ameter/", nil)
	if err != nil {
		logger.Error().Err(err).Msg("create request")

		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("make request")

		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error().Err(err).Msg("read response")

		return nil, err
	}

	return parseReg(accountRegexp, body), nil
}

func parseReg(reg *regexp.Regexp, in []byte) []string {
	matches := reg.FindAllSubmatchIndex(in, -1)
	result := make([]string, 0, len(matches))
	for _, match := range matches {
		result = append(result, string(in[match[2]:match[3]]))
	}
	return result
}

func parseRegString(reg *regexp.Regexp, in []byte) string {
	res := parseReg(reg, in)
	if len(res) == 0 {
		return ""
	}

	return res[0]
}

func (c *Session) getAccountInfo(ctx context.Context, accountID string) (*accountInfo, error) {
	logger := log.Ctx(ctx).With().Str("method", "getAccountInfo").Logger()

	req, err := http.NewRequest(http.MethodGet, c.config.Address+"/ameter/?accountId="+accountID, nil)
	if err != nil {
		logger.Error().Err(err).Msg("create request")

		return nil, err
	}

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("make request")

		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error().Err(err).Msg("read response")

		return nil, err
	}

	return &accountInfo{
		webAccountID: parseRegString(webAccountIDRegexp, body),
		accountID:    parseRegString(accountIDRegexp, body),
		accountPID:   parseRegString(accountPIDRegexp, body),
		dbID:         parseRegString(dbIDRegexp, body),
	}, nil
}

func (c *Session) getMetrics(ctx context.Context, info *accountInfo) (*metersResponse, error) {
	logger := log.Ctx(ctx).With().Str("method", "getMetrics").Logger()

	data := url.Values{
		"webAccountId": []string{info.webAccountID},
		"accountPid":   []string{info.accountPID},
		"dbId":         []string{info.dbID},
	}

	req, err := http.NewRequest(http.MethodPost, c.config.Address+"/ajax/meters.php", strings.NewReader(data.Encode()))
	if err != nil {
		logger.Error().Err(err).Msg("create request")

		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("make request")

		return nil, err
	}

	defer resp.Body.Close()

	var result metersResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Error().Err(err).Msg("decode")

		return nil, err
	}

	return &result, nil
}

func (c *Session) setMetrics(ctx context.Context, accountInfo *accountInfo, meterInfo *metersResponse, cold, hot float64) error {
	logger := log.Ctx(ctx).With().
		Str("method", "setMetrics").
		Str("account", accountInfo.accountID).
		Float64("cold", cold).
		Float64("hot", hot).
		Logger()

	req, err := http.NewRequest(http.MethodPost,
		c.config.Address+"/ajax/write_newmetering.php",
		strings.NewReader(makePayload(accountInfo, meterInfo, cold, hot).Encode()))
	if err != nil {
		logger.Error().Err(err).Msg("create request")

		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("make request")

		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error().Int("status", resp.StatusCode).Msg("status")

		return errors.New("invalid status")
	}

	logger.Debug().Msg("Send success")

	return nil
}

func (c *Session) SendMetrics(ctx context.Context, cold, hot float64) error {
	logger := log.Ctx(ctx).With().Str("method", "SendMetrics").Logger()

	accounts, err := c.getAccounts(ctx)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		logger := logger.With().Str("account", account).Logger()
		ctx := logger.WithContext(ctx)

		info, err := c.getAccountInfo(ctx, account)
		if err != nil {
			logger.Error().Err(err).Msg("account info")

			return err
		}

		metricsInfo, err := c.getMetrics(ctx, info)
		if err != nil {
			logger.Error().Err(err).Msg("metrics info")

			return err
		}

		if err := c.setMetrics(ctx, info, metricsInfo, cold, hot); err != nil {
			logger.Error().Err(err).Msg("save metrics info")

			return err
		}
	}

	return nil
}

func makePayload(accountInfo *accountInfo, meterInfo *metersResponse, cold, hot float64) url.Values {
	result := url.Values{
		"webAccountId": []string{accountInfo.webAccountID},
		"accountPid":   []string{accountInfo.accountPID},
		"dbId":         []string{accountInfo.dbID},
		"trans":        []string{meterInfo.Answer.Trans},
	}
	var meterIndex int
	for _, room := range meterInfo.Answer.Data {
		for _, meter := range room.Data {
			prefixMeter := "meters[" + strconv.Itoa(meterIndex) + "]"
			result.Set(prefixMeter+"[Loc]", room.RemarkIDText+" - "+meter.ServiceName)
			result.Set(prefixMeter+"[ServiceNameId]", meter.ServiceNameID)
			for _, meterData := range meter.Data {
				result.Set(prefixMeter+"[MeterId]", meterData.MeterID)
				result.Set(prefixMeter+"[RateNumbers]", meterData.RateNumbers)
				result.Set(prefixMeter+"[RateNumber]", meterData.RateNumber1)
				result.Set(prefixMeter+"[OldMeter]", meterData.Value1)
				break
			}
			switch meter.ServiceName {
			case coldName:
				result.Set(prefixMeter+"[MeterValue]", strconv.FormatFloat(cold, 'f', 3, 64))
			case hotName:
				result.Set(prefixMeter+"[MeterValue]", strconv.FormatFloat(hot, 'f', 3, 64))
			}
			meterIndex++
		}
	}

	return result
}
