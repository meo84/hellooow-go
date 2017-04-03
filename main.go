package main

import (
  "encoding/json"
  "net/http"
  "strings"
  "os"
)

func main() {
  mw := multiWeatherProvider{
    openWeatherMap{},
    weatherUnderground{apiKey: os.Getenv("WG_API_KEY")},
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
  temps := make(chan float64, len(w))
  errs := make(chan error, len(w))

  for _, provider := range w {
    go func(p weatherProvider) {
      k, err := p.temperature(city)
      if err != nil {
        errs <- err
        return
      }
      temps <- k
    }(provider)
  }

  sum := 0.0

  for i := 0; i < len(w); i++ {
    select {
    case temp := <-temps:
      sum += temp
    case err := <-errs:
      return 0, err
    }
  }
  return sum / float64(len(w)), nil
}
