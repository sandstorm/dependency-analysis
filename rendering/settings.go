package rendering

type RenderingSettings struct {
	GraphLabel     string
	ShowLegend     bool // TODO: implement legend
	ShowNodeLabels bool
}

func NewRenderingSettings() *RenderingSettings {
	return &RenderingSettings{
		GraphLabel:     "rendered by github.com/sandstorm/dependency-analysis",
		ShowLegend:     true,
		ShowNodeLabels: true,
	}
}
