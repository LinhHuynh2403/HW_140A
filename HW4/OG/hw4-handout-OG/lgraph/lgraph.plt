:- initialization main.

main :-
    consult(['definitions.pl', 'lgraph.pl']),
    (show_coverage(run_tests) ; true),
    halt.

:- begin_tests(lgraph).

test(test_find_sequence_01, [fail]) :- find_sequence(g1, g2, 17, 0, 2, _).
test(test_find_sequence_02, [fail]) :- find_sequence(g1, g2, 7, 6, -100, _).
test(test_find_sequence_03, [set(S == [[]])]) :- find_sequence(g1, g2, 7, 7, 0, S).
test(test_find_sequence_04, [set(S == [[l]])]) :- find_sequence(g1, g2, 7, 6, 1, S).
test(test_find_sequence_05, [fail]) :- find_sequence(g1, g2, 0, 0, 0, _).
test(test_find_sequence_06, [fail]) :- find_sequence(g1, g2, 2, 0, 1, _).
test(test_find_sequence_07, [fail]) :- find_sequence(g1, g2, 4, 0, 3, _).
test(test_find_sequence_08, [set(S == [[e]])]) :- find_sequence(g2, g1, 2, 0, 1, S).
test(test_find_sequence_09, [set(S == [[h]])]) :- find_sequence(g2, g1, 4, 1, 1, S).
test(test_find_sequence_10, [fail]) :- find_sequence(g2, g1, 3, 6, 2, _).
test(test_find_sequence_11, [set(S == [[f, k]])]) :- find_sequence(g1, g2, 1, 4, 2, S).
test(test_find_sequence_12, [fail]) :- find_sequence(g1, g2, 0, 2, 2, _).
test(test_find_sequence_13, [fail]) :- find_sequence(g2, g1, 1, 2, 1, _).
test(test_find_sequence_14, [fail]) :- find_sequence(g2, g1, 0, 3, 3, _).
test(test_find_sequence_15, [set(S == [[a, b, f, g]])]) :- find_sequence(g2, g1, 0, 4, 4, S).
test(test_find_sequence_16, [set(S == [[m, l]])]) :- find_sequence(g1, g2, 4, 6, 2, S).
test(test_find_sequence_17, [set(S ==  [[h, f, j, l]])]) :- find_sequence(g2, g1, 4, 6, 4, S).
test(test_find_sequence_18, [set(S == [[e, b, f, j, l]])]) :-find_sequence(g2, g1, 2, 6, 5, S).
test(test_find_sequence_19, [set(S == [[b, f, k, m, l]])]) :- find_sequence(g1, g2, 0, 6, 5, S).
test(test_find_sequence_20, [set(S ==  [[a]])]) :- find_sequence(g5, g4, 14, 14, 1, S).
test(test_find_sequence_21, [fail]) :- find_sequence(g5, g4, 14, 14, 8, _).
test(test_find_sequence_22, [set(S ==  [[a, a, a]])]) :- find_sequence(g5, g4, 14, 14, 3, S).
test(test_find_sequence_23, [fail]) :- find_sequence(g4, g5, 11, 199, 1, _).
test(test_find_sequence_24, [set(S == [[y, n]])]) :- find_sequence(g4, g5, 11, 199, 2, S).
test(test_find_sequence_25, [fail]) :- find_sequence(g4, g5, 11, 199, 3, _).
test(test_find_sequence_26, [fail]) :- find_sequence(g4, g5, 200, 202, 2, _).
test(test_find_sequence_27, [set(S == [[x, a]])]) :- find_sequence(g5, g4, 200, 202, 2, S).
test(test_find_sequence_28, [set(S == [[c, e, b, f], [f, g, h, f]]) ]) :- find_sequence(g2, g1, 1, 3, 4, S).
test(test_find_sequence_29, [set(S == [[c, d, e], [c, e, a]])]) :-  find_sequence(g2, g1, 1, 0, 3, S).
test(test_find_sequence_30, [set(S == [[f, g], [n, o]])]) :- find_sequence(g3, g1, 1, 4, 2, S).
test(test_find_sequence_31, [set(S ==[[b, n, o], [a, a, o]])]) :- find_sequence(g3, g2, 0, 4, 3, S).
test(test_find_sequence_32, [set(S == [[d, d, d], [e, b, c]])]) :- find_sequence(g2, g1, 2, 2, 3, S).
test(test_find_sequence_33, [set(S == [[c, d, d, d, d],
                                        [c, d, e, b, c],
                                        [c, e, a, b, c],
                                        [c, e, b, c, d],
                                        [f, g, h, c, d]
                                        ])]) :- find_sequence(g2, g1, 1, 2, 5, S).
