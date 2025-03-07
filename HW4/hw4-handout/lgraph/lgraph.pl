% A node is “in” the graph if it appears as a source or target in an edge.
node(G, N) :- edge(G, N, _, _).
node(G, N) :- edge(G, _, _, N).

% The predicate path works when the length K is given.
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


% The predicate path_seq works when the length K is not given.
% S is a sequence (list of labels) for a path from U to V in graph G.
path_seq(G, U, V, []) :-
    U = V,
    node(G, U).
path_seq(G, U, V, [L|Ls]) :-
    edge(G, U, L, W),
    path_seq(G, W, V, Ls).

% find_sequence/6: S is a sequence (of K labels) for a path from U to V in graph G1,
% and S is NOT a sequence for a path from U to V in graph G2.
%
% If K is provided (nonvar) we use path/5; otherwise we generate S with path_seq/4
% and then determine K as its length.
find_sequence(G1, G2, U, V, K, S) :-
    ( nonvar(K) ->
         path(G1, U, V, K, S)
    ;
         path_seq(G1, U, V, S),
         length(S, K)
    ),
    not(path_seq(G2, U, V, S)).
