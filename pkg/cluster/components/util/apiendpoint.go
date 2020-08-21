package util

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/criticalstack/crit/pkg/log"
	fmtutil "github.com/criticalstack/crit/pkg/util/fmt"
	netutil "github.com/criticalstack/crit/pkg/util/net"
)

// APIEndpoint represents a reachable Kubernetes API endpoint using scheme-less
// host:port (port is optional).
type APIEndpoint struct {
	// The hostname on which the API server is serving.
	Host string `json:"host"`

	// The port on which the API server is serving.
	Port int32 `json:"port"`
}

// String returns a formatted version HOST:PORT of this APIEndpoint.
func (v APIEndpoint) String() string {
	return fmt.Sprintf("%s:%d", v.Host, v.Port)
}

// IsZero returns true if host and the port are zero values.
func (v APIEndpoint) IsZero() bool {
	return v.Host == "" && v.Port == 0
}

func (v *APIEndpoint) UnmarshalJSON(data []byte) error {
	type alias APIEndpoint
	if reterr := json.Unmarshal(data, (*alias)(v)); reterr != nil {
		if err := v.tryUnmarshalText(fmtutil.Unquote(string(data))); err != nil {
			log.Debug("tryUnmarshalText", zap.Error(err))
			return reterr
		}
	}
	return nil
}

func (v *APIEndpoint) tryUnmarshalText(s string) (err error) {
	v.Host, v.Port, err = netutil.SplitHostPort(s)
	return err
}
