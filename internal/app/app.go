package app

import (
	"fmt"
	"github.com/joho/godotenv"
	delivery "github.com/weazyexe/fonto-server/internal/delivery/grpc"
	"github.com/weazyexe/fonto-server/internal/delivery/grpc/interceptors"
	"github.com/weazyexe/fonto-server/internal/repository"
	"github.com/weazyexe/fonto-server/internal/service"
	"github.com/weazyexe/fonto-server/pkg/crypto"
	"github.com/weazyexe/fonto-server/pkg/domain"
	"github.com/weazyexe/fonto-server/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func Run(configPath string) {
	// Initializing logger
	logger.InitializeLogger()
	defer logger.Zap.Sync()

	// Loading environment variables
	if err := godotenv.Load(); err != nil {
		logger.Zap.Fatalf("Failed to load enviroment variables: %v", err)
	}

	// Reading config files
	config, err := ReadConfig(configPath)
	if err != nil {
		logger.Zap.Fatalf("Failed to read config: %v\n%v", configPath, err)
	}
	logger.Zap.Info("Configs was read")

	// Establishing connection to postgres database
	db := establishPostgresConnection(config.Postgres)
	defer func(db *gorm.DB) {
		sqlDb, err := db.DB()
		if err != nil {
			logger.Zap.Fatalf("Failed to open database connection: %v", err)
		}
		if err := sqlDb.Close(); err != nil {
			logger.Zap.Fatalf("Failed to close database connection: %v", err)
		}
	}(db)

	// Listening to the port
	lis := listenToPort(config.Port)

	// Initializing gRPC server & receivers
	server := initializeGrpc(db, config.Auth)

	// RUN
	logger.Zap.Info("Here we go!")
	if err := server.Serve(lis); err != nil {
		logger.Zap.Fatalf("❌ Failed to serve: %v", err)
	}
}

func establishPostgresConnection(config Postgres) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Host,
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		config.DbName,
		config.Port,
	)
	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN:                  dsn,
				PreferSimpleProtocol: true, // disables implicit prepared statement usage
			},
		),
		&gorm.Config{},
	)
	if err != nil {
		logger.Zap.Fatalf("Failed to open connection to the database: %v", err)
	}

	logger.Zap.Info("Connected to PostgreSQL")

	return db
}

func listenToPort(port string) net.Listener {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Zap.Fatalf("❌ Failed to listen: %v", err)
	}
	logger.Zap.Infof("Server is listening at %v", lis.Addr())
	return lis
}

func initializeGrpc(db *gorm.DB, authConfig Auth) *grpc.Server {
	jwtManager := crypto.NewJwtManager(
		[]byte(os.Getenv("ACCESS_TOKEN_SECRET")),
		[]byte(os.Getenv("REFRESH_TOKEN_SECRET")),
		&domain.JwtConfig{
			Issuer:               authConfig.Issuer,
			ExpireTimeForAccess:  authConfig.ExpireTimeForAccess,
			ExpireTimeForRefresh: authConfig.ExpireTimeForRefresh,
		},
	)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.NewAuthInterceptor(jwtManager).Intercept()),
	)

	// Authentication feature
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, jwtManager)
	delivery.NewAuthReceiver(authService).Register(s)

	// Greeter feature
	delivery.NewGreeterReceiver().Register(s)

	logger.Zap.Info("Initialized gRPC server")
	return s
}
