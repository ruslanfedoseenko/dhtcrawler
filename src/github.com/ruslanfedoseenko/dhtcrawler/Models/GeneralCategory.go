package Models

type GeneralCategory struct {
	Id            int32             `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Name          string            `gorm:"column:name" json:"value"`
	Icon          string            `gorm:"column:icon" json:"icon"`
	TrainingData  []TrainingPortion `gorm:"many2many:training_data_to_gemeral_groups" json:",omitempty"`
	OrderPosition int               `gorm:"column:parent_id" json:"-"`
	Children      []GeneralCategory `gorm:"ForeignKey:ParentId" json:"data,omitempty"`
	ParentId      int32             `gorm:"column:parent_id" json:"-"`
}

func (GeneralCategory) TableName() string {
	return "general_groups"
}
