package domain

import "time"

type Store struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID      string    `gorm:"type:uuid;not null" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	StoreName   string    `gorm:"not null;unique" json:"store_name"`
	StoreLogo   string    `json:"store_logo"`
	StoreBanner string    `json:"store_banner"`
	StoreDesc   string    `gorm:"type:text" json:"store_desc"`
	Follower    int       `gorm:"default:0" json:"follower"`
	StoreRating float64   `gorm:"default:0" json:"store_rating"`
	CreatedAt   time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}

// type StoreFinance struct {
// 	ID            string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
// 	StoreID       string    `gorm:"type:uuid;not null" json:"store_id"`
// 	TotalRevenue  int64     `gorm:"default:0" json:"total_revenue"`
// 	TotalExpenses int64     `gorm:"default:0" json:"total_expenses"`
// 	NetProfit     int64     `gorm:"default:0" json:"net_profit"`
// 	OrdersCount   int       `gorm:"default:0" json:"orders_count"`
// 	CreatedAt     time.Time `gorm:"default:current_timestamp" json:"created_at"`
// 	UpdatedAt     time.Time `gorm:"default:current_timestamp" json:"updated_at"`
// }

// type Order struct {
// 	ID             string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
// 	StoreID        string    `gorm:"type:uuid;not null" json:"store_id"`
// 	UserID         string    `gorm:"type:uuid;not null" json:"user_id"`
// 	TotalPrice     int64     `gorm:"not null" json:"total_price"`
// 	OrderStatus    string    `gorm:"not null" json:"order_status"`
// 	ShippingStatus string    `gorm:"not null" json:"shipping_status"`
// 	PaymentStatus  string    `gorm:"not null" json:"payment_status"`
// 	CreatedAt      time.Time `gorm:"default:current_timestamp" json:"created_at"`
// 	UpdatedAt      time.Time `gorm:"default:current_timestamp" json:"updated_at"`
// }

// type OrderDetail struct {
// 	ID         string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
// 	OrderID    string    `gorm:"type:uuid;not null" json:"order_id"`
// 	ProductID  string    `gorm:"type:uuid;not null" json:"product_id"`
// 	Quantity   int       `gorm:"not null" json:"quantity"`
// 	UnitPrice  int64     `gorm:"not null" json:"unit_price"`
// 	TotalPrice int64     `gorm:"not null" json:"total_price"`
// 	CreatedAt  time.Time `gorm:"default:current_timestamp" json:"created_at"`
// 	UpdatedAt  time.Time `gorm:"default:current_timestamp" json:"updated_at"`
// }

// type Shipping struct {
// 	ID             string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
// 	OrderID        string    `gorm:"type:uuid;not null" json:"order_id"`
// 	ShippingDate   time.Time `gorm:"default:current_timestamp" json:"shipping_date"`
// 	ShippingCost   int64     `gorm:"not null" json:"shipping_cost"`
// 	TrackingNo     string    `gorm:"not null" json:"tracking_no"`
// 	ShippingStatus string    `gorm:"not null" json:"shipping_status"`
// 	CreatedAt      time.Time `gorm:"default:current_timestamp" json:"created_at"`
// 	UpdatedAt      time.Time `gorm:"default:current_timestamp" json:"updated_at"`
// }

// type Cancellation struct {
// 	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
// 	OrderID     string    `gorm:"type:uuid;not null" json:"order_id"`
// 	Reason      string    `gorm:"type:text" json:"reason"`
// 	CancelledAt time.Time `gorm:"default:current_timestamp" json:"cancelled_at"`
// 	CreatedAt   time.Time `gorm:"default:current_timestamp" json:"created_at"`
// 	UpdatedAt   time.Time `gorm:"default:current_timestamp" json:"updated_at"`
// }

// type Return struct {
// 	ID         string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
// 	OrderID    string    `gorm:"type:uuid;not null" json:"order_id"`
// 	Reason     string    `gorm:"type:text" json:"reason"`
// 	ReturnedAt time.Time `gorm:"default:current_timestamp" json:"returned_at"`
// 	CreatedAt  time.Time `gorm:"default:current_timestamp" json:"created_at"`
// 	UpdatedAt  time.Time `gorm:"default:current_timestamp" json:"updated_at"`
// }
