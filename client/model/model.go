package model

import (
	"encoding/json"
	"github.com/Xrefullx/YanDip/client/pkg"
	"github.com/google/uuid"
	"log"
)

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"

	SecretTypes = map[string]int{
		"CARD":   1,
		"AUTH":   2,
		"TEXT":   3,
		"BINARY": 4,
	}

	SecretStatuses = map[string]int{
		"NEW":     1,
		"EDITED":  2,
		"ACTUAL":  3,
		"DELETED": 4,
	}
)

type Info struct {
	TypeID      int    `json:"type_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Informer interface {
	GetInfo() Info
}

// Cardholder name
// PAN
// Expiration date
// Service code
type Card struct {
	Info
	CardholderName  string `json:"cardholder"`
	PAN             string `json:"pan"`
	ExpirationMonth int    `json:"expiration_month"`
	ExpirationYear  int    `json:"expiration_year"`
	ServiceCode     int    `json:"code"`
}

func (c *Card) GetInfo() Info {
	return c.Info
}

var _ Informer = (*Card)(nil)

type Auth struct {
	Info
	Login    string `json:"login"`
	Password string `json:"password"`
}
type Text struct {
	Info
	Text string `json:"text"`
}

func (a *Auth) GetInfo() Info {
	return a.Info
}

type Binary struct {
	Info
	Data        []byte
	ContentType string
	Filename    string
}

type Secret struct {
	Info

	ID        int64
	SecretID  uuid.UUID
	SecretVer int
	StatusID  int
	TimeStamp int64

	SecretData string
}
type SecretMeta struct {
	ID        int64
	SecretID  uuid.UUID
	SecretVer int
	StatusID  int
	TimeStamp int64
}

func (s *Info) FromEncodedData(enc string, masterKey string) error {
	decData, err := pkg.Decode(enc, masterKey)
	if err != nil {
		log.Fatal("error:", err)
	}

	return json.Unmarshal(decData, s)
}
