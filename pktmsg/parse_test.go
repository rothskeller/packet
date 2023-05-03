package pktmsg

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var parseTests = []struct {
	name    string
	raw     string
	wantErr bool
	msg     *Message
}{
	{
		name:    "unparseable",
		raw:     "nothing\n",
		wantErr: true,
	},
	{
		name:    "no plain text",
		raw:     "Content-Type: text/html\n\n<div>nothing</div>",
		wantErr: true,
	},
	{
		name: "sent",
		raw:  "From: <nobody@nowhere>\nTo: <somebody@somewhere>\nSubject: Hello, World\nDate: Wed, 1 Dec 2021 08:04:29 +0000\n\nnothing\n",
		msg: &Message{
			From:          "<nobody@nowhere>",
			To:            "<somebody@somewhere>",
			SentDate:      "Wed, 01 Dec 2021 08:04:29 +0000",
			Subject:       "Hello, World",
			SubjectHeader: "Hello, World",
			Body:          "nothing\n",
		},
	},
	{
		name: "multiple recipients",
		raw:  "From: <nobody@nowhere>\nTo: <somebody@somewhere>\nCc: <number2@somewhere>\nBcc: <number3@somewhere>\nSubject: Hello, World\nDate: Wed, 1 Dec 2021 08:04:29 +0000\n\nnothing\n",
		msg: &Message{
			From:          "<nobody@nowhere>",
			To:            "<somebody@somewhere>, <number2@somewhere>, <number3@somewhere>",
			SentDate:      "Wed, 01 Dec 2021 08:04:29 +0000",
			Subject:       "Hello, World",
			SubjectHeader: "Hello, World",
			Body:          "nothing\n",
		},
	},
	{
		name: "XSC subject",
		raw:  "From: <nobody@nowhere>\nTo: <somebody@somewhere>\nSubject: AAA-111P_R_Hello, World\nDate: Wed, 1 Dec 2021 08:04:29 +0000\n\nnothing\n",
		msg: &Message{
			From:          "<nobody@nowhere>",
			To:            "<somebody@somewhere>",
			SentDate:      "Wed, 01 Dec 2021 08:04:29 +0000",
			Subject:       "Hello, World",
			SubjectHeader: "AAA-111P_R_Hello, World",
			Body:          "nothing\n",
			OriginMsgID:   "AAA-111P",
			Handling:      "ROUTINE",
		},
	},
	{
		name: "XSC subject with severity",
		raw:  "From: <nobody@nowhere>\nTo: <somebody@somewhere>\nSubject: AAA-111P_O/R_Hello, World\nDate: Wed, 1 Dec 2021 08:04:29 +0000\n\nnothing\n",
		msg: &Message{
			From:          "<nobody@nowhere>",
			To:            "<somebody@somewhere>",
			SentDate:      "Wed, 01 Dec 2021 08:04:29 +0000",
			Subject:       "Hello, World",
			SubjectHeader: "AAA-111P_O/R_Hello, World",
			Body:          "nothing\n",
			OriginMsgID:   "AAA-111P",
			Handling:      "ROUTINE",
			Severity:      "OTHER",
		},
	},
	{
		name: "received",
		raw:  "Received: FROM bbs.ampr.org BY pktmsg.local FOR area; Wed, 01 Dec 2021 08:04:29 +0000\nFrom: <nobody@nowhere>\nTo: <somebody@somewhere>\nSubject: Hello, World\nDate: Wed, 1 Dec 2021 08:04:29 +0000\n\nnothing\n",
		msg: &Message{
			From:          "<nobody@nowhere>",
			To:            "<somebody@somewhere>",
			SentDate:      "Wed, 01 Dec 2021 08:04:29 +0000",
			Subject:       "Hello, World",
			SubjectHeader: "Hello, World",
			Body:          "nothing\n",
			RxBBS:         "bbs",
			RxArea:        "area",
			RxDate:        "Wed, 01 Dec 2021 08:04:29 +0000",
		},
	},
	{
		name:    "received with bad header",
		raw:     "Received: blah\nFrom: <nobody@nowhere>\nTo: <somebody@somewhere>\nSubject: Hello, World\nDate: Wed, 1 Dec 2021 08:04:29 +0000\n\nnothing\n",
		wantErr: true,
	},
	{
		name: "Outpost flags",
		raw:  "From: <nobody@nowhere>\n\n!RDR!!RRR!!URG!nothing\n",
		msg: &Message{
			From:                   "<nobody@nowhere>",
			Body:                   "nothing\n",
			OutpostUrgent:          true,
			RequestDeliveryReceipt: true,
			RequestReadReceipt:     true,
		},
	},
	{
		name: "base64 (Outpost)",
		raw:  "From: <nobody@nowhere>\n\n\n!B64!bm90aGluZwo=\n",
		msg: &Message{
			From: "<nobody@nowhere>",
			Body: "nothing\n",
		},
	},
	{
		name: "minimal valid form",
		raw:  "From: <nobody@nowhere>\n\n!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [x]\n!/ADDON!\n",
		msg: &Message{
			From:         "<nobody@nowhere>",
			PIFOVersion:  "1",
			FormHTML:     "tt.html",
			FormVersion:  "2",
			TaggedFields: []TaggedField{{"A", "x"}},
		},
	},
	{
		name: "form with XSC subject",
		raw:  "From: <nobody@nowhere>\nSubject: AAA-111P_R_FTag_Subject\n\n!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [x]\n!/ADDON!\n",
		msg: &Message{
			From:          "<nobody@nowhere>",
			PIFOVersion:   "1",
			FormHTML:      "tt.html",
			FormVersion:   "2",
			SubjectHeader: "AAA-111P_R_FTag_Subject",
			Subject:       "Subject",
			OriginMsgID:   "AAA-111P",
			Handling:      "ROUTINE",
			FormTag:       "FTag",
			TaggedFields:  []TaggedField{{"A", "x"}},
		},
	},
	{
		name: "form with stuff before it",
		raw:  "From: <nobody@nowhere>\n\nHello, world!\n!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [x]\n!/ADDON!\n",
		msg: &Message{
			From:         "<nobody@nowhere>",
			PIFOVersion:  "1",
			FormHTML:     "tt.html",
			FormVersion:  "2",
			TaggedFields: []TaggedField{{"A", "x"}},
		},
	},
	{
		name: "form with stuff after it",
		raw:  "From: <nobody@nowhere>\n\n!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [x]\n!/ADDON!\nGoodbye, cruel world!\n",
		msg: &Message{
			From:         "<nobody@nowhere>",
			PIFOVersion:  "1",
			FormHTML:     "tt.html",
			FormVersion:  "2",
			TaggedFields: []TaggedField{{"A", "x"}},
		},
	},
	{
		name: "missing header",
		raw:  "From: <nobody@nowhere>\n\n#T: tt.html\n#V: 1-2\nA: [x]\n!/ADDON!\n",
		msg: &Message{
			From: "<nobody@nowhere>",
			Body: "#T: tt.html\n#V: 1-2\nA: [x]\n!/ADDON!\n",
		},
	},
	{
		name: "missing type",
		raw:  "From: <nobody@nowhere>\n\n!SCCoPIFO!\n#V: 1-2\nA: [x]\n!/ADDON!\n",
		msg: &Message{
			From: "<nobody@nowhere>",
			Body: "!SCCoPIFO!\n#V: 1-2\nA: [x]\n!/ADDON!\n",
		},
	},
	{
		name: "invalid type",
		raw:  "From: <nobody@nowhere>\n\n!SCCoPIFO!\n#T: t\n#V: 1-2\nA: [x]\n!/ADDON!\n",
		msg: &Message{
			From: "<nobody@nowhere>",
			Body: "!SCCoPIFO!\n#T: t\n#V: 1-2\nA: [x]\n!/ADDON!\n",
		},
	},
	{
		name: "missing version",
		raw:  "From: <nobody@nowhere>\n\n!SCCoPIFO!\n#T: tt.html\n!/ADDON!\n",
		msg: &Message{
			From: "<nobody@nowhere>",
			Body: "!SCCoPIFO!\n#T: tt.html\n!/ADDON!\n",
		},
	},
	{
		name: "invalid version",
		raw:  "From: <nobody@nowhere>\n\n!SCCoPIFO!\n#T: tt.html\n#V: X\nA: [x]\n!/ADDON!\n",
		msg: &Message{
			From: "<nobody@nowhere>",
			Body: "!SCCoPIFO!\n#T: tt.html\n#V: X\nA: [x]\n!/ADDON!\n",
		},
	},
	{
		name: "invalid field",
		raw:  "From: <nobody@nowhere>\n\n!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA\n!/ADDON!\n",
		msg: &Message{
			From: "<nobody@nowhere>",
			Body: "!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA\n!/ADDON!\n",
		},
	},
	{
		name: "missing footer",
		raw:  "From: <nobody@nowhere>\n\n!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [x]\n",
		msg: &Message{
			From: "<nobody@nowhere>",
			Body: "!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [x]\n",
		},
	},
	{
		name: "bracket quoting",
		raw:  "From: <nobody@nowhere>\n\n!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [nl\\nbs\\\\rb`]et`]]]\n!/ADDON!\n",
		msg: &Message{
			From:         "<nobody@nowhere>",
			PIFOVersion:  "1",
			FormHTML:     "tt.html",
			FormVersion:  "2",
			TaggedFields: []TaggedField{{"A", "nl\nbs\\rb]et`"}},
		},
	},
	{
		name: "end of input inside brackets",
		raw:  "From: <nobody@nowhere>\n\n!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [x",
		msg: &Message{
			From: "<nobody@nowhere>",
			Body: "!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [x",
		},
	},
	{
		name: "extra stuff after brackets",
		raw:  "From: <nobody@nowhere>\n\n!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [x] \n!/ADDON!\n",
		msg: &Message{
			From: "<nobody@nowhere>",
			Body: "!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [x] \n!/ADDON!\n",
		},
	},
	{
		name: "line continuation",
		raw:  "From: <nobody@nowhere>\n\n!SCCoPIFO!\n#T: tt.html\n#V: 1-2\nA: [this is \na test]\n!/ADDON!\n",
		msg: &Message{
			From:         "<nobody@nowhere>",
			PIFOVersion:  "1",
			FormHTML:     "tt.html",
			FormVersion:  "2",
			TaggedFields: []TaggedField{{"A", "this is a test"}},
		},
	},
}

func TestParse(t *testing.T) {
	for _, tt := range parseTests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := ParseMessage(tt.raw)
			if err != nil && !tt.wantErr {
				t.Fatalf("unexpected error %s", err)
			}
			if err == nil && tt.wantErr {
				t.Fatal("unexpected success")
			}
			if !reflect.DeepEqual(m, tt.msg) {
				spew.Fdump(os.Stderr, "actual", m)
				spew.Fdump(os.Stderr, "expected", tt.msg)
				t.Fatal("incorrect result")
			}
		})
	}
}

var receiveTests = []struct {
	name    string
	raw     string
	bbs     string
	area    string
	wantErr bool
	msg     *Message
}{
	{
		name:    "unparseable",
		raw:     "nothing\n",
		bbs:     "bbs",
		wantErr: true,
	},
	{
		name:    "bounce",
		raw:     actualBounceMessage,
		bbs:     "bbs",
		wantErr: false,
		msg: &Message{

			From:          "Microsoft Outlook <MicrosoftExchange329e71ec88ae4615bbc36ab6ce41109e@cityofsunnyvale.onmicrosoft.com>",
			To:            "<cert@sunnyvale.ca.gov>",
			Subject:       "Undeliverable: SERV Volunteer Hours for November 2021",
			SubjectHeader: "Undeliverable: SERV Volunteer Hours for November 2021",
			SentDate:      "Wed, 01 Dec 2021 08:04:29 +0000",
			Body:          expectedBounceBody,
			RxBBS:         "bbs",
			RxDate:        "Sun, 01 Jan 2023 00:00:00 -0800",
			Autoresponse:  true,
			BBSRxDate:     "Wed, 01 Dec 2021 08:04:29 -0800",
			NotPlainText:  true,
		},
	},
	{
		name:    "no plain text",
		raw:     "Content-Type: text/html\n\n<div>nothing</div>",
		bbs:     "bbs",
		wantErr: true,
		msg: &Message{

			RxBBS:  "bbs",
			RxDate: "Sun, 01 Jan 2023 00:00:00 -0800",

			NotPlainText: true,
		},
	},
	{
		name: "quoted-printable",
		raw:  "Content-Transfer-Encoding: quoted-printable\n\nn=6fthing\n",
		bbs:  "bbs",
		msg: &Message{

			Body: "nothing\n",

			RxBBS:  "bbs",
			RxDate: "Sun, 01 Jan 2023 00:00:00 -0800",

			NotPlainText: true,
		},
	},
	{
		name: "base64 (MIME)",
		raw:  "Content-Transfer-Encoding: base64\n\nbm90aGluZwo=\n",
		bbs:  "bbs",
		msg: &Message{

			Body: "nothing\n",

			RxBBS:  "bbs",
			RxDate: "Sun, 01 Jan 2023 00:00:00 -0800",

			NotPlainText: true,
		},
	},
	{
		name: "unparseable content type",
		raw:  "Content-Type: //bogus\n\nnothing\n",
		bbs:  "bbs",
		msg: &Message{

			RxBBS:  "bbs",
			RxDate: "Sun, 01 Jan 2023 00:00:00 -0800",
		},
		wantErr: true,
	},
	{
		name: "non-plain-text content type",
		raw:  "Content-Type: text/html\n\nnothing\n",
		bbs:  "bbs",
		msg: &Message{

			RxBBS:  "bbs",
			RxDate: "Sun, 01 Jan 2023 00:00:00 -0800",

			NotPlainText: true,
		},
		wantErr: true,
	},
	{
		name: "multipart with plain text",
		raw:  "Content-Type: multipart/alternative; boundary=\"X\"\n\n\n--X\nContent-Type: text/plain\n\nnothing\n\n--X--\n",
		bbs:  "bbs",
		msg: &Message{

			Body: "nothing\n",

			RxBBS:  "bbs",
			RxDate: "Sun, 01 Jan 2023 00:00:00 -0800",

			NotPlainText: true,
		},
	},
	{
		name: "nested multipart with plain text",
		raw:  "Content-Type: multipart/mixed; boundary=\"Y\"\n\n\n--Y\nContent-Type: multipart/alternative; boundary=\"X\"\n\n\n--X\nContent-Type: text/plain\n\nnothing\n\n--X--\n\n--Y--\n",
		bbs:  "bbs",
		msg: &Message{

			Body: "nothing\n",

			RxBBS:  "bbs",
			RxDate: "Sun, 01 Jan 2023 00:00:00 -0800",

			NotPlainText: true,
		},
	},
	{
		name: "nested multipart ill-formed",
		raw:  "Content-Type: multipart/mixed; boundary=\"Y\"\n\n\n--Y\nContent-Type: multipart/alternative; boundary=\"X\"\n\n\n--X\nContent-Type: text/plain\n\nnothing\n\n--Y--\n",
		bbs:  "bbs",
		msg: &Message{

			RxBBS:  "bbs",
			RxDate: "Sun, 01 Jan 2023 00:00:00 -0800",
		},
		wantErr: true,
	},
	{
		name: "multipart with no plain text",
		raw:  "Content-Type: multipart/alternative; boundary=\"X\"\n\n\n--X\nContent-Type: text/html\n\n<div>nothing</div>\n\n--X--\n",
		bbs:  "bbs",
		msg: &Message{
			RxBBS:        "bbs",
			RxDate:       "Sun, 01 Jan 2023 00:00:00 -0800",
			NotPlainText: true,
		},
		wantErr: true,
	},
	{
		name: "envelope line",
		raw:  "From nobody@nowhere Wed Dec  1 08:04:29 2021\n\nnothing\n",
		bbs:  "bbs",
		msg: &Message{
			Body:       "nothing\n",
			RxBBS:      "bbs",
			RxDate:     "Sun, 01 Jan 2023 00:00:00 -0800",
			ReturnAddr: "nobody@nowhere",
			BBSRxDate:  "Wed, 01 Dec 2021 08:04:29 -0800",
		},
	},
}

