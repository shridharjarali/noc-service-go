package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"digit-oss/noc-services/internal/config"
	nocpostgres "digit-oss/noc-services/internal/repository/postgres"
	"digit-oss/noc-services/internal/service"
	"digit-oss/noc-services/internal/service/notification"
	nochttp "digit-oss/noc-services/internal/transport/http"
	"digit-oss/noc-services/internal/transport/kafka"
	"digit-oss/noc-services/internal/validator"
	wfintegrator "digit-oss/noc-services/internal/workflow"
)

func main() {
	// ── 1. Load config via viper ─────────────────────────────────────────
	log.Println("[main] loading configuration...")
	cfg := loadConfig()

	// ── 2. Connect to PostgreSQL via GORM ────────────────────────────────
	log.Println("[main] connecting to PostgreSQL via GORM...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dsn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable search_path=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser, cfg.DBPassword, cfg.DBSchema,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("[main] failed to connect to PostgreSQL via GORM: %v", err)
	}
	// Configure underlying connection pool.
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[main] failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()
	log.Printf("[main] PostgreSQL connected via GORM: %s:%d/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)

	// ── 3. Init repositories ─────────────────────────────────────────────
	log.Println("[main] initializing repositories...")
	httpClient := &http.Client{Timeout: 30 * time.Second}

	svcRequestRepo := &nocpostgres.ServiceRequestRepository{Client: httpClient}
	idGenRepo := &nocpostgres.IdGenRepository{Client: httpClient, Cfg: cfg}

	// ── 4. Init Kafka producer ───────────────────────────────────────────
	log.Printf("[main] connecting to Kafka brokers: %v", cfg.KafkaBrokers)
	kafkaProducer, err := kafka.NewKafkaProducer(cfg.KafkaBrokers)
	if err != nil {
		log.Fatalf("[main] failed to create Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()
	log.Println("[main] Kafka producer ready")

	nocRepo := &nocpostgres.NocRepository{
		DB:       db,
		Cfg:      cfg,
		Producer: kafkaProducer,
	}

	// ── 5. Init services ─────────────────────────────────────────────────
	log.Println("[main] initializing services...")

	enrichmentSvc := &service.EnrichmentService{
		Cfg:       cfg,
		IdGenRepo: idGenRepo,
	}

	userSvc := &service.UserService{
		Cfg:        cfg,
		SvcRequest: svcRequestRepo,
	}

	wfIntegrator := &wfintegrator.WorkflowIntegrator{
		Cfg:        cfg,
		SvcRequest: svcRequestRepo,
	}

	wfService := &wfintegrator.WorkflowService{
		Cfg:        cfg,
		SvcRequest: svcRequestRepo,
	}

	nocValidator := &validator.NOCValidator{
		Cfg: cfg,
	}

	nocService := &service.NOCServiceImpl{
		Cfg:          cfg,
		Repo:         nocRepo,
		Enrichment:   enrichmentSvc,
		Validator:    nocValidator,
		WfIntegrator: wfIntegrator,
		WfService:    wfService,
		SvcRequest:   svcRequestRepo,
	}

	notificationSvc := &notification.NOCNotificationService{
		Cfg:         cfg,
		Producer:    kafkaProducer,
		SvcRequest:  svcRequestRepo,
		UserService: userSvc,
	}

	// ── 6. Init Kafka consumer group ─────────────────────────────────────
	log.Println("[main] starting Kafka consumer...")
	nocConsumer := &kafka.NOCConsumer{
		Cfg:                 cfg,
		NotificationService: notificationSvc,
	}

	// ── 7. Init HTTP router (gin) ────────────────────────────────────────
	log.Println("[main] setting up HTTP router...")
	router := nochttp.NewRouter(nocService, cfg.ServerContext)

	// ── 8. Start Kafka consumer in goroutine ─────────────────────────────
	if err := nocConsumer.Start(ctx, cfg.KafkaBrokers, cfg.KafkaConsumerGroup); err != nil {
		log.Fatalf("[main] failed to start Kafka consumer: %v", err)
	}

	// ── 9. Start HTTP server ─────────────────────────────────────────────
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Graceful shutdown: listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("[main] NOC service listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[main] HTTP server error: %v", err)
		}
	}()

	// Block until shutdown signal
	sig := <-quit
	log.Printf("[main] received signal %v, shutting down...", sig)

	// Cancel context → stops Kafka consumer
	cancel()

	// Shutdown HTTP server with 10s timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("[main] HTTP server shutdown error: %v", err)
	}

	// Close Kafka producer (deferred above)
	// Close DB pool (deferred above)

	log.Println("[main] shutdown complete")

	// TODO: Verify Kafka broker connectivity in staging environment
	// TODO: Verify MDMS/workflow/idgen service URLs in deployment config
	// TODO: Add Prometheus metrics endpoint if required
	// TODO: Add request tracing (OpenTelemetry) if required
}

