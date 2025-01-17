package tg

// DisplayMode defines how a message should be displayed in the UI.
type DisplayMode uint8

const (
	// DisplayModal indicates the message should be displayed in a modal dialog.
	DisplayModal DisplayMode = iota + 1
	// DisplayInline indicates the message should be displayed in a box.
	DisplayInline
	// DisplayScroller indicates the message should be displayed in a scrolling box.
	DisplayScroller
)

// MessageKind defines the message to be displayed.
type MessageKind uint16

const (
	MsgSessionExpires MessageKind = iota + 1
	MsgSessionExpired
	MsgComplexRoundTimeout
)

var messageCodes = []string{
	"message.session-expires",
	"message.session-expired",
	"message.complex-round-timeout",
}

// String implements the Stringer interface and returns the i18n message code.
func (m MessageKind) String() string {
	if m > 0 && int(m) <= len(messageCodes) {
		return messageCodes[m-1]
	}
	return "message.unsupported-message"
}
