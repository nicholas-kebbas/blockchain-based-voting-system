package p1

import (
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"reflect"
	"strings"
)

/* Place into leaf and extension node */
type Flag_value struct {
	encoded_prefix []uint8
	value string
}

/* Only use branch value if Branch type Node */
/* Don't use flag_value if Branch type Node */
/* If leaf, flag_value is the actual value of the key value pair */
/* If extension, flag_value is the hashed string of the node branch or leaf node it's pointing to */
type Node struct {
	node_type int // 0: Null, 1: Branch, 2: Ext or Leaf
	branch_value [17]string
	flag_value Flag_value
}

// MerklePatriciaTrie class has map of strings called db
// and root of string value. This is the encoded string value I think
// the db maps the string to a node
type MerklePatriciaTrie struct {
	db map[string]Node
	stringDb map[string]string
	root string
}


func (mpt *MerklePatriciaTrie) GetRoot() string {
	return mpt.root
}

func (mpt *MerklePatriciaTrie) GetDb() map[string]Node {
	return mpt.db
}

func (mpt *MerklePatriciaTrie) GetStringDb() map[string]string {
	return mpt.stringDb
}

func (mpt *MerklePatriciaTrie) buildNodeExt(commonPath []uint8, branchHashedValue string) Node {
	newNode := Node{}
	newNode.node_type = 2
	newNode.flag_value.encoded_prefix = compact_encode(commonPath)
	newNode.flag_value.value = branchHashedValue
	return newNode
}
func (mpt *MerklePatriciaTrie) buildNodeLeaf(path []uint8, new_value string) Node {
	newNode := Node{}
	newNode.node_type = 2
	path = append(path, 16)
	newNode.flag_value = Flag_value{compact_encode(path), new_value}
	return newNode
}
func (mpt *MerklePatriciaTrie) buildNodeBranch() Node {
	newNode := Node{}

	return newNode
}
// takes as input a string key
// returns a string. The Value of the node (leaf.value or branch.branch_value)
func (mpt *MerklePatriciaTrie) Get(key string) (string, error) {
	/* Convert the string to a hex array */
	byteArray := []uint8(key)
	hexArray := []uint8{}
	/* Then, encode the hex array */

	for i := 0; i < len(byteArray); i+=1 {
		hex1 := byteArray[i]/16
		hex2 := byteArray[i]%16
		hexArray = append(hexArray, hex1)
		hexArray = append(hexArray, hex2)
	}
	commonPath := []uint8{}

	nodeThatContainsValue := mpt.getHelper(mpt.root, commonPath, hexArray)
	return mpt.db[nodeThatContainsValue].flag_value.value, errors.New("path_not_found")
}

func (mpt *MerklePatriciaTrie) getHelper(parentHash string, commonPath, newKey []uint8) string {

	if len(mpt.db) == 0 {
		return ""
	}

	/* If there's nothing left in the new key, we found it, so return whatever value is at the common path */

	currNode := mpt.db[parentHash]
	if currNode.node_type == 0 {
		return ""
	}

	/* Current is branch  */
	if currNode.node_type == 1 {
		/* Check what's at currNode.branch_value[16] */
		if len(newKey) == 0 {
			return mpt.getHelper(currNode.branch_value[16], commonPath, newKey)
		} else {
			/* Slice off for index */
			if len(currNode.branch_value[newKey[0]]) > 0 {
				indexToCheck := newKey[0]
				newKey = newKey[1:]
				returnValue := mpt.getHelper(currNode.branch_value[indexToCheck], commonPath, newKey)
				return returnValue
			}
			return ""
		}
	}

	if currNode.node_type == 2 {
		arrayToCheckAgainst := compact_decode(currNode.flag_value.encoded_prefix)
		if !(is_ext_node(currNode.flag_value.encoded_prefix)) {
			commonPath := []uint8{}
			commonPath = mpt.findCommonPath(newKey, arrayToCheckAgainst, commonPath, newKey)
			remainingPath := newKey[len(commonPath):]
			/* It's a leaf so we can't recurse down. Need to check if remaining path */
			if len(remainingPath) == 0 {
				return currNode.hash_node()
			}

		} else if is_ext_node(currNode.flag_value.encoded_prefix) {
			commonPath := []uint8{}
			commonPath = mpt.findCommonPath(newKey, arrayToCheckAgainst, commonPath, newKey)
			/* Recurse Down */
			remainingPath := newKey[len(commonPath):]
			returnValue := mpt.getHelper(currNode.flag_value.value, commonPath, remainingPath)
			return returnValue
		}
	}

	return ""
}

func (mpt *MerklePatriciaTrie) findCommonPath(hexArray []uint8, arrayInTree []uint8, commonPath []uint8, remainingPath []uint8) []uint8 {

	commonArray := commonPath
	if len(hexArray) == 0 || len(arrayInTree) == 0 {
		return commonArray
	}

	if hexArray[0] == arrayInTree[0] {
		commonPath = append(commonPath, hexArray[0])
		remainingPath = hexArray[1:]
		return mpt.findCommonPath(hexArray[1:], arrayInTree[1:], commonPath, remainingPath)
	}
	return commonPath
}
// Insert a key, value pair into the trie
func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {

	mpt.stringDb[key] = new_value
	/* Convert the string to a hex array */
	byteArray := []uint8(key)
	hexArray := []uint8{}
	/* Then, encode the hex array */

	for i := 0; i < len(byteArray); i+=1 {
		hex1 := byteArray[i]/16
		hex2 := byteArray[i]%16
		hexArray = append(hexArray, hex1)
		hexArray = append(hexArray, hex2)
	}

	commonPath := []uint8{}
	root := mpt.insertHelper(mpt.root, commonPath, hexArray, new_value)
	mpt.root = root
}

/* Traverse the trie and get to the place you need to be */
/* This needs to check the type of the node, then go down to the next node */
/* Need to check the prefix (0,1,2,3) and not 0,16,32,48 because decode removes that */

