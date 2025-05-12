package model

type Collector struct {
	ID        string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	UserID    string  `gorm:"not null" json:"userId"`
	User      User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	JobStatus string  `gorm:"default:nonactive" json:"jobstatus"`
	Rating    float32 `gorm:"default:5" json:"rating"`
	AddressId string  `gorm:"not null" json:"address_id"`
	Address   Address `gorm:"foreignKey:AddressId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"address"`
}
