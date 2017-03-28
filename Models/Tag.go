package Models


type Tag struct {
	Id uint32
	Tag string
}



func (Tag) TableName() string {
	return "tags"
}