func createNewRoot(mpt *MerklePatriciaTrie, newKey []uint8, new_value string) string {
	newNode := Node{}
	newKey = append(newKey, 16)
	newNode.node_type = 2
	/* Set the flag_value parameters */
	newNode.flag_value = Flag_value{compact_encode(newKey), new_value}
	hashedValue := newNode.hash_node()
	mpt.db[hashedValue] = newNode
	mpt.root = hashedValue
	return hashedValue
}

func (mpt *MerklePatriciaTrie) insertHelper(parentHash string, commonPath, newKey []uint8, new_value string) string {

	/* CurrNode is empty, so empty tree */
	if len(mpt.db) == 0 {
		newRootValue := createNewRoot(mpt, newKey, new_value)
		return newRootValue
	}

	/* Handle case for empty key */
	if len(newKey) == 0 {
		newNode := mpt.db[parentHash]
		if newNode.node_type == 2 {
			if !is_ext_node(newNode.flag_value.encoded_prefix) {
				newNode.flag_value.value = new_value
				newHash := newNode.hash_node()
				mpt.db[newHash] = newNode
				return newHash
			}
		} else {
			newNode := Node{}
			newKey = append(newKey, 16)
			newNode.node_type = 2
			/* Set the flag_value parameters */
			newNode.flag_value = Flag_value{compact_encode(newKey), new_value}
			hashedValue := newNode.hash_node()
			mpt.db[hashedValue] = newNode
			return hashedValue
		}
	}

	currNode := mpt.db[parentHash]

	if currNode.node_type == 0 {
		newNode := Node{}
		newKey = append(newKey, 16)
		newNode.node_type = 2
		/* Set the flag_value parameters */
		newNode.flag_value = Flag_value{compact_encode(newKey), new_value}
		hashedValue := newNode.hash_node()
		mpt.db[hashedValue] = newNode
		return hashedValue
	}

	/* Branch Case */
	if currNode.node_type == 1 {
		if len(newKey) == 0 {
			currNode.branch_value[16] = new_value
			hashedValue := currNode.hash_node()
			mpt.db[hashedValue] = currNode
			/* Need to update the parent */
			return hashedValue
		}

		if len(newKey) == 0 {
			currNode.branch_value[16] = new_value
			hashedValue := currNode.hash_node()
			mpt.db[hashedValue] = currNode
			/* Need to update the parent */
			return hashedValue
		}
		trimmedPath := newKey[1:]
		if len(newKey) == 1 {
			childNode := mpt.db[currNode.branch_value[newKey[0]]]
			if childNode.node_type == 0 {
				childNode := Node{}
				trimmedPath  = append(trimmedPath , 16)
				childNode.node_type = 2
				/* Set the flag_value parameters */
				childNode.flag_value = Flag_value{compact_encode(trimmedPath ), new_value}
				childHashValue := childNode.hash_node()
				mpt.db[childHashValue] = childNode
				currNode.branch_value[newKey[0]] = childHashValue
				newNewHashValue := currNode.hash_node()
				mpt.db[newNewHashValue] = currNode
				if parentHash == mpt.root {
					mpt.root = newNewHashValue
				}
				return newNewHashValue
			} else if childNode.node_type == 1 {
				childNode.branch_value[16] = new_value
				newHashValue := childNode.hash_node()
				mpt.db[newHashValue] = childNode
				currNode.branch_value[newKey[0]] = newHashValue
				newNewHashValue := currNode.hash_node()
				mpt.db[newNewHashValue] = currNode
				if parentHash == mpt.root {
					mpt.root = newHashValue
				}
				return newNewHashValue
			} else if childNode.node_type == 2 {
				if is_ext_node(childNode.flag_value.encoded_prefix) {
					childNode.node_type = 1

				} else {
					childNode.node_type = 1
					oldValue := childNode.flag_value.value
					oldEncodedPrefixPath := compact_decode(childNode.flag_value.encoded_prefix)
					if len(oldEncodedPrefixPath) > 0 {
						newEncodedPrefixPath := oldEncodedPrefixPath[1:]
						newLeafNode := mpt.buildNodeLeaf(newEncodedPrefixPath, oldValue)
						newLeafNodeHash := newLeafNode.hash_node()
						mpt.db[newLeafNodeHash] = newLeafNode

						childNode.branch_value[oldEncodedPrefixPath[0]] = newLeafNodeHash
						childNode.branch_value[16] = new_value
						childNodeHash := childNode.hash_node()
						mpt.db[childNodeHash] = childNode

						currNode.branch_value[newKey[0]] = childNodeHash
						newNewHashValue := currNode.hash_node()
						mpt.db[newNewHashValue] = currNode
						if parentHash == mpt.root {
							mpt.root = newNewHashValue
						}
						return newNewHashValue
					}

				}
			}

		}

		newHashValue := mpt.insertHelper(currNode.branch_value[newKey[0]], commonPath, trimmedPath, new_value)
		currNode.branch_value[newKey[0]] = newHashValue
		newNewHashValue := currNode.hash_node()
		mpt.db[newNewHashValue] = currNode
		if parentHash == mpt.root {
			mpt.root = newHashValue
		}
		return newNewHashValue
	}

	/* Leaf or Extension */
	if currNode.node_type == 2 {
		/* Need to decide whether the node is leaf or ext */
		/* Need to decode the node's flag_value's encoded prefix */
		/* So we want to use this to compare */
		arrayToCheckAgainst := compact_decode(currNode.flag_value.encoded_prefix)
		/* Bug where branch we're checking does not have correct prefix */
		/* If it's a leaf */

		/* If existing encoded prefix is empty, need to do something */
		if len(currNode.flag_value.encoded_prefix) == 0 {
			return parentHash
		}

		if !(is_ext_node(currNode.flag_value.encoded_prefix)) {
			/* Don't encode yet, we need this path */
			commonPath := []uint8{}
			commonPath = mpt.findCommonPath(newKey, arrayToCheckAgainst, commonPath, newKey)
			remainingPathKey := newKey[len(commonPath):]
			remainingPathExisting := arrayToCheckAgainst[len(commonPath):]
			oldValue := currNode.flag_value.value
			/* If no commonPath, just convert leaf to branch. If not, convert to extension */
			updatedHashedValue := ""

			/* So create branch that points to the new leaves, no extension */
			if len(commonPath) == 0 {
				currNode.node_type = 1
				currNode.flag_value.value = ""

				if len(remainingPathKey) == 0 {
					currNode.branch_value[16] = new_value
				} else {
					branchIndex := remainingPathKey[0]
					arrayToInsert := remainingPathKey[1:]
					newLeafNodeKey := mpt.buildNodeLeaf(arrayToInsert, new_value)
					leafHashedValueKey := newLeafNodeKey.hash_node()
					mpt.db[leafHashedValueKey] = newLeafNodeKey
					currNode.branch_value[branchIndex] = leafHashedValueKey
					// return leafHashedValueKey
				}

				/* This needs to retain the old value */
				if len(remainingPathExisting) == 0 {
					currNode.branch_value[16] = oldValue
				} else {
					branchIndex := remainingPathExisting[0]
					arrayToInsert := remainingPathExisting[1:]
					newLeafNodeExisting := mpt.buildNodeLeaf(arrayToInsert, oldValue)
					leafHashedValueKey := newLeafNodeExisting.hash_node()
					mpt.db[leafHashedValueKey] = newLeafNodeExisting
					currNode.branch_value[branchIndex] = leafHashedValueKey
					/* Hash the key */
					mpt.db[leafHashedValueKey] = newLeafNodeExisting
					// return leafHashedValueKey
				}

				updatedHashedValue = currNode.hash_node()
				mpt.db[updatedHashedValue] = currNode
				if parentHash == mpt.root {
					mpt.root = updatedHashedValue
				}
				return updatedHashedValue
				/* Else common path is not 0 */
			} else {
				newBranchNode := Node{}
				newBranchNode.node_type = 1
				/* Stick the hashed value in the correct array. */
				if len(remainingPathKey) == 0 && len(remainingPathExisting) == 0 {
					currNode.flag_value.value = new_value
					updatedHash := currNode.hash_node()
					mpt.db[updatedHash] = currNode
					if parentHash == mpt.root {
						mpt.root = updatedHash
					}
					return updatedHash
				}
				/* Only these should have the new value */
				if len(remainingPathKey) == 0 {
					newBranchNode.branch_value[16] = new_value
				} else {
					branchIndex := remainingPathKey[0]
					arrayToInsert := remainingPathKey[1:]
					/* Make an extension instead */

					newLeafNodeKey := mpt.buildNodeLeaf(arrayToInsert, new_value)
					leafHashedValueKey := newLeafNodeKey.hash_node()
					mpt.db[leafHashedValueKey] = newLeafNodeKey
					newBranchNode.branch_value[branchIndex] = leafHashedValueKey

					if len(remainingPathExisting) == 0 {
						newBranchNode.branch_value[16] = oldValue
					} else {
						arrayToInsertExisting := remainingPathExisting[1:]
						branchIndexExisting := remainingPathExisting[0]
						newLeafNodeExisting := mpt.buildNodeLeaf(arrayToInsertExisting, oldValue)
						leafHashedValueExisting := newLeafNodeExisting.hash_node()
						mpt.db[leafHashedValueExisting] = newLeafNodeExisting
						newBranchNode.branch_value[branchIndexExisting] = leafHashedValueExisting
					}

					branchHashedValue := newBranchNode.hash_node()
					mpt.db[branchHashedValue] = newBranchNode
					currNode.flag_value.value = branchHashedValue
					currNode.flag_value.encoded_prefix = compact_encode(commonPath)
					updatedHashedValue = currNode.hash_node()
					mpt.db[updatedHashedValue] = currNode
					/* Need to send back up the string */
					if parentHash == mpt.root {
						mpt.root = updatedHashedValue
					}
					return updatedHashedValue
				}

				/* This needs to retain the old value */
				if len(remainingPathExisting) == 0 {
					newBranchNode.branch_value[16] = oldValue
				} else {
					branchIndex := remainingPathExisting[0]
					arrayToInsert := remainingPathExisting[1:]
					newLeafNodeExisting := mpt.buildNodeLeaf(arrayToInsert, oldValue)
					leafHashedValueExisting := newLeafNodeExisting.hash_node()
					mpt.db[leafHashedValueExisting] = newLeafNodeExisting
					newBranchNode.branch_value[branchIndex] = leafHashedValueExisting
				}
				/* Add the branch to the DB. This branch has connections to the above leaves */
				branchHashedValue := newBranchNode.hash_node()
				mpt.db[branchHashedValue] = newBranchNode
				currNode.flag_value.encoded_prefix = compact_encode(commonPath)
				currNode.flag_value.value = branchHashedValue
				updatedHashedValue = currNode.hash_node()
				mpt.db[updatedHashedValue] = currNode
				/* Need to send back up the string */
				if parentHash == mpt.root {
					mpt.root = updatedHashedValue
				}
				return updatedHashedValue
				/* This return should send it back up the chain */
				/* Need to update if we're on root */
			}


		} else if is_ext_node(currNode.flag_value.encoded_prefix) {
			/* So we need to check whether we keep going down, if the extension is contained within the hexArray
			we're checking, or if we need to shorten the hex array and change everything else.
			  */
			commonPath := []uint8{}
			commonPath = mpt.findCommonPath(newKey, arrayToCheckAgainst, commonPath, newKey)
			remainingPathKey := newKey[len(commonPath):]
			remainingPathExisting := arrayToCheckAgainst[len(commonPath):]

			/* If there's nothing in common, need to turn into branch. Don't convert to ext node */
			/* But need to add extenstion below it IF current length of extension is greater than 1 */
			if len(commonPath) == 0 {

				valueOfChild := currNode.flag_value.value
				currentNibbles := compact_decode(currNode.flag_value.encoded_prefix)
				branchIndex := remainingPathKey[0]
				arrayToInsert := remainingPathKey[1:]
				arrayToInsertLeaf := append(arrayToInsert, 16)

				/* Create a new extension that takes old data, splice 1 from nibbles, if there's a nibble */
				if len(currentNibbles) > 1 {
					nibblesToInsert := currentNibbles[1:]
					newExtNode := Node{}
					newExtNode.node_type = 2
					newExtNode.flag_value.encoded_prefix = compact_encode(nibblesToInsert)
					newExtNode.flag_value.value = currNode.flag_value.value
					newExtNodeHash := newExtNode.hash_node()
					mpt.db[newExtNodeHash] = newExtNode

					/* Create the new leaf node */
					newLeafNode := Node{}
					newLeafNode.node_type = 2
					newLeafNode.flag_value.encoded_prefix = compact_encode(arrayToInsertLeaf)
					newLeafNode.flag_value.value = new_value
					newLeafNodeHash := newLeafNode.hash_node()
					mpt.db[newLeafNodeHash] = newLeafNode

					/* Turn current into branch */
					currNode.node_type = 1
					currNode.branch_value[currentNibbles[0]] = newExtNodeHash
					currNode.branch_value[branchIndex] = newLeafNodeHash
					newCurrNodeHash := currNode.hash_node()
					mpt.db[newCurrNodeHash] = currNode
					if parentHash == mpt.root {
						mpt.root = newCurrNodeHash
					}
					return newCurrNodeHash
				} else {
					nibblesToInsert := compact_decode(currNode.flag_value.encoded_prefix)
					newExtNode := Node{}
					newExtNode.node_type = 2
					newExtNode.flag_value.encoded_prefix = compact_encode(nibblesToInsert)
					newExtNode.flag_value.value = currNode.flag_value.value
					newExtNodeHash := newExtNode.hash_node()
					mpt.db[newExtNodeHash] = newExtNode

					/* Create the new leaf node */
					newLeafNode := Node{}
					newLeafNode.node_type = 2
					newLeafNode.flag_value.encoded_prefix = compact_encode(arrayToInsertLeaf)
					newLeafNode.flag_value.value = new_value
					newLeafNodeHash := newLeafNode.hash_node()
					mpt.db[newLeafNodeHash] = newLeafNode

					/* Turn current into branch */
					currNode.node_type = 1
					currNode.branch_value[currentNibbles[0]] = valueOfChild
					currNode.branch_value[branchIndex] = newLeafNodeHash
					newCurrNodeHash := currNode.hash_node()
					mpt.db[newCurrNodeHash] = currNode
					if parentHash == mpt.root {
						mpt.root = newCurrNodeHash
					}
					return newCurrNodeHash
				}
				/* Add the newly inserted leaf node to the branch with new value*/
			}
			/* there's a common path so shorten the existing ext node and create a branch node */

			/* Do this if only one in common for now */
			if len(commonPath) == 1 && len(remainingPathExisting) == 1 && len(remainingPathKey) > 0 {
				destinationHashValue := currNode.flag_value.value
				newBranchNode := Node{}
				newBranchNode.node_type = 1
				newLeafNode := Node{}
				newLeafNode.node_type = 2
				newLeafNodeInsert := remainingPathKey[1:]
				newLeafNodeNibbles := append(newLeafNodeInsert, 16)
				newLeafNodeEncodedPrefix := compact_encode(newLeafNodeNibbles)
				newLeafNode.flag_value.encoded_prefix = newLeafNodeEncodedPrefix
				newLeafNode.flag_value.value = new_value
				newLeafNodeHash := newLeafNode.hash_node()
				newBranchNode.branch_value[remainingPathExisting[0]] = destinationHashValue
				newBranchNode.branch_value[remainingPathKey[0]] = newLeafNodeHash
				newBranchNodeHash := newBranchNode.hash_node()

				currNode.flag_value.value = newBranchNodeHash
				currNode.flag_value.encoded_prefix = compact_encode(commonPath)
				updatedCurrNodeHash := currNode.hash_node()

				/* Add everything to DB */
				mpt.db[newLeafNodeHash] = newLeafNode
				mpt.db[newBranchNodeHash] = newBranchNode
				mpt.db[updatedCurrNodeHash] = currNode

				if parentHash == mpt.root {
					mpt.root = updatedCurrNodeHash
				}
				return updatedCurrNodeHash
			}

			/* Special Case where need to change the extension to branch and add a value at the end */
			if len(remainingPathKey) == 0 && len(remainingPathExisting)== 0 {

				/* Need to change */
				nextNodeHash := currNode.flag_value.value
				nextNode := mpt.db[nextNodeHash]

				if nextNode.node_type == 1 {
					nextNode.branch_value[16] = new_value
					// nextNode.flag_value.value = new_value
					newHash := nextNode.hash_node()
					mpt.db[newHash] = nextNode
					currNode.flag_value.value = newHash
					/* Have to rehash the currNode then too */
					updatedHashedValue := currNode.hash_node()
					mpt.db[updatedHashedValue] = currNode
					if parentHash == mpt.root {
						mpt.root = updatedHashedValue
					}
					return updatedHashedValue
				}


			}

			newBranchNode := Node{}
			newBranchNode.node_type = 1
			/* Stick the hashed value in the correct array. */

			/* Convert the existing Node to an Extension Node */
			if len(remainingPathKey) == 0 {
				newBranchNode.branch_value[16] = new_value
			} else {
				branchIndex := remainingPathKey[0]
				arrayToInsert := remainingPathKey[1:]
				newLeafNodeKey := mpt.buildNodeLeaf(arrayToInsert, new_value)
				leafHashedValueKey := newLeafNodeKey.hash_node()
				mpt.db[leafHashedValueKey] = newLeafNodeKey
				newBranchNode.branch_value[branchIndex] = leafHashedValueKey
			}

			if len(remainingPathExisting) == 0 {
				if len(remainingPathKey) > 0 {
					newHash := mpt.insertHelper(currNode.flag_value.value, commonPath, remainingPathKey, new_value)
					currNode.flag_value.value = newHash
					newCurrNodeHash := currNode.hash_node()
					mpt.db[newCurrNodeHash] = currNode
					if parentHash == mpt.root {
						mpt.root = newCurrNodeHash
					}
					return newCurrNodeHash
				} else {
					/* Just get out */
					branchHashedValue := newBranchNode.hash_node()
					mpt.db[branchHashedValue] = newBranchNode
					currNode.flag_value.encoded_prefix = compact_encode(commonPath)
					currNode.flag_value.value = branchHashedValue
					updatedHashedValue := currNode.hash_node()
					mpt.db[updatedHashedValue] = currNode

					/* Need to update if we're on root */
					if parentHash == mpt.root {
						mpt.root = updatedHashedValue
					}

					/* Don't know why i need to run this here. Don't think i do. */

					// mpt.insertHelper(currNode.flag_value.value, commonPath, remainingPathKey, new_value)
					return updatedHashedValue

				}

			} else {
				branchIndex := remainingPathExisting[0]
				arrayToInsert := remainingPathExisting[1:]
				/* Need to check whether this should be a leaf or extension */
				/* Check if it's currently an extension. might not need to do this */
				if is_ext_node(currNode.flag_value.encoded_prefix) {
					newExtNode := Node{}
					newExtNode.node_type = 2
					newExtNode.flag_value.encoded_prefix = compact_encode(arrayToInsert)
					/* Point to whatever the parent extension node was pointing to before */
					destinationNode := currNode.flag_value.value
					newExtNode.flag_value.value = destinationNode
					hashedValueExisting := newExtNode.hash_node()
					mpt.db[hashedValueExisting] = newExtNode

					/* Hook this up with the new branch */
					newBranchNode.branch_value[branchIndex] = hashedValueExisting
					branchHashedValue := newBranchNode.hash_node()
					mpt.db[branchHashedValue] = newBranchNode
					currNode.flag_value.encoded_prefix = compact_encode(commonPath)
					currNode.flag_value.value = branchHashedValue
					updatedHashedValue := currNode.hash_node()
					mpt.db[updatedHashedValue] = currNode
					if parentHash == mpt.root {
						mpt.root = updatedHashedValue
					}
					return updatedHashedValue
				} else {
					newLeafNodeExisting := mpt.buildNodeLeaf(arrayToInsert, new_value)
					leafHashedValueExisting := newLeafNodeExisting.hash_node()
					mpt.db[leafHashedValueExisting] = newLeafNodeExisting
					newBranchNode.branch_value[branchIndex] = leafHashedValueExisting
				}
			}
		}
	}
	return parentHash
}

