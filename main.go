package main

import (
  "encoding/json"
  "net/http"
  "strings"
)

func main() {
  mw := multiWeatherProvider{
    openWeatherMap{},
    weatherUnderground{apiKey: "need-key"},
  }

  http.HandleFunc("/weather/", weather(mw))
  http.ListenAndServe(":8080", nil)
}

func weather(mw multiWeatherProvider) func(w http.ResponseWriter, r *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    city := strings.SplitN(r.URL.Path, "/", 3)[2]

    temp, err := mw.temperature(city)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    json.NewEncoder(w).Encode(map[string]interface{}{
    "city": city,
    "temp": temp,
  })
  }
}

type weatherProvider interface {
  temperature(city string) (float64, error) //in Kelvin
}

type multiWeatherProvider []weatherProvider

func (w multiWeatherProvider) temperature(city string) (float64, error) {
  sum := 0.0

  for _, provider := range w {
    k, err := provider.temperature(city)
    if err != nil {
      return 0, err
    }

    sum += k
  }

  return sum / float64(len(w)), nil
}
