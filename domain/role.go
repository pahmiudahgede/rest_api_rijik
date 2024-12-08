package domain

type UserRole struct {
	ID       string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique;not null" json:"id"`
	RoleName string `gorm:"unique;not null" json:"roleName"`
	Users    []User `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"users"`
}
