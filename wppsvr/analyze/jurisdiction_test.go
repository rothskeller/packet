package analyze

// Leave this test disabled most of the time.  We can uncomment it when we need it.
/*
import "testing"

func TestFetchJurisdiction(t *testing.T) {
	var a Analysis
	a.sm.FromCallSign = "KC6RSC"
	a.fetchJurisdiction()
	if a.sm.Jurisdiction != "SNY" {
		t.Errorf("Expected SNY, got %q", a.sm.Jurisdiction)
	}
	a.sm.FromCallSign = "KE6TIM"
	a.sm.Jurisdiction = ""
	a.fetchJurisdiction()
	if a.sm.Jurisdiction != "MLP" {
		t.Errorf("Expected MLP, got %q", a.sm.Jurisdiction)
	}
}
*/
