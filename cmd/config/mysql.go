package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/go-sql-driver/mysql"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"log"
	"time"
)

const (
	// 64MB
	MysqlMaxAllowedPacket = 64 * 1024 * 1024
)

func InitMysql() {
	m := Cfg.Mysql
	var err error
	dsn := GetMysqlDSN(m)
	//dsn := m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
	if Db, err = gorm.Open(gmysql.New(gmysql.Config{
		DriverName: "mysql",
		DSN:        dsn,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.LogLevel(m.LogLevel)),
	}); err != nil {
		log.Panicln("init db error", err.Error())
	} else {
		sqlDB, _ := Db.DB()
		sqlDB.SetMaxIdleConns(Cfg.Mysql.MaxIdleConns)
		sqlDB.SetMaxOpenConns(Cfg.Mysql.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}
}

func GetMysqlDSN(cfg Mysql) string {

	isTLS := false
	if cfg.ClientKey != "" && cfg.CACert != "" && cfg.ClientCert != "" {
		isTLS = true
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(cfg.CACert)
		if err != nil {
			log.Fatal(err)
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			log.Fatal("Failed to append PEM.")
		}
		clientCert := make([]tls.Certificate, 0, 1)
		certs, err := tls.LoadX509KeyPair(cfg.ClientCert, cfg.ClientKey)
		if err != nil {
			log.Fatal(err)
		}
		clientCert = append(clientCert, certs)
		err = mysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs:            rootCertPool,
			Certificates:       clientCert,
			InsecureSkipVerify: true,
		})
	}

	config := mysql.Config{
		User:                 cfg.Username,
		Passwd:               cfg.Password,
		Addr:                 cfg.Path,
		Net:                  "tcp",
		DBName:               cfg.Dbname,
		MaxAllowedPacket:     MysqlMaxAllowedPacket,
		AllowNativePasswords: true,
		Params: map[string]string{
			"charset":   "utf8mb4",
			"parseTime": "True",
			"loc":       "Asia/Shanghai",
		},
	}
	if isTLS == true {
		fmt.Printf("mysql TLS mode on, using ca cert: %s, client cert: %s, client key: %s \n",
			cfg.CACert, cfg.ClientCert, cfg.ClientKey)
		config.TLSConfig = "custom"
	}

	dsn := config.FormatDSN()
	return dsn
}
