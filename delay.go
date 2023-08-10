package wiremock

import "encoding/json"

type DelayInterface interface {
	ParseDelay() map[string]interface{}
}

type fixedDelay struct {
	milliseconds int64
}

func (d fixedDelay) ParseDelay() map[string]interface{} {
	return map[string]interface{}{
		"type":         "fixed",
		"milliseconds": d.milliseconds,
	}
}

type logNormalRandomDelay struct {
	median int64
	sigma  float64
}

func (d logNormalRandomDelay) ParseDelay() map[string]interface{} {
	return map[string]interface{}{
		"type":   "lognormal",
		"median": d.median,
		"sigma":  d.sigma,
	}
}

type uniformRandomDelay struct {
	lower int64
	upper int64
}

func (d uniformRandomDelay) ParseDelay() map[string]interface{} {
	return map[string]interface{}{
		"type":  "uniform",
		"lower": d.lower,
		"upper": d.upper,
	}
}

type chunkedDribbleDelay struct {
	numberOfChunks int64
	totalDuration  int64
}

func (d chunkedDribbleDelay) MarshalJSON() ([]byte, error) {
	jsonMap := map[string]interface{}{
		"numberOfChunks": d.numberOfChunks,
		"totalDuration":  d.totalDuration,
	}

	return json.Marshal(jsonMap)
}
