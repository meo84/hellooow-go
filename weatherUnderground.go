package main

import (
  "encoding/json"
  "net/http"
  "fmt"
)

type weatherUnderground struct{
  apiKey string
}

func (w weatherUnderground) temperature(city string) (float64, error) {
  resp, err := http.Get("http://api.wunderground.com/api/" + w.apiKey + "/conditions/q/" + city + ".json")
  if err != nil {
    return 0, err
  }

  defer resp.Body.Close()

  var temp weatherUndergroundData

  if err := json.NewDecoder(resp.Body).Decode(&temp); err != nil {
    return 0, err
  }

  kelvin := temp.Observation.Celsius + 273.15
  fmt.Printf("weatherUndergroundMap: %s: %.2f", city, kelvin)
  return kelvin, nil
}

type weatherUndergroundData struct{
  Observation struct {
    Celsius float64 `json:"temp_c"`
  } `json:"current_observation"`
}
