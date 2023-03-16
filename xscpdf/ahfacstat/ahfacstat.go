package ahfacstat

import (
	_ "embed" // oh, please
	"encoding/hex"
	"regexp"

	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/ahfacstat"
	"github.com/rothskeller/packet/xscpdf"
)

//go:embed Allied_Health_Facility_Status_DEOC-9_v20180200_with_XSC_RACES_Routing_Slip_Fillable_v20190527.pdf
var basePDF []byte

var idFormAndRouting, _ = hex.DecodeString("3d6c1e15e585a50f819d74dc910f04c5")
var idFormOnly, _ = hex.DecodeString("725f4fc4fac15e4da1759c47ac3a2186")

func init() {
	xscpdf.RegisterReader(xscpdf.ReaderMap{XSCTag: ahfacstat.Tag, PDFID: idFormAndRouting, Fields: fieldMap})
	xscpdf.RegisterReader(xscpdf.ReaderMap{XSCTag: ahfacstat.Tag, PDFID: idFormOnly, Fields: fieldMap})
	xscpdf.RegisterWriter(xscpdf.WriterMap{XSCTag: ahfacstat.Tag, BasePDF: basePDF, Fields: fieldMap})
}

var dateRE = regexp.MustCompile(`\s*((?:0?[1-9]|1[012])/(?:0?[1-9]|[12][0-9]|3[01])/(?:20)?[0-9][0-9])$`)
var fieldMap = []xscpdf.FieldMap{
	{PDFName: "Origin Msg Nbr", XSCTag: "MsgNo"},          // 1
	{PDFName: "Destination Msg Nbr", XSCTag: "DestMsgNo"}, // 2
	{PDFName: "Date Created", XSCTag: "1a."},              // 3
	{PDFName: "Time Created", XSCTag: "1b."},              // 4
	{PDFName: "Handling", XSCTag: "5.", Values: []xscpdf.ValueMap{ // 5
		{PDFValue: "Off", XSCValue: ""},
		{PDFValue: "Immediate", XSCValue: "IMMEDIATE"},
		{PDFValue: "Priority", XSCValue: "PRIORITY"},
		{PDFValue: "Routine", XSCValue: "ROUTINE"},
	}},
	{PDFName: "To ICS Position", XSCTag: "7a."},   // 6
	{PDFName: "To Location", XSCTag: "7b."},       // 7
	{PDFName: "To Name", XSCTag: "7c."},           // 8
	{PDFName: "To Contact Info", XSCTag: "7d."},   // 9
	{PDFName: "From ICS Position", XSCTag: "8a."}, // 10
	{PDFName: "From Location", XSCTag: "8b."},     // 11
	{PDFName: "From Name", XSCTag: "8c."},         // 12
	{PDFName: "From Contact Info", XSCTag: "8d."}, // 13
	{PDFName: "Form Type", FromXSC: func(*xscmsg.Message) string { return "Allied Health Facility Status" }}, // 14
	{PDFName: "Form Topic", XSCTag: "21."},                           // 15
	{PDFName: "Relay Rcvd", XSCTag: "OpRelayRcvd"},                   // 16
	{PDFName: "Relay Sent", XSCTag: "OpRelaySent"},                   // 17
	{PDFName: "Op Name", XSCTag: "OpName"},                           // 18
	{PDFName: "Op Call Sign", XSCTag: "OpCall"},                      // 19
	{PDFName: "Op Date", XSCTag: "OpDate"},                           // 20
	{PDFName: "Op Time", XSCTag: "OpTime"},                           // 21
	{PDFName: "FACILTY TYPE TIME DATE FACILITY NAME", XSCTag: "20."}, // 22
	{PDFName: "Contact Name", XSCTag: "23."},                         // 23
	{PDFName: "Phone", XSCTag: "23p."},                               // 24
	{PDFName: "Fax", XSCTag: "23f."},                                 // 25
	{PDFName: "Other Phone Fax Cell Phone Radio", XSCTag: "24."},     // 26
	{PDFName: "Incident Name and Date", FromXSC: func(xsc *xscmsg.Message) string {
		return xsc.Field("25.").Value + " " + xsc.Field("25d.").Value
	}}, // 27
	{XSCTag: "25.", FromPDF: func(pfields map[string]string) string {
		ind := pfields["Incident Name and Date"]
		if match := dateRE.FindString(ind); match != "" {
			return ind[:len(ind)-len(match)]
		}
		return ind
	}},
	{XSCTag: "25d.", FromPDF: func(pfields map[string]string) string {
		ind := pfields["Incident Name and Date"]
		if match := dateRE.FindStringSubmatch(ind); match != nil {
			return match[1]
		}
		return ""
	}},
	{PDFName: "CHECK ONEGREEN FULLY FUNCTIONAL", FromXSC: func(xsc *xscmsg.Message) string {
		if xsc.Field("35.").Value == "Fully Functional" {
			return "X"
		}
		return ""
	}}, // 28
	{XSCTag: "35.", FromPDF: func(pfields map[string]string) string {
		switch {
		case pfields["CHECK ONEGREEN FULLY FUNCTIONAL"] != "":
			return "Fully Functional"
		case pfields["CHECK ONERED LIMITED SERVICES"] != "":
			return "Limited Services"
		case pfields["CHECK ONEBLACK IMPAIREDCLOSED"] != "":
			return "Impaired/Closed"
		default:
			return ""
		}
	}},
	{PDFName: "YesNoNHICSICS ORGANIZATION CHART", XSCTag: "26a."}, // 29
	{PDFName: "CHECK ONERED LIMITED SERVICES", FromXSC: func(xsc *xscmsg.Message) string {
		if xsc.Field("35.").Value == "Limited Services" {
			return "X"
		}
		return ""
	}}, // 30
	{PDFName: "YesNoDEOC9A RESOURCE REQUEST FORMS", XSCTag: "26b."}, // 31
	{PDFName: "CHECK ONEBLACK IMPAIREDCLOSED", FromXSC: func(xsc *xscmsg.Message) string {
		if xsc.Field("35.").Value == "Impaired/Closed" {
			return "X"
		}
		return ""
	}}, // 32
	{PDFName: "YesNoNHICSICS STATUS REPORT FORM  STANDARD", XSCTag: "26c."},                           // 33
	{PDFName: "YesNoNHICSICS INCIDENT ACTION PLAN", XSCTag: "26d."},                                   // 34
	{PDFName: "EOC MAIN CONTACT NUMBER", XSCTag: "27p."},                                              // 35
	{PDFName: "YesNoPHONECOMMUNICATIONS DIRECTORY", XSCTag: "26e."},                                   // 36
	{PDFName: "EOC MAIN CONTACT FAX", XSCTag: "27f."},                                                 // 37
	{PDFName: "NAME LIAISON TO PUBLIC HEALTHMEDICAL HEALTH BRANCH", XSCTag: "28."},                    // 38
	{PDFName: "CONTACT NUMBER", XSCTag: "28p."},                                                       // 39
	{PDFName: "INFORMATION OFFICER NAME", XSCTag: "29."},                                              // 40
	{PDFName: "CONTACT NUMBER_2", XSCTag: "29p."},                                                     // 41
	{PDFName: "CONTACT EMAIL", XSCTag: "29e."},                                                        // 42
	{PDFName: "GENERAL SUMMARY OF SITUATIONCONDITIONSRow1", XSCTag: "34."},                            // 43
	{PDFName: "IF EOC IS NOT ACTIVATED WHO SHOULD BE CONTACTED FOR QUESTIONSREQUESTS", XSCTag: "30."}, // 44
	{PDFName: "CONTACT NUMBER_3", XSCTag: "30p."},                                                     // 45
	{PDFName: "Staffed Bed MSKILLED NURSING", XSCTag: "40a."},                                         // 46
	{PDFName: "Staffed BedFSKILLED NURSING", XSCTag: "40b."},                                          // 47
	{PDFName: "Vacant BedsMSKILLED NURSING", XSCTag: "40c."},                                          // 48
	{PDFName: "Vacant BedFSKILLED NURSING", XSCTag: "40d."},                                           // 49
	{PDFName: "Surge SKILLED NURSING", XSCTag: "40e."},                                                // 50
	{PDFName: "CONTACT EMAIL_2", XSCTag: "30e."},                                                      // 51
	{PDFName: "Staffed Bed MASSISTED LIVING", XSCTag: "41a."},                                         // 52
	{PDFName: "Staffed BedFASSISTED LIVING", XSCTag: "41b."},                                          // 53
	{PDFName: "Vacant BedsMASSISTED LIVING", XSCTag: "41c."},                                          // 54
	{PDFName: "Vacant BedFASSISTED LIVING", XSCTag: "41d."},                                           // 55
	{PDFName: "Surge ASSISTED LIVING", XSCTag: "41e."},                                                // 56
	// Yes, PDF fields 57-61, 63-67, 69-73, 75-79, and 81-85 have incorrect
	// names in the PDF.  Sigh.
	{PDFName: "Staffed Bed MSURGICAL BEDS", XSCTag: "42a."},             // 57
	{PDFName: "Staffed BedFSURGICAL BEDS", XSCTag: "42b."},              // 58
	{PDFName: "Vacant BedsMSURGICAL BEDS", XSCTag: "42c."},              // 59
	{PDFName: "Vacant BedFSURGICAL BEDS", XSCTag: "42d."},               // 60
	{PDFName: "Surge SURGICAL BEDS", XSCTag: "42e."},                    // 61
	{PDFName: "TOTALPATIENTS TO EVACUATE", XSCTag: "31a."},              // 62
	{PDFName: "Staffed Bed MSUBACUTE", XSCTag: "43a."},                  // 63
	{PDFName: "Staffed BedFSUBACUTE", XSCTag: "43b."},                   // 64
	{PDFName: "Vacant BedsMSUBACUTE", XSCTag: "43c."},                   // 65
	{PDFName: "Vacant BedFSUBACUTE", XSCTag: "43d."},                    // 66
	{PDFName: "Surge SUBACUTE", XSCTag: "43e."},                         // 67
	{PDFName: "TOTALPATIENTS  INJURED  MINOR", XSCTag: "31b."},          // 68
	{PDFName: "Staffed Bed MALZEIMERSDIMENTIA", XSCTag: "44a."},         // 69
	{PDFName: "Staffed BedFALZEIMERSDIMENTIA", XSCTag: "44b."},          // 70
	{PDFName: "Vacant BedsMALZEIMERSDIMENTIA", XSCTag: "44c."},          // 71
	{PDFName: "Vacant BedFALZEIMERSDIMENTIA", XSCTag: "44d."},           // 72
	{PDFName: "Surge ALZEIMERSDIMENTIA", XSCTag: "44e."},                // 73
	{PDFName: "TOTALPATIENTS TRANSFERED OUT OF COUNTY", XSCTag: "31c."}, // 74
	{PDFName: "Staffed Bed MPEDIATRICSUB ACUTE", XSCTag: "45a."},        // 75
	{PDFName: "Staffed BedFPEDIATRICSUB ACUTE", XSCTag: "45b."},         // 76
	{PDFName: "Vacant BedsMPEDIATRICSUB ACUTE", XSCTag: "45c."},         // 77
	{PDFName: "Vacant BedFPEDIATRICSUB ACUTE", XSCTag: "45d."},          // 78
	{PDFName: "Surge PEDIATRICSUB ACUTE", XSCTag: "45e."},               // 79
	{PDFName: "OTHER PATIENT CARE INFORMATION", XSCTag: "33."},          // 80
	{PDFName: "Staffed Bed MPSYCHIATRIC", XSCTag: "46a."},               // 81
	{PDFName: "Staffed BedFPSYCHIATRIC", XSCTag: "46b."},                // 82
	{PDFName: "Vacant BedsMPSYCHIATRIC", XSCTag: "46c."},                // 83
	{PDFName: "Vacant BedFPSYCHIATRIC", XSCTag: "46d."},                 // 84
	{PDFName: "Surge PSYCHIATRIC", XSCTag: "46e."},                      // 85
	{PDFName: "CHAIRS ROOMSDIALYSIS", XSCTag: "50a."},                   // 87
	{PDFName: "VANCANT CHAIRS ROOMDIALYSIS", XSCTag: "50b."},            // 88
	{PDFName: "FRONT DESK STAFFDIALYSIS", XSCTag: "50c."},               // 89
	{PDFName: "MEDICAL SUPPORT STAFFDIALYSIS", XSCTag: "50d."},          // 90
	{PDFName: "PROVIDER STAFFDIALYSIS", XSCTag: "50e."},                 // 91
	{PDFName: "CHAIRS ROOMSSURGICAL", XSCTag: "51a."},                   // 92
	{PDFName: "VANCANT CHAIRS ROOMSURGICAL", XSCTag: "51b."},            // 93
	{PDFName: "FRONT DESK STAFFSURGICAL", XSCTag: "51c."},               // 94
	{PDFName: "MEDICAL SUPPORT STAFFSURGICAL", XSCTag: "51d."},          // 95
	{PDFName: "PROVIDER STAFFSURGICAL", XSCTag: "51e."},                 // 96
	{PDFName: "CHAIRS ROOMSCLINIC", XSCTag: "52a."},                     // 97
	{PDFName: "VANCANT CHAIRS ROOMCLINIC", XSCTag: "52b."},              // 98
	{PDFName: "FRONT DESK STAFFCLINIC", XSCTag: "52c."},                 // 99
	{PDFName: "MEDICAL SUPPORT STAFFCLINIC", XSCTag: "52d."},            // 100
	{PDFName: "PROVIDER STAFFCLINIC", XSCTag: "52e."},                   // 101
	{PDFName: "CHAIRS ROOMSHOMEHEALTH", XSCTag: "53a."},                 // 102
	{PDFName: "VANCANT CHAIRS ROOMHOMEHEALTH", XSCTag: "53b."},          // 103
	{PDFName: "FRONT DESK STAFFHOMEHEALTH", XSCTag: "53c."},             // 104
	{PDFName: "MEDICAL SUPPORT STAFFHOMEHEALTH", XSCTag: "53d."},        // 105
	{PDFName: "PROVIDER STAFFHOMEHEALTH", XSCTag: "53e."},               // 106
	{PDFName: "CHAIRS ROOMSADULT DAY CENTER", XSCTag: "54a."},           // 107
	{PDFName: "VANCANT CHAIRS ROOMADULT DAY CENTER", XSCTag: "54b."},    // 108
	{PDFName: "FRONT DESK STAFFADULT DAY CENTER", XSCTag: "54c."},       // 109
	{PDFName: "MEDICAL SUPPORT STAFFADULT DAY CENTER", XSCTag: "54d."},  // 110
	{PDFName: "PROVIDER STAFFADULT DAY CENTER", XSCTag: "54e."},         // 111
	{PDFName: "facility type", XSCTag: "21."},                           // 112
	{PDFName: "date", XSCTag: "22d."},                                   // 113
	{PDFName: "time", XSCTag: "22t."},                                   // 114
}
