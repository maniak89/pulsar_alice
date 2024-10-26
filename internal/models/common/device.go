package common

import (
	"fmt"
	"time"
)

type Meter struct {
	ID           string
	UserID       string
	SerailNumber string
	Name         string
	Model        string
	SWVersion    string
	Manufacturer string
	Cold         bool
	Value        float64
	Updated      time.Time
	Changed      time.Time
}

func (r *Meter) String() string {
	return fmt.Sprintf("%s (%s %s)", r.Name, r.Model, r.ID)
}

func (r *Meter) Clone() *Meter {
	if r == nil {
		return nil
	}
	result := Meter{
		ID:           r.ID,
		UserID:       r.UserID,
		Name:         r.Name,
		Model:        r.Model,
		SWVersion:    r.SWVersion,
		Manufacturer: r.Manufacturer,
		SerailNumber: r.SerailNumber,
		Value:        r.Value,
		Cold:         r.Cold,
		Updated:      r.Updated,
		Changed:      r.Changed,
	}

	return &result
}
