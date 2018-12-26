package migateway

import (
	"time"
)

const (
	MODEL_MOTIONAQ2 = "sensor_motion.aq2"

	MOTIONAQ2_FIELD_NO_MOTION = "no_motion"
	MOTIONAQ2_FIELD_LUX       = "lux"
)

type MotionAq2 struct {
	*Device
	State MotionStateAq2
}

type MotionStateAq2 struct {
	Battery    float32
	HasMotion  bool
	LastMotion time.Time
	Lux        int
}

type MotionStateChangeAq2 struct {
	ID   string
	From MotionStateAq2
	To   MotionStateAq2
}

func (m MotionStateChangeAq2) IsChanged() bool {
	return m.From.HasMotion != m.To.HasMotion || m.From.LastMotion != m.To.LastMotion
}

func (m *MotionAq2) GetData() interface{} {
	return m.Data
}

func NewMotionAq2(dev *Device) *MotionAq2 {
	return &MotionAq2{
		Device: dev,
		State: MotionStateAq2{
			Battery:   dev.GetBatteryLevel(0),
			Lux:       dev.GetDataAsInt(MOTIONAQ2_FIELD_LUX),
			HasMotion: dev.GetDataAsBool(FIELD_STATUS),
		},
	}
}

func (m *MotionAq2) Set(dev *Device) {
	dev.SetLastUpdate()
	timestamp := time.Now()
	change := &MotionStateChangeAq2{ID: m.Sid, From: m.State, To: m.State}
	if dev.hasField(FIELD_STATUS) {
		m.State.HasMotion = dev.GetDataAsBool(FIELD_STATUS)
		if m.State.HasMotion {
			m.State.LastMotion = timestamp
		}
	} else if dev.hasField(FIELD_NO_MOTION) {
		m.State.HasMotion = false
		nomotionInSeconds := int64(dev.GetDataAsInt(FIELD_NO_MOTION)) * -1
		timestamp.Add(time.Duration(nomotionInSeconds) * time.Second)
		m.State.LastMotion = timestamp
	} else if dev.hasField(MOTIONAQ2_FIELD_LUX) {
		m.State.Lux = dev.GetDataAsInt(MOTIONAQ2_FIELD_LUX)
	}

	m.State.Battery = dev.GetBatteryLevel(m.State.Battery)
	change.To = m.State
	if change.IsChanged() || m.shouldPushUpdates() {
		m.Aqara.StateMessages <- change
	}
	if dev.Token != "" {
		m.Token = dev.Token
	}
}
