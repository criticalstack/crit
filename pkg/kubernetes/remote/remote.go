package remote

import (
	"context"
	"io"
	"io/ioutil"
	"net"
	"strings"

	"github.com/hpcloud/tail"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"

	"github.com/criticalstack/crit/pkg/log"
)

func dial(ctx context.Context, addr string) (net.Conn, error) {
	return (&net.Dialer{}).DialContext(ctx, "unix", addr)
}

type RuntimeServiceClient struct {
	runtimeapi.RuntimeServiceClient
}

func NewRuntimeServiceClient(ctx context.Context, endpoint string) (*RuntimeServiceClient, error) {
	addr := strings.TrimPrefix(endpoint, "unix://")
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithContextDialer(dial))
	if err != nil {
		return nil, err
	}
	r := &RuntimeServiceClient{
		RuntimeServiceClient: runtimeapi.NewRuntimeServiceClient(conn),
	}
	return r, nil
}

func (r *RuntimeServiceClient) GetContainerByName(ctx context.Context, name string) (*runtimeapi.Container, error) {
	resp, err := r.ListContainers(ctx, &runtimeapi.ListContainersRequest{
		Filter: &runtimeapi.ContainerFilter{
			LabelSelector: map[string]string{
				"io.kubernetes.container.name": name,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Containers) == 0 {
		return nil, errors.Errorf("cannot find container: %q", name)
	}
	return resp.Containers[0], nil
}

func (r *RuntimeServiceClient) GetContainerStatus(ctx context.Context, id string) (*runtimeapi.ContainerStatus, error) {
	resp, err := r.ContainerStatus(ctx, &runtimeapi.ContainerStatusRequest{
		ContainerId: id,
	})
	if err != nil {
		return nil, err
	}
	return resp.Status, nil
}

func (r *RuntimeServiceClient) ReadLogs(ctx context.Context, path string, w io.Writer) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (r *RuntimeServiceClient) TailLogs(ctx context.Context, path string, w io.Writer) error {
	t, err := tail.TailFile(path, tail.Config{Follow: true})
	if err != nil {
		return err
	}

	for {
		select {
		case l := <-t.Lines:
			if _, err := w.Write([]byte(l.Text + "\n")); err != nil {
				return err
			}
		case <-ctx.Done():
			log.Debug("finished tailing log", zap.String("path", path), zap.Error(ctx.Err()))
			return t.Stop()
		}
	}
}
