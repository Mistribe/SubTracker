package exch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Oleexo/config-go"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/pkg/slicesx"
)

const (
	exchangeRateHost = "https://api.exchangerate.host/"
)

type Client interface {
	GetLiveExchangeRates(currencies []currency.Unit) (LiveExchangeRates, error)
}

type client struct {
	apiKey   string
	fakeData bool
}

func NewClient(cfg config.Configuration) Client {
	apiKey := cfg.GetString("EXCHANGERATE_HOST_API_KEY")
	fakeData := cfg.GetBoolOrDefault("ECHANGERATE_FAKE_DATA", false)
	return &client{
		apiKey:   apiKey,
		fakeData: fakeData,
	}
}

type liveExchangeRatesResponse struct {
	Success   bool               `json:"success"`
	Timestamp int64              `json:"timestamp"`
	Source    string             `json:"source"`
	Quotes    map[string]float64 `json:"quotes"`
}

type LiveExchangeRates struct {
	Timestamp time.Time
	Rates     map[currency.Unit]float64
}

func (c *client) GetLiveExchangeRates(currencies []currency.Unit) (LiveExchangeRates, error) {
	if c.fakeData {
		return getLiveFakeData(), nil
	}
	supportedCurrencies := slicesx.Select(currencies, func(in currency.Unit) string {
		return in.String()
	})
	url := fmt.Sprintf("%s/live?access_key=%s&format=1&currencies=%s", exchangeRateHost, c.apiKey,
		strings.Join(supportedCurrencies, ","))

	resp, err := http.Get(url)
	if err != nil {
		return LiveExchangeRates{}, fmt.Errorf("failed to fetch exchange rates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LiveExchangeRates{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var response liveExchangeRatesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return LiveExchangeRates{}, fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.Success {
		return LiveExchangeRates{}, fmt.Errorf("API request failed")
	}

	exchangeRates := make(map[currency.Unit]float64)
	for quote, rate := range response.Quotes {
		isoCurrency, err := currency.ParseISO(quote)
		if err != nil {
			return LiveExchangeRates{}, err
		}

		exchangeRates[isoCurrency] = rate
	}

	return LiveExchangeRates{
		Timestamp: time.Unix(response.Timestamp, 0),
		Rates:     exchangeRates,
	}, nil
}
