package aviation

// FacilityControlPosition describes a controller at a facility in the codex
// format.
type FacilityControlPosition struct {
	SectorRoutingID string    `json:"sector_routing_id"`
	Frequency       Frequency `json:"frequency"`
	RadioName       string    `json:"radio_name"`
}

// Facility describes an air traffic facility and its controllers in the codex
// scenario format.
type Facility struct {
	FacilityType          string                    `json:"facility_type"`
	FacilityName          string                    `json:"facility_name"`
	AdjacentTRACONID      int                       `json:"adjacent_tracon_id"`
	AbbreviatedFacilityID string                    `json:"abbreviated_facility_id"`
	TerminalSectors       map[string]string         `json:"terminal_sectors"`
	ControlPositions      []FacilityControlPosition `json:"control_positions"`
}