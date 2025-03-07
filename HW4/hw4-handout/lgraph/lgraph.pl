% A node is “in” the graph if it appears as a source or target in an edge.
node(G, N) :- edge(G, N, _, _).
node(G, N) :- edge(G, _, _, N).

% A path of length 0 exists from a node to itself provided the node is in the graph.
path(G, U, V, 0, []) :-
    node(G, U),
    U = V.

% A path of length K>0 exists from U to V if there is an edge from U to some intermediate W
% labeled L and a path of length K-1 from W to V with the rest of the labels.
path(G, U, V, K, [L|Ls]) :-
    K > 0,
    edge(G, U, L, W),
    K1 is K - 1,
    path(G, W, V, K1, Ls).

% find_sequence(G1, G2, U, V, K, S) is true if S is a sequence of K labels for a path from U to V in graph G1
% and S is NOT a sequence for a path from U to V in graph G2.
find_sequence(G1, G2, U, V, K, S) :-
    path(G1, U, V, K, S),
    not(path(G2, U, V, K, S)).