test(test_find_sequence_34, [set(S ==  [[c, d, d, d, e, b],
                                        [c, d, d, e, a, b],
                                        [c, d, e, a, a, b],
                                        [c, e, a, a, a, b],
                                        [c, e, b, c, e, b],
                                        [c, e, b, f, g, h],
                                        [f, g, h, c, e, b],
                                        [f, g, h, f, g, h]
                                        ])]) :- find_sequence(g2, g1, 1, 1, 6, S).

test(test_find_source_01, [fail]) :- find_sequence(g1, g2, _, 7, 5, [b, f, k, m, l]).
test(test_find_source_02, [fail]) :- find_sequence(g1, g2, _, 6, 6, [b, f, k, m, l]).
test(test_find_source_03, [set(Start == [0])]) :- find_sequence(g1, g2, Start, 6, 5, [b, f, k, m, l]).
test(test_find_source_04, [set(Start == [2])]) :- find_sequence(g2, g1, Start, 2, 1, [d]).
test(test_find_source_05, [fail]) :- find_sequence(g2, g1, _, 2, -2, [d]).
test(test_find_source_06, [fail]) :- find_sequence(g2, g1, _, -2, -100, [d]).
test(test_find_source_07, [fail]) :- find_sequence(g4, g5, _, 202, 2, [x, a]).
test(test_find_source_08, [set(Start == [300, 301])]) :- find_sequence(g4, g5, Start, 302, 1, [z]).

test(test_find_target_01, [fail]) :- find_sequence(g2, g1, -1, _, 4, [a, b, f, g]).
test(test_find_target_02, [fail]) :- find_sequence(g2, g1, 0, _, 5, [a, b, f, g]).
test(test_find_target_03, [fail]) :- find_sequence(g5, g4, 200, _, 2, [x, y]).
test(test_find_target_04, [set(Target == [4])]) :- find_sequence(g2, g1, 0, Target, 4, [a, b, f, g]).
test(test_find_target_05, [set(Target == [0])]) :-find_sequence(g2, g1, 2, Target, 3, [d, d, e]).
test(test_find_target_06, [set(Target == [2])]) :- find_sequence(g2, g1, 2, Target, 11, [d, d, d, e, b, f, g, h, c, d, d]).
test(test_find_target_07, [set(Target == [203])]) :- find_sequence(g5, g4, 200, Target, 1, [x]).
test(test_find_target_08, [set(Target == [301, 302])]) :- find_sequence(g5, g4, 300, Target, 1, [u]).

test(test_misc_01, [fail]) :- find_sequence(g1, g2, _, 7, _, [b, f, k, m, l]).
test(test_misc_02, [set(K == [5])]) :- find_sequence(g1, g2, 0, 6, K, [_, _, _, _, _]).
test(test_misc_03, [set(Start == [2])]) :- find_sequence(g2, g1, Start, 2, _, [_]).
test(test_misc_04, [set(Start == [300, 301])]) :- find_sequence(g4, g5, Start, 302, _, [_]).
test(test_misc_05, [set(K == [0])]) :- find_sequence(g1, g2, 7, 7, K, []).
test(test_misc_06, [set(Target == [2])]) :- find_sequence(g2, g1, 2, Target, _, [d, d, d, e, _, f, g, h, c, _, d]).
test(test_misc_07, [set(K == [4])]) :- find_sequence(g2, g1, 0, 4, K, [_,_,_,_]).

:- end_tests(lgraph).