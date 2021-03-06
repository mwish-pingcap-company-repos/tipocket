// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package core

// Model specifies the behavior of a data object.
type Model interface {
	// Prepare the initial state of the data object.
	Prepare(state interface{})

	// Initial state of the data object.
	Init() interface{}

	// Step function for the data object. Returns whether or not the system
	// could take this step with the given inputs and outputs and also
	// returns the new state. This should not mutate the existing state.
	//
	// state must support encoding to and decoding from json.
	Step(state interface{}, input interface{}, output interface{}) (bool, interface{})

	// Equality on states.
	Equal(state1, state2 interface{}) bool

	// Name returns the unique name for the model.
	Name() string
}

// Operation action
const (
	InvokeOperation = "call"
	ReturnOperation = "return"
)

// Operation of a data object.
type Operation struct {
	Action string      `json:"action"`
	Proc   int64       `json:"proc"`
	Data   interface{} `json:"data"`
}

// NoopModel is noop model.
type NoopModel struct {
	perparedState interface{}
}

// Prepare the initial state of the data object.
func (n *NoopModel) Prepare(state interface{}) {
	n.perparedState = state
}

// Init initials state of the data object.
func (n *NoopModel) Init() interface{} {
	return n.perparedState
}

// Step function for the data object.
func (*NoopModel) Step(state interface{}, input interface{}, output interface{}) (bool, interface{}) {
	return true, state
}

// Equal returns equality on states.
func (*NoopModel) Equal(state1, state2 interface{}) bool {
	return true
}

// Name returns the unique name for the model.
func (*NoopModel) Name() string {
	return "NoopModel"
}