func (mpt *MerklePatriciaTrie) Delete(key string) (string, error) {
	delete(mpt.stringDb, key)
	/* Convert the string to a hex array */
	byteArray := []uint8(key)
	hexArray := []uint8{}
	/* Then, encode the hex array */

	for i := 0; i < len(byteArray); i+=1 {
		hex1 := byteArray[i]/16
		hex2 := byteArray[i]%16
		hexArray = append(hexArray, hex1)
		hexArray = append(hexArray, hex2)
	}

	commonPath := []uint8{}
	root := mpt.deleteHelper(mpt.root, commonPath, hexArray)
	if root != "-1" {
		mpt.root = root
	} else {
	}

	return "", errors.New("path_not_found")


}

func (mpt *MerklePatriciaTrie) deleteHelper(parentHash string, commonPath, newKey []uint8) string {
	if len(mpt.db) == 0 {
		return parentHash
	}

	/* Handle case for empty key */
	currNode := mpt.db[parentHash]

	if len(newKey) < len(compact_decode(mpt.db[parentHash].flag_value.encoded_prefix)) {
		if len(currNode.flag_value.value) > 0 {
			currNode.flag_value.value = ""
		}
	}



	if currNode.node_type == 0 {
		return "-1"
	}

	/* Branch Case */
	if currNode.node_type == 1 {
		var counter = 0
		/* Current is branch  */

		/* If length is 0, check the end of the branch array to see if it's there */
		if len(newKey) == 0 {
			if len(currNode.branch_value[16]) > 0 {
				currNode.branch_value[16] = ""
				newHash := currNode.hash_node()
				mpt.db[newHash] = currNode
				var i = 0
				for i = 0; i < len(currNode.branch_value); i++ {
					if len(currNode.branch_value[i]) > 0 {
						counter++
					}
				}

				if counter == 1 {
					for i = 0; i < len(currNode.branch_value); i++ {
						if currNode.branch_value[i] != "" {
							/* add i as a nibble to the extension we turn the current branch into */
							nibbleOfExt := uint8(i)
							nibbleOfExtAsArray := []uint8{nibbleOfExt}
							currNode.node_type = 2
							currNode.flag_value.value = currNode.branch_value[i]
							childNode := mpt.db[currNode.flag_value.value]
							if  childNode.node_type == 1{
							} else if  childNode.node_type == 0{
							} else if is_ext_node(childNode.flag_value.encoded_prefix) {
								if counter == 1 {
									newValue := childNode.flag_value.value
									newNibbles := compact_decode(childNode.flag_value.encoded_prefix)
									currentNibbles := []uint8{}
									currentNibbles = append(currentNibbles, uint8(i))
									currNode.flag_value.value = newValue
									totalPrefix := append(currentNibbles, newNibbles...)
									newEncodedPrefix := compact_encode(totalPrefix)
									currNode.flag_value.encoded_prefix = newEncodedPrefix
									newHashedNode := currNode.hash_node()
									mpt.db[newHashedNode] = currNode
									return newHashedNode
								} else {
								}
							} else if !is_ext_node(childNode.flag_value.encoded_prefix) {
								nibbleOfExtAsArray = append(nibbleOfExtAsArray, compact_decode(childNode.flag_value.encoded_prefix)...)
								currNode.flag_value.value = childNode.flag_value.value
								nibbleOfExtAsArray = append(nibbleOfExtAsArray, 16)
							}
							currNode.flag_value.encoded_prefix = compact_encode(nibbleOfExtAsArray)
							hashedOldBranchNowExt := currNode.hash_node()
							mpt.db[hashedOldBranchNowExt] = currNode
							return hashedOldBranchNowExt
						}
					}
				}
				return newHash
			} else {
				return "-1"
			}
		} else  {
			/* Deleting the key if we find it */
			trimmedKey := newKey[1:]
			if len(currNode.branch_value[newKey[0]]) == 0 {
				return "-1"
			}
			/* we're here for 011 */
			newHashValue := mpt.deleteHelper(currNode.branch_value[newKey[0]], commonPath, trimmedKey)
			currNode.branch_value[newKey[0]] = newHashValue
			newCurrhash := currNode.hash_node()
			mpt.db[newCurrhash] = currNode
			checkingNode := mpt.db[currNode.branch_value[newKey[0]]]
			if newHashValue == "-1" {
				/* If we couldn't find it, this will be -1 and return here */
				return "-1"
			}

			/* Found it, so decrement the counter */
			if checkingNode.node_type == 0 {
				// counter = counter-1
				currNode.branch_value[newKey[0]] = newHashValue
				/* Need to rebalance still, don't return */
				var i = 0
				for i = 0; i < len(currNode.branch_value); i++ {
					if len(currNode.branch_value[i]) > 0 {
						counter++
					}
				}

				if counter == 1 {
					for i = 0; i < len(currNode.branch_value); i++ {
						if currNode.branch_value[i] != "" {
							/* add i as a nibble to the extension we turn the current branch into */
							nibbleOfExt := uint8(i)
							nibbleOfExtAsArray := []uint8{nibbleOfExt}
							currNode.node_type = 2
							currNode.flag_value.value = currNode.branch_value[i]
							childNode := mpt.db[currNode.flag_value.value]
							if  childNode.node_type == 1{
							} else if  childNode.node_type == 0{
							} else if is_ext_node(childNode.flag_value.encoded_prefix) {
								if counter == 1 {
									newValue := childNode.flag_value.value
									newNibbles := compact_decode(childNode.flag_value.encoded_prefix)
									currentNibbles := []uint8{}
									currentNibbles = append(currentNibbles, uint8(i))
									currNode.flag_value.value = newValue
									totalPrefix := append(currentNibbles, newNibbles...)
									newEncodedPrefix := compact_encode(totalPrefix)
									currNode.flag_value.encoded_prefix = newEncodedPrefix
									newHashedNode := currNode.hash_node()
									mpt.db[newHashedNode] = currNode
									return newHashedNode
								}
							} else if !is_ext_node(childNode.flag_value.encoded_prefix) {
								nibbleOfExtAsArray = append(nibbleOfExtAsArray, compact_decode(childNode.flag_value.encoded_prefix)...)
								currNode.flag_value.value = childNode.flag_value.value
								nibbleOfExtAsArray = append(nibbleOfExtAsArray, 16)
							}
							currNode.flag_value.encoded_prefix = compact_encode(nibbleOfExtAsArray)
							hashedOldBranchNowExt := currNode.hash_node()
							mpt.db[hashedOldBranchNowExt] = currNode
							return hashedOldBranchNowExt
						}
					}
				}
			}
			var i = 0
			for i = 0; i < len(currNode.branch_value); i++ {
				if len(currNode.branch_value[i]) > 0 {
					counter++
					if i==16 && len(newKey) == 0 {
						currNode.branch_value[i]=""
						newHash := currNode.hash_node()
						mpt.db[newHash] = currNode
						return newHash
					}
				}
			}

			/* if branch is still necessary, keep it there and rehash. otherwise change to leaf */
			if counter > 1 {
				return newCurrhash
			} else {
				/* Need to add the array back into the extension node nibble (at the front) */
				/* Remove the branch and point parent to whatever is stored as branch's child */
				for i = 0; i < len(currNode.branch_value); i++ {
					if currNode.branch_value[i] != "" {
						/* This should work eventually, we want to get the value */
						updatedHashedValue := currNode.branch_value[i]
						nodeToGet := mpt.db[updatedHashedValue]
						encodedPrefixToGet := nodeToGet.flag_value.encoded_prefix
						valueToGet:= nodeToGet.flag_value.value
						if len(encodedPrefixToGet) == 0 {
							return "-1"
						}
						/* If it's extension, we have to continue down. If not, we can move everything up */
						if is_ext_node(encodedPrefixToGet) {
							nibbles := compact_decode(encodedPrefixToGet)
							appendThis := []uint8{}
							appendThis = append(appendThis, uint8(i))
							newNibbles := append(appendThis, nibbles...)
							nodeToGet.flag_value.encoded_prefix = compact_encode(newNibbles)
							updatedHashedValue = nodeToGet.hash_node()
							mpt.db[updatedHashedValue] = nodeToGet
							return updatedHashedValue
						} else {
							/* else we have to make sure we find the key */
							currNode.node_type = 2
							prefix := []uint8{}
							prefix = append(prefix, uint8(i))
							prefix = append(commonPath, prefix...)
							decodedIndex := compact_decode(encodedPrefixToGet)
							decodedIndex = append(prefix, decodedIndex...)
							/* add 16 to make it a leaf */
							decodedIndex = append(decodedIndex, 16)
							finalEncodedIndex := compact_encode(decodedIndex)
							currNode.flag_value.encoded_prefix = finalEncodedIndex
							currNode.flag_value.value = valueToGet

							updatedHashedValue = currNode.hash_node()
							mpt.db[updatedHashedValue] = currNode
							return updatedHashedValue
						}
					}
				}
			}
		}
		updatedBranchNode := currNode.hash_node()
		mpt.db[updatedBranchNode] = currNode
		return updatedBranchNode
		/* Maybe at the end of this, I have to decide whether to convert this branch to an extension or leaf node */
	}

	if currNode.node_type == 2 {
		arrayToCheckAgainst := compact_decode(currNode.flag_value.encoded_prefix)

		if  !(is_ext_node(currNode.flag_value.encoded_prefix)) {
			/* Don't encode yet, we need this path */
			commonPath := []uint8{}
			commonPath = mpt.findCommonPath(newKey, arrayToCheckAgainst, commonPath, newKey)
			remainingPathKey := newKey[len(commonPath):]
			remainingPathExisting := arrayToCheckAgainst[len(commonPath):]

			/* We've found the end of the key, make sure it's the end of currNode as well */
			if len(remainingPathKey) == 0 {
				if len(remainingPathExisting) == 0 {
					/* We found a match so remove the value and rehash, then update trie db */
					currNode.flag_value.value = ""
					currNode.node_type = 0
					updatedHashedValue := currNode.hash_node()
					mpt.db[updatedHashedValue] = currNode
					return ""

				} else {
					/* It's not a match so return, unable to delete */
					return parentHash
				}

			} else {
				return parentHash
			}

		} else if is_ext_node(currNode.flag_value.encoded_prefix) {
			commonPath := []uint8{}
			commonPath = mpt.findCommonPath(newKey, arrayToCheckAgainst, commonPath, newKey)
			remainingPathKey := newKey[len(commonPath):]

			/* We're in extension, so trim and go down a level to check branch*/
			/* Set to parent hash because we need to pass the value up */

			if len(commonPath) > 0 {
				newHashValue := mpt.deleteHelper(currNode.flag_value.value, commonPath, remainingPathKey)
				if newHashValue == "-1" {
					return "-1"
				}

				/* In case we made changes to the currNode, we need to reassign it here */
				currNode.flag_value.value = newHashValue

				if newHashValue != parentHash && len(remainingPathKey) == 0 {
					currNode.flag_value.value = newHashValue
					childNode := mpt.db[currNode.flag_value.value]
					if childNode.node_type == 1{
					} else if childNode.node_type == 0{
					} else if is_ext_node(childNode.flag_value.encoded_prefix) {
						newValue := childNode.flag_value.value
						newNibbles := compact_decode(childNode.flag_value.encoded_prefix)
						currentNibbles := compact_decode(currNode.flag_value.encoded_prefix)
						currNode.flag_value.value = newValue
						totalPrefix := append(currentNibbles, newNibbles...)
						newEncodedPrefix := compact_encode(totalPrefix)
						currNode.flag_value.encoded_prefix = newEncodedPrefix
						newHashedNode := currNode.hash_node()
						mpt.db[newHashedNode] = currNode
						return newHashedNode
					} else if !is_ext_node(childNode.flag_value.encoded_prefix) {
					}
					newHash := currNode.hash_node()
					mpt.db[newHash] = currNode
					return newHash
				}

				/* Check to see whether child is leaf now */
				nodeToCheck := mpt.db[newHashValue]

				/* Check if it's a leaf. If it's a leaf, we can merge up into parent */
				if len(nodeToCheck.flag_value.encoded_prefix) == 0 {
					childNode := mpt.db[currNode.flag_value.value]
					if  childNode.node_type == 1{
					} else if  childNode.node_type == 0{
					} else if is_ext_node(childNode.flag_value.encoded_prefix) {
						newValue := childNode.flag_value.value
						newNibbles := compact_decode(childNode.flag_value.encoded_prefix)
						currentNibbles := compact_decode(currNode.flag_value.encoded_prefix)
						currNode.flag_value.value = newValue
						totalPrefix := append(currentNibbles, newNibbles...)
						newEncodedPrefix := compact_encode(totalPrefix)
						currNode.flag_value.encoded_prefix = newEncodedPrefix
						newHashedNode := currNode.hash_node()
						mpt.db[newHashedNode] = currNode
						return newHashedNode
					} else if !is_ext_node(childNode.flag_value.encoded_prefix) {
					}
					currNode.flag_value.value = newHashValue
					newHash := currNode.hash_node()
					mpt.db[newHash] = currNode
					return newHash
				}
				/* First check if encoded prefix is empty */
				if !is_ext_node(nodeToCheck.flag_value.encoded_prefix) {
					/* get the info of the leaf node */
					decodedPrefixToGet := compact_decode(nodeToCheck.flag_value.encoded_prefix)
					valueToGet := nodeToCheck.flag_value.value
					currentDecodedPrefix := compact_decode(currNode.flag_value.encoded_prefix)
					mergedPrefix := append(currentDecodedPrefix, decodedPrefixToGet...)
					finalPrefix := append(mergedPrefix, 16)
					finalEncodedPrefix := compact_encode(finalPrefix)
					currNode.flag_value.encoded_prefix = finalEncodedPrefix
					currNode.flag_value.value = valueToGet
					hashedValueOfNewCurrNode := currNode.hash_node()
					mpt.db[hashedValueOfNewCurrNode] = currNode
					return hashedValueOfNewCurrNode
				}
				childNode := mpt.db[currNode.flag_value.value]
				if childNode.node_type == 1 {
				} else if childNode.node_type == 0 {
				} else if is_ext_node(childNode.flag_value.encoded_prefix) {
					newValue := childNode.flag_value.value
					newNibbles := compact_decode(childNode.flag_value.encoded_prefix)
					currentNibbles := compact_decode(currNode.flag_value.encoded_prefix)
					currNode.flag_value.value = newValue
					totalPrefix := append(currentNibbles, newNibbles...)
					newEncodedPrefix := compact_encode(totalPrefix)
					currNode.flag_value.encoded_prefix = newEncodedPrefix
					newHashedNode := currNode.hash_node()
					mpt.db[newHashedNode] = currNode
					return newHashedNode
				} else if !is_ext_node(childNode.flag_value.encoded_prefix) {
				}
				currNode.flag_value.value = newHashValue
				newHash := currNode.hash_node()
				mpt.db[newHash] = currNode
				return newHash
			} else {
				newHashValue := mpt.deleteHelper(currNode.flag_value.value, commonPath, remainingPathKey)
				currNode.flag_value.value = newHashValue
				if newHashValue == "-1" {
					fmt.Println("Value is -1 in ext")
				}
				childNode := mpt.db[currNode.flag_value.value]
				if  childNode.node_type == 1 {
				} else if  childNode.node_type == 0 {
				} else if is_ext_node(childNode.flag_value.encoded_prefix) {

					newValue := childNode.flag_value.value
					newNibbles := compact_decode(childNode.flag_value.encoded_prefix)
					currentNibbles := compact_decode(currNode.flag_value.encoded_prefix)
					currNode.flag_value.value = newValue
					totalPrefix := append(currentNibbles, newNibbles...)
					newEncodedPrefix := compact_encode(totalPrefix)
					currNode.flag_value.encoded_prefix = newEncodedPrefix
					newHashedNode := currNode.hash_node()
					mpt.db[newHashedNode] = currNode
					return newHashedNode
				} else if !is_ext_node(childNode.flag_value.encoded_prefix) {
				}
				mpt.db[newHashValue] = currNode
				return newHashValue
			}

		}
	}
	return parentHash
}

