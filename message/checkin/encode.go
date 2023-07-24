package checkin

import (
	"fmt"

	"github.com/rothskeller/packet/message/common"
)

// EncodeSubject encodes the message subject.
func (m *CheckIn) EncodeSubject() string {
	if m.TacticalCallSign != "" {
		return common.EncodeSubject(m.OriginMsgID, m.Handling, "",
			fmt.Sprintf("Check-In %s, %s", m.TacticalCallSign, m.TacticalStationName))
	}
	return common.EncodeSubject(m.OriginMsgID, m.Handling, "",
		fmt.Sprintf("Check-In %s, %s", m.OperatorCallSign, m.OperatorName))
}

// EncodeBody encodes the message body.
func (m *CheckIn) EncodeBody() string {
	if m.TacticalCallSign != "" {
		return fmt.Sprintf("Check-In %s, %s\n%s, %s\n",
			m.TacticalCallSign, m.TacticalStationName, m.OperatorCallSign, m.OperatorName)
	}
	return fmt.Sprintf("Check-In %s, %s\n", m.OperatorCallSign, m.OperatorName)
}
