// Code generated by gopkg.in/reform.v1. DO NOT EDIT.

package storage

import (
	"fmt"
	"strings"

	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/parse"
)

type meterTableType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *meterTableType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("meters").
func (v *meterTableType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *meterTableType) Columns() []string {
	return []string{
		"id",
		"user_id",
		"name",
		"address",
		"serail_number",
		"is_cold",
		"created_at",
		"updated_at",
	}
}

// NewStruct makes a new struct for that view or table.
func (v *meterTableType) NewStruct() reform.Struct {
	return new(Meter)
}

// NewRecord makes a new record for that table.
func (v *meterTableType) NewRecord() reform.Record {
	return new(Meter)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *meterTableType) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// MeterTable represents meters view or table in SQL database.
var MeterTable = &meterTableType{
	s: parse.StructInfo{
		Type:    "Meter",
		SQLName: "meters",
		Fields: []parse.FieldInfo{
			{Name: "ID", Type: "string", Column: "id"},
			{Name: "UserID", Type: "string", Column: "user_id"},
			{Name: "Name", Type: "string", Column: "name"},
			{Name: "Address", Type: "string", Column: "address"},
			{Name: "SerialNumber", Type: "string", Column: "serail_number"},
			{Name: "Cold", Type: "bool", Column: "is_cold"},
			{Name: "CreatedAt", Type: "time.Time", Column: "created_at"},
			{Name: "UpdatedAt", Type: "time.Time", Column: "updated_at"},
		},
		PKFieldIndex: 0,
	},
	z: new(Meter).Values(),
}

// String returns a string representation of this struct or record.
func (s Meter) String() string {
	res := make([]string, 8)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "UserID: " + reform.Inspect(s.UserID, true)
	res[2] = "Name: " + reform.Inspect(s.Name, true)
	res[3] = "Address: " + reform.Inspect(s.Address, true)
	res[4] = "SerialNumber: " + reform.Inspect(s.SerialNumber, true)
	res[5] = "Cold: " + reform.Inspect(s.Cold, true)
	res[6] = "CreatedAt: " + reform.Inspect(s.CreatedAt, true)
	res[7] = "UpdatedAt: " + reform.Inspect(s.UpdatedAt, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Meter) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.UserID,
		s.Name,
		s.Address,
		s.SerialNumber,
		s.Cold,
		s.CreatedAt,
		s.UpdatedAt,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Meter) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.UserID,
		&s.Name,
		&s.Address,
		&s.SerialNumber,
		&s.Cold,
		&s.CreatedAt,
		&s.UpdatedAt,
	}
}

// View returns View object for that struct.
func (s *Meter) View() reform.View {
	return MeterTable
}

// Table returns Table object for that record.
func (s *Meter) Table() reform.Table {
	return MeterTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *Meter) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *Meter) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *Meter) HasPK() bool {
	return s.ID != MeterTable.z[MeterTable.s.PKFieldIndex]
}

// SetPK sets record primary key, if possible.
//
// Deprecated: prefer direct field assignment where possible: s.ID = pk.
func (s *Meter) SetPK(pk interface{}) {
	reform.SetPK(s, pk)
}

// check interfaces
var (
	_ reform.View   = MeterTable
	_ reform.Struct = (*Meter)(nil)
	_ reform.Table  = MeterTable
	_ reform.Record = (*Meter)(nil)
	_ fmt.Stringer  = (*Meter)(nil)
)

type logTableType struct {
	s parse.StructInfo
	z []interface{}
}

// Schema returns a schema name in SQL database ("").
func (v *logTableType) Schema() string {
	return v.s.SQLSchema
}

// Name returns a view or table name in SQL database ("logs").
func (v *logTableType) Name() string {
	return v.s.SQLName
}

// Columns returns a new slice of column names for that view or table in SQL database.
func (v *logTableType) Columns() []string {
	return []string{
		"id",
		"meter_id",
		"time",
		"level",
		"message",
	}
}

// NewStruct makes a new struct for that view or table.
func (v *logTableType) NewStruct() reform.Struct {
	return new(Log)
}

// NewRecord makes a new record for that table.
func (v *logTableType) NewRecord() reform.Record {
	return new(Log)
}

// PKColumnIndex returns an index of primary key column for that table in SQL database.
func (v *logTableType) PKColumnIndex() uint {
	return uint(v.s.PKFieldIndex)
}

// LogTable represents logs view or table in SQL database.
var LogTable = &logTableType{
	s: parse.StructInfo{
		Type:    "Log",
		SQLName: "logs",
		Fields: []parse.FieldInfo{
			{Name: "ID", Type: "string", Column: "id"},
			{Name: "MeterID", Type: "string", Column: "meter_id"},
			{Name: "Time", Type: "time.Time", Column: "time"},
			{Name: "Level", Type: "LogLevel", Column: "level"},
			{Name: "Message", Type: "string", Column: "message"},
		},
		PKFieldIndex: 0,
	},
	z: new(Log).Values(),
}

// String returns a string representation of this struct or record.
func (s Log) String() string {
	res := make([]string, 5)
	res[0] = "ID: " + reform.Inspect(s.ID, true)
	res[1] = "MeterID: " + reform.Inspect(s.MeterID, true)
	res[2] = "Time: " + reform.Inspect(s.Time, true)
	res[3] = "Level: " + reform.Inspect(s.Level, true)
	res[4] = "Message: " + reform.Inspect(s.Message, true)
	return strings.Join(res, ", ")
}

// Values returns a slice of struct or record field values.
// Returned interface{} values are never untyped nils.
func (s *Log) Values() []interface{} {
	return []interface{}{
		s.ID,
		s.MeterID,
		s.Time,
		s.Level,
		s.Message,
	}
}

// Pointers returns a slice of pointers to struct or record fields.
// Returned interface{} values are never untyped nils.
func (s *Log) Pointers() []interface{} {
	return []interface{}{
		&s.ID,
		&s.MeterID,
		&s.Time,
		&s.Level,
		&s.Message,
	}
}

// View returns View object for that struct.
func (s *Log) View() reform.View {
	return LogTable
}

// Table returns Table object for that record.
func (s *Log) Table() reform.Table {
	return LogTable
}

// PKValue returns a value of primary key for that record.
// Returned interface{} value is never untyped nil.
func (s *Log) PKValue() interface{} {
	return s.ID
}

// PKPointer returns a pointer to primary key field for that record.
// Returned interface{} value is never untyped nil.
func (s *Log) PKPointer() interface{} {
	return &s.ID
}

// HasPK returns true if record has non-zero primary key set, false otherwise.
func (s *Log) HasPK() bool {
	return s.ID != LogTable.z[LogTable.s.PKFieldIndex]
}

// SetPK sets record primary key, if possible.
//
// Deprecated: prefer direct field assignment where possible: s.ID = pk.
func (s *Log) SetPK(pk interface{}) {
	reform.SetPK(s, pk)
}

// check interfaces
var (
	_ reform.View   = LogTable
	_ reform.Struct = (*Log)(nil)
	_ reform.Table  = LogTable
	_ reform.Record = (*Log)(nil)
	_ fmt.Stringer  = (*Log)(nil)
)

func init() {
	parse.AssertUpToDate(&MeterTable.s, new(Meter))
	parse.AssertUpToDate(&LogTable.s, new(Log))
}
