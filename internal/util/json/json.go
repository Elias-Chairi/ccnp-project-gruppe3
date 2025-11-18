package json

import (
	"encoding/json"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/cmd/node/greenhouse"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/entity"
)

type RootDTO struct {
	Node              NodeDTO              `json:"node"`
	OutdoorConditions OutdoorConditionsDTO `json:"outdoorConditions"`
}

func ParseGreenhouseJSON(input string) (entity.Node, greenhouse.OutdoorConditions, error) {
	var root RootDTO

	if err := json.Unmarshal([]byte(input), &root); err != nil {
		return entity.Node{}, greenhouse.OutdoorConditions{}, err
	}

	// Convert Node
	node := entity.Node{
		ID: root.Node.ID,
	}

	for _, a := range root.Node.Actuators {
		node.Actuators = append(node.Actuators, ToActuator(a))
	}

	for _, s := range root.Node.Sensors {
		node.Sensors = append(node.Sensors, ToSensor(s))
	}

	// Convert outdoor conditions
	outdoor := ToOutdoorConditions(root.OutdoorConditions)

	return node, outdoor, nil
}
