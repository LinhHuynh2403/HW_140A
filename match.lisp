
(defun len (a)
  (apply '+ (mapcar (lambda (n) 1) a))
)

(defun regMatch (pattern assertion)
   (cond
      ((and (eq 0 (len pattern)) (not (eq 0 (len assertion)))) ;; len miss match
         (eval nil)
      )
      ((and (not (eq 0 (len pattern))) (eq 0 (len assertion))) ;; len miss match
         (eval nil)
      )
      ((and (eq 0 (len pattern)) (eq 0 (len assertion))) ;; end of list
         (eval t)
      )
      ((eq (car pattern) (car assertion)) ;; match
         (checkMatch (cdr pattern) (cdr assertion))
      )
      ((not (eq (car pattern) (car assertion))) ;; miss match
         (eval nil)
      )
   )
)

(defun queMatch (pattern assertion)
   (cond
      ((and (not (null pattern)) (not (null assertion)))
         (checkMatch (cdr pattern) (cdr assertion))
      )
      ((or (null pattern) (null assertion))
         (eval nil)
      )
   )
)

(defun exMatch (pattern assertion stopWord)
   (cond
      ((null pattern) ;; ! is the last symbol in pattern
         (eval t)
      )
      ((eq '! (car pattern)) ;;no last !
         (setq newP (cons '? (cdr pattern)))
         (queMatch newP assertion)
      )
      ((eq '? (car pattern)) ;; next is ?
         (setq newP (cons '? (cdr pattern)))
         (queMatch newP assertion)
      )
      ((eq '! stopWord)
         (exMatch (cdr pattern) assertion (car (cdr pattern)))
      )
      ((not (eq '! stopWord))
         (cond
            ((eq stopWord (car assertion))
               (checkMatch pattern assertion)
            )
            ((not (eq stopWord (car assertion)))
               (exMatch pattern (cdr assertion) stopWord)
            )
         )
      )
   )
)

(defun checkMatch (pattern assertion) ;; deal with normal and ?

   (setq q "?")
   (setq e "!")

   (cond
      ((and (not (eq '? (car pattern))) (not (eq '! (car pattern))))
         (regMatch pattern assertion)
      )
      ((eq '? (car pattern))
         (queMatch pattern assertion)
      )
      ((eq '! (car pattern))
         (exMatch (cdr pattern) assertion (car (cdr pattern)))
      )
   )
)


(defun match (pattern assertion)
   ;; TODO: incomplete function.
   ;; The next line should not be in your solution.
   (checkMatch pattern assertion)
)

