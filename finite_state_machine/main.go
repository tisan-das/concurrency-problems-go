package main

import (
	"concurrency_problems/finite_state_machine/fsm"
	"log"
)

func main() {

	stateMachine := fsm.NewFiniteStateMachine("locked", &map[fsm.State]fsm.Transition{
		"locked": fsm.Transition{
			"coin": fsm.Action{Operation: coinInserted, Args: []interface{}{"locked"},
				Destination: "unlocked"},
			"push": fsm.Action{Operation: pushedButton, Args: []interface{}{"locked"},
				Destination: "locked"},
		},
		"unlocked": fsm.Transition{
			"coin": fsm.Action{Operation: coinInserted, Args: []interface{}{"unlocked"},
				Destination: "unlocked"},
			"push": fsm.Action{Operation: pushedButton, Args: []interface{}{"unlocked"},
				Destination: "locked"},
		},
	})

	orderedEvents := []fsm.Event{"coin", "coin", "push", "push", "coin", "touch"}
	log.Printf("Triggering the state conversion with the events %+v for the state machine: %+v",
		orderedEvents, stateMachine)

	for i, event := range orderedEvents {
		log.Printf("Before handing event %s at index %d Current state of FSM: %s",
			event, i, stateMachine.GetCurrentState())
		err := stateMachine.HandleEvent(event)
		if err != nil {
			log.Printf("Error: %s", err)
		}
		log.Printf("After handing event %s at index %d Current state of FSM: %s",
			event, i, stateMachine.GetCurrentState())
	}
	log.Print("Execution Completed!")
}

func coinInserted(args ...interface{}) error {
	msg := "Method coinInserted: Coin Inserted"
	if len(args) > 0 {
		stringVal, _ := args[0].(string)
		msg += " with current state " + stringVal
	}
	log.Print(msg)
	return nil
}

func pushedButton(args ...interface{}) error {
	msg := "Method pushedButton: Button pushed"
	if len(args) > 0 {
		stringVal, _ := args[0].(string)
		msg += " with current state " + stringVal
	}
	log.Print(msg)
	return nil
}
