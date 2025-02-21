(defun find-sequence (g1 g2 start target k)
  (labels (
           ;; Helper to check if a node exists in the graph
           ;; A node exists if (funcall graph node) does not return nil
           (node-exists-p (graph node)
             (not (null (funcall graph node))))

           ;; Helper to determine if there is a path following sequence S
           ;; from node to target in graph
           (can-follow-to (graph node S target)
             (if (null S)
                 ;; If sequence is empty, check if current node is target
                 (equal node target)
                 (let ((label (first S)))
                   (let ((edges (funcall graph node)))
                     (and edges
                          ;; Find all edges with the given label
                          (let ((next-nodes (mapcar #'second
                                                   (remove-if-not
                                                    (lambda (e) (equal (first e) label))
                                                    edges))))
                            ;; Recurse on each possible next node
                            (some (lambda (next)
                                    (can-follow-to graph next (rest S) target))
                                  next-nodes)))))))

           ;; DFS to find a sequence of length k from current to target in g1
           ;; sequence is the list of labels collected so far (in reverse order)
           (dfs (current k sequence)
             (if (= k 0)
                 (if (equal current target)
                     (let ((S (reverse sequence)))
                       ;; Check if S does not lead from start to target in g2
                       (if (not (can-follow-to g2 start S target))
                           (cons S t)
                           nil))
                     nil)
                 (let ((edges (funcall g1 current)))
                   (and edges
                        (some (lambda (edge)
                                (let ((label (first edge))
                                      (next (second edge)))
                                  (dfs next (- k 1) (cons label sequence))))
                              edges))))))
    ;; Handle k=0 case separately
    (if (= k 0)
        (if (and (equal start target)
                 (node-exists-p g1 start)
                 (not (node-exists-p g2 start)))
            (cons nil t)
            nil)
        (dfs start k nil))))
