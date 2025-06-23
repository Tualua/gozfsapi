package server

import (
	"net/http"
	"slices"
	"sync"

	. "github.com/Tualua/gozfsapi/internal/models"
	"github.com/Tualua/gozfsapi/pkgs/zfs"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

//go:generate go tool oapi-codegen -config server.cfg.yaml ../../api/api.yaml
//go:generate go tool oapi-codegen -config models.cfg.yaml ../../api/api.yaml
type ApiServer struct {
	lock sync.Mutex
	Log  *zap.Logger
}

// NewApiServer creates a new ApiServer instance.
func NewApiServer(log *zap.Logger) *ApiServer {
	return &ApiServer{
		Log: log,
	}
}

func (s *ApiServer) ListZfsDatasets(ctx echo.Context, params ListZfsDatasetsParams) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	datasets := slices.Collect(func(yield func(ZFSDataset) bool) {
		for _, name := range zfs.ListDatasets(string(*params.Type)) {
			if !yield(ZFSDataset{Name: name}) {
				return
			}
		}
	})

	if datasets == nil {
		datasets = make([]ZFSDataset, 0)
	}
	s.Log.Info("Listed ZFS datasets", zap.Int("count", len(datasets)))
	return ctx.JSON(http.StatusOK, datasets)
}

func (s *ApiServer) GetReadyz(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Ready")
}
