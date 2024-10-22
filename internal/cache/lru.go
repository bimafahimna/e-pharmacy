package cache

type EvictAlgo interface {
	get(key string)
	set(key string)
}

type Node struct {
	key  string
	next *Node
	prev *Node
}

type LRU struct {
	capacity int
	counter  int
	cache    map[string](*Node)
	head     *Node
	tail     *Node
}

func constructor(capacity int) EvictAlgo {
	headNode := Node{}
	tailNode := Node{}
	headNode.next = &tailNode
	tailNode.prev = &headNode
	return &LRU{
		capacity,
		0,
		make(map[string](*Node)),
		&headNode,
		&tailNode,
	}
}

func (e *LRU) get(key string) {
	if node, exists := e.cache[key]; exists {
		removeNode(node)
		e.insertNodeToHead(node)
		return
	} else {
		return
	}
}

func (e *LRU) set(key string) {
	if node, exists := e.cache[key]; exists {
		removeNode(node)
		e.insertNodeToHead(node)
	} else {
		newNode := Node{key, nil, nil}
		e.insertNodeToHead(&newNode)
		e.cache[key] = &newNode

		if e.counter == e.capacity {
			evictedNode := e.tail.prev
			nextTail := evictedNode.prev
			nextTail.next = e.tail
			e.tail.prev = nextTail
			delete(e.cache, evictedNode.key)
		} else {
			e.counter += 1
		}
	}

}

func (e *LRU) insertNodeToHead(node *Node) {
	currHead := e.head.next
	node.next = currHead
	node.prev = e.head
	e.head.next = node
	currHead.prev = node
}

func removeNode(node *Node) {
	prevNode := node.prev
	nextNode := node.next
	prevNode.next = nextNode
	nextNode.prev = prevNode
}
