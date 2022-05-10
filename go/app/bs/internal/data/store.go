package data

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

var ErrNoRowsUpdated = errors.New("no rows updated")
var ErrNoRowsDeleted = errors.New("no rows deleted")

type Bits uint32

const (
	FLAG_CUSTOMIZED Bits = 1 << iota
	FLAG_DISABLED
	FLAG_FRAUD_DETECTION
	FLAG_NO_REFERRER
)

func (b Bits) Set(flag Bits, enabled bool) Bits {
	if enabled {
		return b | flag
	}
	return b &^ flag
}

func (b Bits) Has(flag Bits) bool {
	return b&flag != 0
}

// StringSlice is a customized database type which is used to
// store string slice as string in database. It implements scanner
// and valuer interface
type StringSlice []string

func (ss *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*ss = StringSlice{}
	} else {
		switch value.(type) {
		case []byte:
			if len(value.([]byte)) == 0 {
				*ss = StringSlice{}
			} else {
				json.Unmarshal(value.([]byte), ss)
			}
		case string:
			if len(value.(string)) == 0 {
				*ss = StringSlice{}
			} else {
				json.Unmarshal([]byte(value.(string)), ss)
			}
		}
	}
	return nil
}

func (ss StringSlice) Value() (driver.Value, error) {
	v, _ := json.Marshal(ss)
	return string(v), nil
}

// create an empty array to mongodb in case StringSlice is nil
func (ss StringSlice) MarshalBSONValue() (bsontype.Type, []byte, error) {
	bsonArr := bsonx.Arr{}
	if ss != nil {
		for _, s := range ss {
			bsonArr = append(bsonArr, bsonx.String(s))
		}
	}
	return bsonArr.MarshalBSONValue()
}

// ShortUrl
type ShortUrl struct {
	Alias       string      `json:"alias"        bson:"alias"`
	Url         string      `json:"url"          bson:"url"`
	Oid         string      `json:"oid"          bson:"oid"`
	Title       string      `json:"title"        bson:"title"`
	Tags        StringSlice `json:"tags"         bson:"tags"`
	Flags       Bits        `json:"flags"        bson:"flags"`
	UtmSource   string      `json:"utm_source"   bson:"utm_source"`
	UtmMedium   string      `json:"utm_medium"   bson:"utm_medium"`
	UtmCampaign string      `json:"utm_campaign" bson:"utm_campaign"`
	UtmTerm     string      `json:"utm_term"     bson:"utm_term"`
	UtmContent  string      `json:"utm_content"  bson:"utm_content"`
	CreatedAt   time.Time   `json:"created_at"   bson:"created_at"`
}

// UpdateShortUrl
type UpdateShortUrl struct {
	Url         string      `json:"url"          bson:"url"`
	Title       string      `json:"title"        bson:"title"`
	Tags        StringSlice `json:"tags"         bson:"tags"`
	Flags       Bits        `json:"flags"        bson:"flags"`
	UtmSource   string      `json:"utm_source"   bson:"utm_source"`
	UtmMedium   string      `json:"utm_medium"   bson:"utm_medium"`
	UtmCampaign string      `json:"utm_campaign" bson:"utm_campaign"`
	UtmTerm     string      `json:"utm_term"     bson:"utm_term"`
	UtmContent  string      `json:"utm_content"  bson:"utm_content"`
}

func (s ShortUrl) Disabled() bool {
	return s.Flags.Has(FLAG_DISABLED)
}

func (s ShortUrl) FraudDetection() bool {
	return s.Flags.Has(FLAG_FRAUD_DETECTION)
}

func (s ShortUrl) NoReferrer() bool {
	return s.Flags.Has(FLAG_NO_REFERRER)
}

// IShortUrlStore interface
type IShortUrlStore interface {
	Open() error
	Close() error
	CreateShortUrl(alias string, customized bool, url string, oid string) error
	DeleteShortUrl(alias string, oid string) error
	GetShortUrl(alias string) (*ShortUrl, error)
	GetShortUrlByOid(alias string, oid string) (*ShortUrl, error)
	ListShortUrl(oid string, limit int64, offset int64) ([]*ShortUrl, error)
	UpdateShortUrl(alias string, oid string, updateShortUrl UpdateShortUrl) error
}
