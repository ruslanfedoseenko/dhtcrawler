package Models

type TrainingPortion struct {
	Id   int32  `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Data string `gorm:"column:token"`
}

func (TrainingPortion) TableName() string {
	return "training_data"
}
