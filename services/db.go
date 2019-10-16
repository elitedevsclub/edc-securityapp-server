package services

import "github.com/jinzhu/gorm"

func Connect(connectUri string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", connectUri)
	if err != nil {
		return nil, err
	}
	if err := db.DB().Ping(); err != nil {
		return nil, err
	}
	runMigration(db)
	return db, err
}

func runMigration(db *gorm.DB) {}
