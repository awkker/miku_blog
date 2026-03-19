package service

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"sort"
	"strings"
	"time"

	"nanamiku-blog/backend/query"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AnalyticsCollectInput struct {
	VisitorID   uuid.UUID
	SessionKey  string
	Path        string
	Title       string
	Referrer    string
	Timezone    string
	Language    string
	OccurredAt  time.Time
	UserAgent   string
	ClientIP    string
	CountryCode string
	Region      string
	City        string
}

type DashboardMetric struct {
	Value  float64 `json:"value"`
	Change float64 `json:"change"`
}

type DashboardAnalyticsSummary struct {
	Visitors      DashboardMetric `json:"visitors"`
	Visits        DashboardMetric `json:"visits"`
	Views         DashboardMetric `json:"views"`
	BounceRate    DashboardMetric `json:"bounce_rate"`
	VisitDuration DashboardMetric `json:"visit_duration"`
}

type DashboardAnalyticsWindow struct {
	Range        string `json:"range"`
	Label        string `json:"label"`
	Start        string `json:"start"`
	End          string `json:"end"`
	PreviousFrom string `json:"previous_from"`
	PreviousTo   string `json:"previous_to"`
	Granularity  string `json:"granularity"`
}

type DashboardTrendItem struct {
	Bucket   string `json:"bucket"`
	Visitors int64  `json:"visitors"`
	Views    int64  `json:"views"`
}

type DashboardPageItem struct {
	Path     string `json:"path"`
	Visitors int64  `json:"visitors"`
	Views    int64  `json:"views"`
	Entries  int64  `json:"entries"`
	Exits    int64  `json:"exits"`
}

type DashboardNamedCount struct {
	Name     string `json:"name"`
	Visitors int64  `json:"visitors"`
}

type DashboardCountryCount struct {
	Code     string `json:"code"`
	Visitors int64  `json:"visitors"`
}

type DashboardSources struct {
	Referrers []DashboardNamedCount `json:"referrers"`
	Channels  []DashboardNamedCount `json:"channels"`
}

type DashboardEnvironment struct {
	Browsers []DashboardNamedCount `json:"browsers"`
	OS       []DashboardNamedCount `json:"os"`
	Devices  []DashboardNamedCount `json:"devices"`
}

type DashboardLocation struct {
	Countries []DashboardCountryCount `json:"countries"`
	Regions   []DashboardNamedCount   `json:"regions"`
	Cities    []DashboardNamedCount   `json:"cities"`
}

type DashboardTrafficPoint struct {
	Dow   int   `json:"dow"`
	Hour  int   `json:"hour"`
	Value int64 `json:"value"`
}

type DashboardAnalyticsOverview struct {
	Window      DashboardAnalyticsWindow  `json:"window"`
	Summary     DashboardAnalyticsSummary `json:"summary"`
	Trend       []DashboardTrendItem      `json:"trend"`
	Pages       []DashboardPageItem       `json:"pages"`
	Sources     DashboardSources          `json:"sources"`
	Environment DashboardEnvironment      `json:"environment"`
	Location    DashboardLocation         `json:"location"`
	Traffic     []DashboardTrafficPoint   `json:"traffic"`
}

type analyticsWindow struct {
	rangeKey     string
	label        string
	start        time.Time
	end          time.Time
	previousFrom time.Time
	previousTo   time.Time
	granularity  string
}

