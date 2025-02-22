(defun match (pattern assertion)
(cond
  ((and (null pattern) (null assertion)) t)  ; Both empty: match
  ((or (null pattern) (null assertion)) nil) ; One empty: no match
  ((equal (car pattern) '?)                  ; Handle '?'
    (and (consp assertion)                    ; Assertion has at least one element
        (match (cdr pattern) (cdr assertion))))
  ((equal (car pattern) '!)                  ; Handle '!'
    (or (and (consp assertion) 
            (match (cdr pattern) (cdr assertion)))  ; Consume one element
        (and (consp assertion) 
            (match pattern (cdr assertion)))))      ; Try skipping one element
  (t                                         ; Regular atom comparison
    (and (equal (car pattern) (car assertion))
        (match (cdr pattern) (cdr assertion))))))
