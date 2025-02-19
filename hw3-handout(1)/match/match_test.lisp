;;;;;;;;;;;;;;;;;;;;;;;;
;;; pre-testing prep ;;;
;;;;;;;;;;;;;;;;;;;;;;;;

(load "../lisp-unit.lisp")

(use-package :lisp-unit)

(load "match.lisp")

(remove-tests :all)

(setq *print-failures* t)

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;
;;; match test definitions ;;;
;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(define-test test-match
    (assert-equal T (match nil nil))
    (assert-equal NIL (match nil '( a)))
    (assert-equal NIL (match '(a) nil))
    (assert-equal NIL (match '(a) '(b)))
    (assert-equal NIL (match '(a b c) '(a b d)))
    (assert-equal T (match '(?) '(a)))
    (assert-equal NIL (match '(?) nil))
    (assert-equal T (match '(!) '(a)))
    (assert-equal T (match '(!) '(a b c)))
    (assert-equal T (match '(! c) '(a b c)))
    (assert-equal NIL (match '(! d) '(a b c)))
    (assert-equal T (match '(? ! c) '(a a b c)))
    (assert-equal T (match '(? ! c d) '(a a b c d)))
    (assert-equal NIL (match '(? ! c d) '(a a b c e)))
    (assert-equal NIL (match '(? ? ! ? A) '(A A A A)))
    (assert-equal T (match '(? ? ! ? A !) '(A A A A A NIL)))
    (assert-equal T (match '(! B ? B A B A ! ? !) '(B B A B A B A B A B A)))
    (assert-equal NIL (match '(! B ? A B A B ! ? !) '(B B A B A B A B A B)))
)

;;;;;;;;;;;;;;;;;
;;; run tests ;;;
;;;;;;;;;;;;;;;;;

(run-tests :all)
