package pktmsg

import (
	"net/textproto"
	"testing"
	"time"
)

var encodeBodyTests = []struct {
	name     string
	msg      *Message
	human    bool
	wantBody string
}{
	{
		"empty",
		&Message{},
		false,
		"",
	},
	{
		"plain body",
		&Message{
			Body: "nothing\n",
		},
		false,
		"nothing\n",
	},
	{
		"Outpost flags",
		&Message{
			Body:  "nothing\n",
			Flags: RequestDeliveryReceipt | RequestReadReceipt | OutpostUrgent,
		},
		false,
		"!URG!!RDR!!RRR!nothing\n",
	},
	{
		"base64",
		&Message{
			Body: "nöthing\n",
		},
		false,
		"!B64!bsO2dGhpbmcK\n",
	},
	{
		"base64 human",
		&Message{
			Body: "nöthing\n",
		},
		true,
		"nöthing\n",
	},
}

func TestEncodeBody(t *testing.T) {
	for _, tt := range encodeBodyTests {
		t.Run(tt.name, func(t *testing.T) {
			if gotBody := tt.msg.EncodeBody(tt.human); gotBody != tt.wantBody {
				t.Errorf("Message.EncodeBody() = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}

var encodeTests = []struct {
	name  string
	msg   *Message
	human bool
	want  string
}{
	{
		"plain",
		&Message{Body: "nothing\n"},
		false,
		"\nnothing\n",
	},
	{
		"headers",
		&Message{
			Header: textproto.MIMEHeader{
				"To":   []string{"<nobody@nowhere>", "<somebody@somewhere>"},
				"From": []string{"<me@here>"},
			},
			Body: "nothing\n",
		},
		false,
		"To: <nobody@nowhere>, <somebody@somewhere>\nFrom: <me@here>\n\nnothing\n",
	},
	{
		"envelope",
		&Message{
			EnvelopeAddress: "me@here",
			Body:            "nothing\n",
		},
		false,
		"From me@here\n\nnothing\n",
	},
	{
		"everything",
		&Message{
			EnvelopeAddress: "me@here",
			EnvelopeDate:    time.Date(2021, 12, 1, 8, 4, 29, 0, time.Local),
			Header:          textproto.MIMEHeader{"To": []string{"<nobody@nowhere>"}},
			Body:            "nothing\n",
		},
		false,
		"From me@here Wed Dec  1 08:04:29 2021\nTo: <nobody@nowhere>\n\nnothing\n",
	},
}

func TestEncode(t *testing.T) {
	for _, tt := range encodeTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.Encode(tt.human); got != tt.want {
				t.Errorf("Message.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
