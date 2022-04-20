package pktmsg_test

import (
	"testing"
	"time"

	"steve.rothskeller.net/packet/pktmsg"
)

func TestPlainEncode(t *testing.T) {
	var msg pktmsg.TxMessage
	msg.MessageNumber = "AAA-111"
	msg.HandlingOrder = pktmsg.HandlingImmediate
	msg.RequestDeliveryReceipt = true
	msg.RequestReadReceipt = true
	msg.Subject = "subject"
	msg.Body = "This is a test.\n"
	var subject, body, err = msg.Encode()
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if subject != "AAA-111_I_subject" {
		t.Errorf("subject: %q", subject)
	}
	if body != "!URG!!RDR!!RRR!This is a test.\n" {
		t.Errorf("body: %q", body)
	}
}

func TestBase64Encode(t *testing.T) {
	var msg pktmsg.TxMessage
	msg.MessageNumber = "AAA-111"
	msg.HandlingOrder = pktmsg.HandlingRoutine
	msg.Subject = "subject"
	msg.Body = "¡Esto es una prueba!\n"
	var subject, body, err = msg.Encode()
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if subject != "AAA-111_R_subject" {
		t.Errorf("subject: %q", subject)
	}
	if body != "!B64!wqFFc3RvIGVzIHVuYSBwcnVlYmEhCg==\n" {
		t.Errorf("body: %q", body)
	}
}

func TestDeliveryReceiptEncode(t *testing.T) {
	var msg pktmsg.TxDeliveryReceipt
	msg.DeliveredSubject = "Test 1"
	msg.DeliveredTo = "A6AAA"
	msg.LocalMessageID = "TUE-2111P"
	msg.DeliveredTime = time.Date(2021, 11, 30, 19, 4, 7, 0, time.Local)
	var subject, body, err = msg.Encode()
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if subject != "DELIVERED: Test 1" {
		t.Errorf("subject: %q", subject)
	}
	if body != "!LMI!TUE-2111P!DR!2021-11-30 19:04:07\nYour Message\nTo: A6AAA\nSubject: Test 1\nwas delivered on 2021-11-30 19:04:07\nRecipient's Local Message ID: TUE-2111P\n" {
		t.Errorf("body: %q", body)
	}
}

func TestReadReceiptEncode(t *testing.T) {
	var msg pktmsg.TxReadReceipt
	msg.ReadSubject = "Test 1"
	msg.ReadTo = "A6AAA"
	msg.ReadTime = time.Date(2021, 12, 16, 10, 38, 10, 0, time.Local)
	var subject, body, err = msg.Encode()
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if subject != "READ: Test 1" {
		t.Errorf("subject: %q", subject)
	}
	if body != "!RR!2021-12-16 10:38:10\nYour Message\n\nTo: A6AAA\nSubject: Test 1\n\nwas read on 2021-12-16 10:38:10\n" {
		t.Errorf("body: %q", body)
	}
}

func TestICS213Encode(t *testing.T) {
	var msg pktmsg.TxICS213Form
	msg.MessageNumber = "AAA-111"
	msg.Severity = pktmsg.SeverityOther
	msg.HandlingOrder = pktmsg.HandlingRoutine
	msg.TakeAction = "No"
	msg.Reply = "Yes"
	msg.ReplyBy = "never"
	msg.FYI = true
	msg.DateTime = time.Date(2021, 12, 18, 12, 34, 0, 0, time.Local)
	msg.ToICSPosition = "topos"
	msg.FromICSPosition = "frompos"
	msg.ToLocation = "toloc"
	msg.FromLocation = "fromloc"
	msg.ToName = "toname"
	msg.FromName = "fromname"
	msg.ToTelephone = "totel"
	msg.FromTelephone = "fromtel"
	msg.Subject = "subject"
	msg.Reference = "ref"
	msg.MessageBody = "message"
	msg.RelayReceived = "rr"
	msg.RelaySent = "rs"
	msg.ReceiverSender = "sender"
	msg.OperatorCallSign = "A6AAA"
	msg.OperatorName = "Fred Flintstone"
	msg.OperatorMethod = "Packet"
	msg.OperatorDateTime = time.Date(2021, 12, 18, 9, 54, 0, 0, time.Local)
	var subject, body, err = msg.Encode()
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if subject != "AAA-111_R_ICS213_subject" {
		t.Errorf("subject: %q", subject)
	}
	if body != "!SCCoPIFO!\n#T: form-ics213.html\n#V: 3.2-2.1\nMsgNo: [AAA-111]\n1a.: [12/18/2021]\n4.: [OTHER]\n5.: [ROUTINE]\n6a.: [No]\n6b.: [Yes]\n6d.: [never]\n6c.: [checked]\n1b.: [12:34]\n7.: [topos]\n8.: [frompos]\n9a.: [toloc]\n9b.: [fromloc]\nToName: [toname]\nFmName: [fromname]\nToTel: [totel]\nFmTel: [fromtel]\n10.: [subject]\n11.: [ref]\n12.: [message]\nOpRelayRcvd: [rr]\nOpRelaySent: [rs]\nRec-Sent: [sender]\nOpCall: [A6AAA]\nOpName: [Fred Flintstone]\nMethod: [Other]\nOther: [Packet]\nOpDate: [12/18/2021]\nOpTime: [09:54]\n!/ADDON!\n" {
		t.Errorf("body: %q", body)
	}
}

