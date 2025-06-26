package sim

import (
	"slices"
	"strings"

	av "github.com/mmp/vice/pkg/aviation"
)

// TCPSpec describes a terminal control position.
type TCPSpec struct {
	TCPName           string   `json:"tcp_name"`
	FacilityID        string   `json:"facility_id"`
	TerminalSector    string   `json:"terminal_sector"`
	ExcludedExitFixes []string `json:"excluded_exit_fixes"`
}

// SectorConfiguration identifies a set of sectors active for a configuration.
type SectorConfiguration struct {
	ConfigurationID      string `json:"configuration_id"`
	ConfigurationName    string `json:"configuration_name"`
	DefaultConfiguration bool   `json:"default_configuration"`
}

// FixPair describes an entry/exit pair for a flight type.
type FixPair struct {
	TerminalSector string `json:"terminal_sector"`
	FlightType     string `json:"flight_type"`
	EntryFix       string `json:"entry_fix"`
	ExitFix        string `json:"exit_fix"`
}

// FixPairTCPAssignment maps a fix pair to a TCP for one or more configurations.
type FixPairTCPAssignment struct {
	Configuration []string `json:"configuration"`
	FlightType    string   `json:"flight_type"`
	EntryFix      string   `json:"entry_fix"`
	ExitFix       string   `json:"exit_fix"`
	TCP           string   `json:"tcp"`
}

// FixPairConfiguration groups all fix pair data for a facility.
type FixPairConfiguration struct {
	SectorConfigurations  []SectorConfiguration  `json:"sector_configurations"`
	FixPairs              []FixPair              `json:"fix_pairs"`
	FixPairTCPAssignments []FixPairTCPAssignment `json:"fix_pair_tcp_assignments"`
}

// TCPConfiguration contains TCP specifications and fix pair rules.
type TCPConfiguration struct {
	TCPs                 map[string]TCPSpec              `json:"tcps"`
	FixPairConfiguration map[string]FixPairConfiguration `json:"fix_pair_configuration"`
}

func flightTypeString(ft av.TypeOfFlight) string {
	switch ft {
	case av.FlightTypeArrival:
		return "arrival"
	case av.FlightTypeDeparture:
		return "departure"
	case av.FlightTypeOverflight:
		return "overflight"
	default:
		return ""
	}
}

// TCPForFixPair returns the TCP for the given fix pair using the active
// default sector configuration.
func (fa STARSFacilityAdaptation) TCPForFixPair(ft av.TypeOfFlight, entry, exit string) string {
	t := flightTypeString(ft)
	for _, fac := range fa.TCPConfiguration.FixPairConfiguration {
		var defaults []string
		for _, sc := range fac.SectorConfigurations {
			if sc.DefaultConfiguration {
				defaults = append(defaults, sc.ConfigurationID)
			}
		}
		for _, as := range fac.FixPairTCPAssignments {
			if as.FlightType == t && strings.EqualFold(as.EntryFix, entry) && strings.EqualFold(as.ExitFix, exit) {
				for _, cfg := range as.Configuration {
					if slices.Contains(defaults, cfg) {
						return as.TCP
					}
				}
			}
		}
	}
	return ""
}

// TRACONForFixPair returns the TRACON whose configuration contains the given
// fix pair. If none is found, an empty string is returned.
func (fa STARSFacilityAdaptation) TRACONForFixPair(ft av.TypeOfFlight, entry, exit string) string {
	t := flightTypeString(ft)
	for fac, cfg := range fa.TCPConfiguration.FixPairConfiguration {
		for _, fp := range cfg.FixPairs {
			if strings.EqualFold(fp.EntryFix, entry) && strings.EqualFold(fp.ExitFix, exit) &&
				strings.EqualFold(fp.FlightType, t) {
				return fac
			}
		}
	}
	return ""
}