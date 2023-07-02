package checkout

import (
	"fmt"

	"github.com/rothskeller/packet/message/common"
)

// EncodeSubject encodes the message subject.
func (m *CheckOut) EncodeSubject() string {
	if m.TacticalCallSign != "" {
		return common.EncodeSubject(m.OriginMsgID, m.Handling, "",
			fmt.Sprintf("Check-Out %s, %s\n", m.TacticalCallSign, m.TacticalStationName))
	}
	return common.EncodeSubject(m.OriginMsgID, m.Handling, "",
		fmt.Sprintf("Check-Out %s, %s\n", m.OperatorCallSign, m.OperatorName))
}

// EncodeBody encodes the message body.
func (m *CheckOut) EncodeBody() string {
	if m.TacticalCallSign != "" {
		return fmt.Sprintf("Check-Out %s, %s\n%s, %s\n",
			m.TacticalCallSign, m.TacticalStationName, m.OperatorCallSign, m.OperatorName)
	}
	return fmt.Sprintf("Check-Out %s, %s\n", m.OperatorCallSign, m.OperatorName)
}
