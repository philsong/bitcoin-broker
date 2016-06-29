package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"logger"
	"os"
	"sync"
)

var mutex = &sync.Mutex{}

var sqlconn string
var g_ormDB *gorm.DB

func Init_sqlstr(_sqlconn string) {
	sqlconn = _sqlconn
}

func GetDB() (db *gorm.DB, err error) {
	return getORMDB()
}

func getORMDB() (*gorm.DB, error) {
	if g_ormDB != nil {
		return g_ormDB, nil
	}

	mutex.Lock()
	defer mutex.Unlock()

	//check again
	if g_ormDB != nil {
		return g_ormDB, nil
	}

	db, err := open_db_orm()
	if err == nil {
		g_ormDB = db
	}

	return g_ormDB, err
}

func open_db_orm() (db *gorm.DB, err error) {
	db, err = gorm.Open("postgres", sqlconn)
	if err != nil {
		logger.Errorln("open:", err)
		return
	}
	// Get database connection handle [*sql.DB](http://golang.org/pkg/database/sql/#DB)
	db.DB()

	// Then you could invoke `*sql.DB`'s functions with it
	err = db.DB().Ping()
	if err != nil {
		logger.Errorln("ping:", err)
		return
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// Disable table name's pluralization
	db.SingularTable(true)
	db.LogMode(false)
	db.SetLogger(log.New(os.Stdout, "\r\n", 0))

	// migrate(db)

	return
}

func migrate(db *gorm.DB) {
	// Automating Migration
	db.AutoMigrate(&AmountConfig{}, &Ticker{}, &Depth{})

	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&AmountConfig{})
}