// loadConfig reads configs/app.yaml and binds env overrides via viper.
func loadConfig() *config.Config {
	v := viper.New()
	v.SetConfigName("app")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath("/configs")
	v.AddConfigPath(".")

	// Allow env overrides: DB_HOST, KAFKA_BROKERS, etc.
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind alternative environment variable names so the service can pick
	// up Postgres settings coming from scripts that export PG_* instead of DB_*.
	_ = v.BindEnv("db.port", "DB_PORT", "PG_PORT")
	_ = v.BindEnv("db.user", "DB_USER", "PG_USER")
	_ = v.BindEnv("db.password", "DB_PASSWORD", "PG_PASSWORD")
	_ = v.BindEnv("db.host", "DB_HOST", "PG_HOST")
	_ = v.BindEnv("db.name", "DB_NAME", "PG_DATABASE")

	// Kafka bootstrap server names commonly used in shell scripts / Java properties
	_ = v.BindEnv("kafka.brokers", "KAFKA_BROKERS", "KAFKA_BOOTSTRAP_SERVERS", "KAFKA_CONFIG_BOOTSTRAP_SERVER_CONFIG")
	_ = v.BindEnv("kafka.consumer.group", "KAFKA_CONSUMER_GROUP", "KAFKA_CONSUMER")

	// External services that Java core services often set via properties
	_ = v.BindEnv("idgen.host", "EGOV_IDGEN_HOST", "IDGEN_HOST", "EGOV_IDGEN")
	_ = v.BindEnv("workflow.host", "WORKFLOW_CONTEXT_PATH", "WORKFLOW_HOST", "WF_HOST")
	_ = v.BindEnv("mdms.host", "EGOV_MDMS_HOST", "MDMS_HOST")
	_ = v.BindEnv("user.host", "EGOV_USER_HOST", "USER_HOST")
	_ = v.BindEnv("localization.host", "EGOV_LOCALIZATION_HOST", "LOCALIZATION_HOST")

	// Server port override
	_ = v.BindEnv("server.port", "SERVER_PORT")

	if err := v.ReadInConfig(); err != nil {
		log.Printf("[config] warning: %v — falling back to env vars only", err)
	}

	return &config.Config{
		// Server
		ServerPort:    v.GetInt("server.port"),
		ServerContext: v.GetString("server.context"),

		// Database
		DBHost:     v.GetString("db.host"),
		DBPort:     v.GetInt("db.port"),
		DBName:     v.GetString("db.name"),
		DBUser:     v.GetString("db.user"),
		DBPassword: v.GetString("db.password"),
		DBSchema:   v.GetString("db.schema"),

		// Kafka
		KafkaBrokers:       v.GetStringSlice("kafka.brokers"),
		KafkaConsumerGroup: v.GetString("kafka.consumer.group"),

		// Persister topics
		SaveTopic:           v.GetString("persister.save.topic"),
		UpdateTopic:         v.GetString("persister.update.topic"),
		UpdateWorkflowTopic: v.GetString("persister.updateWorkflow.topic"),

		// ID-gen
		IdGenHost:              v.GetString("idgen.host"),
		IdGenPath:              v.GetString("idgen.path"),
		ApplicationNoIdgenName: v.GetString("idgen.applicationId"),

		// User service
		UserHost:           v.GetString("user.host"),
		UserContextPath:    v.GetString("user.contextPath"),
		UserSearchEndpoint: v.GetString("user.searchEndpoint"),

		// Workflow
		WfHost:                      v.GetString("workflow.host"),
		WfTransitionPath:            v.GetString("workflow.transitionPath"),
		WfBusinessServiceSearchPath: v.GetString("workflow.businessServiceSearchPath"),
		WfProcessPath:               v.GetString("workflow.processPath"),

		// MDMS
		MdmsHost:     v.GetString("mdms.host"),
		MdmsEndPoint: v.GetString("mdms.endpoint"),

		// BPA
		BpaHost:           v.GetString("bpa.host"),
		BpaContextPath:    v.GetString("bpa.contextPath"),
		BpaSearchEndpoint: v.GetString("bpa.searchEndpoint"),

		// SMS / Notification
		SmsNotifTopic:              v.GetString("notification.sms.topic"),
		IsSMSEnabled:               v.GetBool("notification.sms.enabled"),
		LocalizationHost:           v.GetString("localization.host"),
		LocalizationContextPath:    v.GetString("localization.contextPath"),
		LocalizationSearchEndpoint: v.GetString("localization.searchEndpoint"),
		IsLocalizationStateLevel:   v.GetBool("localization.stateLevel"),

		// Pagination
		DefaultOffset:  v.GetInt("pagination.defaultOffset"),
		DefaultLimit:   v.GetInt("pagination.defaultLimit"),
		MaxSearchLimit: v.GetInt("pagination.maxLimit"),

		// NOC flags
		NocOfflineDocRequired: v.GetBool("noc.offlineDocRequired"),
		IsFuzzyEnabled:        v.GetBool("noc.fuzzySearchEnabled"),
	}
}
