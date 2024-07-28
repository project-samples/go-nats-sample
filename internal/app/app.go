package app

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/core-go/health"
	hm "github.com/core-go/health/mongo"
	nh "github.com/core-go/health/nats"
	w "github.com/core-go/mongo/writer"
	"github.com/core-go/mq"
	v "github.com/core-go/mq/validator"
	"github.com/core-go/mq/zap"
	"github.com/core-go/nats"
)

type ApplicationContext struct {
	HealthHandler *health.Handler
	Subscribe     func(ctx context.Context, handle func(context.Context, []byte, map[string]string))
	Handle        func(context.Context, []byte, map[string]string)
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	log.Initialize(cfg.Log)
	client, er1 := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.Uri))
	if er1 != nil {
		log.Error(ctx, "Cannot connect to MongoDB: Error: "+er1.Error())
		return nil, er1
	}
	db := client.Database(cfg.Mongo.Database)

	logError := log.ErrorMsg
	var logInfo func(context.Context, string)
	if log.IsInfoEnable() {
		logInfo = log.InfoMsg
	}

	subscriber, er2 := nats.NewSubscriberByConfig(cfg.Subscriber, logError)
	if er2 != nil {
		log.Error(ctx, "Cannot create a new receiver. Error: "+er2.Error())
		return nil, er2
	}
	validator, er3 := v.NewValidator[*User]()
	if er3 != nil {
		return nil, er3
	}
	errorHandler := mq.NewErrorHandler[*User](logError)
	publisher, er4 := nats.NewPublisherByConfig(*cfg.Publisher)
	if er4 != nil {
		return nil, er4
	}
	writer := w.NewWriter[*User](db, "user")
	handler := mq.NewRetryHandlerByConfig[User](cfg.Retry, writer.Write, validator.Validate, errorHandler.RejectWithMap, nil, publisher.Publish, logError, logInfo)
	mongoChecker := hm.NewHealthChecker(client)
	receiverChecker := nh.NewHealthChecker(cfg.Subscriber.Connection.Url, "nats_subscriber")
	senderChecker := nh.NewHealthChecker(cfg.Publisher.Connection.Url, "nats_publisher")
	healthHandler := health.NewHandler(mongoChecker, receiverChecker, senderChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		Subscribe:     subscriber.Subscribe,
		Handle:        handler.Handle,
	}, nil
}
