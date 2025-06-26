package sim

import (
	"testing"

	av "github.com/mmp/vice/pkg/aviation"
)

func TestDefaultTCPAssignment(t *testing.T) {
	adapt := STARSFacilityAdaptation{
		TCPConfiguration: TCPConfiguration{
			FixPairConfiguration: map[string]*FixPairFacilityConfiguration{
				"ABC": {
					SectorConfigurations:  []SectorConfiguration{{ConfigurationID: "CFG", DefaultConfiguration: true}},
					FixPairs:              []FixPair{{TerminalSector: "A", FlightType: "arrival", EntryFix: "AAA", ExitFix: "BBB"}},
					FixPairTCPAssignments: []FixPairTCPAssignment{{Configuration: []string{"CFG"}, FlightType: "arrival", EntryFix: "AAA", ExitFix: "BBB", TCP: "1A"}},
				},
			},
		},
	}
	sc := makeSTARSComputer("ABC", &adapt)
	fp := STARSFlightPlan{EntryFix: "AAA", ExitFix: "BBB", TypeOfFlight: av.FlightTypeArrival}
	fp, err := sc.CreateFlightPlan(fp)
	if err != nil {
		t.Fatalf("CreateFlightPlan: %v", err)
	}
	if fp.TrackingController != "1A" {
		t.Errorf("expected TCP 1A, got %s", fp.TrackingController)
	}
}