func addPrefix(hex_array []uint8) []uint8 {
	var term int

	if len(hex_array) == 0 {
		hex_array = append(hex_array, 0)
		return hex_array
	}

	if hex_array[len(hex_array)-1] == 16 {
		term = 1
	} else {
		term = 0
	}
	// slice off the last int
	if term == 1 {
		hex_array = hex_array[:len(hex_array)-1]
	}
	// see if hex_array is odd
	oddlen := len(hex_array) % 2
	flags := 2 * term + oddlen
	uint8Flags := uint8(flags)
	flagsArray := []uint8{}
	flagsArray = append(flagsArray, uint8Flags)
	if oddlen != 0 {
		hex_array = append(flagsArray, hex_array...)
	} else {
		/* Note: his is the correct OOO, don't touch this! */
		flagsArray = append(flagsArray, 0)
		hex_array = append(flagsArray, hex_array...)
	}

	return hex_array
}
// Encode []uint8
func compact_encode(hex_array []uint8) []uint8 {

	o := make([]uint8, 0)
	hex_array = addPrefix(hex_array)
	if len(hex_array) < 2 {
		return hex_array
	}
	// hexarray now has an even length whose first nibble is the flags.

	for i := 0; i < len(hex_array); i+=2 {
		uintvalue := uint8(16 * hex_array[i] + hex_array[i+1])
		o = append(o, uintvalue)
	}
	return o
}
// If Leaf, ignore 16 at the end
/* Remove any prefix and convert back from ASCII to hex */
func compact_decode(encoded_arr []uint8) []uint8 {
	/* Skip the first part since we've removed the 16 */

	// hexarray now has an even length whose first nibble is the flags.
	// encoded_arr = encoded_arr[1:]
	o := make([]uint8, 0)

	if len(encoded_arr) == 0 {
		return o
	}
	/* do the opposite of encode, so go through each value and split it up */
	/* if odd, do only the first thing */

	/* Then do this */
	for i := 0; i < len(encoded_arr); i += 1 {
		uint1 := uint8(encoded_arr[i] % 16)
		uint2 := uint8(encoded_arr[i] / 16)
		o = append(o, uint2)
		o = append(o, uint1)

	}

	/* if odd, slice off 1. If even slice off 2. no */
	if o[0] == 2 || o[0] == 0 {
		o = o[2:]
	} else if o[0] == 1 || o[0] == 3 {
		o = o[1:]
	}

	return o
}

