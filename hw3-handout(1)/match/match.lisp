(defun match (pattern assertion)
  (cond
    ((and (null pattern) (null assertion)) t)  ; Both empty: match
    ((or (null pattern) (null assertion)) nil) ; One empty: no match
    ((equal (car pattern) '?)                  ; Handle '?'
     (and (consp assertion)                    ; Assertion has at least one element
          (match (cdr pattern) (cdr assertion))))
    ((equal (car pattern) '!)                  ; Handle '!'
     (loop for i from 1 to (length assertion)  ; Try all possible splits
           thereis (match (cdr pattern) (nthcdr i assertion))))
    (t                                         ; Regular atom comparison
     (and (equal (car pattern) (car assertion))
          (match (cdr pattern) (cdr assertion)))))
)