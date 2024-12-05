package schema

type RouteName struct {
	Name string `json:"name"`
}

type Route struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	ExecutableExists bool   `json:"executableExists"`
}
