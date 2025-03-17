package lgraph

import (
	"fmt"
	"reflect"
	"testing"
)

func mkCopy(edges []edge, exists bool) ([]edge, bool) {
	if exists {
		edgescopy := make([]edge, len(edges))
		copy(edgescopy, edges)
		return edgescopy, exists
	} else {
		return edges, exists
	}
}

/*
Graph G1:

	0 --a--> 0
	0 --b--> 1
	1 --c--> 2
	1 --f--> 3
	3 --k--> 4
	3 --j--> 7
	4 --m--> 7
	7 --l--> 6
*/
func g1(source node) ([]edge, bool) {
	graph := map[node][]edge{
		0: {
			{destination: 0, label: 'a'},
			{destination: 1, label: 'b'},
		},
		1: {
			{destination: 2, label: 'c'},
			{destination: 3, label: 'f'},
		},
		2: {},
		3: {
			{destination: 4, label: 'k'},
			{destination: 7, label: 'j'},
		},
		4: {
			{destination: 7, label: 'm'},
		},
		6: {},
		7: {
			{destination: 6, label: 'l'},
		},
	}
	edges, exists := graph[source]
	return mkCopy(edges, exists)
}

/*
Graph G2:

	0 --a--> 0
	0 --b--> 1
	1 --c--> 2
	1 --f--> 3
	2 --d--> 2
	2 --e--> 0
	3 --g--> 4
	3 --j--> 5
	4 --h--> 1
	5 --l--> 6
*/
func g2(source node) ([]edge, bool) {
	graph := map[node][]edge{
		0: {
			{destination: 0, label: 'a'},
			{destination: 1, label: 'b'},
		},
		1: {
			{destination: 2, label: 'c'},
			{destination: 3, label: 'f'},
		},
		2: {
			{destination: 2, label: 'd'},
			{destination: 0, label: 'e'},
		},
		3: {
			{destination: 4, label: 'g'},
			{destination: 5, label: 'j'},
		},
		4: {
			{destination: 1, label: 'h'},
		},
		5: {
			{destination: 6, label: 'l'},
		},
		6: {},
	}
	edges, exists := graph[source]
	return mkCopy(edges, exists)
}

/*
Graph G3:

	0 --a--> 0
	0 --b--> 1
	0 --o--> 4
	1 --n--> 0
	1 --c--> 2
	1 --f--> 3
	2 --d--> 2
	2 --e--> 0
	3 --g--> 4
	4 --h--> 1
*/
func g3(source node) ([]edge, bool) {
	graph := map[node][]edge{
		0: {
			{destination: 0, label: 'a'},
			{destination: 1, label: 'b'},
			{destination: 4, label: 'o'},
		},
		1: {
			{destination: 0, label: 'n'},
			{destination: 2, label: 'c'},
			{destination: 3, label: 'f'},
		},
		2: {
			{destination: 2, label: 'd'},
			{destination: 0, label: 'e'},
		},
		3: {
			{destination: 4, label: 'g'},
		},
		4: {
			{destination: 1, label: 'h'},
		},
	}
	edges, exists := graph[source]
	return mkCopy(edges, exists)
}

/*
Graph G4:

   11 --y--> 34
   11 --y--> 2
   34 --n--> 199
    2 --n--> 199
   14 --a--> 8
    8 --y--> 8 
    8 --a--> 14
  200 --x--> 201
  201 --y--> 202
*/
func g4(source node) ([]edge, bool) {
	graph := map[node][]edge{
		11:  {{34, 'y'}, {2, 'y'}},
		34:  {{199, 'n'}},
		2:   {{199, 'n'}},
		199: {},
		14:  {{8, 'a'}},
		8:   {{8, 'y'}, {14, 'a'}},
		200: {{201, 'x'}},
		201: {{202, 'y'}},
		202: {},
	}
	edges, exists := graph[source]
	return mkCopy(edges, exists)
}

/*
Graph G5:

   14 --a--> 14 
   11 --n--> 3
   11 --n--> 2
    3 --y--> 199
    2 --y--> 199
  200 --x--> 201
  200 --x--> 203
  201 --y--> 202
  203 --a--> 202
*/
func g5(source node) ([]edge, bool) {
	graph := map[node][]edge{
		14:  {{14, 'a'}},
		11:  {{3, 'n'}, {2, 'n'}},
		3:   {{199, 'y'}},
		2:   {{199, 'y'}},
		199: {},
		200: {{201, 'x'}, {203, 'x'}},
		201: {{202, 'y'}},
		202: {},
		203: {{202, 'a'}},
	}
	edges, exists := graph[source]
	return mkCopy(edges, exists)
}

