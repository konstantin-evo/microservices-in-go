package data

type EventSeverity string

const (
	SeverityLog     EventSeverity = "log.INFO"
	SeverityWarning EventSeverity = "log.WARNING"
	SeverityInfo    EventSeverity = "log.ERROR"
)

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
