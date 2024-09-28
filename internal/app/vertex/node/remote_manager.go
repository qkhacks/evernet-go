package node

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type RemoteManager struct {
	httpClient *http.Client
}

func NewRemoteManager() *RemoteManager {
	return &RemoteManager{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (m *RemoteManager) Get(ctx context.Context, nodeVertex string, nodeIdentifier string) (*Node, error) {
	url := fmt.Sprintf("https://%s/api/v1/nodes/%s", nodeVertex, nodeIdentifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := m.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			zap.L().Error("failed to close response body", zap.Error(err))
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var node Node

	if err := json.Unmarshal(body, &node); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &node, nil
}
