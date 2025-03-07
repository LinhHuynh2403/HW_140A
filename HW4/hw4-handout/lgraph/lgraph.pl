path_sequence(_, U, U, 0, []).

% Recursive case: build a path of length K by extending a path of length K-1.
path_sequence(Graph, U, V, K, [L|Labels]) :-
    K > 0,
    edge(Graph, U, L, Next),
    K1 is K - 1,
    path_sequence(Graph, Next, V, K1, Labels).

% Helper predicate to check if a sequence is not present in another graph.
sequence_not_in_graph(Graph, U, V, K, Sequence) :-
    \+ path_sequence(Graph, U, V, K, Sequence).

% Main predicate to find sequences in G1 that are not in G2.
find_sequence(G1, G2, U, V, K, S) :-
    path_sequence(G1, U, V, K, S),
    sequence_not_in_graph(G2, U, V, K, S).

find_sequence(G1, G2, U, V, 0, []) :-
    U = V,
    edge(G1, U, _, _),
    not(edge(G2, U, _, _).