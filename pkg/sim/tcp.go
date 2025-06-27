package sim

import (
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
func tcpForFixPair(cfg *FixPairFacilityConfiguration, ft av.TypeOfFlight, entry, exit string) string {
	var defaults []string
	for _, sc := range cfg.SectorConfigurations {
		if sc.DefaultConfiguration {
			defaults = append(defaults, sc.ConfigurationID)
		}
	}
	t := flightTypeString(ft)
	for _, as := range cfg.FixPairTCPAssignments {
		if as.FlightType == t && strings.EqualFold(as.EntryFix, entry) && strings.EqualFold(as.ExitFix, exit) {
			if len(as.Configuration) == 0 {
				return as.TCP
			}
		}
	}
	return ""
}

// TCPForFixPair returns the TCP for the given fix pair using the active
// default sector configuration.
func (fa *STARSFacilityAdaptation) TCPForFixPair(tracon string, ft av.TypeOfFlight, entry, exit string) string {
	cfg, ok := fa.TCPConfiguration.FixPairConfiguration[tracon]
	if !ok {
		return ""
	}
	return tcpForFixPair(cfg, ft, entry, exit)
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
