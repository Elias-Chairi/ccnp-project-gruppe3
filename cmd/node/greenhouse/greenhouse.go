package greenhouse

import (
	"log"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
)

const (
	DefaultFanRPM     int32   = 1000  // 1000 RPM when fan state is true
	DefaultLightLUX   int32   = 10000 // 10 000 lux when light state is true
	DefaultWindowOpen float32 = 1.0   // 100% open when window state is true
	DefaultHeaterTemp float32 = 30.0  // 30 °C when heater state is true

	TemperatureChangePerFanStep    float32 = 0.50 // 50% convergence towards outdoor temperature per simulation step per fan when MaxWindowOpenness is reached
	TemperatureChangePerHeaterStep float32 = 0.03 // 3% convergence rate towards heater temperature per simulation step
	TemperatureChangePerWindowStep float32 = 0.05 // 5% convergence rate towards outdoor temperature per simulation step per fully open window (i.e., moves 5% closer each step)

	MaxWindowOpenness float32 = 6.0 // Maximum cumulative openness from all windows. Fan effects scale from 0 to TemperatureChangePerFanStep based on openness

	HumidityChangePerWindowStep float32 = 0.05  // 5% convergence toward outdoor humidity per fully open window
	HumidityChangePerFanStep    float32 = 0.02  // 2% additional convergence per fan
	HumidityChangePerHeaterStep float32 = -0.03 // heater dries air: -3% relative change per simulation step

	DefaultHumidifierHumidity float32 = 80.0 // default humidifier target humidity (80%)

	LightChangePerWindowStep float32 = 0.05  // 5% of outdoor light passes per fully open window
)

type OutdoorConditions struct {
	Temperature float32
	Humidity    float32
	LightLevel  int32
}

type Greenhouse struct {
	Outdoor        OutdoorConditions
	Node           entity.Node
	onSensorUpdate func(sensor *entity.Sensor[any])
}

func NewGreenhouse(node entity.Node, outdoor OutdoorConditions, onSensorUpdate func(sensor *entity.Sensor[any])) *Greenhouse {
	return &Greenhouse{
		Outdoor:        outdoor,
		Node:           node,
		onSensorUpdate: onSensorUpdate,
	}
}

func (g *Greenhouse) SimulateStep() {
	for _, s := range g.Node.Sensors {
		switch s.Type {
		case "TEMPERATURE":
			g.updateTemperature(s)
		case "HUMIDITY":
			 g.updateHumidity(s)
		case "LIGHT":
			 g.updateLightLevel(s)
		default:
			log.Println("Unknown sensor type:", s.Type)
		}
	}
}

func (g *Greenhouse) updateTemperature(sensor *entity.Sensor[any]) {
	// Get current temperature value of the sensor
	var currentTemp float32
	if sensor.Value != nil {
		var ok bool
		currentTemp, ok = sensor.Value.(float32)
		if !ok {
			log.Printf("Sensor ID %d has invalid temperature value type of %T\n", sensor.ID, sensor.Value)
			return
		}
	}

	actuators := make(map[string][]*entity.Actuator[any])
	for _, a := range g.Node.Actuators {
		actuators[a.Type] = append(actuators[a.Type], a)
	}

	// Heater effects
	for _, a := range actuators["HEATER"] {
		switch a.Unit {
		case "":
			if on, ok := a.State.(bool); ok {
				if on {
					if currentTemp < DefaultHeaterTemp { // heater only heats
						currentTemp += (DefaultHeaterTemp - currentTemp) * TemperatureChangePerHeaterStep
					}
				}
			} else {
				log.Printf("HEATER %d has invalid state type %T for the '' unit\n", a.ID, a.State)
			}
		case "°C":
			if targetTemp, ok := a.State.(float32); ok {
				if currentTemp < targetTemp { // heater only heats
					currentTemp += (targetTemp - currentTemp) * TemperatureChangePerHeaterStep
				}
			} else {
				log.Printf("HEATER %d has invalid state type %T for the '°C' unit\n", a.ID, a.State)
			}
		default:
			log.Printf("Unknown actuator unit for HEATER: %s\n", a.Unit)
		}
	}

	// Window effects
	windowsOpen := float32(0.0)
	for _, a := range actuators["WINDOW"] {
		switch a.Unit {
		case "":
			switch open := a.State.(type) {
			case bool:
				if open {
					windowsOpen += DefaultWindowOpen
					currentTemp += (g.Outdoor.Temperature - currentTemp) * (DefaultWindowOpen * TemperatureChangePerWindowStep)
				}
			case float32:
				windowsOpen += open
				currentTemp += (g.Outdoor.Temperature - currentTemp) * (open * TemperatureChangePerWindowStep)
			default:
				log.Printf("WINDOW %d has invalid state type %T for the [no unit] unit\n", a.ID, a.State)
			}
		default:
			log.Printf("Unknown actuator unit for WINDOW: %s\n", a.Unit)
		}
	}

	// Fan effects
	if windowsOpen > 0 { // fans only have effect if windows are open
		windowsOpenEffect := min(1.0, windowsOpen/MaxWindowOpenness)

		for _, a := range actuators["FAN"] {
			switch a.Unit {
			case "":
				if on, ok := a.State.(bool); ok {
					if on { // assumes running at DefaultFanRPM
						currentTemp += (g.Outdoor.Temperature - currentTemp) * TemperatureChangePerFanStep * windowsOpenEffect
					}
				} else {
					log.Printf("FAN %d has invalid state type %T for the [no unit] unit\n", a.ID, a.State)
				}
			case "RPM":
				if rpm, ok := a.State.(int32); ok {
					if rpm > 0 {
						fanEffect := (float32(rpm) / float32(DefaultFanRPM)) * TemperatureChangePerFanStep * windowsOpenEffect
						currentTemp += (g.Outdoor.Temperature - currentTemp) * fanEffect
					}
				} else {
					log.Printf("FAN %d has invalid state type %T for the 'RPM' unit\n", a.ID, a.State)
				}
			default:
				log.Printf("Unknown actuator unit for FAN: %s\n", a.Unit)
			}
		}
	}

	// Update sensor value
	sensor.Value = currentTemp
	if g.onSensorUpdate != nil {
		g.onSensorUpdate(sensor)
	}
}



