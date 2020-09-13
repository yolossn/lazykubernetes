package gui

type namespaceState struct {
	SelectedLine int
}

type resourceState struct {
	SelectedLine int
	TabIndex     int
}

type infoState struct {
	SelectedLine int
	TabIndex     int
}

type panelStates struct {
	Namespace *namespaceState
	Resource  *resourceState
	Info      *infoState
}

func NewPanelStates() *panelStates {
	ns := &namespaceState{}
	rs := &resourceState{}
	return &panelStates{Namespace: ns, Resource: rs}
}
