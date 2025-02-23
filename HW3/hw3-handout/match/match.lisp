; You may define helper functions here

(defun match (pattern assertion)
  ;; TODO: incomplete function. 
  ;; The next line should not be in your solution.
  ;; (list 'incomplete)
  (cond
    ;; Check for exact match when no special atoms are present
    ((null pattern) (null assertion))

    ;; Handling pattern containing ?
    ((equal (car pattern) '?)
     (and (not (null assertion))
          (match (cdr pattern) (cdr assertion))))

    ;; Handling pattern containing !
    ((equal (car pattern) '!)
     (or (match (cdr pattern) (cdr assertion))
         (and (not (null assertion))
              (match pattern (cdr assertion)))))

    ;; Matching corresponding elements
    ((equal (car pattern) (car assertion))
     (match (cdr pattern) (cdr assertion)))

    ;; No match
    (t nil)))