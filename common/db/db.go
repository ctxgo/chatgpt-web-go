package db

import (
	"bufio"
	"bytes"
	"chatgpt-web-new-go/common/config"
	"chatgpt-web-new-go/dao"
	"chatgpt-web-new-go/pkgs/retry"

	sqlFilr "chatgpt-web-new-go/model/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbTypeInitializer = map[string]func(){
	"mysql": initMysql,
}

func Init() {
	dbType := config.Config.Db.Type
	dbInitializer := dbTypeInitializer[dbType]
	dbInitializer()

	// gorm gen init
	dao.SetDefault(config.DB)
}

func initMysql() {
	dbConfig := config.Config.Db
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Name)

	var err error
	retry.Retry(
		func() error {
			config.DB, err = openDB(dsn)
			return err
		},
		retry.WithAttempts(10),
	)
	if err != nil {
		panic(err)
	}
	sqlDB, err := config.DB.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := initializeDBData(config.DB); err != nil {
		panic(err)
	}
	// migrate
	//err = config.DB.AutoMigrate(&user.User{})
	//if err != nil {
	//	panic(err)
	//}
}

func openDB(dsn string) (db *gorm.DB, err error) {

	return gorm.Open(
		mysql.New(mysql.Config{
			DSN:                       dsn,   // data source name
			DefaultStringSize:         256,   // default size for string fields
			DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
			DontSupportRenameIndex:    true,  // drop & create when rename messageDao, rename messageDao not supported before MySQL 5.7, MariaDB
			DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
			SkipInitializeWithVersion: false, // autoconfigure based on currently MySQL version
		}),
		&gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold: time.Second,
					LogLevel:      logger.Info,
					Colorful:      true,
				},
			),
		})
}

func initializeDBData(db *gorm.DB) error {
	err := db.AutoMigrate(&DBInitStatus{})
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		// 尝试获取一个全局锁
		var lockStatus int
		err := tx.Raw("SELECT GET_LOCK('db_init_lock', 10)").Scan(&lockStatus).Error
		if err != nil || lockStatus != 1 {
			return fmt.Errorf("failed to acquire lock: %v", err)
		}
		defer tx.Exec("SELECT RELEASE_LOCK('db_init_lock')") // 确保锁被释放

		var initStatus DBInitStatus
		err = tx.First(&initStatus, 1).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err == nil && initStatus.Initialized {
			fmt.Println("Initialized", config.Config.Db)
			return nil
		}

		if err := executeEmbeddedSQL(tx); err != nil {
			return err
		}

		return tx.Save(&DBInitStatus{ID: 1, Initialized: true}).Error
	})
}

func executeEmbeddedSQL(db *gorm.DB) error {
	scanner := bufio.NewScanner(bytes.NewReader(sqlFilr.EmbeddedSQLData))
	var statement strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "--") { // 忽略空行和注释
			continue
		}
		statement.WriteString(line)
		if strings.HasSuffix(strings.TrimSpace(line), ";") {
			execSQL := statement.String()[:statement.Len()-1] // 移除末尾的分号
			if err := db.Exec(execSQL).Error; err != nil {
				return err
			}
			statement.Reset() // 重置字符串构建器
		}
	}
	return scanner.Err()
}
