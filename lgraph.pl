:- use_module(library(lists)).

% Base case: a path of length 0 from a node to itself with no labels.
path_sequence(_, U, U, 0, []).

% Recursive case: build a path of length K by extending a path of length K-1.
path_sequence(Graph, U, V, K, [L|Labels]) :-
    K > 0,
    edge(Graph, U, L, Next),
    K1 is K - 1,
    path_sequence(Graph, Next, V, K1, Labels).

% Main predicate to find sequences in G1 but not in G2.
find_sequence(G1, G2, U, V, K, S) :-
    K >= 0,
    path_sequence(G1, U, V, K, S),
    \+ path_sequence(G2, U, V, K, S).