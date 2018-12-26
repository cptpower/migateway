package migateway

const (
	MODEL_WEATHERV1 = "weather.v1"

	FIELD_WEATHERV1_TEMPERATURE = "temperature"
	FIELD_WEATHERV1_HUMIDITY    = "humidity"
	FIELD_WEATHERV1_PRESSURE    = "pressure"
)

type WeatherV1 struct {
	*Device
	State WeatherV1State
}

type WeatherV1State struct {
	Temperature float64
	Humidity    float64
	Pressure    float64
	Battery     float32
}

type WeatherV1StateChange struct {
	ID   string
	From WeatherV1State
	To   WeatherV1State
}

func (s WeatherV1StateChange) IsChanged() bool {
	return s.From.Temperature != s.To.Temperature || s.From.Humidity != s.To.Humidity || s.From.Battery != s.To.Battery || s.From.Pressure != s.To.Pressure
}

func NewWeatherV1(dev *Device) *WeatherV1 {
	return &WeatherV1{
		Device: dev,
		State: WeatherV1State{
			Temperature: dev.GetDataAsFloat64(FIELD_WEATHERV1_TEMPERATURE) / 100,
			Humidity:    dev.GetDataAsFloat64(FIELD_WEATHERV1_HUMIDITY) / 100,
			Pressure:    dev.GetDataAsFloat64(FIELD_WEATHERV1_PRESSURE) / 100,
			Battery:     dev.GetBatteryLevel(0),
		},
	}
}

func (s *WeatherV1) Set(dev *Device) {
	dev.SetLastUpdate()

	change := &WeatherV1StateChange{ID: s.Sid, From: s.State, To: s.State}
	if dev.hasField(FIELD_WEATHERV1_TEMPERATURE) {
		s.State.Temperature = dev.GetDataAsFloat64(FIELD_WEATHERV1_TEMPERATURE) / 100
	}
	if dev.hasField(FIELD_WEATHERV1_HUMIDITY) {
		s.State.Humidity = dev.GetDataAsFloat64(FIELD_WEATHERV1_HUMIDITY) / 100
	}
	if dev.hasField(FIELD_WEATHERV1_PRESSURE) {
		s.State.Humidity = dev.GetDataAsFloat64(FIELD_WEATHERV1_PRESSURE) / 100
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