func TestEOC213RREncode(t *testing.T) {
	var msg pktmsg.TxEOC213RRForm
	msg.OriginMessageNumber = "AAA-111"
	msg.DateTime = time.Date(2021, 12, 18, 12, 34, 0, 0, time.Local)
	msg.HandlingOrder = pktmsg.HandlingPriority
	msg.ToICSPosition = "topos"
	msg.FromICSPosition = "frompos"
	msg.ToLocation = "toloc"
	msg.FromLocation = "fromloc"
	msg.ToName = "toname"
	msg.FromName = "fromname"
	msg.ToContactInfo = "tocon"
	msg.FromContactInfo = "fromcon"
	msg.RelayReceived = "rr"
	msg.RelaySent = "rs"
	msg.OperatorName = "Fred Flintstone"
	msg.OperatorCallSign = "A6AAA"
	msg.OperatorDateTime = time.Date(2021, 12, 18, 10, 25, 0, 0, time.Local)
	msg.IncidentName = "incident"
	msg.DateTimeInitiated = time.Date(2021, 12, 18, 12, 30, 0, 0, time.Local)
	msg.RequestedBy = "req"
	msg.PreparedBy = "prep"
	msg.ApprovedBy = "app"
	msg.QtyUnit = "qty"
	msg.ResourceDescription = "desc"
	msg.Arrival = "arr"
	msg.Priority = "High"
	msg.EstimatedCost = "cost"
	msg.DeliverTo = "delto"
	msg.DeliverToLocation = "delloc"
	msg.SuggestedSources = "sub"
	msg.RequireEquipmentOperator = true
	msg.RequireLodging = true
	msg.RequireFuel = true
	msg.FuelRequirement = "fuel"
	msg.RequirePower = true
	msg.RequireMeals = true
	msg.RequireMaintenance = true
	msg.RequireWater = true
	msg.RequireOther = true
	msg.SpecialInstructions = "instr"
	var subject, body, err = msg.Encode()
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if subject != "AAA-111_P_EOC213RR_incident" {
		t.Errorf("subject: %q", subject)
	}
	if body != "!SCCoPIFO!\n#T: form-scco-eoc-213rr.html\n#V: 3.2-2.3\nMsgNo: [AAA-111]\n1a.: [12/18/2021]\n1b.: [12:34]\n5.: [PRIORITY]\n7a.: [topos]\n8a.: [frompos]\n7b.: [toloc]\n8b.: [fromloc]\n7c.: [toname]\n8c.: [fromname]\n7d.: [tocon]\n8d.: [fromcon]\n21.: [incident]\n22.: [12/18/2021]\n23.: [12:30]\n25.: [req]\n26.: [prep]\n27.: [app]\n28.: [qty]\n29.: [desc]\n30.: [arr]\n31.: [High]\n32.: [cost]\n33.: [delto]\n34.: [delloc]\n35.: [sub]\n36a.: [checked]\n36b.: [checked]\n36c.: [checked]\n36d.: [fuel]\n36e.: [checked]\n36f.: [checked]\n36g.: [checked]\n36h.: [checked]\n36i.: [checked]\n37.: [instr]\nOpRelayRcvd: [rr]\nOpRelaySent: [rs]\nOpName: [Fred Flintstone]\nOpCall: [A6AAA]\nOpDate: [12/18/2021]\nOpTime: [10:25]\n!/ADDON!\n" {
		t.Errorf("body: %q", body)
	}
}