func TestFindSequence(t *testing.T) {
	tests := []struct {
		testID                 string
		graph1                 LGraph
		graph2                 LGraph
		source, target         node
		sequenceLength         uint
		expectedSequenceExists bool
		expectedSequence       []rune
	}{
		{
			"test_1_00", g1, g2, 17, 0, 2,
			false, nil,
		},
		{
			"test_1_01", g1, g2, 7, 7, 0,
			true, []rune{},
		},
		{
			"test_1_02", g1, g2, 7, 6, 1,
			true, []rune{'l'},
		},
		{
			"test_1_03", g1, g2, 0, 0, 0,
			false, nil,
		},
		{
			"test_1_04", g1, g2, 2, 0, 1,
			false, nil,
		},
		{
			"test_1_05", g1, g2, 4, 0, 3,
			false, nil,
		},
		{
			"test_1_06", g2, g1, 2, 0, 1,
			true, []rune{'e'},
		},
		{
			"test_1_07", g2, g1, 4, 1, 1,
			true, []rune{'h'},
		},
		{
			"test_1_08", g2, g1, 3, 6, 2,
			false, nil,
		},
		{
			"test_1_09", g1, g2, 1, 4, 2,
			true, []rune{'f', 'k'},
		},
		{
			"test_1_10", g1, g2, 0, 2, 2,
			false, nil,
		},
		{
			"test_1_11", g2, g1, 1, 2, 1,
			false, nil,
		},
		{
			"test_1_12", g2, g1, 0, 3, 3,
			false, nil,
		},
		{
			"test_1_13", g2, g1, 0, 4, 4,
			true, []rune{'a', 'b', 'f', 'g'},
		},
		{
			"test_1_14", g1, g2, 4, 6, 2,
			true, []rune{'m', 'l'},
		},
		{
			"test_1_15", g2, g1, 4, 6, 4,
			true, []rune{'h', 'f', 'j', 'l'},
		},
		{
			"test_1_16", g2, g1, 2, 6, 5,
			true, []rune{'e', 'b', 'f', 'j', 'l'},
		},
		{
			"test_1_17", g1, g2, 0, 6, 5,
			true, []rune{'b', 'f', 'k', 'm', 'l'},
		},
		{
			"test_1_18", g5, g4, 14, 14, 1,
			true, []rune{'a'},
		},
		{
			"test_1_19", g5, g4, 14, 14, 8,
			false, nil,
		},
		{
			"test_1_20", g5, g4, 14, 14, 3,
			true, []rune{'a', 'a', 'a'},
		},
		{
			"test_1_21", g4, g5, 11, 199, 1,
			false, nil,
		},
		{
			"test_1_22", g4, g5, 11, 199, 2,
			true, []rune{'y', 'n'},
		},
		{
			"test_1_23", g4, g5, 11, 199, 3,
			false, nil,
		},
		{
			"test_1_24", g4, g5, 200, 202, 2,
			false, nil,
		},
		{
			"test_1_25", g5, g4, 200, 202, 2,
			true, []rune{'x', 'a'},
		},
	}

	for _, test := range tests {
		func() {
			defer func() {
				if recover() != nil {
					t.Errorf("FindSequence panicked on (%s, %d, %d, %d)",
						test.testID, test.source, test.target, test.sequenceLength)
				}
			}()
			ans_sequence, ans_sequence_exists := FindSequence(test.graph1, test.graph2, test.source, test.target, test.sequenceLength)
			if ans_sequence_exists != test.expectedSequenceExists || !reflect.DeepEqual(ans_sequence, test.expectedSequence) {
				formatSequence := func(val []rune) string {
					if val == nil {
						return "nil"
					}
					return fmt.Sprintf("%q", val)
				}
				message := fmt.Sprintf(
					"FindSequence failed on (%s, %d, %d, %d); expected: (%s, %t), got (%s, %t).",
					test.testID, test.source, test.target, test.sequenceLength,
					formatSequence(test.expectedSequence), test.expectedSequenceExists,
					formatSequence(ans_sequence), ans_sequence_exists,
				)
				t.Errorf(message)
			}
		}()
	}
}

