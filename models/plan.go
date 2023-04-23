package models

import "gorm.io/gorm"

type Plan struct {
	ID     uint   `gorm:"primarykey"`
	UID    string `json:"UID" form:"UID" query:"UID"`
	Plan   string `json:"Plan" form:"Plan" query:"Plan"`
	Date   string `json:"Date" form:"Date" query:"Date"`
	Start  string `json:"Start" form:"Start" query:"Start"`
	End    string `json:"End" form:"End" query:"End"`
	Status int    `json:"Status" form:"Status" query:"Status"`
}

type DateTime struct {
	Start string
	End   string
}

type Week struct {
	Week int
}

type Month struct {
	Month int
}

type Year struct {
	Year int
}

type Username struct {
	Username string
}

func CreatePlan(db *gorm.DB, Plan *Plan) (err error) {
	err = db.Table("plans").Create(Plan).Error
	if err != nil {
		return err
	}
	return nil
}

func GetPlan(db *gorm.DB, Plan *Plan, id int) (err error) {
	err = db.Table("plans").Where("id = ?", id).First(Plan).Error
	if err != nil {
		return err
	}
	return nil
}

func GetPlansByUser(db *gorm.DB, Plan *[]Plan, id string) (err error) {
	err = db.Table("plans").Where("UID = ?", id).Find(Plan).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdatePlan(db *gorm.DB, Plan *Plan) (err error) {
	db.Table("plans").Save(Plan)
	return nil
}

func DeletePlan(db *gorm.DB, Plan *Plan, id int) (err error) {
	db.Table("plans").Where("id = ?", id).Delete(Plan)
	return nil
}

func GetDateTime(db *gorm.DB, date string, id int, uid string) ([]DateTime, error) {
	var results []DateTime
	err := db.Table("plans").Select("start, end").Where("date = ? AND id != ? AND UID = ?", date, id, uid).Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func GetWeekly(db *gorm.DB, Plan *[]Plan, week, year int, uid string) (err error) {
	err = db.Table("plans").Where("WEEK(STR_TO_DATE(date, '%Y-%m-%d'), 1) = ? AND YEAR(STR_TO_DATE(date, '%Y-%m-%d')) = ? AND UID = ?", week, year, uid).Find(Plan).Error
	if err != nil {
		return err
	}
	return nil
}

func GetMonthly(db *gorm.DB, Plan *[]Plan, month, year int, uid string) (err error) {
	err = db.Table("plans").Where("MONTH(STR_TO_DATE(date, '%Y-%m-%d')) = ? AND YEAR(STR_TO_DATE(date, '%Y-%m-%d')) = ? AND UID = ?", month, year, uid).Find(Plan).Error
	if err != nil {
		return err
	}
	return nil
}

func GetYears(db *gorm.DB, Year *[]Year, uid string) (err error) {
	err = db.Table("plans").Select("DATE_FORMAT(`date`, '%Y') as `year`").Group("`year`").Where("UID = ?", uid).Order("`year` ASC").Find(Year).Error
	if err != nil {
		return err
	}
	return nil
}

func GetMonthsByYear(db *gorm.DB, Month *[]Month, year int, uid string) (err error) {
	err = db.Table("plans").Select("DATE_FORMAT(`date`, '%m') as `month`").Group("`month`").Where("YEAR(STR_TO_DATE(date, '%Y-%m-%d')) = ? AND UID = ?", year, uid).Order("`month` ASC").Find(Month).Error
	if err != nil {
		return err
	}
	return nil
}

func GetWeeksByYear(db *gorm.DB, Week *[]Week, year int, uid string) (err error) {
	err = db.Table("plans").Select("DATE_FORMAT(`date`, '%u') as `week`").Group("`week`").Where("YEAR(STR_TO_DATE(date, '%Y-%m-%d')) = ? AND UID = ?", year, uid).Order("`week` ASC").Find(Week).Error
	if err != nil {
		return err
	}
	return nil
}
