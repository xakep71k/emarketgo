package emarket

type Magazine struct {
	ID          string `bson:"_id,omitempty" json:"id"`
	Title       string `bson:"title" json:"title"`
	Price       int    `bson:"price" json:"price"`
	Thumb       []byte `bson:"thumb" json:"thumb"`
	Enable      bool   `bson:"enable" json:"enable"`
	Description string `bson:"description" json:"description"`
	Quantity    int    `bson:"quantity" json:"quantity"`
	OldID       int    `bson:"oldid" json:"oldid"`
	OldImgName  string `bson:"oldimgfile" json:"oldimgfile"`
	PageNum     int    `bson:"-" json:"-"`
}

type MagazieService interface {
	Find() ([]*Magazine, error)
}
