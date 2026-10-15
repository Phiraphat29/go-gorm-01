package main

import (
	"fmt"

	"gorm.io/gorm"
)

// * Create (เพิ่มประเภทสินค้า)
func createCategoryDB(db *gorm.DB, name string, description string) (*Category, error) {
	category := Category{
		Name:        name,
		Description: description,
	}

	if err := db.Create(&category).Error; err != nil {
		return nil, fmt.Errorf("เพิ่มประเภทสินค้าไม่สำเร็จ: %v", err)
	}

	fmt.Printf("✅ เพิ่ม Category สำเร็จ: ID %d - %s\n", category.ID, category.Name)
	return &category, nil
}

// * Delete (ลบประเภทสินค้า + เช็กความสัมพันธ์)
// * ห้ามลบประเภทสินค้าที่ยังมีสินค้าค้างอยู่
func deleteCategoryByIDDB(db *gorm.DB, id uint) error {
	// 1. เช็กก่อนว่ามีสินค้าผูกอยู่กับ Category นี้ไหม
	var productCount int64
	db.Model(&Product{}).Where("category_id = ?", id).Count(&productCount)

	if productCount > 0 {
		return fmt.Errorf("❌ ลบไม่ได้! หมวดหมู่นี้ยังมีสินค้าค้างอยู่ %d ชิ้น ต้องเคลียร์สินค้าออกก่อน", productCount)
	}

	// 2. ถ้าไม่มีสินค้า ให้ทำการลบได้เลย
	if err := db.Delete(&Category{}, id).Error; err != nil {
		return fmt.Errorf("ลบ Category ไม่สำเร็จ: %v", err)
	}

	fmt.Printf("🗑️ ลบ Category ID %d สำเร็จแล้ว\n", id)
	return nil
}