func (s *DashboardService) CollectPageview(ctx context.Context, in AnalyticsCollectInput) error {
	now := time.Now()
	occurredAt := normalizeOccurredAt(in.OccurredAt, now)
	path := normalizePath(in.Path)
	if path == "" {
		path = "/"
	}

	sessionKey := normalizeSessionKey(in.SessionKey)
	if sessionKey == "" {
		sessionKey = fmt.Sprintf("auto-%d", now.UnixNano())
	}

	referrer := normalizeReferrer(in.Referrer)
	referrerHost := normalizeHost(referrer)
	channel := classifyChannel(referrerHost)
	browser, osName, device := parseUserAgent(in.UserAgent)

	countryCode := normalizeCountryCode(in.CountryCode)
	region := normalizeGeoName(in.Region)
	city := normalizeGeoName(in.City)
	if s.geoResolver != nil && (countryCode == "ZZ" || region == "Unknown" || city == "Unknown") {
		resolvedCountry, resolvedRegion, resolvedCity, ok := s.geoResolver.Lookup(in.ClientIP)
		if ok {
			if countryCode == "ZZ" && resolvedCountry != "" {
				countryCode = normalizeCountryCode(resolvedCountry)
			}
			if region == "Unknown" && strings.TrimSpace(resolvedRegion) != "" {
				region = normalizeGeoName(resolvedRegion)
			}
			if city == "Unknown" && strings.TrimSpace(resolvedCity) != "" {
				city = normalizeGeoName(resolvedCity)
			}
		}
	}
	timezone := trimRunes(in.Timezone, 64)
	language := trimRunes(in.Language, 32)
	title := trimRunes(strings.TrimSpace(in.Title), 180)

	visitorID := toPgUUID(in.VisitorID)

	session, err := s.q.UpsertAnalyticsSession(ctx, query.UpsertAnalyticsSessionParams{
		SessionKey:   sessionKey,
		VisitorID:    visitorID,
		StartedAt:    occurredAt,
		LastSeenAt:   occurredAt,
		EntryPath:    path,
		ExitPath:     path,
		Referrer:     referrer,
		ReferrerHost: referrerHost,
		Channel:      channel,
		Browser:      browser,
		Os:           osName,
		Device:       device,
		CountryCode:  countryCode,
		Region:       region,
		City:         city,
		Timezone:     timezone,
		Language:     language,
	})
	if err != nil {
		return err
	}

	return s.q.CreateAnalyticsPageview(ctx, query.CreateAnalyticsPageviewParams{
		SessionID:    session.ID,
		VisitorID:    visitorID,
		Path:         path,
		Title:        title,
		Referrer:     referrer,
		ReferrerHost: referrerHost,
		Channel:      channel,
		Browser:      browser,
		Os:           osName,
		Device:       device,
		CountryCode:  countryCode,
		Region:       region,
		City:         city,
		Timezone:     timezone,
		Language:     language,
		OccurredAt:   occurredAt,
	})
}

