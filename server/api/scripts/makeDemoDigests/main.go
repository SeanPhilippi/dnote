package main

import (
	"github.com/dnote/dnote/server/api/helpers"
	"github.com/dnote/dnote/server/database"
	"time"
)

func main() {
	database.InitDB()

	db := database.DBConn
	tx := db.Begin()

	userID, err := helpers.GetDemoUserID()
	if err != nil {
		panic(err)
	}

	var d1Notes []database.Note
	var d2Notes []database.Note
	var d3Notes []database.Note
	var d4Notes []database.Note
	var d5Notes []database.Note
	var d6Notes []database.Note
	var d7Notes []database.Note
	var d8Notes []database.Note
	var d9Notes []database.Note
	var d10Notes []database.Note
	var d11Notes []database.Note
	var d12Notes []database.Note

	if err := db.Order("random()").Limit(5).Where("user_id = ?", userID).Find(&d1Notes).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := db.Order("random()").Limit(5).Where("user_id = ?", userID).Find(&d2Notes).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := db.Order("random()").Limit(5).Where("user_id = ?", userID).Find(&d3Notes).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := db.Order("random()").Limit(5).Where("user_id = ?", userID).Find(&d4Notes).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := db.Order("random()").Limit(5).Where("user_id = ?", userID).Find(&d5Notes).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := db.Order("random()").Limit(5).Where("user_id = ?", userID).Find(&d6Notes).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := db.Order("random()").Limit(5).Where("user_id = ?", userID).Find(&d7Notes).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := db.Order("random()").Limit(5).Where("user_id = ?", userID).Find(&d8Notes).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := db.Order("random()").Limit(5).Where("user_id = ?", userID).Find(&d9Notes).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := db.Order("random()").Limit(5).Where("user_id = ?", userID).Find(&d10Notes).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := db.Order("random()").Limit(5).Where("user_id = ?", userID).Find(&d11Notes).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := db.Order("random()").Limit(5).Where("user_id = ?", userID).Find(&d12Notes).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	d1Date := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)
	d1 := database.Digest{
		UserID:    userID,
		Notes:     d1Notes,
		CreatedAt: d1Date,
		UpdatedAt: d1Date,
	}
	if err := tx.Save(&d1).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	d2Date := time.Date(2019, time.February, 4, 0, 0, 0, 0, time.UTC)
	d2 := database.Digest{
		UserID:    userID,
		Notes:     d2Notes,
		CreatedAt: d2Date,
		UpdatedAt: d2Date,
	}
	if err := tx.Save(&d2).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	d3Date := time.Date(2019, time.February, 12, 0, 0, 0, 0, time.UTC)
	d3 := database.Digest{
		UserID:    userID,
		Notes:     d3Notes,
		CreatedAt: d3Date,
		UpdatedAt: d3Date,
	}
	if err := tx.Save(&d3).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	d4Date := time.Date(2019, time.May, 12, 0, 0, 0, 0, time.UTC)
	d4 := database.Digest{
		UserID:    userID,
		Notes:     d4Notes,
		CreatedAt: d4Date,
		UpdatedAt: d4Date,
	}
	if err := tx.Save(&d4).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	d5Date := time.Date(2019, time.March, 10, 0, 0, 0, 0, time.UTC)
	d5 := database.Digest{
		UserID:    userID,
		Notes:     d5Notes,
		CreatedAt: d5Date,
		UpdatedAt: d5Date,
	}
	if err := tx.Save(&d5).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	d6Date := time.Date(2019, time.February, 20, 0, 0, 0, 0, time.UTC)
	d6 := database.Digest{
		UserID:    userID,
		Notes:     d6Notes,
		CreatedAt: d6Date,
		UpdatedAt: d6Date,
	}
	if err := tx.Save(&d6).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	d7Date := time.Date(2019, time.April, 24, 0, 0, 0, 0, time.UTC)
	d7 := database.Digest{
		UserID:    userID,
		Notes:     d7Notes,
		CreatedAt: d7Date,
		UpdatedAt: d7Date,
	}
	if err := tx.Save(&d7).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	d8Date := time.Date(2018, time.December, 6, 0, 0, 0, 0, time.UTC)
	d8 := database.Digest{
		UserID:    userID,
		Notes:     d8Notes,
		CreatedAt: d8Date,
		UpdatedAt: d8Date,
	}
	if err := tx.Save(&d8).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	d9Date := time.Date(2018, time.November, 2, 0, 0, 0, 0, time.UTC)
	d9 := database.Digest{
		UserID:    userID,
		Notes:     d9Notes,
		CreatedAt: d9Date,
		UpdatedAt: d9Date,
	}
	if err := tx.Save(&d9).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	d10Date := time.Date(2018, time.October, 12, 0, 0, 0, 0, time.UTC)
	d10 := database.Digest{
		UserID:    userID,
		Notes:     d10Notes,
		CreatedAt: d10Date,
		UpdatedAt: d10Date,
	}
	if err := tx.Save(&d10).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	d11Date := time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC)
	d11 := database.Digest{
		UserID:    userID,
		Notes:     d11Notes,
		CreatedAt: d11Date,
		UpdatedAt: d11Date,
	}
	if err := tx.Save(&d11).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	d12Date := time.Date(2018, time.May, 17, 0, 0, 0, 0, time.UTC)
	d12 := database.Digest{
		UserID:    userID,
		Notes:     d12Notes,
		CreatedAt: d12Date,
		UpdatedAt: d12Date,
	}
	if err := tx.Save(&d12).Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()
}
