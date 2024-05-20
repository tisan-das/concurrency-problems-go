package fsm

import "fmt"

type Action struct {
	Operation   func(...interface{}) error
	Args        []interface{}
	Destination State
}

type Transition map[Event]Action

type FiniteStateMachine struct {
	initialState         State
	currentState         State
	stateTransitionTable *map[State]Transition
}

func NewFiniteStateMachine(initState State, transitionTable *map[State]Transition) IStateMachine {
	return &FiniteStateMachine{
		initialState:         initState,
		currentState:         initState,
		stateTransitionTable: transitionTable,
	}
}

func (stateMachine *FiniteStateMachine) GetCurrentState() State {
	return stateMachine.currentState
}

func (stateMachine *FiniteStateMachine) HandleEvent(event Event) error {

	if stateMachine.stateTransitionTable == nil {
		return nil
	}

	transitionsFromCurrentState, found := (*stateMachine.stateTransitionTable)[stateMachine.currentState]
	if !found {
		return fmt.Errorf(
			"Transitions not found for the event %+v from the current state %s",
			event, stateMachine.currentState)
	}
	actionForEvent, found := transitionsFromCurrentState[event]
	if !found {
		return fmt.Errorf("Action not found for event %+v from the current state %s",
			event, stateMachine.currentState)
	}
	if actionForEvent.Operation != nil {
		actionForEvent.Operation(actionForEvent.Args...)
	}
	stateMachine.currentState = actionForEvent.Destination
	return nil
}