func (s *DashboardService) GetAnalyticsOverview(ctx context.Context, rangeKey string, offset int) (*DashboardAnalyticsOverview, error) {
	window := resolveAnalyticsWindow(rangeKey, time.Now(), offset)

	visitors, err := s.q.CountAnalyticsVisitors(ctx, query.CountAnalyticsVisitorsParams{StartedAt: window.start, StartedAt_2: window.end})
	if err != nil {
		return nil, err
	}
	prevVisitors, err := s.q.CountAnalyticsVisitors(ctx, query.CountAnalyticsVisitorsParams{StartedAt: window.previousFrom, StartedAt_2: window.previousTo})
	if err != nil {
		return nil, err
	}

	visits, err := s.q.CountAnalyticsVisits(ctx, query.CountAnalyticsVisitsParams{StartedAt: window.start, StartedAt_2: window.end})
	if err != nil {
		return nil, err
	}
	prevVisits, err := s.q.CountAnalyticsVisits(ctx, query.CountAnalyticsVisitsParams{StartedAt: window.previousFrom, StartedAt_2: window.previousTo})
	if err != nil {
		return nil, err
	}

	views, err := s.q.CountAnalyticsViews(ctx, query.CountAnalyticsViewsParams{OccurredAt: window.start, OccurredAt_2: window.end})
	if err != nil {
		return nil, err
	}
	prevViews, err := s.q.CountAnalyticsViews(ctx, query.CountAnalyticsViewsParams{OccurredAt: window.previousFrom, OccurredAt_2: window.previousTo})
	if err != nil {
		return nil, err
	}

	bounceVisits, err := s.q.CountAnalyticsBouncedVisits(ctx, query.CountAnalyticsBouncedVisitsParams{StartedAt: window.start, StartedAt_2: window.end})
	if err != nil {
		return nil, err
	}
	prevBounceVisits, err := s.q.CountAnalyticsBouncedVisits(ctx, query.CountAnalyticsBouncedVisitsParams{StartedAt: window.previousFrom, StartedAt_2: window.previousTo})
	if err != nil {
		return nil, err
	}

	avgDuration, err := s.q.GetAnalyticsAverageVisitDurationSeconds(ctx, query.GetAnalyticsAverageVisitDurationSecondsParams{StartedAt: window.start, StartedAt_2: window.end})
	if err != nil {
		return nil, err
	}
	prevAvgDuration, err := s.q.GetAnalyticsAverageVisitDurationSeconds(ctx, query.GetAnalyticsAverageVisitDurationSecondsParams{StartedAt: window.previousFrom, StartedAt_2: window.previousTo})
	if err != nil {
		return nil, err
	}

	trendRows, err := s.q.GetAnalyticsTrend(ctx, query.GetAnalyticsTrendParams{OccurredAt: window.start, OccurredAt_2: window.end, Column3: window.granularity})
	if err != nil {
		return nil, err
	}

	pageRows, err := s.q.ListAnalyticsPageStats(ctx, query.ListAnalyticsPageStatsParams{OccurredAt: window.start, OccurredAt_2: window.end, Limit: 12})
	if err != nil {
		return nil, err
	}
	entryRows, err := s.q.ListAnalyticsEntryPaths(ctx, query.ListAnalyticsEntryPathsParams{StartedAt: window.start, StartedAt_2: window.end, Limit: 12})
	if err != nil {
		return nil, err
	}
	exitRows, err := s.q.ListAnalyticsExitPaths(ctx, query.ListAnalyticsExitPathsParams{StartedAt: window.start, StartedAt_2: window.end, Limit: 12})
	if err != nil {
		return nil, err
	}

	referrersRows, err := s.q.ListAnalyticsReferrers(ctx, query.ListAnalyticsReferrersParams{StartedAt: window.start, StartedAt_2: window.end, Limit: 10})
	if err != nil {
		return nil, err
	}
	channelRows, err := s.q.ListAnalyticsChannels(ctx, query.ListAnalyticsChannelsParams{StartedAt: window.start, StartedAt_2: window.end, Limit: 10})
	if err != nil {
		return nil, err
	}

	browserRows, err := s.q.ListAnalyticsBrowsers(ctx, query.ListAnalyticsBrowsersParams{StartedAt: window.start, StartedAt_2: window.end, Limit: 10})
	if err != nil {
		return nil, err
	}
	osRows, err := s.q.ListAnalyticsOperatingSystems(ctx, query.ListAnalyticsOperatingSystemsParams{StartedAt: window.start, StartedAt_2: window.end, Limit: 10})
	if err != nil {
		return nil, err
	}
	deviceRows, err := s.q.ListAnalyticsDevices(ctx, query.ListAnalyticsDevicesParams{StartedAt: window.start, StartedAt_2: window.end, Limit: 10})
	if err != nil {
		return nil, err
	}

	countryRows, err := s.q.ListAnalyticsCountries(ctx, query.ListAnalyticsCountriesParams{StartedAt: window.start, StartedAt_2: window.end, Limit: 200})
	if err != nil {
		return nil, err
	}
	regionRows, err := s.q.ListAnalyticsRegions(ctx, query.ListAnalyticsRegionsParams{StartedAt: window.start, StartedAt_2: window.end, Limit: 10})
	if err != nil {
		return nil, err
	}
	cityRows, err := s.q.ListAnalyticsCities(ctx, query.ListAnalyticsCitiesParams{StartedAt: window.start, StartedAt_2: window.end, Limit: 10})
	if err != nil {
		return nil, err
	}

	heatRows, err := s.q.GetAnalyticsTrafficHeatmap(ctx, query.GetAnalyticsTrafficHeatmapParams{OccurredAt: window.start, OccurredAt_2: window.end})
	if err != nil {
		return nil, err
	}

	bounceRate := percentage(float64(bounceVisits), float64(visits))
	prevBounceRate := percentage(float64(prevBounceVisits), float64(prevVisits))

	summary := DashboardAnalyticsSummary{
		Visitors:      DashboardMetric{Value: float64(visitors), Change: percentChange(float64(visitors), float64(prevVisitors))},
		Visits:        DashboardMetric{Value: float64(visits), Change: percentChange(float64(visits), float64(prevVisits))},
		Views:         DashboardMetric{Value: float64(views), Change: percentChange(float64(views), float64(prevViews))},
		BounceRate:    DashboardMetric{Value: round1(bounceRate), Change: percentChange(bounceRate, prevBounceRate)},
		VisitDuration: DashboardMetric{Value: round1(avgDuration), Change: percentChange(avgDuration, prevAvgDuration)},
	}

	result := &DashboardAnalyticsOverview{
		Window: DashboardAnalyticsWindow{
			Range:        window.rangeKey,
			Label:        window.label,
			Start:        window.start.Format(time.RFC3339),
			End:          window.end.Format(time.RFC3339),
			PreviousFrom: window.previousFrom.Format(time.RFC3339),
			PreviousTo:   window.previousTo.Format(time.RFC3339),
			Granularity:  window.granularity,
		},
		Summary: summary,
		Trend:   fillTrendBuckets(window, trendRows),
		Pages:   mergePageRows(pageRows, entryRows, exitRows),
		Sources: DashboardSources{
			Referrers: mapReferrers(referrersRows),
			Channels:  mapChannelRows(channelRows),
		},
		Environment: DashboardEnvironment{
			Browsers: mapBrowserRows(browserRows),
			OS:       mapOperatingSystemRows(osRows),
			Devices:  mapDeviceRows(deviceRows),
		},
		Location: DashboardLocation{
			Countries: mapCountries(countryRows),
			Regions:   mapRegionRows(regionRows),
			Cities:    mapCityRows(cityRows),
		},
		Traffic: mapTraffic(heatRows),
	}

	return result, nil
}

