(defun find-sequence (g1 g2 u v k)
  (let ((sequences (find-sequences g1 u v k)))
    (check-sequences sequences g2 u v k)))

(defun find-sequences (graph u v k)
  (cond
    ;; Base case: If u == v and k == 0, return an empty sequence.
    ((and (eql u v) (zerop k)) '(nil))
    ;; Base case: If k < 0, return nil (invalid sequence).
    ((< k 0) nil)
    ;; Recursive case: Explore all edges from u.
    (t
     (let ((edges (funcall graph u)))
       (find-sequences-helper edges graph v (1- k))))))

(defun find-sequences-helper (edges graph v k)
  (cond
    ;; Base case: If no more edges, return nil.
    ((null edges) nil)
    ;; Recursive case: Process the current edge and explore further.
    (t
     (let* ((edge (car edges))
            (label (car edge))
            (next-node (cadr edge))
            ;; Find all sequences from next-node to v with length k.
            (sequences (find-sequences graph next-node v k)))
       ;; Prepend the current label to each sequence and combine results.
       (append (mapcar (lambda (seq) (cons label seq)) sequences)
              ;; Recurse on the remaining edges.
              (find-sequences-helper (cdr edges) graph v k))))))

(defun check-sequences (sequences g2 u v k)
  (cond
    ((and (null sequences) (equal u v) (zerop k) (null (funcall g2 u))) (cons nil t))
    ;; Base case: If no more sequences to check, return nil.
    ((null sequences) nil)
    ;; Check if the current sequence does not exist in g2.
    ((not (sequence-exists-in-graph g2 u v (car sequences) k))
     ;; If it doesn't exist, return the sequence wrapped in (seq . t).
     (cons (car sequences) t))
    ;; Otherwise, recurse on the remaining sequences.
    (t (check-sequences (cdr sequences) g2 u v k))))

(defun sequence-exists-in-graph (graph u v seq k)
  (let ((sequences (find-sequences graph u v k)))
    (sequence-exists-helper sequences seq)))

(defun sequence-exists-helper (sequences seq)
  (cond
    ;; Base case: If no more sequences, return nil.
    ((null sequences) nil)
    ;; If the current sequence matches, return t.
    ((equal seq (car sequences)) t)
    ;; Otherwise, recurse on the remaining sequences.
    (t (sequence-exists-helper (cdr sequences) seq))))