func TestReceive(t *testing.T) {
	now = func() time.Time { return time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local) }
	for _, tt := range receiveTests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := ReceiveMessage(tt.raw, tt.bbs, tt.area)
			if err != nil && !tt.wantErr {
				t.Fatalf("unexpected error %s", err)
			}
			if err == nil && tt.wantErr {
				t.Fatal("unexpected success")
			}
			if !reflect.DeepEqual(m, tt.msg) {
				spew.Fdump(os.Stderr, "actual", m)
				spew.Fdump(os.Stderr, "expected", tt.msg)
				t.Fatal("incorrect result")
			}
		})
	}
}

const actualBounceMessage = `From  Wed Dec 01 08:04:29 2021
Received: from BY5PR09MB5379.namprd09.prod.outlook.com (2603:10b6:a03:24d::11)
 by BY5PR09MB5090.namprd09.prod.outlook.com with HTTPS; Wed, 1 Dec 2021
 08:04:30 +0000
MIME-Version: 1.0
From: Microsoft Outlook
	<MicrosoftExchange329e71ec88ae4615bbc36ab6ce41109e@cityofsunnyvale.onmicrosoft.com>
To: <cert@sunnyvale.ca.gov>
Date: Wed, 1 Dec 2021 08:04:29 +0000
X-MS-Exchange-Organization-SCL: -1
X-MS-Exchange-Message-Is-Ndr:
Content-Language: en-US
Message-ID:
 <1b8c300e-e3e3-4ed0-8102-c2d5e64f0c39@BY5PR09MB5379.namprd09.prod.outlook.com>
In-Reply-To:
 <BY5PR09MB5090C0A1798249A18D3F34A3E7689@BY5PR09MB5090.namprd09.prod.outlook.com>
References:
 <BY5PR09MB5090C0A1798249A18D3F34A3E7689@BY5PR09MB5090.namprd09.prod.outlook.com>
Subject: Undeliverable: SERV Volunteer Hours for November 2021
Auto-Submitted: auto-replied
X-MS-PublicTrafficType: Email
X-MS-Exchange-Organization-AuthSource: BY5PR09MB5379.namprd09.prod.outlook.com
X-MS-Exchange-Organization-AuthAs: Internal
X-MS-Exchange-Organization-AuthMechanism: 05
X-MS-Exchange-Organization-Network-Message-Id:
 1a4b7726-319e-440f-f657-08d9b4a13380
X-MS-TrafficTypeDiagnostic: BY5PR09MB5379:
X-MS-Exchange-Organization-ExpirationStartTime: 01 Dec 2021 08:04:29.6410
 (UTC)
X-MS-Exchange-Organization-ExpirationStartTimeReason: SideEffectMessage
X-MS-Exchange-Organization-ExpirationInterval: 1:00:00:00.0000000
X-MS-Exchange-Organization-ExpirationIntervalReason: SideEffectMessage
X-MS-Oob-TLC-OOBClassifiers: OLM:849;
X-Microsoft-Antispam: BCL:0;
X-Forefront-Antispam-Report:
 CIP:255.255.255.255;CTRY:;LANG:en;SCL:-1;SRV:;IPV:NLI;SFV:SKI;H:;PTR:;CAT:NONE;SFS:;DIR:INB;
X-MS-Exchange-CrossTenant-OriginalArrivalTime: 01 Dec 2021 08:04:29.6439
 (UTC)
X-MS-Exchange-CrossTenant-FromEntityHeader: Hosted
X-MS-Exchange-CrossTenant-AuthSource: BY5PR09MB5379.namprd09.prod.outlook.com
X-MS-Exchange-CrossTenant-AuthAs: Internal
X-MS-Exchange-CrossTenant-Network-Message-Id:
 1a4b7726-319e-440f-f657-08d9b4a13380
X-MS-Exchange-Transport-CrossTenantHeadersStamped: BY5PR09MB5379
X-MS-Exchange-Organization-MessageDirectionality: Originating
X-MS-Exchange-Transport-EndToEndLatency: 00:00:00.5040709
X-MS-Exchange-Processed-By-BccFoldering: 15.20.4755.014
X-Microsoft-Antispam-Mailbox-Delivery:
	ucf:0;jmr:0;auth:0;dest:I;ENG:(910001)(944506458)(944626604)(4710097)(4711095)(920097)(425001)(930097);
X-Microsoft-Antispam-Message-Info:
	wuunwctO6u4e1IwT3boFqQAQ+O6n62lWv1aQgPi8u8IQgqfoEVm/uHp+Gmk1LpyRb1hnKgw+i77ySs6Ge3VtFK+lc4BcDYG82y+dBTmbCMC7fzuu8eJlMQtB5EvVtseaxcQzCusXuv4LOGGoNZJPtrHP0rJu5AD7vByvSjtbhEbm0tMeQ7akRiVwKnyoKHi7RfYhbD1yRpT3Ag4Y1Cj1/MlYYMp/MjtazOuTahLDBX8UN3Z0edRGm7sxjQ1qjfZAdBT2P2gR452TK/oeKlN0Rb6X4+uApGzaiuZF2Tf+tJBZBvLhAebzAPqk+swZz4LI38nsQsFlK0bhCgPCkVNuxWisL/wZMdrp6183JqRJuPxuOwvXfDHJRn25xVJdI3BccaKrAbRhyRu4g0WagVrhIY3KTk8xfoibk07tQ8the0ZSMfrVE3WNoH2m0QT2VW6DaRbTKLfTguq9K8anwz2Y68ZEwthIuPcTaZ1VJfwZq8ftyyQoKzivrYwOnx7NuG0E0QQiIDS809T5Mb0bajvdYOOGvHKA1KyXWh9VUuJYerBsZ4Jo7M/O2unh+Sa3B0wHqiS1apbXuyWlwvaYf5VqhNHHS/9zlbRXS15UPRJ1pheLg6Oo1pREgDeU6bV9ZAhnRqeYop1J581bRP+BHGez3wX6ocLhtbUN6cktcGvps2TYFvOKfUKE1tet5jcQeGABvkhCF6YjXSQ7f81H6QunWSTKmFWG4M4p3JcUaO4o2oLheiK5mpunWtXMIp4yeD4HqKlfVZhBv41AmNbfzBq9ng==
Content-type: multipart/mixed;
	boundary="B_3721553817_1320843075"

> This message is in MIME format. Since your mail reader does not understand
this format, some or all of this message may not be legible.

--B_3721553817_1320843075
Content-type: multipart/alternative;
	boundary="B_3721553817_622503736"


--B_3721553817_622503736
Content-type: text/plain;
	charset="UTF-8"
Content-transfer-encoding: 7bit


Your message to gordon@amateur-radio.org couldn't be delivered.
amateur-radio.org suspects your message is spam and rejected it.
cert Office 365 amateur-radio.org
Sender Action Required
Messages suspected as spam

How to Fix It
Try to modify your message, or change how you're sending the message, using the guidance in this article: Bulk E-mailing Best Practices for Senders Using Forefront Online Protection for Exchange. Then resend your message.
If you continue to experience the problem, contact the recipient by some other means (by phone, for example) and ask them to ask their email admin to add your email address, or your domain name, to their allowed senders list.

Was this helpful? Send feedback to Microsoft.

More Info for Email Admins
Status code: 550 5.7.350

When Office 365 tried to send the message to the recipient (outside Office 365), the recipient's email server (or email filtering service) suspected the sender's message is spam.

If the sender can't fix the problem by modifying their message, contact the recipient's email admin and ask them to add your domain name, or the sender's email address, to their list of allowed senders.

Although the sender may be able to alter the message contents to fix this issue, it's likely that only the recipient's email admin can fix this problem. Unfortunately, Office 365 Support is unlikely to be able to help fix these kinds of externally reported errors.

Original Message Details
Created Date:12/1/2021 8:03:44 AM
Sender Address:cert@sunnyvale.ca.gov
Recipient Address:gordon@amateur-radio.org
Subject:SERV Volunteer Hours for November 2021

Error Details
Reported error:550 5.7.350 Remote server returned message detected as spam -> 554 5.7.1 The message from (<cert@sunnyvale.ca.gov>) with the subject of (SERV Volunteer Hours for November 2021) matches a profile the Internet community may consider spam. Please revise your message before resending.
DSN generated by:BY5PR09MB5379.namprd09.prod.outlook.com
Message Hops
HOPTIME (UTC)FROMTOWITHRELAY TIME
112/1/2021
8:04:26 AMlocalhostBL1P223CA0025.NAMP223.PROD.OUTLOOK.COMMicrosoft SMTP Server (version=TLS1_2, cipher=TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384)42 sec
212/1/2021
8:04:27 AMBY5PR09MB5090.namprd09.prod.outlook.comBY5PR09MB5090.namprd09.prod.outlook.commapi1 sec
312/1/2021
8:04:27 AMBY5PR09MB5090.namprd09.prod.outlook.comBY5PR09MB5379.namprd09.prod.outlook.comMicrosoft SMTP Server (version=TLS1_2, cipher=TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384)*

Original Message Headers
ARC-Seal: i=1; a=rsa-sha256; s=arcselector9901; d=microsoft.com; cv=none;
 b=ElDXLNsNor1ieyeILBeH5YK1xaNSUH9k/3DrvN57MlvdD47bwqMqyrCF/aq579lWFnXnV4lxQ7CqmLgzj1JOhLfjGbkhwg64xDwCufjPj19znP1pthkPrbC3ASDOMSM1uVY5jjtdRGBz2KHnJ5gGyn9muTvqqpAP9xNLe4eJWX5Q5k5XsbLCimuLBJoFEF4aTJVB+6WD1wV0cBgwCWiG83NER5cYrWFvm9E0tZAW2ZJwZ7XMs7lrAHG6B8f5aroV0/1wOqLAkbS/31mZH2naMxhL3XsjX90KDUePESyTnclD3bj5jiX62Z6sE1E/DilGz+IBenWmvAoXpD6/Q2TSsg==
ARC-Message-Signature: i=1; a=rsa-sha256; c=relaxed/relaxed; d=microsoft.com;
 s=arcselector9901;
 h=From:Date:Subject:Message-ID:Content-Type:MIME-Version:X-MS-Exchange-AntiSpam-MessageData-ChunkCount:X-MS-Exchange-AntiSpam-MessageData-0:X-MS-Exchange-AntiSpam-MessageData-1;
 bh=xiMYKUzhNVXMjKgeNC0NBpIE1l7F6ew8A35kP50faOY=;
 b=S1+82Rkj0k6Twql7XKI8WsGwxXC/WQHMmqCvfUm+Z80KwmTAMOExeMfg62vES1aToAHjE+Y0eXmBxQBKfo5crXkU5p8/AJDoZPcztwniN4AjJVOS6jAuqsoZvdDwPFa7J/r0L9VoJzNwekS6PXT+YtelGaysf+iu8xwG8mlyZFkx2thBYz5mYUX6qaFaqGbW3oIVPnCr0V/4SydMn6Dqu4BUWzKlvuWcFUkfETyMAhd2gNXj4aDL+3JPmOu/bncFce19PnSBWbvgkNAiZIU3ypCd41v1NgEgyWeNjGeLrys0SRV42QIjTtcjwTf1/AE5dr9FFvK0t8ikjC74uhb57w==
ARC-Authentication-Results: i=1; mx.microsoft.com 1; spf=pass
 smtp.mailfrom=sunnyvale.ca.gov; dmarc=pass action=none
 header.from=sunnyvale.ca.gov; dkim=pass header.d=sunnyvale.ca.gov; arc=none
DKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;
 d=cityofsunnyvale.onmicrosoft.com;
 s=selector2-cityofsunnyvale-onmicrosoft-com;
 h=From:Date:Subject:Message-ID:Content-Type:MIME-Version:X-MS-Exchange-SenderADCheck;
 bh=xiMYKUzhNVXMjKgeNC0NBpIE1l7F6ew8A35kP50faOY=;
 b=0q2m3qItWh5cXQJd+TvjJfR/rqd20b2kN3+hNbTiY57BtFfo6FJkyu2bbOrlgWZR7pe9/rojV2txhcXOg2ulr3nGxetRJ6EhgJF1C2ekRiz0IWE3KGTMOk40K8kZmmGXo3XVJ3bRw4+y2XpB59ZLSWfwsoFed+FmOPOExiE7lx4=
Authentication-Results: dkim=none (message not signed)
 header.d=none;dmarc=none action=none header.from=sunnyvale.ca.gov;
Received: from BY5PR09MB5090.namprd09.prod.outlook.com (2603:10b6:a03:24a::14)
 by BY5PR09MB5379.namprd09.prod.outlook.com (2603:10b6:a03:24d::11) with
 Microsoft SMTP Server (version=TLS1_2,
 cipher=TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384) id 15.20.4734.22; Wed, 1 Dec
 2021 08:04:27 +0000
Received: from BY5PR09MB5090.namprd09.prod.outlook.com
 ([fe80::3ca3:b9c4:5e3d:58f6]) by BY5PR09MB5090.namprd09.prod.outlook.com
 ([fe80::3ca3:b9c4:5e3d:58f6%3]) with mapi id 15.20.4755.014; Wed, 1 Dec 2021
 08:04:27 +0000
From: SunnyvaleSERV.org <cert@sunnyvale.ca.gov>
Content-Type: multipart/alternative; boundary="BOUNDARY"
Subject: SERV Volunteer Hours for November 2021
Date: Wed, 01 Dec 2021 00:03:44 -0800
To: "Gordon Girton" <gordon@amateur-radio.org>
X-ClientProxiedBy: BL1P223CA0025.NAMP223.PROD.OUTLOOK.COM
 (2603:10b6:208:2c4::30) To BY5PR09MB5090.namprd09.prod.outlook.com
 (2603:10b6:a03:24a::14)
Return-Path: cert@sunnyvale.ca.gov
Message-ID: <BY5PR09MB5090C0A1798249A18D3F34A3E7689@BY5PR09MB5090.namprd09.prod.outlook.com>
MIME-Version: 1.0
Received: from localhost (2607:f298:5:100f::320:4a05) by BL1P223CA0025.NAMP223.PROD.OUTLOOK.COM (2603:10b6:208:2c4::30) with Microsoft SMTP Server (version=TLS1_2, cipher=TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384) id 15.20.4734.23 via Frontend Transport; Wed, 1 Dec 2021 08:04:26 +0000
X-MS-PublicTrafficType: Email
X-MS-Office365-Filtering-Correlation-Id: d270b94b-ec7a-41e9-6661-08d9b4a131ee
X-MS-TrafficTypeDiagnostic: BY5PR09MB5379:
X-Microsoft-Antispam-PRVS:
	<BY5PR09MB53799E7833AAA8B89EC78C75E7689@BY5PR09MB5379.namprd09.prod.outlook.com>
X-MS-Oob-TLC-OOBClassifiers: OLM:6790;
X-MS-Exchange-SenderADCheck: 1
X-MS-Exchange-AntiSpam-Relay: 0
X-Microsoft-Antispam: BCL:0;
X-Microsoft-Antispam-Message-Info:
	xDd6uy4T4vl0VR3tODu+h8hmVT8n/RR0LsMnRjE1JkCqQjNbDvZCW//Yc74Yu0AGTolRxrk87iiNr60ari8xnFHroywalti12y5ZgKg+99C3P5FN4LjcX2YIdFqlOllMqziDoyTOu9kA4D8H00nCY/RixXb1oBbssQLuHNuem+OknTSeb3WVrBYtpXrTyJMmQEeuOd5wHAyCR3Wde5/mGfWMLsqLc9Ii2cXyM0NeVZSfP+zmeq5TPi9f5cSLQEpjgvUbiP4HST9O1zkebBU+uXf9d3cvGR17iiLkSz/0V+4wU+ecqdwoY1WiHWdWNIaYQcomQxGD5BE8z8Atyg8IIJI1K0QN+gzaNfkA8joWGMRnpmhptHyFBYd0BLd+ZZoqF+jG5C+zwKhkLMQOij9+CZf9/pTgYjmdybvOO1ICMSjkPzkr1fDoCRmnDNEb9J9vOIsHXs1ijpRN4LAKO43yVmZk/8tWPHsIFIGjAY3PNuhK12AR9abkgWYToFyapyLwHlOytdbXedZdfrV/kPRGk4X6Q5bdqjtX6BjIlqgyjXAXkFFlE947CMNhn5vNEIv6slS4PwwDqP1PuyYBTFGUZN5LO/hOw7P457JTrHICUlPCsixRzLB4kdppg+tKToIQkDKTqR5+VM5bWH60fBHseEpFg+WVb/wEk2gBm0g3gqSOStaJcXx4lpO5lf5F6fcX6NG2Nh6iixZPUD/c2CxIWUeAsUc9aAtouNLNlZFC7N8=
X-Forefront-Antispam-Report:
	CIP:255.255.255.255;CTRY:;LANG:en;SCL:1;SRV:;IPV:NLI;SFV:NSPM;H:BY5PR09MB5090.namprd09.prod.outlook.com;PTR:;CAT:NONE;SFS:(366004)(186003)(4001150100001)(8936002)(38100700002)(5660300002)(83380400001)(6666004)(6916009)(8676002)(508600001)(316002)(9686003)(966005)(52536014)(6496006)(6486002)(2906002)(66476007)(52116002)(33964004)(166002)(33656002)(66556008)(66946007)(86362001);DIR:OUT;SFP:1102;
X-MS-Exchange-AntiSpam-MessageData-ChunkCount: 1
X-MS-Exchange-AntiSpam-MessageData-0:
	=?us-ascii?Q?NbKYDOktWeK55k5Jod3Eh0cXtj/BpM9uyUI3AqGbM/X35OnV/mzFdDopk3xu?=
 =?us-ascii?Q?kckxY0rMyxo87Vfh1HQAObmFYo0rpOMiQkgQkpCfkbIjtOfTz5djTHpd3c8V?=
 =?us-ascii?Q?Whoy/b13NKtH4ksh9KUtGre5sl5y35BThCsJISs/HZKc45JlajC9dFLbXmJk?=
 =?us-ascii?Q?B17SJqoK8aWoXp0sjsX8PwI1SDUBlvQFRrU39yTtf2ogCnp+J9BIZjudmrLo?=
 =?us-ascii?Q?SxtPR4pY9HcMnQWvgW89Gvjs0a5S+OmUuv2dJXWQKF1QvMhePWk/NVw1y83T?=
 =?us-ascii?Q?AT8+/zRPdpv/sHJxEQOWdjyVnjQwG5LXsoPDRoIkR4PQscRcYCtfDzldkR0R?=
 =?us-ascii?Q?K5k41wVViL3Bj2Yt9rRr4f4f3ZY3EfptXUF79hCYwurkDasekmtZEe+W6517?=
 =?us-ascii?Q?uP9qrAGBpBoIfReNdG5qSLVme1vFBziAtB866b4iHBpQJyhxifj9FxfeWuWH?=
 =?us-ascii?Q?V75MfO5G2R0ImZnCXzCD8++mXMKtClF0zeEnTlAxzWg8EJi/nlzXWa1niETU?=
 =?us-ascii?Q?G/eQg//FHfAX9jlVcDMdYlrDi+Hjj1cDsNBVzE4nu7ij20SQfkg3K6odEac9?=
 =?us-ascii?Q?4zi3U03GaEFR7CCR+tI/acVFgAaS7TYSGpNHVlthjhsVPlfwM/H0xUSVEact?=
 =?us-ascii?Q?+lK4SFy6Zc75cdQIjbwLr3x8/A2CDlLxEnctCReahtJCv02RjLAlAtMFkd7j?=
 =?us-ascii?Q?DzXy2sDhOYN3sbXnBGMdUD/F15QG4Rg+YWZXJwZXgPsPLl0eClIe1ygeFfJ8?=
 =?us-ascii?Q?xEB/OdDRbIeGvw0zLW0g1IAMhZCEqQax8BEQtxwU2KBsQVxU/k+tDjiRIgcK?=
 =?us-ascii?Q?aO1UV5M0D9k4nzoIgXWFnbkQFqDdtLqrzU5Vnipgh72VnpxufkOfM2xDIs67?=
 =?us-ascii?Q?Y8ShpLfbWElbdwc0IzxJ3WfS5/CFiYkN4QMHWUzV9reC0F8Ov7uu82alLvLD?=
 =?us-ascii?Q?02+9OVkXZiw4h6TO3czbyj239WhuQoCHZivGvbsU31N7JVuLTSBzBKcfdgWl?=
 =?us-ascii?Q?p7WmccNFOKP/GY4vLaDev9KUUCqGy3DHFqCp/PeJaWNbZlQFyVBgELPmbZbc?=
 =?us-ascii?Q?7mL3iCrFrLF65oYg/6QfLaUregT5f5JENi3S+tKUQitQ1BoE9f5KpbqDDf//?=
 =?us-ascii?Q?PPVdLs1JDSmdCJyIHkmr0I7fpl/PYPf/hRbKeVjxELuHY9fOjpEZQkS/ZDOM?=
 =?us-ascii?Q?g7fE+ebolUFU9dpsPFM3tu7FtO7arSmfB4oyFx7Up3FqrvDA+EnnWZ8Gmyd1?=
 =?us-ascii?Q?CZXfWBBgKTazj/u1xrMw2sJFB/Xv+FnPdnVdz+yS/CRwgdoMbHKQTg60R1UH?=
 =?us-ascii?Q?VHvZ7ldvkZVh1138NhcahYslKeP33e53t7dVkslDbqONpEAbtTMkERgqm+tW?=
 =?us-ascii?Q?Ex9R+pw=3D?=
X-OriginatorOrg: sunnyvale.ca.gov
X-MS-Exchange-CrossTenant-Network-Message-Id: d270b94b-ec7a-41e9-6661-08d9b4a131ee
X-MS-Exchange-CrossTenant-AuthSource: BY5PR09MB5090.namprd09.prod.outlook.com
X-MS-Exchange-CrossTenant-AuthAs: Internal
X-MS-Exchange-CrossTenant-OriginalArrivalTime: 01 Dec 2021 08:04:27.0714
 (UTC)
X-MS-Exchange-CrossTenant-FromEntityHeader: Hosted
X-MS-Exchange-CrossTenant-Id: 63dc83d7-8dcb-489b-bce2-0bf7c37d8f39
X-MS-Exchange-Transport-CrossTenantHeadersStamped: BY5PR09MB5379


--B_3721553817_622503736
Content-type: text/html;
	charset="UTF-8"
Content-transfer-encoding: quoted-printable

<html><head>
<meta http-equiv=3D"Content-Type" content=3D"text/html; charset=3Dutf-8"><meta na=
me=3D"viewport" content=3D"width=3Ddevice-width, initial-scale=3D1"><title>DSN</titl=
e></head><body style=3D"        background-color: white;      "><table style=3D"=
        background-color: white;         max-width: 548px;         color: bl=
ack;         border-spacing: 0px 0px;         padding-top: 0px;         padd=
ing-bottom: 0px;        border-collapse: collapse;      " width=3D"548" cellsp=
acing=3D"0" cellpadding=3D"0"><tbody><tr><td style=3D"        text-align: left;   =
     padding-bottom: 20px;      "><img height=3D"28" width=3D"126" style=3D"      =
  max-width: 100%;      " src=3D"https://products.office.com/en-us/CMSImages/O=
ffice365Logo_Orange.png?version=3Db8d100a9-0a8b-8e6a-88e1-ef488fee0470"> </td>=
</tr><tr><td style=3D"        font-family: 'Segoe UI', Frutiger, Arial, sans-s=
erif;        font-size: 16px;         padding-bottom: 10px;         -ms-text=
-size-adjust: 100%;        text-align: left;      ">Your message to <span st=
yle=3D"        color: #0072c6;      ">gordon@amateur-radio.org</span> couldn't=
 be delivered.<br></td></tr><tr><td style=3D"        font-family: 'Segoe UI', =
Frutiger, Arial, sans-serif;         font-size: 24px;         padding-top: 0=
px;         padding-bottom: 20px;         text-align: center;         -ms-te=
xt-size-adjust: 100%;      "><span style=3D"        color: #0072c6;      ">ama=
teur-radio.org</span> suspects your message is spam and rejected it.<br></td=
></tr><tr><td style=3D"        padding-bottom: 15px;         padding-left: 0px=
;         padding-right: 0px;         border-spacing: 0px 0px;      "><table=
 style=3D"        max-width: 548px;         font-weight: 600;        border-sp=
acing: 0px 0px;         padding-top: 0px;         padding-bottom: 0px;      =
  border-collapse: collapse;      "><tbody><tr><td style=3D"        font-famil=
y: 'Segoe UI', Frutiger, Arial, sans-serif;        font-size: 15px;        f=
ont-weight: 600;        text-align: left;        width: 181px;        -ms-te=
xt-size-adjust: 100%;        vertical-align: bottom;      "><font color=3D"#ff=
ffff"><span style=3D"color:#000000">cert</span> </font></td><td style=3D"       =
 font-family: 'Segoe UI', Frutiger, Arial, sans-serif;        font-size: 15p=
x;        font-weight: 600;        text-align: center;        width: 186px; =
       -ms-text-size-adjust: 100%;        vertical-align: bottom;      "><fo=
nt color=3D"#ffffff"><span style=3D"color:#000000">Office 365</span> </font></td=
><td style=3D"        font-family: 'Segoe UI', Frutiger, Arial, sans-serif;   =
      -ms-text-size-adjust: 100%;         font-size: 15px;         font-weig=
ht: 600;        text-align: right;         width: 181px;         vertical-al=
ign: bottom;      "><font color=3D"#ffffff"><span style=3D"color:#000000">amateu=
r-radio.org</span> </font></td></tr><tr><td style=3D"        font-family: 'Seg=
oe UI', Frutiger, Arial, sans-serif;        -ms-text-size-adjust: 100%;     =
   font-size: 14px;        font-weight: 400;        text-align: left;       =
 padding-top: 0px;        padding-bottom: 0px;        vertical-align: middle=
;        width: 181px;      "><font color=3D"#ffffff"><span style=3D"color:#0000=
00">Sender</span> </font></td><td style=3D"        font-family: 'Segoe UI', Fr=
utiger, Arial, sans-serif;        -ms-text-size-adjust: 100%;        font-si=
ze: 14px;        font-weight: 400;        text-align: center;        padding=
-top: 0px;        padding-bottom: 0px;        vertical-align: middle;       =
 width: 186px;      "></td><td style=3D"        font-family: 'Segoe UI', Fruti=
ger, Arial, sans-serif;        -ms-text-size-adjust: 100%;        font-size:=
 14px;        font-weight: 400;        text-align: right;        padding-top=
: 0px;        padding-bottom: 0px;        vertical-align: middle;        wid=
th: 181px;      "><font color=3D"#ffffff"><span style=3D"        color: #c00000;=
      "><b>Action Required</b> </span></font></td></tr><tr><td colspan=3D"3" s=
tyle=3D"        padding-top:0;        padding-bottom:0;        padding-left:0;=
        padding-right:0      "><table cellspacing=3D"0" cellpadding=3D"0" style=3D=
"        border-spacing: 0px 0px;        padding-top: 0px;        padding-bo=
ttom: 0px;        padding-left: 0px;        padding-right: 0px;        borde=
r-collapse: collapse;      "><tbody><tr height=3D"10"><td width=3D"180" height=3D"=
10" bgcolor=3D"#cccccc" style=3D"        width: 180px;        line-height: 10px;=
        height: 10px;        font-size: 6px;        padding-top: 0;        p=
adding-bottom: 0;        padding-left: 0;        padding-right: 0;      "><!=
--[if gte mso 15]>&nbsp;<![endif]--></td><td width=3D"4" height=3D"10" bgcolor=3D"=
#ffffff" style=3D"        width: 4px;        line-height: 10px;        height:=
 10px;        font-size: 6px;        padding-top: 0;        padding-bottom: =
0;        padding-left: 0;        padding-right: 0;      "><!--[if gte mso 1=
5]>&nbsp;<![endif]--></td><td width=3D"180" height=3D"10" bgcolor=3D"#cccccc" styl=
e=3D"        width: 180px;        line-height: 10px;        height: 10px;     =
   font-size: 6px;        padding-top: 0;        padding-bottom: 0;        p=
adding-left: 0;        padding-right: 0;      "><!--[if gte mso 15]>&nbsp;<!=
[endif]--></td><td width=3D"4" height=3D"10" bgcolor=3D"#ffffff" style=3D"        wi=
dth: 4px;        line-height: 10px;        height: 10px;        font-size: 6=
px;        padding-top: 0;        padding-bottom: 0;        padding-left: 0;=
        padding-right: 0;      "><!--[if gte mso 15]>&nbsp;<![endif]--></td>=
<td width=3D"180" height=3D"10" bgcolor=3D"#c00000" style=3D"        width: 180px;  =
      line-height: 10px;        height: 10px;        font-size: 6px;        =
padding-top: 0;        padding-bottom: 0;        padding-left: 0;        pad=
ding-right: 0;      "><!--[if gte mso 15]>&nbsp;<![endif]--></td></tr></tbod=
y></table></td></tr><tr><td style=3D"        font-family: 'Segoe UI', Frutiger=
, Arial, sans-serif;        -ms-text-size-adjust: 100%;        font-size: 14=
px;        text-align: left;        width: 181px;        line-height: 20px; =
       font-weight: 400;        padding-top: 0px;        padding-left: 0px; =
       padding-right: 0px;        padding-bottom: 0px;      "></td><td style=
=3D"        font-family: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-t=
ext-size-adjust: 100%;        font-size: 14px;        text-align: center;   =
     width: 186px;        line-height: 20px;        font-weight: 400;       =
 padding-top: 0px;        padding-left: 0px;        padding-right: 0px;     =
   padding-bottom: 0px;      "></td><td style=3D"        font-family: 'Segoe U=
I', Frutiger, Arial, sans-serif;        -ms-text-size-adjust: 100%;        f=
ont-size: 14px;        text-align: right;        width: 181px;        line-h=
eight: 20px;        font-weight: 400;        padding-top: 0px;        paddin=
g-left: 0px;        padding-right: 0px;        padding-bottom: 0px;      "><=
font color=3D"#ffffff"><span style=3D"        color: #c00000;      ">Messages su=
spected as spam</span> </font></td></tr></tbody></table></td></tr><tr><td st=
yle=3D"        width: 100%;        padding-top: 0px;        padding-right: 10p=
x;        padding-left: 10px;      "><br><table style=3D"        width: 100%; =
       padding-right: 0px;        padding-left: 0px;        padding-top: 0px=
;        padding-bottom: 0px;        background-color: #f2f5fa;        margi=
n-left: 0px;      "><tbody><tr><td style=3D"        font-family: 'Segoe UI', F=
rutiger, Arial, sans-serif;        -ms-text-size-adjust: 100%;        font-s=
ize: 21px;        font-weight: 500;        background-color: #f2f5fa;       =
 padding-top: 0px;        padding-bottom: 0px;        padding-left: 10px;   =
     padding-right: 10px;      ">How to Fix It</td></tr><tr><td style=3D"     =
   font-family: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-text-siz=
e-adjust: 100%;        font-size: 16px;        font-weight: 400;        padd=
ing-top: 0px;        padding-bottom: 6px;        padding-left: 10px;        =
padding-right: 10px;        background-color: #f2f5fa;      ">Try to modify =
your message, or change how you're sending the message, using the guidance i=
n this article: <a href=3D"https://go.microsoft.com/fwlink/?LinkID=3D526654">Bul=
k E-mailing Best Practices for Senders Using Forefront Online Protection for=
 Exchange</a>. Then resend your message.</td></tr><tr><td style=3D"        fon=
t-family: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-text-size-adju=
st: 100%;        font-size: 16px;        font-weight: 400;        padding-to=
p: 0px;        padding-bottom: 6px;        padding-left: 10px;        paddin=
g-right: 10px;        background-color: #f2f5fa;      ">If you continue to e=
xperience the problem, contact the recipient by some other means (by phone, =
for example) and ask them to ask their email admin to add your email address=
, or your domain name, to their allowed senders list.</td></tr></tbody></tab=
le></td></tr><tr><td style=3D"        font-family: 'Segoe UI', Frutiger, Arial=
, sans-serif;        -ms-text-size-adjust: 100%;        font-size: 14px;    =
    font-weight: 400;        padding-top: 10px;        padding-bottom: 0px; =
       padding-bottom: 4px;      "><br><em>Was this helpful? <a href=3D"https:=
//go.microsoft.com/fwlink/p/?LinkID=3D717204">Send feedback to Microsoft</a>.<=
/em> </td></tr><tr><td style=3D"        -ms-text-size-adjust: 100%;        fon=
t-size: 0px;        line-height: 0px;        padding-top: 0px;        paddin=
g-bottom: 0px;      "><hr></td></tr><tr><td style=3D"        font-family: 'Seg=
oe UI', Frutiger, Arial, sans-serif;        -ms-text-size-adjust: 100%;     =
   font-size: 21px;        font-weight: 500;      "><br>More Info for Email =
Admins</td></tr><tr><td style=3D"        font-family: 'Segoe UI', Frutiger, Ar=
ial, sans-serif;        -ms-text-size-adjust: 100%;        font-size: 14px; =
     "><em>Status code: 550 5.7.350</em> <br><br>When Office 365 tried to se=
nd the message to the recipient (outside Office 365), the recipient's email =
server (or email filtering service) suspected the sender's message is spam.<=
br><br>If the sender can't fix the problem by modifying their message, conta=
ct the recipient's email admin and ask them to add your domain name, or the =
sender's email address, to their list of allowed senders.<br><br>Although th=
e sender may be able to alter the message contents to fix this issue, it's l=
ikely that only the recipient's email admin can fix this problem. Unfortunat=
ely, Office 365 Support is unlikely to be able to help fix these kinds of ex=
ternally reported errors.<br><br></td></tr><tr><td style=3D"        font-famil=
y: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-text-size-adjust: 100=
%;        font-size: 17px;        font-weight: 500;      ">Original Message =
Details</td></tr><tr><td style=3D"        font-size: 14px;        line-height:=
 20px;        font-family: 'Segoe UI', Frutiger, Arial, sans-serif;        -=
ms-text-size-adjust: 100%;        font-weight: 500;      "><table style=3D"   =
     width: 100%;        border-collapse: collapse;        margin-left: 10px=
;      "><tbody><tr><td valign=3D"top" style=3D"        font-family: 'Segoe UI',=
 Frutiger, Arial, sans-serif;        font-size: 14px;        -ms-text-size-a=
djust: 100%;        white-space: nowrap;        font-weight: 500;        wid=
th: 140px;      ">Created Date:</td><td style=3D"        font-family: 'Segoe U=
I', Frutiger, Arial, sans-serif;        -ms-text-size-adjust: 100%;        f=
ont-size: 14px;        font-weight: 400;      ">12/1/2021 8:03:44 AM</td></t=
r><tr><td valign=3D"top" style=3D"        font-family: 'Segoe UI', Frutiger, Ari=
al, sans-serif;        font-size: 14px;        -ms-text-size-adjust: 100%;  =
      white-space: nowrap;        font-weight: 500;        width: 140px;    =
  ">Sender Address:</td><td style=3D"        font-family: 'Segoe UI', Frutiger=
, Arial, sans-serif;        -ms-text-size-adjust: 100%;        font-size: 14=
px;        font-weight: 400;      ">cert@sunnyvale.ca.gov</td></tr><tr><td s=
tyle=3D"        font-family: 'Segoe UI', Frutiger, Arial, sans-serif;        f=
ont-size: 14px;        -ms-text-size-adjust: 100%;        white-space: nowra=
p;        font-weight: 500;        width: 140px;      ">Recipient Address:</=
td><td style=3D"        font-family: 'Segoe UI', Frutiger, Arial, sans-serif; =
       -ms-text-size-adjust: 100%;        font-size: 14px;        font-weigh=
t: 400;      ">gordon@amateur-radio.org</td></tr><tr><td style=3D"        font=
-family: 'Segoe UI', Frutiger, Arial, sans-serif;        font-size: 14px;   =
     -ms-text-size-adjust: 100%;        white-space: nowrap;        font-wei=
ght: 500;        width: 140px;      ">Subject:</td><td style=3D"        font-f=
amily: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-text-size-adjust:=
 100%;        font-size: 14px;        font-weight: 400;      ">SERV Voluntee=
r Hours for November 2021</td></tr></tbody></table></td></tr><tr><td style=3D"=
        font-family: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-tex=
t-size-adjust: 100%;        font-size: 17px;        font-weight: 500;      "=
><br>Error Details</td></tr><tr><td style=3D"        font-size: 14px;        l=
ine-height: 20px;        font-family: 'Segoe UI', Frutiger, Arial, sans-seri=
f;        -ms-text-size-adjust: 100%;        font-weight: 500;      "><table=
 style=3D"        width: 100%;        border-collapse: collapse;        margin=
-left: 10px;      "><tbody><tr><td valign=3D"top" style=3D"        font-family: =
'Segoe UI', Frutiger, Arial, sans-serif;        font-size: 14px;        -ms-=
text-size-adjust: 100%;        white-space: nowrap;        font-weight: 500;=
        width: 140px;      ">Reported error:</td><td style=3D"        font-fam=
ily: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-text-size-adjust: 1=
00%;        font-size: 14px;        font-weight: 400;      "><em>550 5.7.350=
 Remote server returned message detected as spam -&gt; 554 5.7.1 The message=
 from (&lt;cert@sunnyvale.ca.gov&gt;) with the subject of (SERV Volunteer Ho=
urs for November 2021) matches a profile the Internet community may consider=
 spam. Please revise your message before resending.</em> </td></tr><tr><td s=
tyle=3D"        font-family: 'Segoe UI', Frutiger, Arial, sans-serif;        f=
ont-size: 14px;        -ms-text-size-adjust: 100%;        white-space: nowra=
p;        font-weight: 500;        width: 140px;      ">DSN generated by:</t=
d><td style=3D"        font-family: 'Segoe UI', Frutiger, Arial, sans-serif;  =
      -ms-text-size-adjust: 100%;        font-size: 14px;        font-weight=
: 400;      ">BY5PR09MB5379.namprd09.prod.outlook.com</td></tr></tbody></tab=
le></td></tr></tbody></table><br><table style=3D"width: 880px;" cellspacing=3D"0=
"><tbody><tr><td colspan=3D"6" style=3D"        padding-top: 4px;        border-=
bottom: 1px solid #999999;        padding-bottom: 4px;        line-height: 1=
20%;        font-size: 17px;        font-family: 'Segoe UI', Frutiger, Arial=
, sans-serif;        -ms-text-size-adjust: 100%;        font-weight: 500;   =
   ">Message Hops</td></tr><tr><td style=3D"        font-size: 12px;        fo=
nt-family: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-text-size-adj=
ust: 100%;        font-weight: 500;        background-color: #f2f5fa;       =
 border-bottom: 1px solid #999999;        white-space: nowrap;        paddin=
g: 8px;      ">HOP</td><td style=3D"        font-size: 12px;        font-famil=
y: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-text-size-adjust: 100=
%;        font-weight: 500;        background-color: #f2f5fa;        border-=
bottom: 1px solid #999999;        white-space: nowrap;        padding: 8px; =
       width: 80px;      ">TIME (UTC)</td><td style=3D"        font-size: 12px=
;        font-family: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-te=
xt-size-adjust: 100%;        font-weight: 500;        background-color: #f2f=
5fa;        border-bottom: 1px solid #999999;        white-space: nowrap;   =
     padding: 8px;      ">FROM</td><td style=3D"        font-size: 12px;      =
  font-family: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-text-size=
-adjust: 100%;        font-weight: 500;        background-color: #f2f5fa;   =
     border-bottom: 1px solid #999999;        white-space: nowrap;        pa=
dding: 8px;      ">TO</td><td style=3D"        font-size: 12px;        font-fa=
mily: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-text-size-adjust: =
100%;        font-weight: 500;        background-color: #f2f5fa;        bord=
er-bottom: 1px solid #999999;        white-space: nowrap;        padding: 8p=
x;      ">WITH</td><td style=3D"        font-size: 12px;        font-family: '=
Segoe UI', Frutiger, Arial, sans-serif;        -ms-text-size-adjust: 100%;  =
      font-weight: 500;        background-color: #f2f5fa;        border-bott=
om: 1px solid #999999;        white-space: nowrap;        padding: 8px;     =
 ">RELAY TIME</td></tr><tr><td style=3D"        font-size: 12px;        font-f=
amily: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-text-size-adjust:=
 100%;        font-weight: 500;        border-bottom: 1px solid #999999;    =
    padding: 8px;        text-align: center;      ">1</td><td style=3D"       =
 font-size: 12px;        font-family: 'Segoe UI', Frutiger, Arial, sans-seri=
f;        -ms-text-size-adjust: 100%;        font-weight: 500;        border=
-bottom: 1px solid #999999;        padding: 8px;        text-align: left;   =
     width: 80px;      ">12/1/2021<br>8:04:26 AM</td><td style=3D"        font=
-size: 12px;        font-family: 'Segoe UI', Frutiger, Arial, sans-serif;   =
     -ms-text-size-adjust: 100%;        font-weight: 500;        border-bott=
om: 1px solid #999999;        padding: 8px;        text-align: left;      ">=
localhost</td><td style=3D"        font-size: 12px;        font-family: 'Segoe=
 UI', Frutiger, Arial, sans-serif;        -ms-text-size-adjust: 100%;       =
 font-weight: 500;        border-bottom: 1px solid #999999;        padding: =
8px;        text-align: left;      ">BL1P223CA0025.NAMP223.PROD.OUTLOOK.COM<=
/td><td style=3D"        font-size: 12px;        font-family: 'Segoe UI', Frut=
iger, Arial, sans-serif;        -ms-text-size-adjust: 100%;        font-weig=
ht: 500;        border-bottom: 1px solid #999999;        padding: 8px;      =
  text-align: left;      ">Microsoft SMTP Server (version=3DTLS1_2, cipher=3DTLS=
_ECDHE_RSA_WITH_AES_256_GCM_SHA384)</td><td style=3D"        font-size: 12px; =
       font-family: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-text=
-size-adjust: 100%;        font-weight: 500;        border-bottom: 1px solid=
 #999999;        padding: 8px;        text-align: left;      ">42&nbsp;sec</=
td></tr><tr><td style=3D"        font-size: 12px;        font-family: 'Segoe U=
I', Frutiger, Arial, sans-serif;        -ms-text-size-adjust: 100%;        f=
ont-weight: 500;        border-bottom: 1px solid #999999;        padding: 8p=
x;        text-align: center;      ">2</td><td style=3D"        font-size: 12p=
x;        font-family: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-t=
ext-size-adjust: 100%;        font-weight: 500;        border-bottom: 1px so=
lid #999999;        padding: 8px;        text-align: left;        width: 80p=
x;      ">12/1/2021<br>8:04:27 AM</td><td style=3D"        font-size: 12px;   =
     font-family: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-text-s=
ize-adjust: 100%;        font-weight: 500;        border-bottom: 1px solid #=
999999;        padding: 8px;        text-align: left;      ">BY5PR09MB5090.n=
amprd09.prod.outlook.com</td><td style=3D"        font-size: 12px;        font=
-family: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-text-size-adjus=
t: 100%;        font-weight: 500;        border-bottom: 1px solid #999999;  =
      padding: 8px;        text-align: left;      ">BY5PR09MB5090.namprd09.p=
rod.outlook.com</td><td style=3D"        font-size: 12px;        font-family: =
'Segoe UI', Frutiger, Arial, sans-serif;        -ms-text-size-adjust: 100%; =
       font-weight: 500;        border-bottom: 1px solid #999999;        pad=
ding: 8px;        text-align: left;      ">mapi</td><td style=3D"        font-=
size: 12px;        font-family: 'Segoe UI', Frutiger, Arial, sans-serif;    =
    -ms-text-size-adjust: 100%;        font-weight: 500;        border-botto=
m: 1px solid #999999;        padding: 8px;        text-align: left;      ">1=
&nbsp;sec</td></tr><tr><td style=3D"        font-size: 12px;        font-famil=
y: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-text-size-adjust: 100=
%;        font-weight: 500;        border-bottom: 1px solid #999999;        =
padding: 8px;        text-align: center;      ">3</td><td style=3D"        fon=
t-size: 12px;        font-family: 'Segoe UI', Frutiger, Arial, sans-serif;  =
      -ms-text-size-adjust: 100%;        font-weight: 500;        border-bot=
tom: 1px solid #999999;        padding: 8px;        text-align: left;       =
 width: 80px;      ">12/1/2021<br>8:04:27 AM</td><td style=3D"        font-siz=
e: 12px;        font-family: 'Segoe UI', Frutiger, Arial, sans-serif;       =
 -ms-text-size-adjust: 100%;        font-weight: 500;        border-bottom: =
1px solid #999999;        padding: 8px;        text-align: left;      ">BY5P=
R09MB5090.namprd09.prod.outlook.com</td><td style=3D"        font-size: 12px; =
       font-family: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-text=
-size-adjust: 100%;        font-weight: 500;        border-bottom: 1px solid=
 #999999;        padding: 8px;        text-align: left;      ">BY5PR09MB5379=
.namprd09.prod.outlook.com</td><td style=3D"        font-size: 12px;        fo=
nt-family: 'Segoe UI', Frutiger, Arial, sans-serif;        -ms-text-size-adj=
ust: 100%;        font-weight: 500;        border-bottom: 1px solid #999999;=
        padding: 8px;        text-align: left;      ">Microsoft SMTP Server =
(version=3DTLS1_2, cipher=3DTLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384)</td><td style=
=3D"        font-size: 12px;        font-family: 'Segoe UI', Frutiger, Arial, =
sans-serif;        -ms-text-size-adjust: 100%;        font-weight: 500;     =
   border-bottom: 1px solid #999999;        padding: 8px;        text-align:=
 left;      ">*</td></tr></tbody></table><p style=3D"        font-family: 'Seg=
oe UI', Frutiger, Arial, sans-serif;        -ms-text-size-adjust: 100%;     =
   font-size: 17px;        font-weight: 500;        padding-top: 4px;       =
 padding-bottom: 0;        margin-top: 19px;        margin-bottom: 5px;     =
 ">Original Message Headers</p><pre style=3D"        color: gray;        white=
-space: pre;        padding-top: 0;        margin-top: 5px;      ">ARC-Seal:=
 i=3D1; a=3Drsa-sha256; s=3Darcselector9901; d=3Dmicrosoft.com; cv=3Dnone;
 b=3DElDXLNsNor1ieyeILBeH5YK1xaNSUH9k/3DrvN57MlvdD47bwqMqyrCF/aq579lWFnXnV4lx=
Q7CqmLgzj1JOhLfjGbkhwg64xDwCufjPj19znP1pthkPrbC3ASDOMSM1uVY5jjtdRGBz2KHnJ5gG=
yn9muTvqqpAP9xNLe4eJWX5Q5k5XsbLCimuLBJoFEF4aTJVB+6WD1wV0cBgwCWiG83NER5cYrWFv=
m9E0tZAW2ZJwZ7XMs7lrAHG6B8f5aroV0/1wOqLAkbS/31mZH2naMxhL3XsjX90KDUePESyTnclD=
3bj5jiX62Z6sE1E/DilGz+IBenWmvAoXpD6/Q2TSsg=3D=3D
ARC-Message-Signature: i=3D1; a=3Drsa-sha256; c=3Drelaxed/relaxed; d=3Dmicrosoft.co=
m;
 s=3Darcselector9901;
 h=3DFrom:Date:Subject:Message-ID:Content-Type:MIME-Version:X-MS-Exchange-Ant=
iSpam-MessageData-ChunkCount:X-MS-Exchange-AntiSpam-MessageData-0:X-MS-Excha=
nge-AntiSpam-MessageData-1;
 bh=3DxiMYKUzhNVXMjKgeNC0NBpIE1l7F6ew8A35kP50faOY=3D;
 b=3DS1+82Rkj0k6Twql7XKI8WsGwxXC/WQHMmqCvfUm+Z80KwmTAMOExeMfg62vES1aToAHjE+Y0=
eXmBxQBKfo5crXkU5p8/AJDoZPcztwniN4AjJVOS6jAuqsoZvdDwPFa7J/r0L9VoJzNwekS6PXT+=
YtelGaysf+iu8xwG8mlyZFkx2thBYz5mYUX6qaFaqGbW3oIVPnCr0V/4SydMn6Dqu4BUWzKlvuWc=
FUkfETyMAhd2gNXj4aDL+3JPmOu/bncFce19PnSBWbvgkNAiZIU3ypCd41v1NgEgyWeNjGeLrys0=
SRV42QIjTtcjwTf1/AE5dr9FFvK0t8ikjC74uhb57w=3D=3D
ARC-Authentication-Results: i=3D1; mx.microsoft.com 1; spf=3Dpass
 smtp.mailfrom=3Dsunnyvale.ca.gov; dmarc=3Dpass action=3Dnone
 header.from=3Dsunnyvale.ca.gov; dkim=3Dpass header.d=3Dsunnyvale.ca.gov; arc=3Dnon=
e
DKIM-Signature: v=3D1; a=3Drsa-sha256; c=3Drelaxed/relaxed;
 d=3Dcityofsunnyvale.onmicrosoft.com;
 s=3Dselector2-cityofsunnyvale-onmicrosoft-com;
 h=3DFrom:Date:Subject:Message-ID:Content-Type:MIME-Version:X-MS-Exchange-Sen=
derADCheck;
 bh=3DxiMYKUzhNVXMjKgeNC0NBpIE1l7F6ew8A35kP50faOY=3D;
 b=3D0q2m3qItWh5cXQJd+TvjJfR/rqd20b2kN3+hNbTiY57BtFfo6FJkyu2bbOrlgWZR7pe9/roj=
V2txhcXOg2ulr3nGxetRJ6EhgJF1C2ekRiz0IWE3KGTMOk40K8kZmmGXo3XVJ3bRw4+y2XpB59ZL=
SWfwsoFed+FmOPOExiE7lx4=3D
Authentication-Results: dkim=3Dnone (message not signed)
 header.d=3Dnone;dmarc=3Dnone action=3Dnone header.from=3Dsunnyvale.ca.gov;
Received: from BY5PR09MB5090.namprd09.prod.outlook.com (2603:10b6:a03:24a::=
14)
 by BY5PR09MB5379.namprd09.prod.outlook.com (2603:10b6:a03:24d::11) with
 Microsoft SMTP Server (version=3DTLS1_2,
 cipher=3DTLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384) id 15.20.4734.22; Wed, 1 Dec
 2021 08:04:27 +0000
Received: from BY5PR09MB5090.namprd09.prod.outlook.com
 ([fe80::3ca3:b9c4:5e3d:58f6]) by BY5PR09MB5090.namprd09.prod.outlook.com
 ([fe80::3ca3:b9c4:5e3d:58f6%3]) with mapi id 15.20.4755.014; Wed, 1 Dec 20=
21
 08:04:27 +0000
From: SunnyvaleSERV.org &lt;cert@sunnyvale.ca.gov&gt;
Content-Type: multipart/alternative; boundary=3D&quot;BOUNDARY&quot;
Subject: SERV Volunteer Hours for November 2021
Date: Wed, 01 Dec 2021 00:03:44 -0800
To: &quot;Gordon Girton&quot; &lt;gordon@amateur-radio.org&gt;
X-ClientProxiedBy: BL1P223CA0025.NAMP223.PROD.OUTLOOK.COM
 (2603:10b6:208:2c4::30) To BY5PR09MB5090.namprd09.prod.outlook.com
 (2603:10b6:a03:24a::14)
Return-Path: cert@sunnyvale.ca.gov
Message-ID: &lt;BY5PR09MB5090C0A1798249A18D3F34A3E7689@BY5PR09MB5090.namprd=
09.prod.outlook.com&gt;
MIME-Version: 1.0
Received: from localhost (2607:f298:5:100f::320:4a05) by BL1P223CA0025.NAMP=
223.PROD.OUTLOOK.COM (2603:10b6:208:2c4::30) with Microsoft SMTP Server (ver=
sion=3DTLS1_2, cipher=3DTLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384) id 15.20.4734.23 =
via Frontend Transport; Wed, 1 Dec 2021 08:04:26 +0000
X-MS-PublicTrafficType: Email
X-MS-Office365-Filtering-Correlation-Id: d270b94b-ec7a-41e9-6661-08d9b4a131=
ee
X-MS-TrafficTypeDiagnostic: BY5PR09MB5379:
X-Microsoft-Antispam-PRVS:
	&lt;BY5PR09MB53799E7833AAA8B89EC78C75E7689@BY5PR09MB5379.namprd09.prod.out=
look.com&gt;
X-MS-Oob-TLC-OOBClassifiers: OLM:6790;
X-MS-Exchange-SenderADCheck: 1
X-MS-Exchange-AntiSpam-Relay: 0
X-Microsoft-Antispam: BCL:0;
X-Microsoft-Antispam-Message-Info:
	xDd6uy4T4vl0VR3tODu+h8hmVT8n/RR0LsMnRjE1JkCqQjNbDvZCW//Yc74Yu0AGTolRxrk87i=
iNr60ari8xnFHroywalti12y5ZgKg+99C3P5FN4LjcX2YIdFqlOllMqziDoyTOu9kA4D8H00nCY/=
RixXb1oBbssQLuHNuem+OknTSeb3WVrBYtpXrTyJMmQEeuOd5wHAyCR3Wde5/mGfWMLsqLc9Ii2c=
XyM0NeVZSfP+zmeq5TPi9f5cSLQEpjgvUbiP4HST9O1zkebBU+uXf9d3cvGR17iiLkSz/0V+4wU+=
ecqdwoY1WiHWdWNIaYQcomQxGD5BE8z8Atyg8IIJI1K0QN+gzaNfkA8joWGMRnpmhptHyFBYd0BL=
d+ZZoqF+jG5C+zwKhkLMQOij9+CZf9/pTgYjmdybvOO1ICMSjkPzkr1fDoCRmnDNEb9J9vOIsHXs=
1ijpRN4LAKO43yVmZk/8tWPHsIFIGjAY3PNuhK12AR9abkgWYToFyapyLwHlOytdbXedZdfrV/kP=
RGk4X6Q5bdqjtX6BjIlqgyjXAXkFFlE947CMNhn5vNEIv6slS4PwwDqP1PuyYBTFGUZN5LO/hOw7=
P457JTrHICUlPCsixRzLB4kdppg+tKToIQkDKTqR5+VM5bWH60fBHseEpFg+WVb/wEk2gBm0g3gq=
SOStaJcXx4lpO5lf5F6fcX6NG2Nh6iixZPUD/c2CxIWUeAsUc9aAtouNLNlZFC7N8=3D
X-Forefront-Antispam-Report:
	CIP:255.255.255.255;CTRY:;LANG:en;SCL:1;SRV:;IPV:NLI;SFV:NSPM;H:BY5PR09MB5=
090.namprd09.prod.outlook.com;PTR:;CAT:NONE;SFS:(366004)(186003)(40011501000=
01)(8936002)(38100700002)(5660300002)(83380400001)(6666004)(6916009)(8676002=
)(508600001)(316002)(9686003)(966005)(52536014)(6496006)(6486002)(2906002)(6=
6476007)(52116002)(33964004)(166002)(33656002)(66556008)(66946007)(86362001)=
;DIR:OUT;SFP:1102;
X-MS-Exchange-AntiSpam-MessageData-ChunkCount: 1
X-MS-Exchange-AntiSpam-MessageData-0:
	=3D?us-ascii?Q?NbKYDOktWeK55k5Jod3Eh0cXtj/BpM9uyUI3AqGbM/X35OnV/mzFdDopk3xu?=
=3D
 =3D?us-ascii?Q?kckxY0rMyxo87Vfh1HQAObmFYo0rpOMiQkgQkpCfkbIjtOfTz5djTHpd3c8V?=
=3D
 =3D?us-ascii?Q?Whoy/b13NKtH4ksh9KUtGre5sl5y35BThCsJISs/HZKc45JlajC9dFLbXmJk?=
=3D
 =3D?us-ascii?Q?B17SJqoK8aWoXp0sjsX8PwI1SDUBlvQFRrU39yTtf2ogCnp+J9BIZjudmrLo?=
=3D
 =3D?us-ascii?Q?SxtPR4pY9HcMnQWvgW89Gvjs0a5S+OmUuv2dJXWQKF1QvMhePWk/NVw1y83T?=
=3D
 =3D?us-ascii?Q?AT8+/zRPdpv/sHJxEQOWdjyVnjQwG5LXsoPDRoIkR4PQscRcYCtfDzldkR0R?=
=3D
 =3D?us-ascii?Q?K5k41wVViL3Bj2Yt9rRr4f4f3ZY3EfptXUF79hCYwurkDasekmtZEe+W6517?=
=3D
 =3D?us-ascii?Q?uP9qrAGBpBoIfReNdG5qSLVme1vFBziAtB866b4iHBpQJyhxifj9FxfeWuWH?=
=3D
 =3D?us-ascii?Q?V75MfO5G2R0ImZnCXzCD8++mXMKtClF0zeEnTlAxzWg8EJi/nlzXWa1niETU?=
=3D
 =3D?us-ascii?Q?G/eQg//FHfAX9jlVcDMdYlrDi+Hjj1cDsNBVzE4nu7ij20SQfkg3K6odEac9?=
=3D
 =3D?us-ascii?Q?4zi3U03GaEFR7CCR+tI/acVFgAaS7TYSGpNHVlthjhsVPlfwM/H0xUSVEact?=
=3D
 =3D?us-ascii?Q?+lK4SFy6Zc75cdQIjbwLr3x8/A2CDlLxEnctCReahtJCv02RjLAlAtMFkd7j?=
=3D
 =3D?us-ascii?Q?DzXy2sDhOYN3sbXnBGMdUD/F15QG4Rg+YWZXJwZXgPsPLl0eClIe1ygeFfJ8?=
=3D
 =3D?us-ascii?Q?xEB/OdDRbIeGvw0zLW0g1IAMhZCEqQax8BEQtxwU2KBsQVxU/k+tDjiRIgcK?=
=3D
 =3D?us-ascii?Q?aO1UV5M0D9k4nzoIgXWFnbkQFqDdtLqrzU5Vnipgh72VnpxufkOfM2xDIs67?=
=3D
 =3D?us-ascii?Q?Y8ShpLfbWElbdwc0IzxJ3WfS5/CFiYkN4QMHWUzV9reC0F8Ov7uu82alLvLD?=
=3D
 =3D?us-ascii?Q?02+9OVkXZiw4h6TO3czbyj239WhuQoCHZivGvbsU31N7JVuLTSBzBKcfdgWl?=
=3D
 =3D?us-ascii?Q?p7WmccNFOKP/GY4vLaDev9KUUCqGy3DHFqCp/PeJaWNbZlQFyVBgELPmbZbc?=
=3D
 =3D?us-ascii?Q?7mL3iCrFrLF65oYg/6QfLaUregT5f5JENi3S+tKUQitQ1BoE9f5KpbqDDf//?=
=3D
 =3D?us-ascii?Q?PPVdLs1JDSmdCJyIHkmr0I7fpl/PYPf/hRbKeVjxELuHY9fOjpEZQkS/ZDOM?=
=3D
 =3D?us-ascii?Q?g7fE+ebolUFU9dpsPFM3tu7FtO7arSmfB4oyFx7Up3FqrvDA+EnnWZ8Gmyd1?=
=3D
 =3D?us-ascii?Q?CZXfWBBgKTazj/u1xrMw2sJFB/Xv+FnPdnVdz+yS/CRwgdoMbHKQTg60R1UH?=
=3D
 =3D?us-ascii?Q?VHvZ7ldvkZVh1138NhcahYslKeP33e53t7dVkslDbqONpEAbtTMkERgqm+tW?=
=3D
 =3D?us-ascii?Q?Ex9R+pw=3D3D?=3D
X-OriginatorOrg: sunnyvale.ca.gov
X-MS-Exchange-CrossTenant-Network-Message-Id: d270b94b-ec7a-41e9-6661-08d9b=
4a131ee
X-MS-Exchange-CrossTenant-AuthSource: BY5PR09MB5090.namprd09.prod.outlook.c=
om
X-MS-Exchange-CrossTenant-AuthAs: Internal
X-MS-Exchange-CrossTenant-OriginalArrivalTime: 01 Dec 2021 08:04:27.0714
 (UTC)
X-MS-Exchange-CrossTenant-FromEntityHeader: Hosted
X-MS-Exchange-CrossTenant-Id: 63dc83d7-8dcb-489b-bce2-0bf7c37d8f39
X-MS-Exchange-Transport-CrossTenantHeadersStamped: BY5PR09MB5379
</pre></body></html>

--B_3721553817_622503736--


--B_3721553817_1320843075
Content-type: message/delivery-status; name="ATT00001";
 x-mac-creator="4F50494D"
Content-ID: <0D3695B74C239D4BB4A6C2C7197F1AE0@namprd09.prod.outlook.com>
Content-disposition: attachment;
	filename="ATT00001"
Content-transfer-encoding: base64


UmVwb3J0aW5nLU1UQTogZG5zO0JZNVBSMDlNQjUzNzkubmFtcHJkMDkucHJvZC5vdXRsb29r
LmNvbQ0KUmVjZWl2ZWQtRnJvbS1NVEE6IGRucztCWTVQUjA5TUI1MDkwLm5hbXByZDA5LnBy
b2Qub3V0bG9vay5jb20NCkFycml2YWwtRGF0ZTogV2VkLCAxIERlYyAyMDIxIDA4OjA0OjI3
ICswMDAwDQoNCkZpbmFsLVJlY2lwaWVudDogcmZjODIyO2dvcmRvbkBhbWF0ZXVyLXJhZGlv
Lm9yZw0KQWN0aW9uOiBmYWlsZWQNClN0YXR1czogNS43LjM1MA0KRGlhZ25vc3RpYy1Db2Rl
OiBzbXRwOzU1MCA1LjcuMzUwIFJlbW90ZSBzZXJ2ZXIgcmV0dXJuZWQgbWVzc2FnZSBkZXRl
Y3RlZCBhcyBzcGFtIC0+IDU1NCA1LjcuMSBUaGUgbWVzc2FnZSBmcm9tICg8Y2VydEBzdW5u
eXZhbGUuY2EuZ292Pikgd2l0aCB0aGUgc3ViamVjdCBvZiAoU0VSViBWb2x1bnRlZXIgSG91
cnMgZm9yIE5vdmVtYmVyIDIwMjEpIG1hdGNoZXMgYSBwcm9maWxlIHRoZSBJbnRlcm5ldCBj
b21tdW5pdHkgbWF5IGNvbnNpZGVyIHNwYW0uIFBsZWFzZSByZXZpc2UgeW91ciBtZXNzYWdl
IGJlZm9yZSByZXNlbmRpbmcuDQpYLURpc3BsYXktTmFtZTogR29yZG9uIEdpcnRvbg0KDQo=
--B_3721553817_1320843075
Content-type: message/rfc822
Content-disposition: attachment

ARC-Seal: i=1; a=rsa-sha256; s=arcselector9901; d=microsoft.com; cv=none;
 b=ElDXLNsNor1ieyeILBeH5YK1xaNSUH9k/3DrvN57MlvdD47bwqMqyrCF/aq579lWFnXnV4lxQ7CqmLgzj1JOhLfjGbkhwg64xDwCufjPj19znP1pthkPrbC3ASDOMSM1uVY5jjtdRGBz2KHnJ5gGyn9muTvqqpAP9xNLe4eJWX5Q5k5XsbLCimuLBJoFEF4aTJVB+6WD1wV0cBgwCWiG83NER5cYrWFvm9E0tZAW2ZJwZ7XMs7lrAHG6B8f5aroV0/1wOqLAkbS/31mZH2naMxhL3XsjX90KDUePESyTnclD3bj5jiX62Z6sE1E/DilGz+IBenWmvAoXpD6/Q2TSsg==
ARC-Message-Signature: i=1; a=rsa-sha256; c=relaxed/relaxed; d=microsoft.com;
 s=arcselector9901;
 h=From:Date:Subject:Message-ID:Content-Type:MIME-Version:X-MS-Exchange-AntiSpam-MessageData-ChunkCount:X-MS-Exchange-AntiSpam-MessageData-0:X-MS-Exchange-AntiSpam-MessageData-1;
 bh=xiMYKUzhNVXMjKgeNC0NBpIE1l7F6ew8A35kP50faOY=;
 b=S1+82Rkj0k6Twql7XKI8WsGwxXC/WQHMmqCvfUm+Z80KwmTAMOExeMfg62vES1aToAHjE+Y0eXmBxQBKfo5crXkU5p8/AJDoZPcztwniN4AjJVOS6jAuqsoZvdDwPFa7J/r0L9VoJzNwekS6PXT+YtelGaysf+iu8xwG8mlyZFkx2thBYz5mYUX6qaFaqGbW3oIVPnCr0V/4SydMn6Dqu4BUWzKlvuWcFUkfETyMAhd2gNXj4aDL+3JPmOu/bncFce19PnSBWbvgkNAiZIU3ypCd41v1NgEgyWeNjGeLrys0SRV42QIjTtcjwTf1/AE5dr9FFvK0t8ikjC74uhb57w==
ARC-Authentication-Results: i=1; mx.microsoft.com 1; spf=pass
 smtp.mailfrom=sunnyvale.ca.gov; dmarc=pass action=none
 header.from=sunnyvale.ca.gov; dkim=pass header.d=sunnyvale.ca.gov; arc=none
DKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;
 d=cityofsunnyvale.onmicrosoft.com;
 s=selector2-cityofsunnyvale-onmicrosoft-com;
 h=From:Date:Subject:Message-ID:Content-Type:MIME-Version:X-MS-Exchange-SenderADCheck;
 bh=xiMYKUzhNVXMjKgeNC0NBpIE1l7F6ew8A35kP50faOY=;
 b=0q2m3qItWh5cXQJd+TvjJfR/rqd20b2kN3+hNbTiY57BtFfo6FJkyu2bbOrlgWZR7pe9/rojV2txhcXOg2ulr3nGxetRJ6EhgJF1C2ekRiz0IWE3KGTMOk40K8kZmmGXo3XVJ3bRw4+y2XpB59ZLSWfwsoFed+FmOPOExiE7lx4=
Authentication-Results: dkim=none (message not signed)
 header.d=none;dmarc=none action=none header.from=sunnyvale.ca.gov;
Received: from BY5PR09MB5090.namprd09.prod.outlook.com (2603:10b6:a03:24a::14)
 by BY5PR09MB5379.namprd09.prod.outlook.com (2603:10b6:a03:24d::11) with
 Microsoft SMTP Server (version=TLS1_2,
 cipher=TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384) id 15.20.4734.22; Wed, 1 Dec
 2021 08:04:27 +0000
Received: from BY5PR09MB5090.namprd09.prod.outlook.com
 ([fe80::3ca3:b9c4:5e3d:58f6]) by BY5PR09MB5090.namprd09.prod.outlook.com
 ([fe80::3ca3:b9c4:5e3d:58f6%3]) with mapi id 15.20.4755.014; Wed, 1 Dec 2021
 08:04:27 +0000
From: SunnyvaleSERV.org <cert@sunnyvale.ca.gov>
Content-Type: multipart/alternative; boundary="BOUNDARY"
Subject: SERV Volunteer Hours for November 2021
Date: Wed, 01 Dec 2021 00:03:44 -0800
To: "Gordon Girton" <gordon@amateur-radio.org>
X-ClientProxiedBy: BL1P223CA0025.NAMP223.PROD.OUTLOOK.COM
 (2603:10b6:208:2c4::30) To BY5PR09MB5090.namprd09.prod.outlook.com
 (2603:10b6:a03:24a::14)
Return-Path: cert@sunnyvale.ca.gov
Message-ID: <BY5PR09MB5090C0A1798249A18D3F34A3E7689@BY5PR09MB5090.namprd09.prod.outlook.com>
Received: from localhost (2607:f298:5:100f::320:4a05) by BL1P223CA0025.NAMP223.PROD.OUTLOOK.COM (2603:10b6:208:2c4::30) with Microsoft SMTP Server (version=TLS1_2, cipher=TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384) id 15.20.4734.23 via Frontend Transport; Wed, 1 Dec 2021 08:04:26 +0000
X-MS-PublicTrafficType: Email
X-MS-Office365-Filtering-Correlation-Id: d270b94b-ec7a-41e9-6661-08d9b4a131ee
X-MS-TrafficTypeDiagnostic: BY5PR09MB5379:
X-Microsoft-Antispam-PRVS:
 <BY5PR09MB53799E7833AAA8B89EC78C75E7689@BY5PR09MB5379.namprd09.prod.outlook.com>
X-MS-Oob-TLC-OOBClassifiers: OLM:6790;
X-MS-Exchange-SenderADCheck: 1
X-MS-Exchange-AntiSpam-Relay: 0
X-Microsoft-Antispam: BCL:0;
X-Microsoft-Antispam-Message-Info:
 xDd6uy4T4vl0VR3tODu+h8hmVT8n/RR0LsMnRjE1JkCqQjNbDvZCW//Yc74Yu0AGTolRxrk87iiNr60ari8xnFHroywalti12y5ZgKg+99C3P5FN4LjcX2YIdFqlOllMqziDoyTOu9kA4D8H00nCY/RixXb1oBbssQLuHNuem+OknTSeb3WVrBYtpXrTyJMmQEeuOd5wHAyCR3Wde5/mGfWMLsqLc9Ii2cXyM0NeVZSfP+zmeq5TPi9f5cSLQEpjgvUbiP4HST9O1zkebBU+uXf9d3cvGR17iiLkSz/0V+4wU+ecqdwoY1WiHWdWNIaYQcomQxGD5BE8z8Atyg8IIJI1K0QN+gzaNfkA8joWGMRnpmhptHyFBYd0BLd+ZZoqF+jG5C+zwKhkLMQOij9+CZf9/pTgYjmdybvOO1ICMSjkPzkr1fDoCRmnDNEb9J9vOIsHXs1ijpRN4LAKO43yVmZk/8tWPHsIFIGjAY3PNuhK12AR9abkgWYToFyapyLwHlOytdbXedZdfrV/kPRGk4X6Q5bdqjtX6BjIlqgyjXAXkFFlE947CMNhn5vNEIv6slS4PwwDqP1PuyYBTFGUZN5LO/hOw7P457JTrHICUlPCsixRzLB4kdppg+tKToIQkDKTqR5+VM5bWH60fBHseEpFg+WVb/wEk2gBm0g3gqSOStaJcXx4lpO5lf5F6fcX6NG2Nh6iixZPUD/c2CxIWUeAsUc9aAtouNLNlZFC7N8=
X-Forefront-Antispam-Report:
 CIP:255.255.255.255;CTRY:;LANG:en;SCL:1;SRV:;IPV:NLI;SFV:NSPM;H:BY5PR09MB5090.namprd09.prod.outlook.com;PTR:;CAT:NONE;SFS:(366004)(186003)(4001150100001)(8936002)(38100700002)(5660300002)(83380400001)(6666004)(6916009)(8676002)(508600001)(316002)(9686003)(966005)(52536014)(6496006)(6486002)(2906002)(66476007)(52116002)(33964004)(166002)(33656002)(66556008)(66946007)(86362001);DIR:OUT;SFP:1102;
X-MS-Exchange-AntiSpam-MessageData-ChunkCount: 1
X-MS-Exchange-AntiSpam-MessageData-0:
 =?us-ascii?Q?NbKYDOktWeK55k5Jod3Eh0cXtj/BpM9uyUI3AqGbM/X35OnV/mzFdDopk3xu?=
 =?us-ascii?Q?kckxY0rMyxo87Vfh1HQAObmFYo0rpOMiQkgQkpCfkbIjtOfTz5djTHpd3c8V?=
 =?us-ascii?Q?Whoy/b13NKtH4ksh9KUtGre5sl5y35BThCsJISs/HZKc45JlajC9dFLbXmJk?=
 =?us-ascii?Q?B17SJqoK8aWoXp0sjsX8PwI1SDUBlvQFRrU39yTtf2ogCnp+J9BIZjudmrLo?=
 =?us-ascii?Q?SxtPR4pY9HcMnQWvgW89Gvjs0a5S+OmUuv2dJXWQKF1QvMhePWk/NVw1y83T?=
 =?us-ascii?Q?AT8+/zRPdpv/sHJxEQOWdjyVnjQwG5LXsoPDRoIkR4PQscRcYCtfDzldkR0R?=
 =?us-ascii?Q?K5k41wVViL3Bj2Yt9rRr4f4f3ZY3EfptXUF79hCYwurkDasekmtZEe+W6517?=
 =?us-ascii?Q?uP9qrAGBpBoIfReNdG5qSLVme1vFBziAtB866b4iHBpQJyhxifj9FxfeWuWH?=
 =?us-ascii?Q?V75MfO5G2R0ImZnCXzCD8++mXMKtClF0zeEnTlAxzWg8EJi/nlzXWa1niETU?=
 =?us-ascii?Q?G/eQg//FHfAX9jlVcDMdYlrDi+Hjj1cDsNBVzE4nu7ij20SQfkg3K6odEac9?=
 =?us-ascii?Q?4zi3U03GaEFR7CCR+tI/acVFgAaS7TYSGpNHVlthjhsVPlfwM/H0xUSVEact?=
 =?us-ascii?Q?+lK4SFy6Zc75cdQIjbwLr3x8/A2CDlLxEnctCReahtJCv02RjLAlAtMFkd7j?=
 =?us-ascii?Q?DzXy2sDhOYN3sbXnBGMdUD/F15QG4Rg+YWZXJwZXgPsPLl0eClIe1ygeFfJ8?=
 =?us-ascii?Q?xEB/OdDRbIeGvw0zLW0g1IAMhZCEqQax8BEQtxwU2KBsQVxU/k+tDjiRIgcK?=
 =?us-ascii?Q?aO1UV5M0D9k4nzoIgXWFnbkQFqDdtLqrzU5Vnipgh72VnpxufkOfM2xDIs67?=
 =?us-ascii?Q?Y8ShpLfbWElbdwc0IzxJ3WfS5/CFiYkN4QMHWUzV9reC0F8Ov7uu82alLvLD?=
 =?us-ascii?Q?02+9OVkXZiw4h6TO3czbyj239WhuQoCHZivGvbsU31N7JVuLTSBzBKcfdgWl?=
 =?us-ascii?Q?p7WmccNFOKP/GY4vLaDev9KUUCqGy3DHFqCp/PeJaWNbZlQFyVBgELPmbZbc?=
 =?us-ascii?Q?7mL3iCrFrLF65oYg/6QfLaUregT5f5JENi3S+tKUQitQ1BoE9f5KpbqDDf//?=
 =?us-ascii?Q?PPVdLs1JDSmdCJyIHkmr0I7fpl/PYPf/hRbKeVjxELuHY9fOjpEZQkS/ZDOM?=
 =?us-ascii?Q?g7fE+ebolUFU9dpsPFM3tu7FtO7arSmfB4oyFx7Up3FqrvDA+EnnWZ8Gmyd1?=
 =?us-ascii?Q?CZXfWBBgKTazj/u1xrMw2sJFB/Xv+FnPdnVdz+yS/CRwgdoMbHKQTg60R1UH?=
 =?us-ascii?Q?VHvZ7ldvkZVh1138NhcahYslKeP33e53t7dVkslDbqONpEAbtTMkERgqm+tW?=
 =?us-ascii?Q?Ex9R+pw=3D?=
X-OriginatorOrg: sunnyvale.ca.gov
X-MS-Exchange-CrossTenant-Network-Message-Id: d270b94b-ec7a-41e9-6661-08d9b4a131ee
X-MS-Exchange-CrossTenant-AuthSource: BY5PR09MB5090.namprd09.prod.outlook.com
X-MS-Exchange-CrossTenant-AuthAs: Internal
X-MS-Exchange-CrossTenant-OriginalArrivalTime: 01 Dec 2021 08:04:27.0714
 (UTC)
X-MS-Exchange-CrossTenant-FromEntityHeader: Hosted
X-MS-Exchange-CrossTenant-Id: 63dc83d7-8dcb-489b-bce2-0bf7c37d8f39
X-MS-Exchange-Transport-CrossTenantHeadersStamped: BY5PR09MB5379
MIME-Version: 1.0

--BOUNDARY
Content-Type: text/plain; charset=utf-8
Content-Transfer-Encoding: quoted-printable

Hello, Gordon Girton,

The Sunnyvale Office of Emergency Services would like to be sure that you g=
et credit for all of the volunteer work you may have done for SERV, CERT, L=
istos, SARES, and/or SNAP during November 2021.  Currently, our records sho=
w:
    2021-11-01 SARES Net: 1.0 Hours
    2021-11-08 SARES Net: 1.0 Hours
    2021-11-15 SARES Net: 1.0 Hours
    2021-11-22 SARES Net: 1.0 Hours
    2021-11-29 SARES Net: 1.0 Hours
    2021-11-30 Other SARES Hours: 1.0 Hours
    Total Hours: 6.0 Hours
If these records are incorrect, or you have any additional volunteer time t=
o report, please visit
    https://sunnyvaleserv.org/volunteer-hours/lBwLeBsmrJYX8DpgaNoRdkNM3l2W1=
5ry
and report it prior to December 10.  If you have questions about reporting =
volunteer hours, just reply to this email.

Many thanks,
Sunnyvale OES


--BOUNDARY
Content-Type: text/html; charset=utf-8
Content-Transfer-Encoding: quoted-printable

<html><head>
<meta http-equiv=3D"Content-Type" content=3D"text/html; charset=3Dutf-8"></=
head><body><div>Hello, Gordon Girton,</div><div><br></div><div>The Sunnyval=
e Office of Emergency Services would like to be sure that you get credit fo=
r all of the volunteer work you may have done for SERV, CERT, Listos, SARES=
, and/or SNAP during November 2021.  Currently, our records show:</div><div=
><br></div><table style=3D"margin-left:2em"><tr><td>2021-11-01 SARES Net</t=
d><td style=3D"text-align:right;padding-left:1em">1.0 Hours</td></tr><tr><t=
d>2021-11-08 SARES Net</td><td style=3D"text-align:right;padding-left:1em">=
1.0 Hours</td></tr><tr><td>2021-11-15 SARES Net</td><td style=3D"text-align=
:right;padding-left:1em">1.0 Hours</td></tr><tr><td>2021-11-22 SARES Net</t=
d><td style=3D"text-align:right;padding-left:1em">1.0 Hours</td></tr><tr><t=
d>2021-11-29 SARES Net</td><td style=3D"text-align:right;padding-left:1em">=
1.0 Hours</td></tr><tr><td>2021-11-30 Other SARES Hours</td><td style=3D"te=
xt-align:right;padding-left:1em">1.0 Hours</td></tr><tr><td style=3D"text-a=
lign:right">Total</td><td style=3D"text-align:right;padding-left:1em">6.0 H=
ours</td></tr></table><div><br></div><div>If these records are incorrect, o=
r you have any additional volunteer time to report, please visit <a href=3D=
"https://sunnyvaleserv.org/volunteer-hours/lBwLeBsmrJYX8DpgaNoRdkNM3l2W15ry=
">our web site</a> and report it prior to December 10.</div><div style=3D"m=
argin:1em 0"><a style=3D"color:#fff;background-color:#007bff;border:1px sol=
id #007bff;border-radius:4px;padding:6px 12px;line-height:1.5;text-align:ce=
nter;vertical-align:middle;display:inline-block;cursor:pointer;user-select:=
none;text-decoration:none" href=3D"https://sunnyvaleserv.org/volunteer-hour=
s/lBwLeBsmrJYX8DpgaNoRdkNM3l2W15ry">Report Hours</a><div style=3D"display:i=
nline-block;margin-left:16px;color:#888">This button takes you directly to =
your Activity page without needing to log in.</div></div><div>If you have q=
uestions about reporting volunteer hours, just reply to this email.</div><d=
iv><br></div><div>Many thanks,<br>Sunnyvale OES</div></body></html>

--BOUNDARY--

--B_3721553817_1320843075--

`

