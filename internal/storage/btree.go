// internal/storage/btree.go
//
// This file implements a simple B-tree (order = 4) for string keys that maps
// each key to one or more row indexes ([]int). The B-tree is used by HaruDB to
// accelerate equality lookups and prepare for future range queries.
//
// High-level design (read this first):
// - A B-tree is a multi-way, balanced search tree.
// - Every node has up to (order-1) keys and up to (order) children.
// - All leaves are at the same depth; the tree stays balanced via splits.
// - We store keys in sorted order inside nodes to support binary-search-like
//   navigation (here we use linear search for simplicity).
// - Leaf nodes keep values parallel to keys: each key maps to []int row indexes.
// - Internal nodes keep only keys and child pointers; values live only in leaves.
//
// What is implemented here:
// - Insert(key, rowIndex): O(log n) insertion with node splitting as needed.
// - GetEqual(key): O(log n) lookup that returns []int of row positions.
//
// Not implemented (future work):
// - Range search (e.g., BETWEEN) and ordered traversal APIs.
// - Deletion (we currently rebuild or append as needed in HaruDB flows).

package storage

// btreeOrder sets max children per node. order=4 => up to 3 keys per node.
const btreeOrder = 4

// btreeNode represents a single node in the B-tree.
type btreeNode struct {
	keys     []string     // sorted keys within the node
	children []*btreeNode // child pointers (len = len(keys)+1 for internal nodes)
	values   [][][]int    // only for leaves: parallel to keys; values[i] is []int list for keys[i]
	leaf     bool         // true if node is a leaf
}

// BTree is the main B-tree structure.
type BTree struct {
	root *btreeNode // root pointer
}

// NewBTree creates an empty B-tree with a single leaf root.
func NewBTree() *BTree {
	return &BTree{root: &btreeNode{leaf: true}}
}

// GetEqual returns the row index list for an exact key match.
func (t *BTree) GetEqual(key string) []int {
	n := t.root
	for {
		// Linear search within node (small node sizes keep this simple & fast)
		i := 0
		for i < len(n.keys) && key > n.keys[i] {
			i++
		}

		// If key matches within this node
		if i < len(n.keys) && key == n.keys[i] {
			if n.leaf {
				// values mirror keys one-to-one; take the list stored at this position
				if len(n.values) > i {
					// Flatten the [][]int (we keep one slice per duplicate insertion) to a single []int
					flat := make([]int, 0)
					for _, group := range n.values[i] {
						flat = append(flat, group...)
					}
					return flat
				}
				return nil
			}
			// Internal node: descend to the right child of the matching key
			n = n.children[i+1]
			continue
		}

		// If leaf and not found, it's a miss
		if n.leaf {
			return nil
		}

		// Internal node and not found here; descend to child i
		n = n.children[i]
	}
}

// Insert inserts a (key,rowIndex) into the B-tree.
func (t *BTree) Insert(key string, rowIndex int) {
	root := t.root
	// If root is full (has order-1 keys), split it and grow the tree height
	if len(root.keys) == btreeOrder-1 {
		newRoot := &btreeNode{leaf: false, children: []*btreeNode{root}}
		t.splitChild(newRoot, 0)
		t.root = newRoot
		t.insertNonFull(newRoot, key, rowIndex)
		return
	}
	// Otherwise, insert directly into non-full root
	t.insertNonFull(root, key, rowIndex)
}

// insertNonFull inserts into a node guaranteed to have spare capacity.
func (t *BTree) insertNonFull(n *btreeNode, key string, rowIndex int) {
	if n.leaf {
		// In a leaf: find insertion point
		i := 0
		for i < len(n.keys) && key > n.keys[i] {
			i++
		}
		// If key exists, append rowIndex to its value list
		if i < len(n.keys) && n.keys[i] == key {
			// Ensure parallel structure for values
			for len(n.values) < len(n.keys) {
				n.values = append(n.values, nil)
			}
			n.values[i] = append(n.values[i], []int{rowIndex})
			return
		}
		// Insert key at i and its value
		n.keys = append(n.keys, "")
		copy(n.keys[i+1:], n.keys[i:])
		n.keys[i] = key

		// Ensure values slice is aligned with keys length
		if !n.leaf {
			// defensive; leaf expected here
		} else {
			if len(n.values) < len(n.keys) {
				n.values = append(n.values, nil)
			}
			copy(n.values[i+1:], n.values[i:])
			n.values[i] = [][]int{{rowIndex}}
		}
		return
	}

	// Internal node: find child to descend into
	i := 0
	for i < len(n.keys) && key > n.keys[i] {
		i++
	}
	// If target child is full, split it first, then decide which child to go to
	if len(n.children[i].keys) == btreeOrder-1 {
		t.splitChild(n, i)
		// After split, decide which of the two children to descend into
		if key > n.keys[i] {
			i++
		}
	}
	// Recurse into the (now non-full) child
	t.insertNonFull(n.children[i], key, rowIndex)
}

// splitChild splits child c = n.children[i] into two nodes and promotes the middle key into n.
func (t *BTree) splitChild(n *btreeNode, i int) {
	// c is the full child to split
	c := n.children[i]
	mid := (btreeOrder - 1) / 2 // with order=4, mid=1 (2 keys => promote index 1)

	// Create new node that will receive the upper half of c's keys
	newNode := &btreeNode{leaf: c.leaf}

	// Move keys to newNode: keys after mid
	newNode.keys = append(newNode.keys, c.keys[mid+1:]...)
	c.keys = c.keys[:mid]

	// If leaf, move values parallel to keys
	if c.leaf {
		if len(c.values) > 0 {
			newNode.values = append(newNode.values, c.values[mid+1:]...)
			c.values = c.values[:mid]
		}
	} else {
		// If internal node, split children pointers accordingly
		newNode.children = append(newNode.children, c.children[mid+1:]...)
		c.children = c.children[:mid+1]
	}

	// Insert new child pointer into parent n at position i+1
	n.children = append(n.children, nil)
	copy(n.children[i+2:], n.children[i+1:])
	n.children[i+1] = newNode

	// Promote middle key into parent n at position i
	n.keys = append(n.keys, "")
	copy(n.keys[i+1:], n.keys[i:])
	n.keys[i] = c.keys[mid]

	// For leaves, values for the promoted key stay in the left child (c)
}