func resolveAnalyticsWindow(rangeKey string, now time.Time, offset int) analyticsWindow {
	if offset < 0 {
		offset = 0
	}

	r := strings.ToLower(strings.TrimSpace(rangeKey))
	if r == "" {
		r = "24h"
	}

	if r == "7d" {
		dayEnd := startOfDay(now).Add(24 * time.Hour)
		shiftDays := offset * 7
		dayEnd = dayEnd.AddDate(0, 0, -shiftDays)
		start := dayEnd.AddDate(0, 0, -7)
		return analyticsWindow{
			rangeKey:     "7d",
			label:        "Last 7 days",
			start:        start,
			end:          dayEnd,
			previousFrom: start.AddDate(0, 0, -7),
			previousTo:   start,
			granularity:  "day",
		}
	}

	if r == "30d" {
		dayEnd := startOfDay(now).Add(24 * time.Hour)
		shiftDays := offset * 30
		dayEnd = dayEnd.AddDate(0, 0, -shiftDays)
		start := dayEnd.AddDate(0, 0, -30)
		return analyticsWindow{
			rangeKey:     "30d",
			label:        "Last 30 days",
			start:        start,
			end:          dayEnd,
			previousFrom: start.AddDate(0, 0, -30),
			previousTo:   start,
			granularity:  "day",
		}
	}

	hourEnd := now.Truncate(time.Hour).Add(time.Hour)
	if offset > 0 {
		hourEnd = hourEnd.Add(-time.Duration(offset) * 24 * time.Hour)
	}
	start := hourEnd.Add(-24 * time.Hour)
	return analyticsWindow{
		rangeKey:     "24h",
		label:        "Last 24 hours",
		start:        start,
		end:          hourEnd,
		previousFrom: start.Add(-24 * time.Hour),
		previousTo:   start,
		granularity:  "hour",
	}
}

