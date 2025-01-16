package models

type Conversion struct {
    ID        int64 `gorm:"primary_key;auto_increment"`
    UserID    int64
    Source    string
    Destination string
    Status    string
}
