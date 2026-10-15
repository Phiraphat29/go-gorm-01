// ไฟล์: main.go
package main

import (
	"fmt"

	"gorm.io/gorm"
)

func main() {
	db := connectDB()

	seedData(db)

	fmt.Println("\n--- 🚀 เริ่มทดสอบฟังก์ชันตามใบงาน (เปิดคอมเมนต์ฟังก์ชันที่ต้องการเทสได้เลย) ---")

	// ==========================================
	// * CREATE (เพิ่มข้อมูล)
	// ==========================================

	// -> เพิ่มประเภทสินค้า (Category)
	// createCategoryDB(db, "Smart Home", "Smart devices and IoT")

	// -> เพิ่มสินค้าพร้อมประเภท (Product) [พารามิเตอร์: db, name, price, quantity, categoryID]
	// createProductDB(db, "Smart Bulb", 490, 25, 3)

	// -> เพิ่มลูกค้า (Customer) [พารามิเตอร์: db, name, email, phone]
	// createCustomerDB(db, "Todsapol Chai", "todsapol@email.com", "0899999999")

	// -> เพิ่มคำสั่งซื้อด้วย Transaction (Order) [พารามิเตอร์: db, customerID, productID, quantity]
	// createOrderDB(db, 1, 1, 2)

	// ==========================================
	// * READ / CONDITIONAL QUERIES (ค้นหาและดึงข้อมูล)
	// ==========================================

	// -> สินค้าทั้งหมด
	// getAllProductsDB(db)

	// -> ค้นหาจาก ID
	// getProductByIDDB(db, 1)

	// -> ค้นหาชื่อบางส่วนด้วย LIKE
	// getProductByNameDB(db, "MacBook")

	// -> ตาม CategoryID
	// getProductsByCategoryDB(db, 1)

	// -> ช่วง Min-Max (ช่วงราคา)
	// getProductsByPriceRangeDB(db, 1000, 10000)

	// -> สต็อกน้อยกว่าค่ากำหนด (เช่น น้อยกว่า 5 ชิ้น)
	// getProductsLowStockDB(db, 5)

	// -> ประเภทและช่วงราคา
	// getProductsByCategoryAndPriceDB(db, 1, 10000, 40000)

	// -> เรียงราคา (true = น้อยไปมาก ASC, false =มากไปน้อย DESC)
	// getProductsSortedByPriceDB(db, true)

	// -> สินค้าแพงที่สุด N รายการ (เช่น 3 รายการแรก)
	// getTopExpensiveProductsDB(db, 3)

	// -> ค้นหาลูกค้าด้วย Email
	// getCustomerByEmailDB(db, "somchai@email.com")

	// -> คำสั่งซื้อของลูกค้า (ตาม CustomerID)
	// getOrdersByCustomerDB(db, 1)

	// -> คำสั่งซื้อตามสถานะ ("Pending", "Completed", "Cancelled")
	// getOrdersByStatusDB(db, "Pending")

	// -> คำสั่งซื้อตามช่วงวันที่ (รูปแบบ "YYYY-MM-DD")
	// getOrdersByDateRangeDB(db, "2026-01-01", "2026-12-31")

	// -> คำสั่งซื้อพร้อมรายละเอียดลูกค้าและสินค้า
	// getOrdersWithDetailsDB(db)

	// ==========================================
	// * UPDATE (แก้ไขข้อมูล)
	// ==========================================

	// -> แก้ไขสินค้าตาม ID [พารามิเตอร์: db, id, name, price, quantity, categoryID]
	// updateProductByIDDB(db, 1, "MacBook Air M3", 42900, 8, 1)

	// -> แก้ไขราคา
	// updateProductPriceDB(db, 1, 38900)

	// -> เพิ่ม/ลดสต็อก (ใส่เลขบวกเพื่อเพิ่ม, เลขลบเพื่อลด เช่น -2)
	// updateProductStockDB(db, 1, -2)

	// -> แก้ไขลูกค้าตาม ID
	// updateCustomerByIDDB(db, 1, "Somchai Jaidee (Updated)", "somchai@email.com", "0811111111")

	// -> เปลี่ยนสถานะคำสั่งซื้อ ("Pending", "Completed", "Cancelled")
	// updateOrderStatusDB(db, 1, "Completed")

	// ==========================================
	// * DELETE & RESTORE (ลบและกู้คืน)
	// ==========================================

	// -> Soft Delete สินค้า
	// softDeleteProductDB(db, 15)

	// -> ลบสินค้าถาวร (Hard Delete)
	// hardDeleteProductDB(db, 15)

	// -> ดึงสินค้าที่ถูก Soft Delete
	// getDeletedProductsDB(db)

	// -> กู้คืนสินค้าที่ถูก Soft Delete
	// restoreProductDB(db, 15)

	// -> ลบลูกค้าตาม ID (ติดเงื่อนไขห้ามลบถ้ามีประวัติ Order)
	// deleteCustomerByIDDB(db, 8)

	// -> ลบประเภทสินค้าตาม ID (ติดเงื่อนไขห้ามลบถ้ายังมีสินค้าข้างใน)
	// deleteCategoryByIDDB(db, 4)

	fmt.Println("\n✅ โปรแกรมรันเสร็จสิ้น!")
}

