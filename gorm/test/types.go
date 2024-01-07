package types

import (
	"time"
)

type EventAPIName string

type EventStatus string

const (
	Normal                  EventStatus = "Normal"
	Postponed                           = "Postponed"
	Rescheduled                         = "Rescheduled"
	Cancelled                           = "Cancelled"
	Relocated                           = "Relocated"
	RelocatedAndRescheduled             = "RelocatedAndRescheduled"
	Draft                               = "Draft"
	Contingent                          = "Contingent"
	Deleted                             = "Deleted"
)

type CountryCode string

const (
	US CountryCode = "US"
	CA             = "CA"
)

// g2d tablename:categories
type Category struct {
	ID   int64  `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type Money struct {
	Amount       float64 `json:"amount"`
	CurrencyCode string  `json:"currency_code"`
	Display      string  `json:"display"`
}

// g2d tablename:events
type Event struct {
	ID             int64        `json:"id" gorm:"primaryKey"`
	ExternalID     int64        `json:"external_id" gorm:"index"`
	EventAPI       EventAPIName `json:"event_api" gorm:"index"`
	Name           string       `json:"name" django:"display_list"`
	Link           string       `json:"link"`
	StartDate      time.Time    `json:"start_date" gorm:"index"`
	EndDate        *time.Time   `json:"end_date,omitempty"`
	OnSaleDate     *time.Time   `json:"on_sale_date,omitempty" gorm:"index"`
	DateConfirmed  bool         `json:"date_confirmed"`
	TimeConfirmed  bool         `json:"time_confirmed"`
	Type           string       `json:"type" gorm:"index"`
	MinTicketPrice *Money       `json:"min_ticket_price,omitempty" gorm:"embedded;embeddedPrefix:min_ticket_price_"`
	Status         EventStatus  `json:"status"`
	GenreID        int64        `json:"genre_id,omitempty"`
	Genre          *Genre       `json:"genre,omitempty"`
	VenueID        int64        `json:"venue_id,omitempty"`
	Venue          *Venue       `json:"venue,omitempty"`
	Categories     []*Category  `json:"categories" gorm:"many2many:event_categories;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt      *time.Time   `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt      *time.Time   `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	DeletedAt      *time.Time   `json:"deleted_at,omitempty"`
}

// g2d tablename:venues
type Venue struct {
	ID         int64      `json:"id" gorm:"primaryKey"`
	Name       string     `json:"name"`
	Link       string     `json:"link"`
	PostalCode string     `json:"postal_code"`
	Address    string     `json:"address"`
	CityID     int64      `json:"city_id"`
	City       *City      `json:"city"`
	Latitude   float64    `json:"latitude"`
	Longitude  float64    `json:"longitude"`
	StateID    int64      `json:"state_id"`
	State      *State     `json:"state"`
	CountryID  int64      `json:"country_id"`
	Country    *Country   `json:"country"`
	CreatedAt  *time.Time `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

// g2d tablename:genres
type Genre struct {
	ID        int64      `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// g2d tablename:countries
type Country struct {
	ID             int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name           string     `json:"name"`
	ISO3           string     `json:"iso3"`
	ISO2           string     `json:"iso2"`
	NumericCode    string     `json:"numeric_code"`
	PhoneCode      string     `json:"phone_code"`
	Capital        string     `json:"capital"`
	Currency       string     `json:"currency"`
	CurrencyName   string     `json:"currency_name"`
	CurrencySymbol string     `json:"currency_symbol"`
	TLD            string     `json:"tld"`
	Native         string     `json:"native"`
	Region         string     `json:"region"`
	Subregion      string     `json:"subregion"`
	Timezones      string     `json:"timezones"`
	Latitude       float64    `json:"latitude"`
	Longitude      float64    `json:"longitude"`
	Emoji          string     `json:"emoji"`
	EmojiU         string     `json:"emojiU"`
	Active         bool       `json:"active" gorm:"not null"`
	CreatedAt      *time.Time `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// g2d tablename:states
type State struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string     `json:"name"`
	CountryID int64      `json:"country_id"`
	Country   *Country   `json:"country"`
	Code      string     `json:"code"`
	Type      string     `json:"type"`
	Latitude  float64    `json:"latitude"`
	Longitude float64    `json:"longitude"`
	Active    bool       `json:"active" gorm:"not null"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// g2d tablename:cities
type City struct {
	ID         int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string     `json:"name"`
	StateID    int64      `json:"state_id"`
	State      *State     `json:"state"`
	CountryID  int64      `json:"country_id"`
	Country    *Country   `json:"country"`
	Latitude   float64    `json:"latitude"`
	Longitude  float64    `json:"longitude"`
	WikiDataID string     `json:"wikiDataId"`
	Active     bool       `json:"active" gorm:"not null"`
	EventCount int64      `json:"event_cout"`
	CreatedAt  *time.Time `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}