func TestMuniStatEncode(t *testing.T) {
	var msg pktmsg.TxMuniStatForm
	msg.OriginMessageNumber = "AAA-111"
	msg.DateTime = time.Date(2021, 12, 18, 12, 34, 0, 0, time.Local)
	msg.HandlingOrder = pktmsg.HandlingImmediate
	msg.ToICSPosition = "topos"
	msg.FromICSPosition = "frompos"
	msg.ToLocation = "toloc"
	msg.FromLocation = "fromloc"
	msg.ToName = "toname"
	msg.FromName = "fromname"
	msg.ToContactInfo = "tocon"
	msg.FromContactInfo = "fromcon"
	msg.RelayReceived = "rr"
	msg.RelaySent = "rs"
	msg.OperatorName = "Fred Flintstone"
	msg.OperatorCallSign = "A6AAA"
	msg.OperatorDateTime = time.Date(2021, 12, 18, 10, 51, 0, 0, time.Local)
	msg.ReportType = "Update"
	msg.Jurisdiction = "Sunnyvale"
	msg.EOCPhone = "1111111111"
	msg.EOCFax = "2222222222"
	msg.PrimaryEMContactName = "pecn"
	msg.PrimaryEMContactPhone = "3333333333"
	msg.SecondaryEMContactName = "secn"
	msg.SecondaryEMContactPhone = "4444444444"
	msg.GovtOfficeStatus = "Open"
	msg.GovtOfficeExpectedOpenDateTime = time.Date(2021, 1, 1, 13, 0, 0, 0, time.Local)
	msg.GovtOfficeExpectedCloseDateTime = time.Date(2021, 2, 2, 14, 0, 0, 0, time.Local)
	msg.EOCIsOpen = "Yes"
	msg.EOCActivationLevel = "Normal"
	msg.EOCExpectedOpenDateTime = time.Date(2021, 3, 3, 15, 0, 0, 0, time.Local)
	msg.EOCExpectedCloseDateTime = time.Date(2021, 4, 4, 16, 0, 0, 0, time.Local)
	msg.StateOfEmergency = "Yes"
	msg.StateOfEmergencySent = "pony express"
	msg.CommunicationsStatus = "Unknown"
	msg.CommunicationsComments = "a"
	msg.DebrisStatus = "Normal"
	msg.DebrisComments = "b"
	msg.FloodingStatus = "Problem"
	msg.FloodingComments = "c"
	msg.HazmatStatus = "Failure"
	msg.HazmatComments = "e"
	msg.EmergencyServicesStatus = "Delayed"
	msg.EmergencyServicesComments = "f"
	msg.CasualtiesStatus = "Closed"
	msg.CasualtiesComments = "g"
	msg.GasUtilitiesStatus = "Early Out"
	msg.GasUtilitiesComments = "h"
	msg.ElectricUtilitiesStatus = "Unknown"
	msg.ElectricUtilitiesComments = "i"
	msg.PowerInfraStatus = "Unknown"
	msg.PowerInfraComments = "j"
	msg.WaterInfraStatus = "Unknown"
	msg.WaterInfraComments = "k"
	msg.SewerInfraStatus = "Unknown"
	msg.SewerInfraComments = "l"
	msg.SearchRescueStatus = "Unknown"
	msg.SearchRescueComments = "m"
	msg.TransportRoadStatus = "Unknown"
	msg.TransportRoadComments = "n"
	msg.TransportBridgeStatus = "Unknown"
	msg.TransportBridgeComments = "o"
	msg.CivilUnrestStatus = "Unknown"
	msg.CivilUnrestComments = "p"
	msg.AnimalIssueStatus = "Unknown"
	msg.AnimalIssueComments = "q"
	var subject, body, err = msg.Encode()
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if subject != "AAA-111_I_MuniStat_Sunnyvale" {
		t.Errorf("subject: %q", subject)
	}
	if body != "!URG!!SCCoPIFO!\n#T: form-oa-muni-status.html\n#V: 3.2-2.1\nMsgNo: [AAA-111]\n1a.: [12/18/2021]\n1b.: [12:34]\n5.: [IMMEDIATE]\n7a.: [topos]\n8a.: [frompos]\n7b.: [toloc]\n8b.: [fromloc]\n7c.: [toname]\n8c.: [fromname]\n7d.: [tocon]\n8d.: [fromcon]\n19.: [Update]\n21.: [Sunnyvale]\n23.: [1111111111]\n24.: [2222222222]\n25.: [pecn]\n26.: [3333333333]\n27.: [secn]\n28.: [4444444444]\n29.: [Open]\n30.: [01/01/2021]\n31.: [13:00]\n32.: [02/02/2021]\n33.: [14:00]\n34.: [Yes]\n35.: [Normal]\n36.: [03/03/2021]\n37.: [15:00]\n38.: [04/04/2021]\n39.: [16:00]\n40.: [Yes]\n99.: [pony express]\n41.0.: [Unknown]\n41.1.: [a]\n42.0.: [Normal]\n42.1.: [b]\n43.0.: [Problem]\n43.1.: [c]\n44.0.: [Failure]\n44.1.: [e]\n45.0.: [Delayed]\n45.1.: [f]\n46.0.: [Closed]\n46.1.: [g]\n47.0.: [Early Out]\n47.1.: [h]\n48.0.: [Unknown]\n48.1.: [i]\n49.0.: [Unknown]\n49.1.: [j]\n50.0.: [Unknown]\n50.1.: [k]\n51.0.: [Unknown]\n51.1.: [l]\n52.0.: [Unknown]\n52.1.: [m]\n53.0.: [Unknown]\n53.1.: [n]\n54.0.: [Unknown]\n54.1.: [o]\n55.0.: [Unknown]\n55.1.: [p]\n56.0.: [Unknown]\n56.1.: [q]\nOpRelayRcvd: [rr]\nOpRelaySent: [rs]\nOpName: [Fred Flintstone]\nOpCall: [A6AAA]\nOpDate: [12/18/2021]\nOpTime: [10:51]\n!/ADDON!\n" {
		t.Errorf("body: %q", body)
	}
}
