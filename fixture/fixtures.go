package fixture

const (
	// YFinQuotes are the yfin quote responses.
	YFinQuotes ResourceID = "quote"
	// YFinQuotePath is the yfin basic quote path.
	YFinQuotePath Path = "/v7/finance/quote"

	// YFinChart are the yfin chart responses.
	YFinChart ResourceID = "chart"
	// YFinChartPath is the yfin chart path.
	YFinChartPath Path = "/v8/finance/chart"

	// ServiceYFin is the yfin service.
	ServiceYFin ServiceID = "yfin"
)

// Path is a url path.
type Path string

// ResourceID is just an identifier for a resource.
type ResourceID string

// ServiceID is just an identifier for a service.
type ServiceID string

// Resources alias for resource map.
type Resources map[ResourceID]interface{}

// Fixtures is a collection of resources.
type Fixtures struct {
	Resources map[ServiceID]Resources `json:"resources"`
}

// Spec specification of services.
type Spec struct {
	Services map[ServiceID]*Service `yaml:"services"`
}

// Service is a collection of url paths and resources.
type Service struct {
	Paths map[Path]*Operation `yaml:"paths"`
}

// Operation defines a service operation.
type Operation struct {
	Parameters []*Parameter `yaml:"parameters"`
	ResourceID ResourceID   `yaml:"resource"`
}

// Parameter describes a url parameter.
type Parameter struct {
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	Required    bool   `yaml:"required"`
}
