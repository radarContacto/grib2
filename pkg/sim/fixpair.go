package sim

// FixPair describes an entry and exit fix pair for a particular flight type.
type FixPair struct {
	TerminalSector string `json:"terminal_sector"`
	FlightType     string `json:"flight_type"`
	EntryFix       string `json:"entry_fix"`
	ExitFix        string `json:"exit_fix"`
}

// FixPairTCPAssignment maps a fix pair to a TCP for one or more sector configurations.
type FixPairTCPAssignment struct {
	Configuration []string `json:"configuration"`
	FlightType    string   `json:"flight_type"`
	EntryFix      string   `json:"entry_fix"`
	ExitFix       string   `json:"exit_fix"`
	TCP           string   `json:"tcp"`
}

// FixPairFacilityConfiguration groups all fix pair data for a facility.
type FixPairFacilityConfiguration struct {
	SectorConfigurations  []SectorConfiguration  `json:"sector_configurations"`
	FixPairs              []FixPair              `json:"fix_pairs"`
	FixPairTCPAssignments []FixPairTCPAssignment `json:"fix_pair_tcp_assignments"`
}

// TCPConfiguration contains TCP specifications and facility fix pair rules.
type TCPConfiguration struct {
	TCPs                 map[string]TCPSpec                       `json:"tcps"`
	FixPairConfiguration map[string]*FixPairFacilityConfiguration `json:"fix_pair_configuration"`
}