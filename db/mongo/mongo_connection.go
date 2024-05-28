package mongo

import (
	"context"
	"github.com/dmitriibb/go-common/logging"
	"github.com/dmitriibb/go-common/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var logger = logging.NewLogger("MongoConnection")
var initialized = false
var mongoUri = ""
var dbName = ""

func Init() {
	if initialized {
		logger.Warn("already initialized")
		return
	}
	mongoUri = utils.GetEnvProperty(MongoUriEnv)
	dbName = utils.GetEnvProperty(MongoDbNameEnv)

	initialized = true
	logger.Debug("initialized")
}

func GetDbName() string {
	if len(dbName) == 0 {
		panic("DB name is empty")
	}
	return dbName
}

func TestConnection() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	logger.Debug("test connection to mongo")
}

// Deprecated: you must manually close the clientmodel
func GetClient() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}
	return client
}

func UseClient(ctx context.Context, f func(client *mongo.Client) any) any {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	return f(client)
}
