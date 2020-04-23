package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessGraph(t *testing.T) {
	history, err := ParseHistory(`{:index 0 :process 1 :type :ok}
{:index 1 :process 2 :type :ok}
{:index 2 :process 2 :type :ok}
{:index 3 :process 1 :type :ok}`)

	assert.Equal(t, err, nil, "test process graph, parse history")

	var (
		v0 = Vertex{Value: history[0]}
		v1 = Vertex{Value: history[1]}
		v2 = Vertex{Value: history[2]}
		v3 = Vertex{Value: history[3]}
	)
	processGraph := NewDirectedGraph()
	processGraph.Outs = map[Vertex]map[Vertex][]Rel{
		v0: {
			v3: {Process},
		},
		v1: {
			v2: {Process},
		},
	}

	_, g, _ := ProcessGraph(history)

	var graphOuts map[Vertex]map[Vertex][]Rel = map[Vertex]map[Vertex][]Rel{}
	for k, v := range g.Outs {
		if len(v) != 0 {
			v, e := processGraph.Outs[k]
			if !e {
				continue
			}
			graphOuts[k] = v
		}
	}

	assert.Equal(t, graphOuts, processGraph.Outs)
}

func TestRealtimeGraph(t *testing.T) {
	history, err := ParseHistory(`{:type :invoke :process 1 :f :read :value nil}
{:type :ok      :process 1 :f :read :value 1}
{:type :invoke :process 2 :f :read :value nil}
{:type :ok      :process 2 :f :read :value 2}
{:type :invoke :process 3 :f :read :value nil}
{:type :ok      :process 3 :f :read :value 3}
{:type :invoke :process 4 :f :read :value nil}
{:type :ok      :process 4 :f :read :value 4}
{:type :invoke :process 5 :f :read :value nil}
{:type :ok      :process 5 :f :read :value 5}`)

	assert.Equal(t, err, nil, "test process graph, parse history")
	assert.Equal(t, len(history), 10, "length of history error")

	call1, resp1 := history[0], history[1]
	call2, resp2 := history[2], history[3]
	call3, resp3 := history[4], history[5]
	call4, resp4 := history[6], history[7]
	call5, resp5 := history[8], history[9]
	// Test	empty
	_, g, _ := RealtimeGraph(History{})
	if len(g.Ins) != 0 || len(g.Outs) != 0 {
		assert.Fail(t, "empty history should return empty graph")
	}

	_, g, _ = RealtimeGraph(History{call1, resp1})

	assert.Equal(t, g.Outs, map[Vertex]map[Vertex][]Rel{})

	_, g, _ = RealtimeGraph(History{call1, resp1, call2, resp2})
	dest := map[Vertex]map[Vertex][]Rel{
		Vertex{Value: resp1}: {
			Vertex{Value: resp2}: []Rel{Realtime},
		},
		Vertex{Value: resp2}: {},
	}

	assert.Equal(t, g.Outs, dest)

	_, g, _ = RealtimeGraph(History{call1, resp1, call2, resp2, call3, resp3})

	dest = map[Vertex]map[Vertex][]Rel{
		Vertex{Value: resp1}: {
			Vertex{Value: resp2}: []Rel{Realtime},
		},
		Vertex{Value: resp2}: {
			Vertex{Value: resp3}: []Rel{Realtime},
		},
		Vertex{Value: resp3}: {},
	}

	assert.Equal(t, g.Outs, dest)

	_, g, _ = RealtimeGraph(History{call1, resp1, call2, call3, resp3, resp2})

	dest = map[Vertex]map[Vertex][]Rel{
		Vertex{Value: resp1}: {
			Vertex{Value: resp2}: []Rel{Realtime},
			Vertex{Value: resp3}: []Rel{Realtime},
		},
		Vertex{Value: resp2}: {},
		Vertex{Value: resp3}: {},
	}

	assert.Equal(t, g.Outs, dest)

	_, g, _ = RealtimeGraph(History{call1, call2, resp2, resp1, call3, call4, resp3, resp4})

	dest = map[Vertex]map[Vertex][]Rel{
		Vertex{Value: resp1}: {
			Vertex{Value: resp4}: []Rel{Realtime},
			Vertex{Value: resp3}: []Rel{Realtime},
		},
		Vertex{Value: resp2}: {
			Vertex{Value: resp3}: []Rel{Realtime},
			Vertex{Value: resp4}: []Rel{Realtime},
		},
		Vertex{Value: resp3}: {},
		Vertex{Value: resp4}: {},
	}

	assert.Equal(t, g.Outs, dest)

	_, g, _ = RealtimeGraph(History{call1, resp1, call2, call4, resp2, call3, resp4, resp3, call5, resp5})

	dest = map[Vertex]map[Vertex][]Rel{
		Vertex{Value: resp1}: {
			Vertex{Value: resp4}: []Rel{Realtime},
			Vertex{Value: resp2}: []Rel{Realtime},
		},
		Vertex{Value: resp2}: {
			Vertex{Value: resp3}: []Rel{Realtime},
		},
		Vertex{Value: resp3}: {
			Vertex{Value: resp5}: []Rel{Realtime},
		},
		Vertex{Value: resp4}: {
			Vertex{Value: resp5}: []Rel{Realtime},
		},
		Vertex{Value: resp5}: {},
	}

	assert.Equal(t, g.Outs, dest)

}

// toJson is a debugging function, which can be used like:
// ```
//for k, v := range g.Outs {
//	fmt.Println(toJson(k), len(v))
//}
//
//fmt.Println()
//fmt.Println()
//
//for k, v := range dest {
//	fmt.Println(toJson(k), len(v))
//}
// ```
func toJson(v interface{}) string {
	s, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(s)
}