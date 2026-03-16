package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	weatherCacheKey = "weather:current"
	weatherCacheTTL = 30 * time.Minute
)

type WeatherData struct {
	Temp      string `json:"temp"`
	FeelsLike string `json:"feels_like"`
	Humidity  string `json:"humidity"`
	Desc      string `json:"desc"`
	Icon      string `json:"icon"`
	WindSpeed string `json:"wind_speed"`
	Location  string `json:"location"`
}

type WeatherService struct {
	rdb      *redis.Client
	location string
}

func NewWeatherService(rdb *redis.Client, location string) *WeatherService {
	return &WeatherService{rdb: rdb, location: location}
}

func (s *WeatherService) GetCurrent(ctx context.Context) (*WeatherData, error) {
	// Try cache first
	cached, err := s.rdb.Get(ctx, weatherCacheKey).Result()
	if err == nil {
		var data WeatherData
		if json.Unmarshal([]byte(cached), &data) == nil {
			return &data, nil
		}
	}

	// Fetch from wttr.in
	data, err := s.fetchWeather(ctx)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if b, err := json.Marshal(data); err == nil {
		if cacheErr := s.rdb.Set(ctx, weatherCacheKey, string(b), weatherCacheTTL).Err(); cacheErr != nil {
			slog.Warn("failed to cache weather data", "error", cacheErr)
		}
	}

	return data, nil
}

// wttr.in JSON response structures (only fields we need)
type wttrResponse struct {
	CurrentCondition []wttrCondition `json:"current_condition"`
	NearestArea      []wttrArea      `json:"nearest_area"`
}

type wttrCondition struct {
	TempC         string           `json:"temp_C"`
	FeelsLikeC    string           `json:"FeelsLikeC"`
	Humidity      string           `json:"humidity"`
	WeatherCode   string           `json:"weatherCode"`
	WindspeedKmph string           `json:"windspeedKmph"`
	WeatherDesc   []wttrLangValue  `json:"lang_zh"`
	WeatherDescEN []wttrLangValue  `json:"weatherDesc"`
}

type wttrLangValue struct {
	Value string `json:"value"`
}

type wttrArea struct {
	AreaName []wttrLangValue `json:"areaName"`
}

func (s *WeatherService) fetchWeather(ctx context.Context) (*WeatherData, error) {
	url := fmt.Sprintf("https://wttr.in/%s?format=j1&lang=zh", s.location)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "nanamiku-blog/1.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wttr.in returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var wttr wttrResponse
	if err := json.Unmarshal(body, &wttr); err != nil {
		return nil, fmt.Errorf("decode weather json: %w", err)
	}

	if len(wttr.CurrentCondition) == 0 {
		return nil, fmt.Errorf("no current condition in weather response")
	}

	cc := wttr.CurrentCondition[0]

	desc := ""
	if len(cc.WeatherDesc) > 0 {
		desc = cc.WeatherDesc[0].Value
	} else if len(cc.WeatherDescEN) > 0 {
		desc = cc.WeatherDescEN[0].Value
	}

	location := s.location
	if len(wttr.NearestArea) > 0 && len(wttr.NearestArea[0].AreaName) > 0 {
		location = wttr.NearestArea[0].AreaName[0].Value
	}

	icon := weatherCodeToIcon(cc.WeatherCode)

	return &WeatherData{
		Temp:      cc.TempC,
		FeelsLike: cc.FeelsLikeC,
		Humidity:  cc.Humidity,
		Desc:      desc,
		Icon:      icon,
		WindSpeed: cc.WindspeedKmph,
		Location:  location,
	}, nil
}

func weatherCodeToIcon(code string) string {
	c, _ := strconv.Atoi(code)
	switch {
	case c == 113:
		return "sunny"
	case c == 116:
		return "partly_cloudy"
	case c == 119, c == 122:
		return "cloudy"
	case c >= 176 && c <= 263:
		return "light_rain"
	case c >= 266 && c <= 314:
		return "rain"
	case c >= 317 && c <= 395:
		return "snow"
	case c >= 200 && c <= 202:
		return "thunderstorm"
	default:
		return "cloudy"
	}
}
