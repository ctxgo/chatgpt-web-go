package db

type DBInitStatus struct {
	ID          int  `gorm:"primaryKey"`
	Initialized bool `gorm:"initialized"`
}

func (*DBInitStatus) TableName() string {
	return "db_init_status"
}
