//go:generate reform
package storage

import (
	"time"
)

//reform:meters
type Meter struct {
	ID           string        `reform:"id,pk"`
	UserID       string        `reform:"user_id"`
	Name         string        `reform:"name"`
	Address      string        `reform:"address"`
	SerialNumber string        `reform:"serail_number"`
	PeriodCheck  time.Duration `reform:"period_check"`
	Cold         bool          `reform:"is_cold"`
	CreatedAt    time.Time     `reform:"created_at"`
	UpdatedAt    time.Time     `reform:"updated_at"`
}

func (s *Meter) BeforeUpdate() error {
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Meter) Equal(o *Meter) bool {
	if s.ID != o.ID ||
		s.Address != o.Address ||
		s.Name != o.Name ||
		s.PeriodCheck != o.PeriodCheck {

		return false
	}

	return true
}

type LogLevel string

const (
	Error LogLevel = "Error"
	Info  LogLevel = "Info"
)

//reform:logs
type Log struct {
	ID      string    `reform:"id,pk"`
	MeterID string    `reform:"meter_id"`
	Time    time.Time `reform:"time"`
	Level   LogLevel  `reform:"level"`
	Message string    `reform:"message"`
}
