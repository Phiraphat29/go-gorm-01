package main

import (
	"fmt"

	"gorm.io/gorm"
)

// * Create (เพิ่มลูกค้า)
func createCustomerDB(db *gorm.DB, name string, email string, phone string) (*Customer, error) {
	customer := Customer{
		Name:  name,
		Email: email,
		Phone: phone,
	}
	if err := db.Create(&customer).Error; err != nil {
		return nil, fmt.Errorf("เพิ่มลูกค้าไม่สำเร็จ (อีเมลอาจซ้ำ): %v", err)
	}
	fmt.Printf("✅ เพิ่ม Customer สำเร็จ: ID %d - %s\n", customer.ID, customer.Name)
	return &customer, nil
}

// * Read (ค้นหาลูกค้าด้วย Email)
func getCustomerByEmailDB(db *gorm.DB, email string) (*Customer, error) {
	var customer Customer
	if err := db.Preload("Orders").Where("email = ?", email).First(&customer).Error; err != nil {
		return nil, fmt.Errorf("ไม่พบลูกค้านี้: %v", err)
	}
	return &customer, nil
}

// * Update (แก้ไขข้อมูลลูกค้า)
func updateCustomerByIDDB(db *gorm.DB, id uint, name string, email string, phone string) error {
	result := db.Model(&Customer{}).Where("id = ?", id).Updates(Customer{
		Name:  name,
		Email: email,
		Phone: phone,
	})
	return result.Error
}

// * Delete (ลบลูกค้า + เช็กความสัมพันธ์)
// * ห้ามลบลูกค้าที่มีประวัติคำสั่งซื้อ
func deleteCustomerByIDDB(db *gorm.DB, id uint) error {
	var orderCount int64
	db.Model(&Order{}).Where("customer_id = ?", id).Count(&orderCount)

	if orderCount > 0 {
		return fmt.Errorf("❌ ลบไม่ได้! ลูกค้ารายนี้มีประวัติคำสั่งซื้อค้างอยู่ %d รายการ", orderCount)
	}

	if err := db.Delete(&Customer{}, id).Error; err != nil {
		return fmt.Errorf("ลบ Customer ไม่สำเร็จ: %v", err)
	}

	fmt.Printf("🗑️ ลบ Customer ID %d สำเร็จแล้ว\n", id)
	return nil
}
