package dht

import "time"

type managedDevice struct {
	t           device
	lastUpdate  time.Time
	policy      UpdatePolicy
	initialized bool
}

func (m *managedDevice) Measurements() (temperature int16, humidity uint16, err error) {
	err = m.checkForUpdateOnDataRequest()
	if err != nil {
		return 0, 0, err
	}
	return m.t.Measurements()
}

func (m *managedDevice) Temperature() (temp int16, err error) {
	err = m.checkForUpdateOnDataRequest()
	if err != nil {
		return 0, err
	}
	temp, err = m.t.Temperature()
	return
}

func (m *managedDevice) checkForUpdateOnDataRequest() (err error) {
	// update if necessary
	if m.policy.UpdateAutomatically {
		err = m.ReadMeasurements()
	}
	// ignore error if the data was updated recently
	// interface comparison does not work in tinygo. Therefore need to cast to explicit type
	if code, ok := err.(ErrorCode); ok && code == UpdateError {
		err = nil
	}
	// add error if the data is not initialized
	if !m.initialized {
		err = UninitializedDataError
	}
	return err
}

func (m *managedDevice) TemperatureFloat(scale TemperatureScale) (float32, error) {
	err := m.checkForUpdateOnDataRequest()
	if err != nil {
		return 0, err
	}
	return m.t.TemperatureFloat(scale)
}

func (m *managedDevice) Humidity() (hum uint16, err error) {
	err = m.checkForUpdateOnDataRequest()
	if err != nil {
		return 0, err
	}
	return m.t.Humidity()
}

func (m *managedDevice) HumidityFloat() (float32, error) {
	err := m.checkForUpdateOnDataRequest()
	if err != nil {
		return 0, err
	}
	return m.t.HumidityFloat()
}

func (m *managedDevice) ReadMeasurements() (err error) {
	timestamp := time.Now()
	if !m.initialized || timestamp.Sub(m.lastUpdate) > m.policy.UpdateTime {
		err = m.t.ReadMeasurements()
	} else {
		err = UpdateError
	}
	if err == nil {
		m.initialized = true
		m.lastUpdate = timestamp
	}
	return
}
