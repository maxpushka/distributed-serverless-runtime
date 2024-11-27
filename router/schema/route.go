package schema

type RouteName struct {
	Name string `json:"name"`
}

type Route struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	ConfigExists     bool   `json:"configExists"`
	ExecutableExists bool   `json:"executableExists"`
}