const expectedBounceBody = `
Your message to gordon@amateur-radio.org couldn't be delivered.
amateur-radio.org suspects your message is spam and rejected it.
cert Office 365 amateur-radio.org
Sender Action Required
Messages suspected as spam

How to Fix It
Try to modify your message, or change how you're sending the message, using the guidance in this article: Bulk E-mailing Best Practices for Senders Using Forefront Online Protection for Exchange. Then resend your message.
If you continue to experience the problem, contact the recipient by some other means (by phone, for example) and ask them to ask their email admin to add your email address, or your domain name, to their allowed senders list.

Was this helpful? Send feedback to Microsoft.

More Info for Email Admins
Status code: 550 5.7.350

When Office 365 tried to send the message to the recipient (outside Office 365), the recipient's email server (or email filtering service) suspected the sender's message is spam.

If the sender can't fix the problem by modifying their message, contact the recipient's email admin and ask them to add your domain name, or the sender's email address, to their list of allowed senders.

Although the sender may be able to alter the message contents to fix this issue, it's likely that only the recipient's email admin can fix this problem. Unfortunately, Office 365 Support is unlikely to be able to help fix these kinds of externally reported errors.

Original Message Details
Created Date:12/1/2021 8:03:44 AM
Sender Address:cert@sunnyvale.ca.gov
Recipient Address:gordon@amateur-radio.org
Subject:SERV Volunteer Hours for November 2021

Error Details
Reported error:550 5.7.350 Remote server returned message detected as spam -> 554 5.7.1 The message from (<cert@sunnyvale.ca.gov>) with the subject of (SERV Volunteer Hours for November 2021) matches a profile the Internet community may consider spam. Please revise your message before resending.
DSN generated by:BY5PR09MB5379.namprd09.prod.outlook.com
Message Hops
HOPTIME (UTC)FROMTOWITHRELAY TIME
112/1/2021
8:04:26 AMlocalhostBL1P223CA0025.NAMP223.PROD.OUTLOOK.COMMicrosoft SMTP Server (version=TLS1_2, cipher=TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384)42 sec
212/1/2021
8:04:27 AMBY5PR09MB5090.namprd09.prod.outlook.comBY5PR09MB5090.namprd09.prod.outlook.commapi1 sec
312/1/2021
8:04:27 AMBY5PR09MB5090.namprd09.prod.outlook.comBY5PR09MB5379.namprd09.prod.outlook.comMicrosoft SMTP Server (version=TLS1_2, cipher=TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384)*

Original Message Headers
ARC-Seal: i=1; a=rsa-sha256; s=arcselector9901; d=microsoft.com; cv=none;
 b=ElDXLNsNor1ieyeILBeH5YK1xaNSUH9k/3DrvN57MlvdD47bwqMqyrCF/aq579lWFnXnV4lxQ7CqmLgzj1JOhLfjGbkhwg64xDwCufjPj19znP1pthkPrbC3ASDOMSM1uVY5jjtdRGBz2KHnJ5gGyn9muTvqqpAP9xNLe4eJWX5Q5k5XsbLCimuLBJoFEF4aTJVB+6WD1wV0cBgwCWiG83NER5cYrWFvm9E0tZAW2ZJwZ7XMs7lrAHG6B8f5aroV0/1wOqLAkbS/31mZH2naMxhL3XsjX90KDUePESyTnclD3bj5jiX62Z6sE1E/DilGz+IBenWmvAoXpD6/Q2TSsg==
ARC-Message-Signature: i=1; a=rsa-sha256; c=relaxed/relaxed; d=microsoft.com;
 s=arcselector9901;
 h=From:Date:Subject:Message-ID:Content-Type:MIME-Version:X-MS-Exchange-AntiSpam-MessageData-ChunkCount:X-MS-Exchange-AntiSpam-MessageData-0:X-MS-Exchange-AntiSpam-MessageData-1;
 bh=xiMYKUzhNVXMjKgeNC0NBpIE1l7F6ew8A35kP50faOY=;
 b=S1+82Rkj0k6Twql7XKI8WsGwxXC/WQHMmqCvfUm+Z80KwmTAMOExeMfg62vES1aToAHjE+Y0eXmBxQBKfo5crXkU5p8/AJDoZPcztwniN4AjJVOS6jAuqsoZvdDwPFa7J/r0L9VoJzNwekS6PXT+YtelGaysf+iu8xwG8mlyZFkx2thBYz5mYUX6qaFaqGbW3oIVPnCr0V/4SydMn6Dqu4BUWzKlvuWcFUkfETyMAhd2gNXj4aDL+3JPmOu/bncFce19PnSBWbvgkNAiZIU3ypCd41v1NgEgyWeNjGeLrys0SRV42QIjTtcjwTf1/AE5dr9FFvK0t8ikjC74uhb57w==
ARC-Authentication-Results: i=1; mx.microsoft.com 1; spf=pass
 smtp.mailfrom=sunnyvale.ca.gov; dmarc=pass action=none
 header.from=sunnyvale.ca.gov; dkim=pass header.d=sunnyvale.ca.gov; arc=none
DKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;
 d=cityofsunnyvale.onmicrosoft.com;
 s=selector2-cityofsunnyvale-onmicrosoft-com;
 h=From:Date:Subject:Message-ID:Content-Type:MIME-Version:X-MS-Exchange-SenderADCheck;
 bh=xiMYKUzhNVXMjKgeNC0NBpIE1l7F6ew8A35kP50faOY=;
 b=0q2m3qItWh5cXQJd+TvjJfR/rqd20b2kN3+hNbTiY57BtFfo6FJkyu2bbOrlgWZR7pe9/rojV2txhcXOg2ulr3nGxetRJ6EhgJF1C2ekRiz0IWE3KGTMOk40K8kZmmGXo3XVJ3bRw4+y2XpB59ZLSWfwsoFed+FmOPOExiE7lx4=
Authentication-Results: dkim=none (message not signed)
 header.d=none;dmarc=none action=none header.from=sunnyvale.ca.gov;
Received: from BY5PR09MB5090.namprd09.prod.outlook.com (2603:10b6:a03:24a::14)
 by BY5PR09MB5379.namprd09.prod.outlook.com (2603:10b6:a03:24d::11) with
 Microsoft SMTP Server (version=TLS1_2,
 cipher=TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384) id 15.20.4734.22; Wed, 1 Dec
 2021 08:04:27 +0000
Received: from BY5PR09MB5090.namprd09.prod.outlook.com
 ([fe80::3ca3:b9c4:5e3d:58f6]) by BY5PR09MB5090.namprd09.prod.outlook.com
 ([fe80::3ca3:b9c4:5e3d:58f6%3]) with mapi id 15.20.4755.014; Wed, 1 Dec 2021
 08:04:27 +0000
From: SunnyvaleSERV.org <cert@sunnyvale.ca.gov>
Content-Type: multipart/alternative; boundary="BOUNDARY"
Subject: SERV Volunteer Hours for November 2021
Date: Wed, 01 Dec 2021 00:03:44 -0800
To: "Gordon Girton" <gordon@amateur-radio.org>
X-ClientProxiedBy: BL1P223CA0025.NAMP223.PROD.OUTLOOK.COM
 (2603:10b6:208:2c4::30) To BY5PR09MB5090.namprd09.prod.outlook.com
 (2603:10b6:a03:24a::14)
Return-Path: cert@sunnyvale.ca.gov
Message-ID: <BY5PR09MB5090C0A1798249A18D3F34A3E7689@BY5PR09MB5090.namprd09.prod.outlook.com>
MIME-Version: 1.0
Received: from localhost (2607:f298:5:100f::320:4a05) by BL1P223CA0025.NAMP223.PROD.OUTLOOK.COM (2603:10b6:208:2c4::30) with Microsoft SMTP Server (version=TLS1_2, cipher=TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384) id 15.20.4734.23 via Frontend Transport; Wed, 1 Dec 2021 08:04:26 +0000
X-MS-PublicTrafficType: Email
X-MS-Office365-Filtering-Correlation-Id: d270b94b-ec7a-41e9-6661-08d9b4a131ee
X-MS-TrafficTypeDiagnostic: BY5PR09MB5379:
X-Microsoft-Antispam-PRVS:
	<BY5PR09MB53799E7833AAA8B89EC78C75E7689@BY5PR09MB5379.namprd09.prod.outlook.com>
X-MS-Oob-TLC-OOBClassifiers: OLM:6790;
X-MS-Exchange-SenderADCheck: 1
X-MS-Exchange-AntiSpam-Relay: 0
X-Microsoft-Antispam: BCL:0;
X-Microsoft-Antispam-Message-Info:
	xDd6uy4T4vl0VR3tODu+h8hmVT8n/RR0LsMnRjE1JkCqQjNbDvZCW//Yc74Yu0AGTolRxrk87iiNr60ari8xnFHroywalti12y5ZgKg+99C3P5FN4LjcX2YIdFqlOllMqziDoyTOu9kA4D8H00nCY/RixXb1oBbssQLuHNuem+OknTSeb3WVrBYtpXrTyJMmQEeuOd5wHAyCR3Wde5/mGfWMLsqLc9Ii2cXyM0NeVZSfP+zmeq5TPi9f5cSLQEpjgvUbiP4HST9O1zkebBU+uXf9d3cvGR17iiLkSz/0V+4wU+ecqdwoY1WiHWdWNIaYQcomQxGD5BE8z8Atyg8IIJI1K0QN+gzaNfkA8joWGMRnpmhptHyFBYd0BLd+ZZoqF+jG5C+zwKhkLMQOij9+CZf9/pTgYjmdybvOO1ICMSjkPzkr1fDoCRmnDNEb9J9vOIsHXs1ijpRN4LAKO43yVmZk/8tWPHsIFIGjAY3PNuhK12AR9abkgWYToFyapyLwHlOytdbXedZdfrV/kPRGk4X6Q5bdqjtX6BjIlqgyjXAXkFFlE947CMNhn5vNEIv6slS4PwwDqP1PuyYBTFGUZN5LO/hOw7P457JTrHICUlPCsixRzLB4kdppg+tKToIQkDKTqR5+VM5bWH60fBHseEpFg+WVb/wEk2gBm0g3gqSOStaJcXx4lpO5lf5F6fcX6NG2Nh6iixZPUD/c2CxIWUeAsUc9aAtouNLNlZFC7N8=
X-Forefront-Antispam-Report:
	CIP:255.255.255.255;CTRY:;LANG:en;SCL:1;SRV:;IPV:NLI;SFV:NSPM;H:BY5PR09MB5090.namprd09.prod.outlook.com;PTR:;CAT:NONE;SFS:(366004)(186003)(4001150100001)(8936002)(38100700002)(5660300002)(83380400001)(6666004)(6916009)(8676002)(508600001)(316002)(9686003)(966005)(52536014)(6496006)(6486002)(2906002)(66476007)(52116002)(33964004)(166002)(33656002)(66556008)(66946007)(86362001);DIR:OUT;SFP:1102;
X-MS-Exchange-AntiSpam-MessageData-ChunkCount: 1
X-MS-Exchange-AntiSpam-MessageData-0:
	=?us-ascii?Q?NbKYDOktWeK55k5Jod3Eh0cXtj/BpM9uyUI3AqGbM/X35OnV/mzFdDopk3xu?=
 =?us-ascii?Q?kckxY0rMyxo87Vfh1HQAObmFYo0rpOMiQkgQkpCfkbIjtOfTz5djTHpd3c8V?=
 =?us-ascii?Q?Whoy/b13NKtH4ksh9KUtGre5sl5y35BThCsJISs/HZKc45JlajC9dFLbXmJk?=
 =?us-ascii?Q?B17SJqoK8aWoXp0sjsX8PwI1SDUBlvQFRrU39yTtf2ogCnp+J9BIZjudmrLo?=
 =?us-ascii?Q?SxtPR4pY9HcMnQWvgW89Gvjs0a5S+OmUuv2dJXWQKF1QvMhePWk/NVw1y83T?=
 =?us-ascii?Q?AT8+/zRPdpv/sHJxEQOWdjyVnjQwG5LXsoPDRoIkR4PQscRcYCtfDzldkR0R?=
 =?us-ascii?Q?K5k41wVViL3Bj2Yt9rRr4f4f3ZY3EfptXUF79hCYwurkDasekmtZEe+W6517?=
 =?us-ascii?Q?uP9qrAGBpBoIfReNdG5qSLVme1vFBziAtB866b4iHBpQJyhxifj9FxfeWuWH?=
 =?us-ascii?Q?V75MfO5G2R0ImZnCXzCD8++mXMKtClF0zeEnTlAxzWg8EJi/nlzXWa1niETU?=
 =?us-ascii?Q?G/eQg//FHfAX9jlVcDMdYlrDi+Hjj1cDsNBVzE4nu7ij20SQfkg3K6odEac9?=
 =?us-ascii?Q?4zi3U03GaEFR7CCR+tI/acVFgAaS7TYSGpNHVlthjhsVPlfwM/H0xUSVEact?=
 =?us-ascii?Q?+lK4SFy6Zc75cdQIjbwLr3x8/A2CDlLxEnctCReahtJCv02RjLAlAtMFkd7j?=
 =?us-ascii?Q?DzXy2sDhOYN3sbXnBGMdUD/F15QG4Rg+YWZXJwZXgPsPLl0eClIe1ygeFfJ8?=
 =?us-ascii?Q?xEB/OdDRbIeGvw0zLW0g1IAMhZCEqQax8BEQtxwU2KBsQVxU/k+tDjiRIgcK?=
 =?us-ascii?Q?aO1UV5M0D9k4nzoIgXWFnbkQFqDdtLqrzU5Vnipgh72VnpxufkOfM2xDIs67?=
 =?us-ascii?Q?Y8ShpLfbWElbdwc0IzxJ3WfS5/CFiYkN4QMHWUzV9reC0F8Ov7uu82alLvLD?=
 =?us-ascii?Q?02+9OVkXZiw4h6TO3czbyj239WhuQoCHZivGvbsU31N7JVuLTSBzBKcfdgWl?=
 =?us-ascii?Q?p7WmccNFOKP/GY4vLaDev9KUUCqGy3DHFqCp/PeJaWNbZlQFyVBgELPmbZbc?=
 =?us-ascii?Q?7mL3iCrFrLF65oYg/6QfLaUregT5f5JENi3S+tKUQitQ1BoE9f5KpbqDDf//?=
 =?us-ascii?Q?PPVdLs1JDSmdCJyIHkmr0I7fpl/PYPf/hRbKeVjxELuHY9fOjpEZQkS/ZDOM?=
 =?us-ascii?Q?g7fE+ebolUFU9dpsPFM3tu7FtO7arSmfB4oyFx7Up3FqrvDA+EnnWZ8Gmyd1?=
 =?us-ascii?Q?CZXfWBBgKTazj/u1xrMw2sJFB/Xv+FnPdnVdz+yS/CRwgdoMbHKQTg60R1UH?=
 =?us-ascii?Q?VHvZ7ldvkZVh1138NhcahYslKeP33e53t7dVkslDbqONpEAbtTMkERgqm+tW?=
 =?us-ascii?Q?Ex9R+pw=3D?=
X-OriginatorOrg: sunnyvale.ca.gov
X-MS-Exchange-CrossTenant-Network-Message-Id: d270b94b-ec7a-41e9-6661-08d9b4a131ee
X-MS-Exchange-CrossTenant-AuthSource: BY5PR09MB5090.namprd09.prod.outlook.com
X-MS-Exchange-CrossTenant-AuthAs: Internal
X-MS-Exchange-CrossTenant-OriginalArrivalTime: 01 Dec 2021 08:04:27.0714
 (UTC)
X-MS-Exchange-CrossTenant-FromEntityHeader: Hosted
X-MS-Exchange-CrossTenant-Id: 63dc83d7-8dcb-489b-bce2-0bf7c37d8f39
X-MS-Exchange-Transport-CrossTenantHeadersStamped: BY5PR09MB5379

`
