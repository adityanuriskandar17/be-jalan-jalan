package models

type Setting struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Group string `gorm:"type:varchar(50);index" json:"group"` // payment, email, system, general
	Key   string `gorm:"type:varchar(100);uniqueIndex" json:"key"`
	Value string `gorm:"type:text" json:"value"`       // Store as string
	Type  string `gorm:"type:varchar(20)" json:"type"` // bool, string, number
}
