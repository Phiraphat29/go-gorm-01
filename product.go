package main

import (
	"fmt"

	"gorm.io/gorm"
)

// * Create (เพิ่มสินค้าพร้อมประเภท)
func createProductDB(db *gorm.DB, name string, price float64, quantity int, categoryID uint) (*Product, error) {
	product := Product{
		Name:       name,
		Price:      price,
		Quantity:   quantity,
		CategoryID: categoryID,
	}
	if err := db.Create(&product).Error; err != nil {
		return nil, fmt.Errorf("เพิ่มสินค้าไม่สำเร็จ: %v", err)
	}
	fmt.Printf("✅ เพิ่ม Product สำเร็จ: ID %d - %s\n", product.ID, product.Name)
	return &product, nil
}

// * Read / Conditional Queries
// 1. สินค้าทั้งหมด
func getAllProductsDB(db *gorm.DB) []Product {
	var products []Product
	db.Preload("Category").Find(&products)
	return products
}

// 2. ค้นหาจาก ID
func getProductByIDDB(db *gorm.DB, id uint) (*Product, error) {
	var product Product
	if err := db.Preload("Category").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// 3. ค้นหาชื่อบางส่วนด้วย LIKE
func getProductByNameDB(db *gorm.DB, keyword string) []Product {
	var products []Product
	db.Preload("Category").Where("name LIKE ?", "%"+keyword+"%").Find(&products)
	return products
}

// 4. ตาม CategoryID
func getProductsByCategoryDB(db *gorm.DB, categoryID uint) []Product {
	var products []Product
	db.Preload("Category").Where("category_id = ?", categoryID).Find(&products)
	return products
}

// 5. ช่วง Min-Max
func getProductsByPriceRangeDB(db *gorm.DB, minPrice, maxPrice float64) []Product {
	var products []Product
	db.Preload("Category").Where("price BETWEEN ? AND ?", minPrice, maxPrice).Find(&products)
	return products
}

// 6. สต็อกน้อยกว่าค่ากำหนด
func getProductsLowStockDB(db *gorm.DB, threshold int) []Product {
	var products []Product
	db.Preload("Category").Where("quantity < ?", threshold).Find(&products)
	return products
}

// 7. ประเภทและช่วงราคา
func getProductsByCategoryAndPriceDB(db *gorm.DB, categoryID uint, minPrice, maxPrice float64) []Product {
	var products []Product
	db.Preload("Category").Where("category_id = ? AND price BETWEEN ? AND ?", categoryID, minPrice, maxPrice).Find(&products)
	return products
}

// 8. เรียงราคา ASC/DESC
func getProductsSortedByPriceDB(db *gorm.DB, ascending bool) []Product {
	var products []Product
	orderClause := "price DESC"
	if ascending {
		orderClause = "price ASC"
	}
	db.Preload("Category").Order(orderClause).Find(&products)
	return products
}

// 9. สินค้าแพงที่สุด N รายการ
func getTopExpensiveProductsDB(db *gorm.DB, limit int) []Product {
	var products []Product
	db.Preload("Category").Order("price DESC").Limit(limit).Find(&products)
	return products
}

//* Update
//* (ราคาและสต็อกต้องไม่ติดลบ)

// 1. แก้ไขสินค้าตาม ID
func updateProductByIDDB(db *gorm.DB, id uint, name string, price float64, quantity int, categoryID uint) error {
	if price < 0 || quantity < 0 {
		return fmt.Errorf("ราคาและสต็อกต้องไม่ติดลบ")
	}
	result := db.Model(&Product{}).Where("id = ?", id).Updates(Product{
		Name:       name,
		Price:      price,
		Quantity:   quantity,
		CategoryID: categoryID,
	})
	return result.Error
}

// 2. แก้ไขราคา
func updateProductPriceDB(db *gorm.DB, id uint, newPrice float64) error {
	if newPrice < 0 {
		return fmt.Errorf("ราคาต้องไม่ติดลบ")
	}
	return db.Model(&Product{}).Where("id = ?", id).Update("price", newPrice).Error
}

// 3. เพิ่ม/ลดสต็อก
func updateProductStockDB(db *gorm.DB, id uint, qtyChange int) error {
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return err
	}
	newQty := product.Quantity + qtyChange
	if newQty < 0 {
		return fmt.Errorf("สต็อกคงเหลือติดลบไม่ได้ (ปัจจุบันเหลือ %d)", product.Quantity)
	}
	return db.Model(&Product{}).Where("id = ?", id).Update("quantity", newQty).Error
}

//* Delete & Restore

// 1. Soft Delete สินค้า
func softDeleteProductDB(db *gorm.DB, id uint) error {
	// เช็กความสัมพันธ์ว่ามี Order ค้างอยู่ไหมตามโจทย์
	var orderCount int64
	db.Model(&Order{}).Where("product_id = ?", id).Count(&orderCount)
	if orderCount > 0 {
		return fmt.Errorf("ลบไม่ได้! สินค้านี้มีประวัติคำสั่งซื้อค้างอยู่")
	}
	return db.Delete(&Product{}, id).Error
}

// 2. ลบสินค้าถาวร
func hardDeleteProductDB(db *gorm.DB, id uint) error {
	return db.Unscoped().Delete(&Product{}, id).Error
}

// 3. ดึงสินค้าที่ถูก Soft Delete
func getDeletedProductsDB(db *gorm.DB) []Product {
	var products []Product
	db.Unscoped().Where("deleted_at IS NOT NULL").Find(&products)
	return products
}

// 4. กู้คืนสินค้า
func restoreProductDB(db *gorm.DB, id uint) error {
	return db.Unscoped().Model(&Product{}).Where("id = ?", id).Update("deleted_at", nil).Error
}
