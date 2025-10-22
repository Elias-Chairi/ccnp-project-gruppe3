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

	TemperatureChangePerFanStep    float32 = 0.50 // 50% change towards outdoor temperature per simulation step per fan when MaxWindowOpenness is reached
	TemperatureChangePerHeaterStep float32 = 0.03 // 3% change towards heater temperature per simulation step
	TemperatureChangePerWindowStep float32 = 0.05 // 5% change towards outdoor temperature per simulation step per fully open window

	MaxWindowOpenness float32 = 6.0 // Maximum cumulative openness from all windows. Fan effects scale from 0 to TemperatureChangePerFanStep based on openness
)

type OutdoorConditions struct {
	Temperature float32
	Humidity    float32
	LightLevel  float32
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
			// g.updateHumidity(s)
		case "LIGHT":
			// g.updateLightLevel(s)
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