func test_compact_encode() {
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
}

// prefixLen returns the length of the common prefix of a and b.
func prefixLen(a, b []byte) int {
	var i, length = 0, len(a)
	if len(b) < length {
		length = len(b)
	}
	for ; i < length; i++ {
		if a[i] != b[i] {
			break
		}
	}
	return i
}

/* Additional Skeleton Code */

// Hashes the Node. Called by node of type Node.
func (node *Node) hash_node() string {
	var str string
	switch node.node_type {
	case 0:
		str = ""
	case 1:
		str = "branch_"
		for _, v := range node.branch_value {
			str += v
		}
	case 2:
		str = node.flag_value.value
	}

	sum := sha3.Sum256([]byte(str))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}

func (node *Node) String() string {
	str := "empty string"
	switch node.node_type {
	case 0:
		str = "[Null Node]"
	case 1:
		str = "Branch["
		for i, v := range node.branch_value[:16] {
			str += fmt.Sprintf("%d=\"%s\", ", i, v)
		}
		str += fmt.Sprintf("value=%s]", node.branch_value[16])
	case 2:
		encoded_prefix := node.flag_value.encoded_prefix
		node_name := "Leaf"
		if is_ext_node(encoded_prefix) {
			node_name = "Ext"
		}
		ori_prefix := strings.Replace(fmt.Sprint(compact_decode(encoded_prefix)), " ", ", ", -1)
		str = fmt.Sprintf("%s<%v, value=\"%s\">", node_name, ori_prefix, node.flag_value.value)
	}
	return str
}