func fillTrendBuckets(window analyticsWindow, rows []query.GetAnalyticsTrendRow) []DashboardTrendItem {
	seriesMap := make(map[string]DashboardTrendItem, len(rows))
	for _, row := range rows {
		bucket := toBucketString(row.Bucket)
		seriesMap[bucket] = DashboardTrendItem{Bucket: bucket, Visitors: row.Visitors, Views: row.Views}
	}

	items := make([]DashboardTrendItem, 0)
	step := 24 * time.Hour
	format := "2006-01-02"
	if window.granularity == "hour" {
		step = time.Hour
		format = "2006-01-02 15:00"
	}

	for cursor := window.start; cursor.Before(window.end); cursor = cursor.Add(step) {
		bucket := cursor.Format(format)
		if got, ok := seriesMap[bucket]; ok {
			items = append(items, got)
			continue
		}
		items = append(items, DashboardTrendItem{Bucket: bucket, Visitors: 0, Views: 0})
	}

	return items
}

func mergePageRows(
	pageRows []query.ListAnalyticsPageStatsRow,
	entryRows []query.ListAnalyticsEntryPathsRow,
	exitRows []query.ListAnalyticsExitPathsRow,
) []DashboardPageItem {
	type pageAgg struct {
		path     string
		visitors int64
		views    int64
		entries  int64
		exits    int64
	}

	m := map[string]*pageAgg{}
	ensure := func(path string) *pageAgg {
		p := normalizePath(path)
		if p == "" {
			p = "/"
		}
		item, ok := m[p]
		if ok {
			return item
		}
		item = &pageAgg{path: p}
		m[p] = item
		return item
	}

	for _, row := range pageRows {
		it := ensure(row.Path)
		it.visitors = row.Visitors
		it.views = row.Views
	}
	for _, row := range entryRows {
		it := ensure(row.Path)
		it.entries = row.Visits
	}
	for _, row := range exitRows {
		it := ensure(row.Path)
		it.exits = row.Visits
	}

	items := make([]DashboardPageItem, 0, len(m))
	for _, row := range m {
		items = append(items, DashboardPageItem{
			Path:     row.path,
			Visitors: row.visitors,
			Views:    row.views,
			Entries:  row.entries,
			Exits:    row.exits,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Visitors != items[j].Visitors {
			return items[i].Visitors > items[j].Visitors
		}
		if items[i].Entries != items[j].Entries {
			return items[i].Entries > items[j].Entries
		}
		if items[i].Exits != items[j].Exits {
			return items[i].Exits > items[j].Exits
		}
		return items[i].Path < items[j].Path
	})

	if len(items) > 12 {
		return items[:12]
	}
	return items
}

func mapReferrers(rows []query.ListAnalyticsReferrersRow) []DashboardNamedCount {
	result := make([]DashboardNamedCount, 0, len(rows))
	for _, row := range rows {
		name := strings.TrimSpace(toBucketString(row.Name))
		if name == "" {
			name = "direct"
		}
		result = append(result, DashboardNamedCount{Name: name, Visitors: row.Visitors})
	}
	return result
}

func mapChannelRows(rows []query.ListAnalyticsChannelsRow) []DashboardNamedCount {
	result := make([]DashboardNamedCount, 0, len(rows))
	for _, row := range rows {
		name := strings.TrimSpace(row.Name)
		if name == "" {
			name = "Unknown"
		}
		result = append(result, DashboardNamedCount{Name: name, Visitors: row.Visitors})
	}
	return result
}

func mapBrowserRows(rows []query.ListAnalyticsBrowsersRow) []DashboardNamedCount {
	result := make([]DashboardNamedCount, 0, len(rows))
	for _, row := range rows {
		name := strings.TrimSpace(row.Name)
		if name == "" {
			name = "Unknown"
		}
		result = append(result, DashboardNamedCount{Name: name, Visitors: row.Visitors})
	}
	return result
}

func mapOperatingSystemRows(rows []query.ListAnalyticsOperatingSystemsRow) []DashboardNamedCount {
	result := make([]DashboardNamedCount, 0, len(rows))
	for _, row := range rows {
		name := strings.TrimSpace(row.Name)
		if name == "" {
			name = "Unknown"
		}
		result = append(result, DashboardNamedCount{Name: name, Visitors: row.Visitors})
	}
	return result
}

func mapDeviceRows(rows []query.ListAnalyticsDevicesRow) []DashboardNamedCount {
	result := make([]DashboardNamedCount, 0, len(rows))
	for _, row := range rows {
		name := strings.TrimSpace(row.Name)
		if name == "" {
			name = "Unknown"
		}
		result = append(result, DashboardNamedCount{Name: name, Visitors: row.Visitors})
	}
	return result
}

func mapRegionRows(rows []query.ListAnalyticsRegionsRow) []DashboardNamedCount {
	result := make([]DashboardNamedCount, 0, len(rows))
	for _, row := range rows {
		name := strings.TrimSpace(row.Name)
		if name == "" {
			name = "Unknown"
		}
		result = append(result, DashboardNamedCount{Name: name, Visitors: row.Visitors})
	}
	return result
}

func mapCityRows(rows []query.ListAnalyticsCitiesRow) []DashboardNamedCount {
	result := make([]DashboardNamedCount, 0, len(rows))
	for _, row := range rows {
		name := strings.TrimSpace(row.Name)
		if name == "" {
			name = "Unknown"
		}
		result = append(result, DashboardNamedCount{Name: name, Visitors: row.Visitors})
	}
	return result
}

func mapCountries(rows []query.ListAnalyticsCountriesRow) []DashboardCountryCount {
	result := make([]DashboardCountryCount, 0, len(rows))
	for _, row := range rows {
		code := normalizeCountryCode(row.CountryCode)
		result = append(result, DashboardCountryCount{Code: code, Visitors: row.Visitors})
	}
	return result
}

func mapTraffic(rows []query.GetAnalyticsTrafficHeatmapRow) []DashboardTrafficPoint {
	result := make([]DashboardTrafficPoint, 0, len(rows))
	for _, row := range rows {
		result = append(result, DashboardTrafficPoint{Dow: int(row.Dow), Hour: int(row.Hour), Value: row.Value})
	}
	return result
}

func normalizeOccurredAt(occurredAt, now time.Time) time.Time {
	if occurredAt.IsZero() {
		return now
	}
	if occurredAt.Before(now.AddDate(0, 0, -90)) || occurredAt.After(now.Add(10*time.Minute)) {
		return now
	}
	return occurredAt
}

func normalizePath(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "/"
	}
	if u, err := url.Parse(trimmed); err == nil {
		if u.Path != "" {
			trimmed = u.Path
		}
	}
	if !strings.HasPrefix(trimmed, "/") {
		trimmed = "/" + trimmed
	}
	trimmed = strings.ReplaceAll(trimmed, "//", "/")
	if len(trimmed) > 180 {
		trimmed = trimmed[:180]
	}
	return trimmed
}