func TestFindSequenceTwoSolutions(t *testing.T) {
	tests := []struct {
		testID                 string
		graph1                 LGraph
		graph2                 LGraph
		source, target         node
		sequenceLength         uint
		expectedSequenceExists bool
		expectedSequence1      []rune
		expectedSequence2      []rune
	}{
		{
			"test_2_00", g2, g1, 1, 3, 4,
			true, []rune{'c', 'e', 'b', 'f'}, []rune{'f', 'g', 'h', 'f'},
		},
		{
			"test_2_01", g2, g1, 1, 0, 3,
			true, []rune{'c', 'd', 'e'}, []rune{'c', 'e', 'a'},
		},
		{
			"test_2_02", g3, g1, 1, 4, 2,
			true, []rune{'f', 'g'}, []rune{'n', 'o'},
		},
		{
			"test_2_03", g3, g2, 0, 4, 3,
			true, []rune{'b', 'n', 'o'}, []rune{'a', 'a', 'o'},
		},
		{
			"test_2_04", g2, g1, 2, 2, 3,
			true, []rune{'e', 'b', 'c'}, []rune{'d', 'd', 'd'},
		},
	}
	for _, test := range tests {
		func() {
			defer func() {
				if recover() != nil {
					t.Errorf("FindSequence panicked on (%s, %d, %d, %d)",
						test.testID, test.source, test.target, test.sequenceLength)
				}
			}()
			ans_sequence, ans_sequence_exists := FindSequence(test.graph1, test.graph2, test.source, test.target, test.sequenceLength)
			if !ans_sequence_exists || (!reflect.DeepEqual(ans_sequence, test.expectedSequence1) && !reflect.DeepEqual(ans_sequence, test.expectedSequence2)) {
				formatSequence := func(val []rune) string {
					if val == nil {
						return "nil"
					}
					return fmt.Sprintf("%q", val)
				}
				message := fmt.Sprintf(
					"FindSequence failed on (%s, %d, %d, %d); expected either (%q, %t) or (%q, %t), got (%q, %t).",
					test.testID, test.source, test.target, test.sequenceLength,
					formatSequence(test.expectedSequence1), test.expectedSequenceExists,
					formatSequence(test.expectedSequence2), test.expectedSequenceExists,
					formatSequence(ans_sequence), ans_sequence_exists,
				)
				t.Errorf(message)
			}
		}()
	}
}

func TestFindSequenceFiveSolutions(t *testing.T) {
	tests := []struct {
		testID                 string
		graph1                 LGraph
		graph2                 LGraph
		source, target         node
		sequenceLength         uint
		expectedSequenceExists bool
		expectedSequence1      []rune
		expectedSequence2      []rune
		expectedSequence3      []rune
		expectedSequence4      []rune
		expectedSequence5      []rune
	}{
		{
			"test_5_00", g2, g1, 1, 2, 5,
			true,
			[]rune{'c', 'd', 'd', 'd', 'd'},
			[]rune{'c', 'd', 'e', 'b', 'c'},
			[]rune{'c', 'e', 'a', 'b', 'c'},
			[]rune{'c', 'e', 'b', 'c', 'd'},
			[]rune{'f', 'g', 'h', 'c', 'd'},
		},
	}
	for _, test := range tests {
		func() {
			defer func() {
				if recover() != nil {
					t.Errorf("FindSequence panicked on (%s, %d, %d, %d)",
						test.testID, test.source, test.target, test.sequenceLength)
				}
			}()
			ans_sequence, ans_sequence_exists := FindSequence(test.graph1, test.graph2, test.source, test.target, test.sequenceLength)
			if !ans_sequence_exists ||
				(!reflect.DeepEqual(ans_sequence, test.expectedSequence1) &&
					!reflect.DeepEqual(ans_sequence, test.expectedSequence2) &&
					!reflect.DeepEqual(ans_sequence, test.expectedSequence3) &&
					!reflect.DeepEqual(ans_sequence, test.expectedSequence4) &&
					!reflect.DeepEqual(ans_sequence, test.expectedSequence5)) {
				
				formatSequence := func(val []rune) string {
					if val == nil {
						return "nil"
					}
					return fmt.Sprintf("%q", val)
				}
				message := fmt.Sprintf(
					"FindSequence failed on (%s, %d, %d, %d); expected any one of the following results: (%s, %t) or (%s, %t) or (%s, %t) or (%s, %t) or (%s, %t), "+
						"got (%s, %t).",
					test.testID, test.source, test.target, test.sequenceLength,
					formatSequence(test.expectedSequence1), test.expectedSequenceExists,
					formatSequence(test.expectedSequence2), test.expectedSequenceExists,
					formatSequence(test.expectedSequence3), test.expectedSequenceExists,
					formatSequence(test.expectedSequence4), test.expectedSequenceExists,
					formatSequence(test.expectedSequence5), test.expectedSequenceExists,
					formatSequence(ans_sequence), ans_sequence_exists,
				)
				t.Errorf(message)
			}
		}()
	}
}

