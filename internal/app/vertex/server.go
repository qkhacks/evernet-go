package vertex

import (
	"fmt"
	"github.com/evernetproto/evernet/internal/app/vertex/db"
	"github.com/evernetproto/evernet/internal/app/vertex/health"
	"github.com/evernetproto/evernet/internal/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
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
	Host       string
	Port       string
	Vertex     string
	DataPath   string
	StaticPath string
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

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile(s.config.StaticPath, true)))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	health.NewHandler(router).Register()

	zap.L().Info("starting vertex", zap.String("host", s.config.Host), zap.String("port", s.config.Port))
	err = router.Run(fmt.Sprintf("%s:%s", s.config.Host, s.config.Port))

	if err != nil {
		zap.L().Panic("error while starting vertex", zap.Error(err))
	}
}
