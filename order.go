// ไฟล์: order.go
package main

import (
	"fmt"

	"gorm.io/gorm"
)

// * Create (เพิ่มคำสั่งซื้อด้วย Transaction)
// * ตรวจสอบลูกค้า/สินค้า, เช็กสต็อก, คำนวณ TotalPrice, ลดสต็อก
// * และต้องสำเร็จหรือยกเลิกพร้อมกันทั้งหมด (Transaction)
func createOrderDB(db *gorm.DB, customerID uint, productID uint, quantity int) (*Order, error) {
	var newOrder Order

	err := db.Transaction(func(tx *gorm.DB) error {
		var customer Customer
		if err := tx.First(&customer, customerID).Error; err != nil {
			return fmt.Errorf("ไม่พบลูกค้า ID %d", customerID)
		}

		var product Product
		if err := tx.First(&product, productID).Error; err != nil {
			return fmt.Errorf("ไม่พบสินค้า ID %d", productID)
		}

		if product.Quantity < quantity {
			return fmt.Errorf("สต็อกไม่พอ (สินค้าเหลือ %d ชิ้น แต่ต้องการสั่ง %d ชิ้น)", product.Quantity, quantity)
		}
		product.Quantity -= quantity
		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		totalPrice := product.Price * float64(quantity)

		newOrder = Order{
			CustomerID: customerID,
			ProductID:  productID,
			Quantity:   quantity,
			TotalPrice: totalPrice,
			Status:     "Pending",
		}

		if err := tx.Create(&newOrder).Error; err != nil {
			return err
		}

		return nil // ถ้าไม่มี Error อะไรเลย Transaction จะ Commit บันทึกลง DB จริง
	})

	if err != nil {
		return nil, fmt.Errorf("❌ สร้าง Order ไม่สำเร็จ (ระบบ Rollback แล้ว): %v", err)
	}

	fmt.Printf("✅ สร้าง Order ID %d สำเร็จ (ยอดรวม: %.2f บาท)\n", newOrder.ID, newOrder.TotalPrice)
	return &newOrder, nil
}

//* Read Orders

// 1. คำสั่งซื้อของลูกค้า (ตาม CustomerID)
func getOrdersByCustomerDB(db *gorm.DB, customerID uint) []Order {
	var orders []Order
	db.Preload("Customer").Preload("Product").Where("customer_id = ?", customerID).Find(&orders)
	return orders
}

// 2. ตามสถานะ (Pending, Completed, Cancelled)
func getOrdersByStatusDB(db *gorm.DB, status string) []Order {
	var orders []Order
	db.Preload("Customer").Preload("Product").Where("status = ?", status).Find(&orders)
	return orders
}

// 3. ตามช่วงวันที่ (Date Range)
func getOrdersByDateRangeDB(db *gorm.DB, startDate string, endDate string) []Order {
	var orders []Order
	db.Preload("Customer").Preload("Product").Where("created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59").Find(&orders)
	return orders
}

// 4. คำสั่งซื้อพร้อมลูกค้าและสินค้า (With Details)
func getOrdersWithDetailsDB(db *gorm.DB) []Order {
	var orders []Order
	db.Preload("Customer").Preload("Product.Category").Find(&orders)
	return orders
}

// * Update Order Status
func updateOrderStatusDB(db *gorm.DB, orderID uint, newStatus string) error {
	if newStatus != "Pending" && newStatus != "Completed" && newStatus != "Cancelled" {
		return fmt.Errorf("❌ สถานะไม่ถูกต้อง (ต้องเป็น Pending, Completed หรือ Cancelled เท่านั้น)")
	}

	result := db.Model(&Order{}).Where("id = ?", orderID).Update("status", newStatus)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("ไม่พบ Order ID %d", orderID)
	}

	fmt.Printf("🔄 อัปเดตสถานะ Order ID %d เป็น '%s' สำเร็จ\n", orderID, newStatus)
	return nil
}
