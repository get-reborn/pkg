/**
 * created by genshen on 2018/11/10
 */

package pkg

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	DlStatusEmpty = iota
	DlStatusSkip
	DlStatusOk
)

type DependencyTree struct {
	Context      DepPkgContext
	Dependencies []*DependencyTree
	Builder      []string // outer builder (lib used by others)
	SelfBuild    []string // inner builder (shows how this package is built)
	CMakeLib     string   // outer cmake script to add this lib.
	SelfCMakeLib string   // inner cmake script to add this lib.
	DlStatus     int
	IsPkgPackage bool
}

type DepPkgContext struct {
	PackageName string
	SrcPath     string
}

// marshal dependency tree content to a json file.
func (depTree *DependencyTree) Dump(filename string) error {
	if content, err := json.Marshal(depTree); err != nil { // unmarshal json to struct
		return err
	} else {
		if dumpFile, err := os.Create(filename); err != nil {
			return err
		} else {
			dumpFile.Write(content)
		}
	}
	return nil
}

// recover the dependency tree from a json file.
// the result is saved in variable deps.
func DepTreeRecover(deps *DependencyTree, filename string) (error) {
	if depFile, err := os.Open(filename); err != nil { // file open error or not exists.
		return err
	} else {
		defer depFile.Close()
		if bytes, err := ioutil.ReadAll(depFile); err != nil { // read file contents
			return err
		} else {
			if err := json.Unmarshal(bytes, &deps); err != nil { // unmarshal json to struct
				return err
			}
			return nil
		}
	}
}

// traversal all tree node with pre-order.
// if the return value of callback function is false, it will skip its children nodes.
func (depTree *DependencyTree) Traversal(callback func(*DependencyTree) bool) {
	if r := callback(depTree); r == false {
		return
	}
	// if this node has children
	if depTree.Dependencies == nil || len(depTree.Dependencies) == 0 {
		return
	} else {
		for _, d := range depTree.Dependencies {
			d.Traversal(callback)
		}
	}
}

// traversal all tree node with pre-order.
// if the return value of callback function is false, then the traversal will break.
func (depTree *DependencyTree) TraversalPreOrder(callback func(*DependencyTree) bool) bool {
	if r := callback(depTree); r == false {
		return false
	}
	// if this node has children
	if depTree.Dependencies == nil || len(depTree.Dependencies) == 0 {
		return true
	} else {
		for _, d := range depTree.Dependencies {
			if r := d.TraversalPreOrder(callback); r == false {
				return false
			}
		}
	}
	return true
}

// traversal all tree node(including the root node) by deep first strategy.
// if return value of callback is false, then the traversal will break.
func (depTree *DependencyTree) TraversalDeep(callback func(*DependencyTree) bool) bool {
	// if this node has children
	if depTree.Dependencies == nil || len(depTree.Dependencies) == 0 {
		return true
	} else {
		for _, d := range depTree.Dependencies {
			if r := d.TraversalDeep(callback); r == false {
				return false
			}
		}
	}
	return callback(depTree)
}
