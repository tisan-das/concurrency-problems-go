package fsm

type State string
type Event string

type IStateMachine interface {
	GetCurrentState() State
	HandleEvent(event Event) error
}