func node_to_string(node Node) string {
	return node.String()
}

func (mpt *MerklePatriciaTrie) Initial() {
	mpt.db = make(map[string]Node)
	mpt.stringDb = make(map[string]string)
	mpt.root = ""
}

func is_ext_node(encoded_arr []uint8) bool {
	return encoded_arr[0] / 16 < 2
}

func TestCompact() {
	test_compact_encode()
}

func (mpt *MerklePatriciaTrie) String() string {
	content := fmt.Sprintf("ROOT=%s\n", mpt.root)
	for hash := range mpt.db {
		content += fmt.Sprintf("%s: %s\n", hash, node_to_string(mpt.db[hash]))
	}
	return content
}

func (mpt *MerklePatriciaTrie) Order_nodes() string {
	raw_content := mpt.String()
	content := strings.Split(raw_content, "\n")
	root_hash := strings.Split(strings.Split(content[0], "HashStart")[1], "HashEnd")[0]
	queue := []string{root_hash}
	i := -1
	rs := ""
	cur_hash := ""
	for len(queue) != 0 {
		last_index := len(queue) - 1
		cur_hash, queue = queue[last_index], queue[:last_index]
		i += 1
		line := ""
		for _, each := range content {
			if strings.HasPrefix(each, "HashStart" + cur_hash + "HashEnd") {
				line = strings.Split(each, "HashEnd: ")[1]
				rs += each + "\n"
				rs = strings.Replace(rs, "HashStart" + cur_hash + "HashEnd", fmt.Sprintf("Hash%v", i),  -1)
			}
		}
		temp2 := strings.Split(line, "HashStart")
		flag := true
		for _, each := range temp2 {
			if flag {
				flag = false
				continue
			}
			queue = append(queue, strings.Split(each, "HashEnd")[0])
		}
	}
	return rs
}