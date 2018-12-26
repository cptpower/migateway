package migateway

const (
	MODEL_SENSORHT = "sensor_ht"

	FIELD_SENSORHT_TEMPERATURE = "temperature"
	FIELD_SENSORHT_HUMIDITY    = "humidity"
)

type SensorHT struct {
	*Device
	State SensorHTState
}

type SensorHTState struct {
	Temperature float64
	Humidity    float64
	Battery     float32
}

type SensorHTStateChange struct {
	ID   string
	From SensorHTState
	To   SensorHTState
}

func (s SensorHTStateChange) IsChanged() bool {
	return s.From.Temperature != s.To.Temperature || s.From.Humidity != s.To.Humidity || s.From.Battery != s.To.Battery
}

func NewSensorHt(dev *Device) *SensorHT {
	return &SensorHT{
		Device: dev,
		State: SensorHTState{
			Temperature: dev.GetDataAsFloat64(FIELD_SENSORHT_TEMPERATURE) / 100,
			Humidity:    dev.GetDataAsFloat64(FIELD_SENSORHT_HUMIDITY) / 100,
			Battery:     dev.GetBatteryLevel(0),
		},
	}
}

func (s *SensorHT) Set(dev *Device) {
	dev.SetLastUpdate()

	change := &SensorHTStateChange{ID: s.Sid, From: s.State, To: s.State}
	if dev.hasField(FIELD_SENSORHT_TEMPERATURE) {
		s.State.Temperature = dev.GetDataAsFloat64(FIELD_SENSORHT_TEMPERATURE) / 100
	}
	if dev.hasField(FIELD_SENSORHT_HUMIDITY) {
		s.State.Humidity = dev.GetDataAsFloat64(FIELD_SENSORHT_HUMIDITY) / 100
	}

	s.State.Battery = dev.GetBatteryLevel(s.State.Battery)
	change.To = s.State
	if change.IsChanged() || s.shouldPushUpdates() {
		s.Aqara.StateMessages <- change
	}
	if dev.Token != "" {
		s.Token = dev.Token
	}
}
