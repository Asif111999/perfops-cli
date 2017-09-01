package perfops

import (
	"context"
	"errors"
	"net/http"
)

type (
	// RunService defines the interface for the run API
	RunService service

	// Ping represents the parameters for a ping request.
	Ping struct {
		// Target name
		Target string `json:"target"`
		// List of nodes ids, comma separated
		Nodes string `json:"nodes,omitempty"`
		// Countries names, comma separated
		Location string `json:"location,omitempty"`
		// Max number of nodes
		Limit int `json:"limit,omitempty"`
	}

	// PingID represents the ID of a ping test.
	PingID string

	// PingResult represents the result of a ping.
	PingResult struct {
		NodeID string `json:"nodeId,omitempty"`
		Output string `json:"output,omitempty"`
	}

	// PingItem represents an item of a ping output.
	PingItem struct {
		ID     string      `json:"id,omitempty"`
		Result *PingResult `json:"result,omitempty"`
	}

	// PingOutput represents the response of ping output calls.
	PingOutput struct {
		ID        string      `json:"id,omitempty"`
		Requested string      `json:"requested,omitempty"`
		Items     []*PingItem `json:"items,omitempty"`
	}
)

// Ping runs a ping test.
func (s *RunService) Ping(ctx context.Context, ping *Ping) (PingID, error) {
	body, err := newJSONReader(ping)
	if err != nil {
		return "", err
	}
	u := s.client.BasePath + "/run/ping"
	req, _ := http.NewRequest("POST", u, body)
	req = req.WithContext(ctx)
	var raw struct {
		Error string
		ID    string `json:"id"`
	}
	if err = s.client.do(req, &raw); err != nil {
		return "", err
	}
	if raw.Error != "" {
		return "", errors.New(raw.Error)
	}
	return PingID(raw.ID), nil
}

// PingOutput returns the full ping output under a test ID.
func (s *RunService) PingOutput(ctx context.Context, pingID PingID) (*PingOutput, error) {
	u := s.client.BasePath + "/run/ping/" + string(pingID)
	req, _ := http.NewRequest("GET", u, nil)
	var v *PingOutput
	if err := s.client.do(req, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IsComplete returns a value indicating whether the whole output is
// complete or not.
func (o *PingOutput) IsComplete() bool {
	if o.Requested == "" {
		return false
	}
	var n int
	for _, item := range o.Items {
		if !item.Result.IsComplete() {
			break
		}
		n++
	}
	return n == len(o.Items)
}

// IsComplete returns a value indicating whether the result is complete
// or not.
func (r *PingResult) IsComplete() bool {
	return r.NodeID != ""
}
