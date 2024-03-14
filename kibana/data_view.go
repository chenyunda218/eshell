package kibana

type DataView struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Namespaces []string `json:"namespaces"`
	Title      string   `json:"title"`
	Type       string   `json:"type"`
	TypeMeta   any      `json:"typeMeta"`
}
