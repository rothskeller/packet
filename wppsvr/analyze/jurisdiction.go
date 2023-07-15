package analyze

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// TestForceJurisdiction is normally an empty string.  When set to a nonempty
// value by tests, it disables the web requests to fetch ham information and
// instead sets all jurisdictions to the value of TestForceJurisdiction.
var TestForceJurisdiction string

var jurisdictionMap = map[string]string{
	"Alameda County":                    "XAL",
	"American Red Cross":                "ARC",
	"CalFIRE Santa Clara Unit":          "SCU",
	"CalOES Coastal Region":             "COS",
	"Campbell":                          "CBL",
	"Contra Costa County":               "XCC",
	"Cupertino":                         "CUP",
	"Gilroy":                            "GIL",
	"Hospitals":                         "HOS",
	"Loma Prieta":                       "LMP",
	"Los Altos":                         "LOS",
	"Los Altos Hills":                   "LAH",
	"Los Gatos":                         "LGT",
	"Marin County":                      "XMR",
	"Milpitas":                          "MLP",
	"Monte Sereno":                      "MSO",
	"Monterey County":                   "XMY",
	"Morgan Hill":                       "MRG",
	"Mountain View":                     "MTV",
	"NASA/AMES":                         "NAM",
	"Palo Alto":                         "PAF",
	"San Benito County":                 "XBE",
	"San Francisco County":              "XSF",
	"San Jose":                          "SJC",
	"San Jose Water Co":                 "SJW",
	"San Mateo County":                  "XSM",
	"Santa Clara":                       "SNC",
	"Santa Clara County":                "XSC",
	"Santa Clara Valley Water District": "VWD",
	"Santa Cruz County":                 "XCZ",
	"Saratoga":                          "SAR",
	"Stanford University":               "STU",
	"Sunnyvale":                         "SNY",
	"Unincorporated":                    "XSC",
}

type getHamInfoResponse struct {
	CallSign      string
	Last          string
	First         string
	HomeAgency    string
	OtherAgencies []string
}

func (a *Analysis) fetchJurisdiction() {
	var (
		timeout context.Context
		cancel  context.CancelFunc
		req     *http.Request
		resp    *http.Response
		ghis    []*getHamInfoResponse
		err     error
	)
	if a.sm.FromCallSign == "" {
		return
	}
	if !fccCallSignRE.MatchString(a.sm.FromCallSign) {
		a.sm.Jurisdiction = a.sm.FromCallSign[:3]
		return
	}
	if TestForceJurisdiction != "" {
		a.sm.Jurisdiction = TestForceJurisdiction
		return
	}
	timeout, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err = http.NewRequestWithContext(timeout, http.MethodGet,
		"https://www.scc-ares-races.org/activities/getHamInfo.php?id="+a.sm.FromCallSign, nil)
	if err != nil {
		panic(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("ERROR: unable to fetch SCCo database info for %s: %s", a.sm.FromCallSign, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("ERROR: unable to fetch SCCo database info for %s: status code %d", a.sm.FromCallSign, resp.StatusCode)
		return
	}
	if err = json.NewDecoder(resp.Body).Decode(&ghis); err != nil {
		log.Printf("ERROR: unable to fetch SCCo database info for %s: json.Decode: %s", a.sm.FromCallSign, err)
		return
	}
	if len(ghis) != 1 {
		log.Printf("ERROR: unable to fetch SCCo database info for %s: %d responses", a.sm.FromCallSign, len(ghis))
		return
	}
	a.sm.Jurisdiction = jurisdictionMap[ghis[0].HomeAgency]
	if a.sm.Jurisdiction == "" || a.sm.Jurisdiction == "XSC" || a.sm.Jurisdiction == "HOS" {
		for _, other := range ghis[0].OtherAgencies {
			if juris := jurisdictionMap[other]; juris != "" && juris != "XSC" && juris != "HOS" {
				a.sm.Jurisdiction = juris
				break
			}
		}
	}
}
