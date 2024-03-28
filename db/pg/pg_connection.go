package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dmitriibb/go-common/logging"
	"github.com/dmitriibb/go-common/utils"
	commonInitializer "github.com/dmitriibb/go-common/utils/initializer"
	_ "github.com/lib/pq"
	"os"
	"strconv"
	"strings"
)

var host string
var portString string
var user string
var password string
var dbname string
var dbInitMode string

var logger = logging.NewLogger("PgConnections")
var initializer = commonInitializer.New(logger)

// TODO connection pool
var initFunc = func() error {
	host = utils.GetEnvProperty(DbHostEnv)
	portString = utils.GetEnvProperty(DbPortEnv)
	user = utils.GetEnvProperty(DbUserEnv)
	password = utils.GetEnvProperty(DbPasswordEnv)
	dbname = utils.GetEnvProperty(DbNameEnv)
	dbInitMode = strings.ToLower(utils.GetEnvProperty(DbInitModeEnv, DbInitModeIgnore))
	logger.Debug("db -> host %s, port %s, dbName %s", host, portString, dbname)
	initDbTables()
	return nil
}

func Init() {
	initFunc := func() error {
		logger.Debug("db -> host %s, port %s, dbName %s", host, portString, dbname)
		initDbTables()
		return nil
	}
	initializer.Init(initFunc)
}

func initDbTables() {
	switch dbInitMode {
	case DbInitModeIgnore:
		return
	case DbInitModeRecreate, DbInitModeUpdate:
		logger.Debug("initDbTables - dbInitMode")
		file, err := os.ReadFile("db_scripts.sql")
		if err != nil {
			panic(err)
		}
		fullScript := string(file)
		scripts := strings.Split(fullScript, ";")
		for i, v := range scripts {
			scripts[i] = strings.TrimSpace(v)
		}
		for _, script := range scripts {
			firstLine := firstNotEmptyString(strings.Split(script, "\n"))
			if len(firstLine) == 0 || strings.HasPrefix(firstLine, "--") {
				continue
			}
			logger.Debug("execute - '%s'", firstLine)
			f := func(db *sql.DB) any {
				res, err := db.Exec(script)
				if err != nil {
					panic(err)
				}
				return res
			}
			UseConnection(f)
		}
	}
}

func firstNotEmptyString(lines []string) string {
	for _, v := range lines {
		if len(v) > 0 {
			return v
		}
	}
	return ""
}

func TestConnectPostgres() {
	port, _ := strconv.Atoi(portString)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	connection, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	err = connection.Ping()
	if err != nil {
		panic(err)
	}
	logger.Info(fmt.Sprintf("Successfully connected to '%v'!", dbname))
}

func UseConnection(f func(db *sql.DB) any) any {
	port, _ := strconv.Atoi(portString)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	connection, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//logger.Debug("db connection open")
	defer func() {
		connection.Close()
		//logger.Debug("db connection close")
	}()
	err = connection.Ping()
	if err != nil {
		panic(err)
	}
	//logger.Debug("execute sql")
	return f(connection)
}
func GetConnection() *sql.DB {
	port, _ := strconv.Atoi(portString)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	connection, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	logger.Debug("db connection open")
	return connection
}

func StartTransaction(ctx context.Context) TxWrapperer {
	logger.Debug("start transaction")
	conn := GetConnection()
	tx, _ := conn.BeginTx(ctx, nil)
	return &TxWrapper{conn, tx}
}

type TxWrapperer interface {
	Commit()
	Rollback()
}

type TxWrapper struct {
	connection *sql.DB
	Tx         *sql.Tx
}

func (txWrapper *TxWrapper) Commit() {
	err := txWrapper.Tx.Commit()
	if err != nil {
		logger.Warn("commit transaction - %s", err.Error())
	} else {
		logger.Debug("commit transaction")
	}
	txWrapper.connection.Close()
}

func (txWrapper *TxWrapper) Rollback() {
	err := txWrapper.Tx.Rollback()
	if err != nil {
		logger.Warn("rollback transaction - %s", err.Error())
	} else {
		logger.Debug("rollback transaction")
	}
	txWrapper.connection.Close()
}
