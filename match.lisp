(defun match (pattern assertion)
  (cond
    ((and (null pattern) (null assertion)) t)
    ((null pattern) nil)
    ((null assertion) nil)
    (t
     (let ((p (car pattern))
           (a (car assertion)))
       (cond
         ((not (or (eq p '?) (eq p '!)))
          (if (eq p a)
              (match (cdr pattern) (cdr assertion))
              nil))
         ((eq p '?)
          (match (cdr pattern) (cdr assertion)))
         ((eq p '!)
          (let ((len (length assertion)))
            (loop for k from 1 to len
                  when (match (cdr pattern) (nthcdr k assertion))
                  return t
                  finally (return nil))))
         (t nil))))))