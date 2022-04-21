package analyze

// This file contains the list of references that can be added to the bottom of
// a response message.

type reference uint8

const (
	refWeeklyPractice reference = 1 << iota
	refSubjectLine
	refOutpostConfig
	refFormRouting
	refPacketGroup
)

var allReferences = []reference{refWeeklyPractice, refSubjectLine, refOutpostConfig, refFormRouting, refPacketGroup}

var referenceText = map[reference]string{
	refWeeklyPractice: `
  * The "Weekly SPECS/SVECS Packet Practice" page on the county ARES/RACES
    website gives details of the packet practice exercise, including the net
    practice schedules, the schedule of what type of message to send, the
    schedule of simulated "down" BBS systems, and the format of the subject
    line for practice messages.  It is available at
    https://www.scc-ares-races.org/data/packet/weekly-packet-practice.html`,
	refSubjectLine: `
  * The "Standard Packet Message Subject Line" document describes how to
    compose the subject line of a packet message following county standards.
    It is available from the "Packet BBS Service" page at
    https://www.scc-ares-races.org/data/packet/index.html`,
	refOutpostConfig: `
  * The "Standard Outpost Configuration Instructions" document describes how
    to configure the Outpost messaging software to send messages following
    county standards.  It is available from the "Packet BBS Service" page at
    https://www.scc-ares-races.org/data/packet/index.html`,
	refFormRouting: `
  * The "SCCo ARES/RACES Recommended Form Routing" document gives
    recommendations for, among other things, what handling orders should be
    used for different types of forms, and what positions and locations they
    should be sent to.  It is available from the "Go Kit Forms" page at
    https://www.scc-ares-races.org/operations/go-kit-forms.html`,
	refPacketGroup: `
  * If you need assistance, you can request it in the packet discussion group.
    To sign up for this group, see the Discussion Groups page at
    https://www.scc-ares-races.org/discuss-groups.html
`,
}
