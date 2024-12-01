package store

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/JinHyeokOh01/go-crwl-server/config"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Initialize(config config.DBConfig) error {

	// 데이터소스 이름 (DB 없이 연결)
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/?parseTime=true",
		config.User,     // MYSQL_USER
		config.Password, // MYSQL_PASSWORD
		config.Host,     // MYSQL_HOST
		config.Port,     // MYSQL_PORT
	)

	var err error
	DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return fmt.Errorf("데이터베이스 연결 실패: %v", err)
	}

	// 연결 테스트
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("데이터베이스 핑 테스트 실패: %v", err)
	}

	log.Println("데이터베이스 연결 성공")

	// 데이터베이스 생성
	err = createDatabase(config.Name)
	if err != nil {
		return fmt.Errorf("데이터베이스 생성 실패: %v", err)
	}

	// 데이터베이스 연결 갱신 (생성된 DB 사용)
	dataSourceNameWithDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.User, config.Password, config.Host, config.Port, config.Name)

	DB, err = sql.Open("mysql", dataSourceNameWithDB)
	if err != nil {
		return fmt.Errorf("데이터베이스 연결 실패: %v", err)
	}

	// 연결 테스트
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("데이터베이스 핑 테스트 실패: %v", err)
	}

	log.Println("데이터베이스 준비 완료")

	// 테이블 생성
	err = createTables()
	if err != nil {
		return fmt.Errorf("테이블 생성 실패: %v", err)
	}

	return nil
}

func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS cse_notices (
            number VARCHAR(255) PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            date VARCHAR(255) NOT NULL,
            link VARCHAR(255) NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS sw_notices (
            number VARCHAR(255) PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            date VARCHAR(255) NOT NULL,
            link VARCHAR(255) NOT NULL
        )`,
	}

	for _, query := range queries {
		_, err := DB.Exec(query)
		if err != nil {
			return fmt.Errorf("테이블 생성 실패: %v", err)
		}
	}
	log.Println("테이블 생성 완료")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// 데이터베이스 생성 함수
func createDatabase(dbName string) error {
	_, err := DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		return fmt.Errorf("데이터베이스 생성 실패: %v", err)
	}
	log.Printf("데이터베이스 %s 생성 완료", dbName)
	return nil
}
