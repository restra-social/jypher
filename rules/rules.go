package rules

import "github.com/restra-social/jypher/models"

func FHIRRules() map[string]models.Rules {

	return map[string]models.Rules{

		"Restaurant": models.Rules{
			SkipField: []string{"time", "picture", "social"},
		},
		"Patient": models.Rules{
			Rename: map[string]interface{}{},
		},
		"Encounter": models.Rules{
			Rename: map[string]interface{}{
				"subject":         "patient",
				"serviceProvider": "organization",
			},
		},
		"Condition": models.Rules{
			Rename: map[string]interface{}{
				"subject": "patient",
				"context": "encounter",
			},
		},
		"Observation": models.Rules{
			Rename: map[string]interface{}{
				"subject": "patient",
				"context": "encounter",
			},
		},
		"DiagnosticReport": models.Rules{
			Rename: map[string]interface{}{
				"subject": "patient",
				"context": "encounter",
				"result":  "observation",
			},
		},
		"CarePlan": models.Rules{
			Rename: map[string]interface{}{
				"subject": "patient",
				"context": "encounter",
			},
		},
		"Goal": models.Rules{
			Rename: map[string]interface{}{
				"addresses": "condition",
			},
		},
		"MedicationRequest": models.Rules{
			Rename: map[string]interface{}{
				"subject":             "patient",
				"context":             "encounter",
				"medicationReference": "Medication",
			},
		},
		"Bundle": models.Rules{
			Rename: map[string]interface{}{
				"subject":         "patient",
				"serviceProvider": "organization",
			},
		},
	}
}
