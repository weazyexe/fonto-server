package app

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	delivery "github.com/weazyexe/fonto-server/internal/delivery/grpc"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func Run(configPath string) {
	// Loading environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("❌ Failed to load enviroment variables: %v", err)
	}

	// Reading config files
	config, err := ReadConfig(configPath)
	if err != nil {
		log.Fatalf("❌ Failed to read config: %v\n%v", configPath, err)
	}
	log.Println("📝 Configs was read")

	// Establishing connection to postgres database
	db := establishPostgresConnection(config.Postgres)
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Fatalf("❌ Failed to close database connection: %v", err)
		}
	}(db)

	// Listening to the port
	lis := listenToPort(config.Port)

	// Initializing gRPC server & receivers
	server := initializeGrpc()

	// RUN
	log.Println("🏄 Here we go!")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("❌ Failed to serve: %v", err)
	}
}

func establishPostgresConnection(config Postgres) *sql.DB {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		config.Host,
		config.Port,
		config.DbName,
	)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("❌ Failed to open connection to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("❌ Failed to connect to the database: %v", err)
	}
	log.Printf("🔑 Connected to PostgreSQL")

	return db
}

func listenToPort(port string) net.Listener {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("❌ Failed to listen: %v", err)
	}
	log.Printf("👂 Server is listening at %v", lis.Addr())
	return lis
}

func initializeGrpc() *grpc.Server {
	s := grpc.NewServer()
	delivery.NewGreeterReceiver().Register(s)
	log.Println("💻 Initialized gRPC server")
	return s
}
