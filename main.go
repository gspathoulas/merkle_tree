package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

// struct to support nodes of Merkle tree
type node struct {
	value       []byte //the hash value of the node
	start_index int    //the start of the range of elements indeces that are represented by the descedants of the node
	end_index   int    //the end of the range of elements indeces that are represented by the descedants of the node
	left_child  int    // the id of left child node of this node
	right_child int    // the id of right child node of this node
}

//struct of the Merkle tree
type Merkle_tree struct {
	elements []string // a slice containing all elements
	tree     []node   // a slice containing all nodes of the Merkle tree
	root     node     // the root node of the Merkle tree
}

// funcion to initialise the Merkle tree by adding the elements
func (mt *Merkle_tree) init(elements []string) {
	mt.elements = elements
}

// function to print the elements
func (mt Merkle_tree) print_elements() {
	for i := 0; i < len(mt.elements); i++ {
		fmt.Println(mt.elements[i])
	}
}

//function to print the nodes of the tree
func (mt Merkle_tree) print_nodes() {
	for i := 0; i < len(mt.tree); i++ {
		fmt.Println(i, mt.tree[i])
	}
}

// function to build theMerkle tree upon the elements that have been added
func (mt *Merkle_tree) build_tree() {
	var nodes []node     // all nodes are added to the nodes slice
	nodes_pop := 0       // keep track of all nodes population
	level_nodes_pop := 0 // keep track of nodes population for each level of the tree

	// hash all elements to create the lower level of the tree
	for i := 0; i < len(mt.elements); i++ {
		nodes_pop++
		level_nodes_pop++
		var n node
		n.start_index = i
		n.end_index = i
		n.value = hash_string(mt.elements[i])
		nodes = append(nodes, n)
	}

	// if elements of the level are odd in number add one more element by copying the most left element of the level
	if level_nodes_pop%2 == 1 {
		nodes_pop++
		level_nodes_pop++
		var n node = nodes[nodes_pop-2]
		n.start_index = n.start_index + 1
		n.end_index = n.end_index + 1
		n.left_child = -1  // it is a duplicate node... no child nodes
		n.right_child = -1 // it is a duplicate node... no child nodes
		nodes = append(nodes, n)
	}

	// loop and build level upon level until the last built level has only one element
	for {
		// calculate the range of indices for the last level
		last_level_start := nodes_pop - level_nodes_pop
		last_level_stop := nodes_pop - 1

		// start counting from 0 for the new level
		level_nodes_pop = 0

		//take nodes of the previous level in pairs and concatenate and hash to create the nodes of the new layer
		for i := last_level_start; i < last_level_stop; i = i + 2 {
			level_nodes_pop++
			nodes_pop++
			var n node
			n.start_index = nodes[i].start_index
			n.end_index = nodes[i+1].end_index
			n.value = hash_node(append(nodes[i].value, nodes[i+1].value...))
			n.left_child = i
			n.right_child = i + 1
			nodes = append(nodes, n)
		}

		// if the last level has only one element stop... root of the Merkle tree has been found
		if level_nodes_pop == 1 {
			break
		}

		// if elements of the level are odd in number add one more element by copying the most left element of the level
		if level_nodes_pop%2 == 1 {
			nodes_pop++
			level_nodes_pop++
			var n node = nodes[nodes_pop-2]
			ran := n.end_index - n.start_index + 1
			n.start_index = n.start_index + ran
			n.end_index = n.end_index + ran
			n.left_child = -1  // it is a duplicate node... no child nodes
			n.right_child = -1 // it is a duplicate node... no child nodes
			nodes = append(nodes, n)
		}
	}

	// store the tree and the root to the Merkle tree struct
	mt.tree = nodes
	mt.root = nodes[nodes_pop-1]
}

// function to generate a proof of inclusion for the element at a given index of elements in the construction
func (mt Merkle_tree) get_proof(index int) [][]byte {

	running_node := mt.root // start from the root
	var proof [][]byte      // create a a 2-dimensional slice to hold hashes that will make up the proof

	// loop by moving down through the levels of the tree until we find the element we are looking for
	for {
		// if the running node is a leaf node stop
		if running_node.start_index == running_node.end_index {
			break
		}

		// check the descedants of the node and select the one to the subtree of which the element belongs to
		// make the chosen node the running one
		// keep the hash value of fthe other
		l := mt.tree[running_node.left_child]
		r := mt.tree[running_node.right_child]
		if index >= l.start_index && index <= l.end_index {
			proof = append(proof, r.value)
			running_node = l
		} else {
			proof = append(proof, l.value)
			running_node = r
		}
	}

	//reverse the stored hashes as those have to be used bottom up
	for i, j := 0, len(proof)-1; i < j; i, j = i+1, j-1 {
		proof[i], proof[j] = proof[j], proof[i]
	}

	return proof
}

//function to validate a proof for a given index
func (mt Merkle_tree) validate_proof(index int, proof [][]byte) bool {

	element := mt.elements[index]
	hash := hash_string(element) // hash the element

	//loop throught the hashes of the proofs to concatenate with the element and hash
	for i := 0; i < len(proof); i++ {
		//according to the index of a node in the level of the tree it belongs to we decide the concatenation order
		if index%2 == 0 {
			concat := append(hash, proof[i]...)
			hash = hash_node(concat)
		} else {
			concat := append(proof[i], hash...)
			hash = hash_node(concat)
		}
		index = index / 2 //calculate the index of the element for the next level
	}

	// check if the result is equal to the hash of the root of the Merkle tree to decide the validity of proof
	return bytes.Equal(hash, mt.root.value)
}

//function to add an element

func (mt *Merkle_tree) add_element(s string) {

	//to be added

}

func (mt *Merkle_tree) update_element(ind int, s string) {

	//to be added

}

// function to hash a string
func hash_string(s string) []byte {
	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return bs
}

//function to hash a bytes array
func hash_node(b []byte) []byte {
	h := sha256.New()
	h.Write(b)
	bs := h.Sum(nil)
	return bs
}

// In the main function set a list of strings
// Create a Merkle Tree
// Initialise it with the strings in the list
// Build the Merkle Tree
// Generatea proof for a specific element (through the use of its index)
// Validate the generated proof

func main() {
	arr := []string{"John", "Lily", "Roy", "Suzie", "Jane", "kane"}

	mt := Merkle_tree{}
	mt.init(arr)
	mt.print_elements()
	mt.build_tree()
	mt.print_nodes()
	proof := mt.get_proof(2)
	fmt.Println(mt.validate_proof(2, proof))

}