// ฟังก์ชันเติมข้อมูลทดสอบขั้นต่ำตามโจทย์ข้อ 4 (Category: 4, Product: 15, Customer: 8, Order: 12)
func seedData(db *gorm.DB) {
	var catCount int64
	db.Model(&Category{}).Count(&catCount)
	if catCount > 0 {
		fmt.Println("ℹ️ ข้อมูล Seed ถูกเติมไว้เรียบร้อยแล้ว ข้ามขั้นตอนการสร้างซ้ำ")
		return
	}

	fmt.Println("🌱 กำลัง Seed ข้อมูลทดสอบลงฐานข้อมูล...")

	// 1. Category (4 รายการ)
	categories := []Category{
		{Name: "Laptops", Description: "Notebooks and MacBooks"},
		{Name: "Smartphones", Description: "Mobile phones and phablets"},
		{Name: "Accessories", Description: "Computer mice, keyboards, and hubs"},
		{Name: "Monitors", Description: "Displays and screens"},
	}
	db.Create(&categories)

	// 2. Product (15 รายการ)
	products := []Product{
		{Name: "MacBook Air", Price: 39900, Quantity: 10, CategoryID: categories[0].ID},
		{Name: "ASUS Vivobook", Price: 21900, Quantity: 15, CategoryID: categories[0].ID},
		{Name: "Lenovo ThinkPad", Price: 42000, Quantity: 4, CategoryID: categories[0].ID}, // สต็อกน้อยไว้เทส LowStock
		{Name: "iPhone 16", Price: 32900, Quantity: 20, CategoryID: categories[1].ID},
		{Name: "Samsung Galaxy", Price: 29900, Quantity: 8, CategoryID: categories[1].ID},
		{Name: "Xiaomi Redmi", Price: 6990, Quantity: 25, CategoryID: categories[1].ID},
		{Name: "Logitech Mouse", Price: 990, Quantity: 50, CategoryID: categories[2].ID},
		{Name: "Mechanical Keyboard", Price: 2500, Quantity: 12, CategoryID: categories[2].ID},
		{Name: "USB-C Hub", Price: 790, Quantity: 30, CategoryID: categories[2].ID},
		{Name: "Webcam HD", Price: 1290, Quantity: 18, CategoryID: categories[2].ID},
		{Name: "Mouse Pad", Price: 350, Quantity: 40, CategoryID: categories[2].ID},
		{Name: "Dell UltraSharp", Price: 14900, Quantity: 6, CategoryID: categories[3].ID},
		{Name: "LG Gaming Monitor", Price: 8900, Quantity: 9, CategoryID: categories[3].ID},
		{Name: "MSI Curved Monitor", Price: 11500, Quantity: 3, CategoryID: categories[3].ID}, // สต็อกน้อย
		{Name: "Portable Monitor", Price: 4500, Quantity: 14, CategoryID: categories[3].ID},
	}
	db.Create(&products)

	// 3. Customer (8 รายการ)
	customers := []Customer{
		{Name: "Somchai Jaidee", Email: "somchai@email.com", Phone: "0811111111"},
		{Name: "Suda Rukngan", Email: "suda@email.com", Phone: "0822222222"},
		{Name: "Pracha Meesuk", Email: "pracha@email.com", Phone: "0833333333"},
		{Name: "Mali Sodsai", Email: "mali@email.com", Phone: "0844444444"},
		{Name: "Arthit Chuchart", Email: "arthit@email.com", Phone: "0855555555"},
		{Name: "Nisa Phrom", Email: "nisa@email.com", Phone: "0866666666"},
		{Name: "Kanya Wong", Email: "kanya@email.com", Phone: "0877777777"},
		{Name: "Thana Kit", Email: "thana@email.com", Phone: "0888888888"},
	}
	db.Create(&customers)

	// 4. Order (12 รายการ)
	orders := []Order{
		{CustomerID: customers[0].ID, ProductID: products[0].ID, Quantity: 1, TotalPrice: products[0].Price * 1, Status: "Completed"},
		{CustomerID: customers[0].ID, ProductID: products[6].ID, Quantity: 2, TotalPrice: products[6].Price * 2, Status: "Completed"},
		{CustomerID: customers[1].ID, ProductID: products[3].ID, Quantity: 1, TotalPrice: products[3].Price * 1, Status: "Pending"},
		{CustomerID: customers[2].ID, ProductID: products[1].ID, Quantity: 1, TotalPrice: products[1].Price * 1, Status: "Completed"},
		{CustomerID: customers[3].ID, ProductID: products[7].ID, Quantity: 1, TotalPrice: products[7].Price * 1, Status: "Pending"},
		{CustomerID: customers[4].ID, ProductID: products[4].ID, Quantity: 1, TotalPrice: products[4].Price * 1, Status: "Cancelled"},
		{CustomerID: customers[5].ID, ProductID: products[8].ID, Quantity: 3, TotalPrice: products[8].Price * 3, Status: "Completed"},
		{CustomerID: customers[6].ID, ProductID: products[11].ID, Quantity: 1, TotalPrice: products[11].Price * 1, Status: "Pending"},
		{CustomerID: customers[7].ID, ProductID: products[6].ID, Quantity: 1, TotalPrice: products[6].Price * 1, Status: "Completed"},
		{CustomerID: customers[1].ID, ProductID: products[9].ID, Quantity: 2, TotalPrice: products[9].Price * 2, Status: "Completed"},
		{CustomerID: customers[2].ID, ProductID: products[12].ID, Quantity: 1, TotalPrice: products[12].Price * 1, Status: "Pending"},
		{CustomerID: customers[3].ID, ProductID: products[2].ID, Quantity: 1, TotalPrice: products[2].Price * 1, Status: "Completed"},
	}
	db.Create(&orders)

	fmt.Println("✨ Seed ข้อมูลสำเร็จครบถ้วนตามโจทย์เรียบร้อย!")
}
