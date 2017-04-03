package main

import (
  "encoding/json"
  "net/http"
  "fmt"
  "os"
)

type openWeatherMap struct{}

func (w openWeatherMap) temperature(city string) (float64, error) {
  resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + os.Getenv("OWD_APPID") + "&q=" + city)
  if err != nil {
    return 0, err
  }

  defer resp.Body.Close()

  var temperature openWeatherData

  if err := json.NewDecoder(resp.Body).Decode(&temperature); err != nil {
    return 0, err
  }

  fmt.Printf("openWeatherMap: %s: %.2f", city, temperature.Main.Kelvin)
  return temperature.Main.Kelvin, nil
}

type openWeatherData struct{
  Main struct {
    Kelvin float64 `json:"temp"`
  } `json:"main"`
}
