package mappers

import (
	"pulsar_alice/internal/models/alice"
	"pulsar_alice/internal/models/common"
)

func MeterState(meter *common.Meter) alice.PayloadStateDevice {
	property := alice.PayloadStateDeviceProperties{
		Type: alice.PropertyTypeFloat,
		State: alice.PayloadStateDevicePropertiesState{
			Instance: alice.PropertyParameterInstanceWaterMeter,
			Value:    meter.Value,
		},
	}

	return alice.PayloadStateDevice{
		ID: meter.ID,
		Properties: []alice.PayloadStateDeviceProperties{
			property,
		},
	}
}

func MetersState(meters []*common.Meter) []alice.PayloadStateDevice {
	if len(meters) == 0 {
		return nil
	}

	res := make([]alice.PayloadStateDevice, 0, len(meters))
	for _, m := range meters {
		res = append(res, MeterState(m))
	}

	return res
}

func MeterDevice(device *common.Meter) alice.Device {
	state := alice.PayloadStateDevicePropertiesState{
		Instance: alice.PropertyParameterInstanceWaterMeter,
		Value:    device.Value,
	}

	tp := alice.DeviceTypeMeterHot
	if device.Cold {
		tp = alice.DeviceTypeMeterCold
	}
	return alice.Device{
		ID:   device.ID,
		Name: device.Name,
		Type: tp,
		DeviceInfo: &alice.DeviceInfo{
			Model:        device.Model,
			SWVersion:    device.SWVersion,
			Manufacturer: device.Manufacturer,
		},
		Properties: []alice.Property{
			{
				Type:        alice.PropertyTypeFloat,
				Retrievable: true,
				Reportable:  true,
				Parameters: alice.PropertiesFloatParameters{
					Instance: alice.PropertyParameterInstanceWaterMeter,
					Unit:     alice.UnitCubicMeter,
				},
				State:          state,
				LastUpdated:    device.Updated,
				StateChangedAt: device.Changed,
			},
		},
	}
}