// updateHumidity simulates how the greenhouse's humidity changes during one simulation step.
// It is affected by outdoor humidity, open windows, running fans, and active heaters.
// Windows and fans move the humidity toward the outdoor value, while heaters dry the air slightly.
func (g *Greenhouse) updateHumidity(sensor *entity.Sensor[any]) {
	
	// Get current humidity value from the sensor
	var currentHum float32
	if sensor.Value != nil {
		val, ok := sensor.Value.(float32)
		if !ok {
			log.Printf("Sensor ID %d has invalid humidity value type %T\n", sensor.ID, sensor.Value)
			return
		}
		currentHum = val
	}

	// Group actuators by type for easy access (same structure as updateTemperature)
	actuators := make(map[string][]*entity.Actuator[any])
	for _, a := range g.Node.Actuators {
		actuators[a.Type] = append(actuators[a.Type], a)
	}

	// Heater effect
	// Heaters dry the air a bit every step since warm air reduces relative humidity.
	for _, a := range actuators["HEATER"] {
		if on, ok := a.State.(bool); ok && on {
			currentHum -= currentHum * 0.02 // lowers humidity by ~2%
		}
	}

	// Window effect
	// Open windows let humidity move toward outdoor humidity.
	windowsOpen := float32(0.0)
	for _, a := range actuators["WINDOW"] {
		switch v := a.State.(type) {
		case bool:
			if v {
				windowsOpen += DefaultWindowOpen
				currentHum += (g.Outdoor.Humidity - currentHum) * HumidityChangePerWindowStep
			}
		case float32:
			windowsOpen += v
			currentHum += (g.Outdoor.Humidity - currentHum) * (v * HumidityChangePerWindowStep)
		}
	}

	// Fan effect
	// Fans accelerate humidity equalization, but only if windows are open.
	if windowsOpen > 0 {
		openEffect := min(1.0, windowsOpen/MaxWindowOpenness)
		for _, a := range actuators["FAN"] {
			switch v := a.State.(type) {
			case bool:
				if v {
					currentHum += (g.Outdoor.Humidity - currentHum) * 0.02 * openEffect
				}
			case int32:
				if v > 0 {
					effect := (float32(v) / float32(DefaultFanRPM)) * 0.02 * openEffect
					currentHum += (g.Outdoor.Humidity - currentHum) * effect
				}
			}
		}
	}

	// humidity between 0% and 100%
	if currentHum < 0 {
		currentHum = 0
	} else if currentHum > 100 {
		currentHum = 100
	}

	// Update the sensor value and trigger the callback
	sensor.Value = currentHum
	if g.onSensorUpdate != nil {
		g.onSensorUpdate(sensor)
	}
}



// updateLightLevel updates the greenhouse's light sensor reading (in lux)
// based on outdoor light and artificial lighting.
// The brightest source (natural or artificial) determines the indoor light level.
func (g *Greenhouse) updateLightLevel(sensor *entity.Sensor[any]) {
	if sensor == nil {
		return
	}

	// Get outdoor light intensity
	outdoorLux := g.Outdoor.LightLevel

	// Calculates total artificial light intensity
	var artificialLux int32
	for _, a := range g.Node.Actuators {
		if a.Type != "LIGHT" {
			continue
		}

		switch v := a.State.(type) {
		case bool:
			if v {
				artificialLux += DefaultLightLUX
			}
		case int32:
			if v > 0 {
				artificialLux += v
			}
		default:
			// ignore unsupported actuator types
		}
	}

	// Chooses the dominant light source
	currentLux := outdoorLux
	if artificialLux > outdoorLux {
		currentLux = artificialLux
	}

	// Updates sensor instantly
	sensor.Value = currentLux

	// Notifies listeners if there are any idk
	if g.onSensorUpdate != nil {
		g.onSensorUpdate(sensor)
	}
}




