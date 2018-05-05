# go-utility


## Stack/Queue
Stack/Queue package is general implementation of concurrent and blocking stack/queue.

### Blocking Stack/Queue
 - Pop: When non-empty, remove top element and returns. When empty, waits some other go routine to push into the stack/queue, whenever got something, it'll remove from top and return it from stack/queue.
 - Push: When stack/queue is maxed out, it'll wait for other go routine to free the stack/queue and whenever stack/queue got freed, it'll push into stack/queue. When stack/queue is free, it'll push into stack/queue.

 ### Concurrent Stack/Queue
 This is non-blocking implementation of Stack/Queue.
  - Pop: It is non blocking, is Queue/Stack is empty, return nil. Otherwise, returns the top element.
  - Push: Push it to the queue/stack. For queue, push to the back. For Stack, push on the top.