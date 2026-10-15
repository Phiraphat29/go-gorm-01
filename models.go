package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Products    []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

type Product struct {
	gorm.Model
	Name       string   `gorm:"type:varchar(150);not null" json:"name"`
	Price      float64  `gorm:"not null" json:"price"`
	Quantity   int      `gorm:"not null" json:"quantity"`
	CategoryID uint     `gorm:"not null" json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Orders     []Order  `gorm:"foreignKey:ProductID" json:"orders,omitempty"`
}

type Customer struct {
	gorm.Model
	Name   string  `gorm:"type:varchar(150);not null" json:"name"`
	Email  string  `gorm:"type:varchar(150);uniqueIndex;not null" json:"email"`
	Phone  string  `gorm:"type:varchar(20)" json:"phone"`
	Orders []Order `gorm:"foreignKey:CustomerID" json:"orders,omitempty"`
}

type Order struct {
	gorm.Model
	CustomerID uint     `gorm:"not null" json:"customer_id"`
	Customer   Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	ProductID  uint     `gorm:"not null" json:"product_id"`
	Product    Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity   int      `gorm:"not null" json:"quantity"`
	TotalPrice float64  `gorm:"not null" json:"total_price"`
	Status     string   `gorm:"type:varchar(20);default:'Pending'" json:"status"` // Pending, Completed, Cancelled
}

func connectDB() *gorm.DB {
	// 1. ต่อเข้า MySQL เปล่าๆ ก่อนเพื่อสั่ง CREATE DATABASE IF NOT EXISTS
	rootDSN := "root:@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
	tempDB, err := gorm.Open(mysql.Open(rootDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("ต่อ XAMPP MySQL ไม่ติด เช็กว่ากด Start MySQL หรือยัง:", err)
	}
	tempDB.Exec("CREATE DATABASE IF NOT EXISTS store_db CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;")

	// 2. ต่อเข้าฐานข้อมูล store_db จริง
	storeDSN := "root:@tcp(127.0.0.1:3306)/store_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(storeDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("เข้าฐานข้อมูล store_db ไม่ได้:", err)
	}

	// 3. AutoMigrate สร้างตารางทั้ง 4
	err = db.AutoMigrate(&Category{}, &Product{}, &Customer{}, &Order{})
	if err != nil {
		log.Fatal("AutoMigrate พัง:", err)
	}

	fmt.Println("✅ เชื่อมต่อ store_db และ AutoMigrate ตารางทั้ง 4 สำเร็จ!")
	return db
}
