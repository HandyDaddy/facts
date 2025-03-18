package factclient

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/HandyDaddy/facts/internal/config"
	"github.com/HandyDaddy/facts/internal/domain/entities"
)

type HTTPFactsClient struct {
	cfg    *config.HttpClient
	client *http.Client
	logger *logrus.Logger
}

func NewHTTPClient(cfg *config.HttpClient) *HTTPFactsClient {
	return &HTTPFactsClient{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
		logger: logrus.New(),
	}
}

func (c HTTPFactsClient) Save(ctx context.Context, fact *entities.Fact) error {
	data := url.Values{}
	data.Set("period_start", fact.PeriodStart)
	data.Set("period_end", fact.PeriodEnd)
	data.Set("period_key", fact.PeriodKey)
	data.Set("indicator_to_mo_id", fmt.Sprintf("%d", fact.IndicatorToMoId))
	data.Set("indicator_to_mo_fact_id", fmt.Sprintf("%d", fact.IndicatorToMoFactId))
	data.Set("value", fmt.Sprintf("%d", fact.Value))
	data.Set("fact_time", fact.FactTime)
	data.Set("is_plan", fmt.Sprintf("%d", fact.IsPlan))
	data.Set("auth_user_id", fmt.Sprintf("%d", fact.AuthUserId))
	data.Set("comment", fact.Comment)

	req, err := http.NewRequest("POST", c.cfg.Addr, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.cfg.Token))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	return nil
}

func (c HTTPFactsClient) Get(ctx context.Context, id string) (*entities.Fact, error) {
	//TODO implement me
	panic("implement me")
}
