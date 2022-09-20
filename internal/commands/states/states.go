package states

type chatState string

const (
	StartCreate      chatState = "startCreate"
	AddTitle         chatState = "addTitle"
	AddDescription   chatState = "addDescription"
	SetAttendeeLimit chatState = "setAttendeeLimit"
	SetAttendeeRetry chatState = "setAttendeeRetry"
	SetDate          chatState = "setDate"
	SetDateRetry     chatState = "setDateRetry"
	SetLocation      chatState = "setLocation"
	SetDuration      chatState = "setDuration"
	SetDurationRetry chatState = "setDurationRetry"
	CreateEvent      chatState = "createEvent"

	StartEdit           chatState = "startEdit"
	StartEditRetry      chatState = "startEditRetry"
	ModifyEvent         chatState = "modifyEvent"
	ModifyEventRetry    chatState = "modifyEventRetry"
	RemoveResponse      chatState = "removeResponse"
	RemoveResponseRetry chatState = "removeResponseRetry"
	ContinueEdit        chatState = "continueEdit"
	ContinueEditRetry   chatState = "continueEditRetry"
	ProcessEdit         chatState = "processEdit"

	AddResponse      chatState = "addResponse"
	UnknownUser      chatState = "unknownUser"
	UnknownUserRetry chatState = "unknownUserRetry"
	SignUp           chatState = "signup"
	SignUpRetry      chatState = "signupRetry"

	Idle    chatState = "idle"
	Cancel  chatState = "cancel"
	Timeout chatState = "timeout"

	SelfTransition chatState = "selfTransition"
)

func (c chatState) String() string {
	return string(c)
}
