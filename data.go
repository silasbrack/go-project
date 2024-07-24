package main

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"os"
)

var data []Company
var db *gorm.DB

type Company struct {
	gorm.Model
	ID      uuid.UUID `gorm:"type:uuid;primary_key;"`
	Company string
	Contact string
	Country string
}

func (base *Company) BeforeCreate(tx *gorm.DB) error {
	uuid, _ := uuid.NewV7()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

func getDsn() string {
	host := os.Getenv("PG_HOST")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	dbname := os.Getenv("PG_DBNAME")
	port := os.Getenv("PG_PORT")
	sslmode := os.Getenv("PG_SSLMODE")
	TimeZone := os.Getenv("PG_TIMEZONE")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, password, dbname, port, sslmode, TimeZone)
	return dsn
}

func init() {
	dsn := getDsn()
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "myschema.",
		},
	})
	if err != nil {
		log.Fatalf("error connecting to postgres: %v", err)
	}
	db.AutoMigrate(&Company{})
	db.Find(&data)
}

func getCompanyByID(id uuid.UUID) Company {
	var result Company
	for _, i := range data {
		if i.ID == id {
			result = i
			break
		}
	}
	return result
}

func updateCompany(company Company) {
	result := []Company{}
	for _, i := range data {
		if i.ID == company.ID {
			i.Company = company.Company
			i.Contact = company.Contact
			i.Country = company.Country
		}
		result = append(result, i)
	}
	data = result
	db.Save(&Company{
		ID:      company.ID,
		Company: company.Company,
		Contact: company.Contact,
		Country: company.Country,
	})
}

func addCompany(company Company) {
	// max := 0
	// for _, i := range data {
	// 	n, _ := strconv.Atoi(i.ID)
	// 	if n > max {
	// 		max = n
	// 	}
	// }
	// max++
	// id := strconv.Itoa(max)

	data = append(data, Company{
		ID:      company.ID,
		Company: company.Company,
		Contact: company.Contact,
		Country: company.Country,
	})
	db.Create(&Company{
		ID:      company.ID,
		Company: company.Company,
		Contact: company.Contact,
		Country: company.Country,
	})
}

func deleteCompany(id uuid.UUID) {
	result := []Company{}
	for _, i := range data {
		if i.ID != id {
			result = append(result, i)
		}
	}
	data = result
	db.Delete(&Company{}, id)
}
