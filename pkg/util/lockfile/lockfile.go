package lockfile

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/criticalstack/crit/pkg/log"
	"go.uber.org/zap"
)

type Lock struct {
	name string
}

func New(name string) *Lock {
	err := os.MkdirAll(filepath.Dir(name), 0755)
	if err != nil {
		log.Error("cannot MkdirAll for lockfile", zap.Error(err), zap.String("filename", name))
	}
	return &Lock{name: name}
}

func (l *Lock) Lock(ctx context.Context) (reterr error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	defer func() {
		if reterr != nil {
			l.Unlock()
		}
	}()

	for {
		select {
		case <-ticker.C:
			f, err := os.OpenFile(l.name, os.O_CREATE|os.O_EXCL, 0)
			if err != nil {
				if !os.IsExist(err) {
					return err
				}
				continue
			}
			return f.Close()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (l *Lock) Unlock() {
	_ = os.Remove(l.name)
}