func normalizeReferrer(ref string) string {
	trimmed := strings.TrimSpace(ref)
	if trimmed == "" {
		return ""
	}
	if len(trimmed) > 400 {
		trimmed = trimmed[:400]
	}
	return trimmed
}

func normalizeHost(raw string) string {
	if raw == "" {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	host := strings.ToLower(strings.TrimSpace(u.Hostname()))
	if strings.HasPrefix(host, "www.") {
		host = strings.TrimPrefix(host, "www.")
	}
	return trimRunes(host, 140)
}

func classifyChannel(refHost string) string {
	if refHost == "" {
		return "direct"
	}
	searchSources := []string{"google.", "bing.", "baidu.", "duckduckgo.", "yahoo.", "sogou.", "so.com", "yandex."}
	for _, token := range searchSources {
		if strings.Contains(refHost, token) {
			return "search"
		}
	}
	socialSources := []string{"x.com", "twitter.com", "weibo.com", "t.co", "facebook.com", "instagram.com", "linkedin.com", "reddit.com", "tiktok.com", "bilibili.com", "zhihu.com"}
	for _, token := range socialSources {
		if strings.Contains(refHost, token) {
			return "social"
		}
	}
	return "referral"
}

func parseUserAgent(ua string) (browser, osName, device string) {
	uaLower := strings.ToLower(strings.TrimSpace(ua))

	browser = "Unknown"
	switch {
	case strings.Contains(ua, "Edg/"):
		browser = "Edge"
	case strings.Contains(ua, "OPR/") || strings.Contains(ua, "Opera"):
		browser = "Opera"
	case strings.Contains(ua, "Firefox/"):
		browser = "Firefox"
	case strings.Contains(ua, "Chrome/") && !strings.Contains(ua, "Edg/"):
		browser = "Chrome"
	case strings.Contains(ua, "Safari/") && strings.Contains(ua, "Version/"):
		browser = "Safari"
	case strings.Contains(ua, "MicroMessenger"):
		browser = "WeChat"
	}

	osName = "Unknown"
	switch {
	case strings.Contains(uaLower, "windows"):
		osName = "Windows"
	case strings.Contains(uaLower, "mac os") || strings.Contains(uaLower, "macintosh"):
		osName = "macOS"
	case strings.Contains(uaLower, "android"):
		osName = "Android"
	case strings.Contains(uaLower, "iphone") || strings.Contains(uaLower, "ipad") || strings.Contains(uaLower, "ios"):
		osName = "iOS"
	case strings.Contains(uaLower, "linux"):
		osName = "Linux"
	}

	device = "Desktop"
	switch {
	case strings.Contains(uaLower, "ipad") || strings.Contains(uaLower, "tablet"):
		device = "Tablet"
	case strings.Contains(uaLower, "mobile") || strings.Contains(uaLower, "iphone") || strings.Contains(uaLower, "android"):
		device = "Mobile"
	}

	return browser, osName, device
}

func normalizeCountryCode(code string) string {
	trimmed := strings.ToUpper(strings.TrimSpace(code))
	if len(trimmed) != 2 {
		return "ZZ"
	}
	for _, ch := range trimmed {
		if ch < 'A' || ch > 'Z' {
			return "ZZ"
		}
	}
	return trimmed
}

func normalizeGeoName(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "Unknown"
	}
	return trimRunes(trimmed, 80)
}

func normalizeSessionKey(v string) string {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return ""
	}

	buf := strings.Builder{}
	for _, ch := range trimmed {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '-' || ch == '_' {
			buf.WriteRune(ch)
		}
		if buf.Len() >= 80 {
			break
		}
	}
	return buf.String()
}

func toPgUUID(id uuid.UUID) pgtype.UUID {
	if id == uuid.Nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: id, Valid: true}
}

func toBucketString(v interface{}) string {
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	default:
		return fmt.Sprint(x)
	}
}

func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func percentage(part, total float64) float64 {
	if total <= 0 {
		return 0
	}
	return (part / total) * 100
}

func percentChange(current, previous float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return round1(((current - previous) / previous) * 100)
}

func round1(v float64) float64 {
	return math.Round(v*10) / 10
}

func trimRunes(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if len(s) <= max {
		return s
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max])
}
