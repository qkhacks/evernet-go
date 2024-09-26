package vertex

import (
	"github.com/evernetproto/evernet/internal/app/vertex/db"
	"github.com/evernetproto/evernet/internal/pkg/logger"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

type Server struct {
	config *ServerConfig
}

func NewServer(config *ServerConfig) *Server {
	return &Server{config: config}
}

type ServerConfig struct {
	Host     string
	Port     string
	Vertex   string
	DataPath string
}

const (
	ServiceName      = "vertex"
	MetaDatabaseFile = "meta.db"
	MetaDatabase     = "metabase"
)

func (s *Server) Start() {
	logger.Init(ServiceName)
	defer func() {
		_ = zap.L().Sync()
	}()

	err := os.MkdirAll(s.config.DataPath, os.ModePerm)

	if err != nil {
		zap.L().Error("error creating data directory", zap.Error(err))
	}

	metaDatabasePath := filepath.Join(s.config.DataPath, MetaDatabaseFile)
	db.MigrateDatabase(metaDatabasePath, MetaDatabase)
}
