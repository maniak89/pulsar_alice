package alice

import (
	"time"
)

type DeviceType string

const (
	DeviceTypeMeterCold DeviceType = "devices.types.smart_meter.cold_water"
	DeviceTypeMeterHot  DeviceType = "devices.types.smart_meter.hot_water"
)

type Device struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Room        string      `json:"room,omitempty"`
	Type        DeviceType  `json:"type"`
	Properties  []Property  `json:"properties,omitempty"`
	DeviceInfo  *DeviceInfo `json:"device_info,omitempty"`
}

type PropertyType string

const (
	PropertyTypeEvent PropertyType = "devices.properties.event"
	PropertyTypeFloat PropertyType = "devices.properties.float"
)

type PropertyParameterValue struct {
	Value float32 `json:"value"`
	Name  string  `json:"name"`
}

type PropertyParameterInstance string

const (
	PropertyParameterInstanceWaterMeter = "water_meter"
)

type Property struct {
	Type           PropertyType                      `json:"type"`
	Retrievable    bool                              `json:"retrievable"`
	Reportable     bool                              `json:"reportable"`
	Parameters     PropertiesFloatParameters         `json:"parameters"`
	State          PayloadStateDevicePropertiesState `json:"state"`
	StateChangedAt time.Time                         `json:"state_changed_at"`
	LastUpdated    time.Time                         `json:"last_updated"`
}

type DeviceInfo struct {
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	SWVersion    string `json:"sw_version,omitempty"`
}

type ActionResultStatus string

const (
	ActionResultStatusDone  ActionResultStatus = "DONE"
	ActionResultStatusError ActionResultStatus = "ERROR"
)

type ErrorCode string

const (
	ErrorCodeDeviceUnreachable ErrorCode = "DEVICE_UNREACHABLE"
	ErrorCodeInvalidAction     ErrorCode = "INVALID_ACTION"
	ErrorCodeInvalidValue      ErrorCode = "INVALID_VALUE"
)

type ActionResult struct {
	Status           ActionResultStatus `json:"status"`
	ErrorCode        ErrorCode          `json:"error_code,omitempty"`
	ErrorDescription string             `json:"error_description,omitempty"`
}

type Unit string

const (
	UnitCubicMeter Unit = "unit.cubic_meter"
)

type PropertiesType string

const (
	PropertiesTypeFloat PropertiesType = "devices.properties.float"
)

type PropertiesFloat struct {
	Type        PropertiesType       `json:"type"`
	Retrievable bool                 `json:"retrievable,omitempty"`
	Reportable  bool                 `json:"reportable,omitempty"`
	Parameters  interface{}          `json:"parameters"`
	State       PropertiesFloatState `json:"state"`
}

type PropertiesFloatParametersInstance string

const (
	PropertiesFloatParametersInstanceTemperature PropertiesFloatParametersInstance = "water_meter"
)

type PropertiesFloatParameters struct {
	Instance PropertiesFloatParametersInstance `json:"instance"`
	Unit     Unit                              `json:"unit"`
}

type PropertiesFloatState struct {
	Instance PropertiesFloatParametersInstance `json:"instance"`
	Value    float32                           `json:"value"`
}

type DeviceRequest struct {
	ID           string            `json:"id"`
	CustomData   map[string]string `json:"custom_data"`
	Capabilities []struct {
		State struct {
			Instance string      `json:"instance"`
			Value    interface{} `json:"value"`
			Relative bool        `json:"relative"`
		}
	} `json:"capabilities"`
}