func TestFindSequenceEightSolutions(t *testing.T) {
	tests := []struct {
		testID                 string
		graph1                 LGraph
		graph2                 LGraph
		source, target         node
		sequenceLength         uint
		expectedSequenceExists bool
		expectedSequence1      []rune
		expectedSequence2      []rune
		expectedSequence3      []rune
		expectedSequence4      []rune
		expectedSequence5      []rune
		expectedSequence6      []rune
		expectedSequence7      []rune
		expectedSequence8      []rune
	}{
		{
			"test_8_00", g2, g1, 1, 1, 6,
			true,
			[]rune{'c', 'd', 'd', 'd', 'e', 'b'},
			[]rune{'c', 'd', 'd', 'e', 'a', 'b'},
			[]rune{'c', 'd', 'e', 'a', 'a', 'b'},
			[]rune{'c', 'e', 'a', 'a', 'a', 'b'},
			[]rune{'c', 'e', 'b', 'c', 'e', 'b'},
			[]rune{'c', 'e', 'b', 'f', 'g', 'h'},
			[]rune{'f', 'g', 'h', 'c', 'e', 'b'},
			[]rune{'f', 'g', 'h', 'f', 'g', 'h'},
		},
	}
	for _, test := range tests {
		func() {
			defer func() {
				if recover() != nil {
					t.Errorf("FindSequence panicked on (%s, %d, %d, %d)",
						test.testID, test.source, test.target, test.sequenceLength)
				}
			}()
			ans_sequence, ans_sequence_exists := FindSequence(test.graph1, test.graph2, test.source, test.target, test.sequenceLength)
			if !ans_sequence_exists ||
				(!reflect.DeepEqual(ans_sequence, test.expectedSequence1) &&
					!reflect.DeepEqual(ans_sequence, test.expectedSequence2) &&
					!reflect.DeepEqual(ans_sequence, test.expectedSequence3) &&
					!reflect.DeepEqual(ans_sequence, test.expectedSequence4) &&
					!reflect.DeepEqual(ans_sequence, test.expectedSequence5) &&
					!reflect.DeepEqual(ans_sequence, test.expectedSequence6) &&
					!reflect.DeepEqual(ans_sequence, test.expectedSequence7) &&
					!reflect.DeepEqual(ans_sequence, test.expectedSequence8)) {

				formatSequence := func(val []rune) string {
					if val == nil {
						return "nil"
					}
					return fmt.Sprintf("%q", val)
				}
				message := fmt.Sprintf(
					"FindSequence failed on (%s, %d, %d, %d); expected any one of the following results: (%s, %t) or (%s, %t) or (%s, %t) or (%s, %t) or (%s, %t) or (%s, %t) or (%s, %t) or (%s, %t), "+
						"got (%s, %t).",
					test.testID, test.source, test.target, test.sequenceLength,
					formatSequence(test.expectedSequence1), test.expectedSequenceExists,
					formatSequence(test.expectedSequence2), test.expectedSequenceExists,
					formatSequence(test.expectedSequence3), test.expectedSequenceExists,
					formatSequence(test.expectedSequence4), test.expectedSequenceExists,
					formatSequence(test.expectedSequence5), test.expectedSequenceExists,
					formatSequence(test.expectedSequence6), test.expectedSequenceExists,
					formatSequence(test.expectedSequence7), test.expectedSequenceExists,
					formatSequence(test.expectedSequence8), test.expectedSequenceExists,
					formatSequence(ans_sequence), ans_sequence_exists,
				)
				t.Errorf(message)
			}
		}()
	}
}