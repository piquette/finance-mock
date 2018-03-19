package fixture

type HTTPVerb string

type ResourceID string

type StatusCode string

type Path string

type Fixtures struct {
	Resources map[ResourceID]interface{} `yaml:"resources"`
}

type Spec struct {
	Paths map[Path]map[HTTPVerb]*Operation `yaml:"paths"`
}

type Operation struct {
	Description string                  `yaml:"description"`
	OperationID string                  `yaml:"operation_id"`
	Parameters  []*Parameter            `yaml:"parameters"`
	Responses   map[StatusCode]Response `yaml:"responses"`
}

type Parameter struct {
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	Required    bool   `yaml:"required"`
}

type Response struct {
	Content map[string]ResourceID `yaml:"content"`
}
