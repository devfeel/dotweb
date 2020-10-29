// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package dotweb

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

const maxParamCount uint8 = ^uint8(0)

func countParams(path string) uint8 {
	var n uint
	for i := range []byte(path) {
		switch path[i] {
		case ':', '*':
			n++
		}
	}
	return uint8(n)
}

type nodeType uint8

const (
	static nodeType = iota // default
	root
	param
	catchAll
)

type Node struct {
	path                 string
	fullPath             string
	hasExcludeMiddleware bool
	wildChild            bool
	nType                nodeType
	maxParams            uint8
	indices              string
	children             []*Node
	appMiddlewares       []Middleware
	groupMiddlewares     []Middleware
	middlewares          []Middleware
	handle               RouterHandle
	priority             uint32
}

// Use registers a middleware
func (n *Node) Use(m ...Middleware) *Node {
	if len(m) <= 0 {
		return n
	}
	step := len(n.middlewares) - 1
	for i := range m {
		if m[i] != nil {
			if step >= 0 {
				n.middlewares[step].SetNext(m[i])
			}
			n.middlewares = append(n.middlewares, m[i])
			step++
		}
	}
	return n
}

// AppMiddlewares return AppMiddlewares
func (n *Node) AppMiddlewares() []Middleware {
	return n.appMiddlewares
}

// GroupMiddlewares return GroupMiddlewares
func (n *Node) GroupMiddlewares() []Middleware {
	return n.groupMiddlewares
}

// Middlewares return middlewares
func (n *Node) Middlewares() []Middleware {
	return n.middlewares
}

// Path return full path in node
func (n *Node) Path() string {
	return n.fullPath
}

func (n *Node) Node() *Node {
	return n
}

// Increments priority of the given child and reorders if necessary
func (n *Node) incrementChildPrio(pos int) int {
	cs := n.children
	cs[pos].priority++
	prio := cs[pos].priority

	// Adjust position (move to front)
	newPos := pos
	for ; newPos > 0 && cs[newPos-1].priority < prio; newPos-- {
		// Swap node positions
		cs[newPos-1], cs[newPos] = cs[newPos], cs[newPos-1]
	}

	// Build new index char string
	if newPos != pos {
		n.indices = n.indices[:newPos] + // Unchanged prefix, might be empty
			n.indices[pos:pos+1] + // The index char we move
			n.indices[newPos:pos] + n.indices[pos+1:] // Rest without char at 'pos'
	}

	return newPos
}

// addRoute adds a node with the given handle to the path.
// Not concurrency-safe!
func (n *Node) addRoute(path string, handle RouterHandle, m ...Middleware) (outnode *Node) {
	outnode = n

	fullPath := path
	n.priority++
	numParams := countParams(path)

	// Non-empty tree
	if len(n.path) > 0 || len(n.children) > 0 {
	walk:
		for {
			// Update maxParams of the current node
			if numParams > n.maxParams {
				n.maxParams = numParams
			}

			// Find the longest common prefix.
			// This also implies that the common prefix contains no ':' or '*'
			// since the existing key can't contain those chars.
			i := 0
			max := min(len(path), len(n.path))
			for i < max && path[i] == n.path[i] {
				i++
			}

			// Split edge
			if i < len(n.path) {
				child := Node{
					path:             n.path[i:],
					fullPath:         n.fullPath,
					wildChild:        n.wildChild,
					nType:            static,
					indices:          n.indices,
					children:         n.children,
					handle:           n.handle,
					priority:         n.priority - 1,
					middlewares:      n.middlewares,
					groupMiddlewares: n.groupMiddlewares,
					appMiddlewares:   n.appMiddlewares,
				}

				// Update maxParams (max of all children)
				for i := range child.children {
					if child.children[i].maxParams > child.maxParams {
						child.maxParams = child.children[i].maxParams
					}
				}

				n.children = []*Node{&child}
				// []byte for proper unicode char conversion, see #65
				n.indices = string([]byte{n.path[i]})
				n.path = path[:i]
				n.fullPath = n.path
				n.handle = nil
				n.wildChild = false
				n.middlewares = nil
				n.groupMiddlewares = nil
				n.appMiddlewares = nil
			}

			// Make new node a child of this node
			if i < len(path) {
				path = path[i:]

				if n.wildChild {
					n = n.children[0]
					n.priority++

					// Update maxParams of the child node
					if numParams > n.maxParams {
						n.maxParams = numParams
					}
					numParams--

					// Check if the wildcard matches
					if len(path) >= len(n.path) && n.path == path[:len(n.path)] &&
						// Check for longer wildcard, e.g. :name and :names
						(len(n.path) >= len(path) || path[len(n.path)] == '/') {
						continue walk
					} else {
						// Wildcard conflict
						var pathSeg string
						if n.nType == catchAll {
							pathSeg = path
						} else {
							pathSeg = strings.SplitN(path, "/", 2)[0]
						}
						prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path
						panic("'" + pathSeg +
							"' in new path '" + fullPath +
							"' conflicts with existing wildcard '" + n.path +
							"' in existing prefix '" + prefix +
							"'")
					}
				}

				c := path[0]

				// Slash after param
				if n.nType == param && c == '/' && len(n.children) == 1 {
					n = n.children[0]
					n.priority++
					continue walk
				}

				// Check if a child with the next path byte exists
				for i := 0; i < len(n.indices); i++ {
					if c == n.indices[i] {
						i = n.incrementChildPrio(i)
						n = n.children[i]
						continue walk
					}
				}

				// Otherwise insert it
				if c != ':' && c != '*' {
					// []byte for proper unicode char conversion, see #65
					n.indices += string([]byte{c})
					child := &Node{
						maxParams: numParams,
					}
					n.children = append(n.children, child)
					n.incrementChildPrio(len(n.indices) - 1)
					n = child
				}
				outnode = n.insertChild(numParams, path, fullPath, handle)
				return

			} else if i == len(path) { // Make node a (in-path) leaf
				if n.handle != nil {
					panic("a handle is already registered for path '" + fullPath + "'")
				}
				n.handle = handle
				n.Use(m...)
				outnode = n
			}
			return
		}
	} else { // Empty tree
		outnode = n.insertChild(numParams, path, fullPath, handle, m...)
		n.nType = root
	}
	return
}

func (n *Node) insertChild(numParams uint8, path, fullPath string, handle RouterHandle, m ...Middleware) *Node {
	var offset int // already handled bytes of the path

	// Find prefix until first wildcard (beginning with ':'' or '*'')
	for i, max := 0, len(path); numParams > 0; i++ {
		c := path[i]
		if c != ':' && c != '*' {
			continue
		}

		// Find wildcard end (either '/' or path end) and check the name for
		// invalid characters
		end := i + 1
		invalid := false
		for end < max {
			c := path[end]
			if c == '/' {
				break
			}
			if c == ':' || c == '*' {
				invalid = true
			}
			end++
		}

		// The wildcard name must not contain ':' and '*'
		if invalid {
			panic("only one wildcard per path segment is allowed, has: '" +
				path[i:end] + "' in path '" + fullPath + "'")
		}

		// Check if the wildcard has a name
		if end-i < 2 {
			panic("wildcards must be named with a non-empty name in path '" + fullPath + "'")
		}

		// Check if this node has existing children which would be
		// unreachable if we insert the wildcard here
		if len(n.children) > 0 {
			panic("wildcard route '" + path[i:end] +
				"' conflicts with existing children in path '" + fullPath + "'")
		}

		if c == ':' { // param
			// Split path at the beginning of the wildcard
			if i > 0 {
				n.path = path[offset:i]
				offset = i
			}

			child := &Node{
				nType:     param,
				maxParams: numParams,
			}
			n.children = []*Node{child}
			n.wildChild = true
			n = child
			n.priority++
			numParams--

			// If the path doesn't end with the wildcard, then there
			// will be another non-wildcard subpath starting with '/'
			if end < max {
				n.path = path[offset:end]
				offset = end

				child := &Node{
					maxParams: numParams,
					priority:  1,
				}
				n.children = []*Node{child}
				n = child
			}

		} else { // catchAll
			if end != max || numParams > 1 {
				panic("catch-all routes are only allowed at the end of the path in path '" + fullPath + "'")
			}

			if len(n.path) > 0 && n.path[len(n.path)-1] == '/' {
				panic("catch-all conflicts with existing handle for the path segment root in path '" + fullPath + "'")
			}

			// currently fixed width 1 for '/'
			i--
			if path[i] != '/' {
				panic("no / before catch-all in path '" + fullPath + "'")
			}

			n.path = path[offset:i]

			// first node: catchAll node with empty path
			child := &Node{
				wildChild: true,
				nType:     catchAll,
				maxParams: 1,
			}
			if child.maxParams > n.maxParams {
				n.maxParams = child.maxParams
			}
			n.children = []*Node{child}
			n.indices = string(path[i])
			n = child
			n.priority++

			// second node: node holding the variable
			child = &Node{
				path:      path[i:],
				nType:     catchAll,
				maxParams: 1,
				handle:    handle,
				priority:  1,
			}
			child.Use(m...)
			n.children = []*Node{child}

			return n
		}
	}

	// Insert remaining path part and handle to the leaf
	n.path = path[offset:]
	n.handle = handle
	n.Use(m...)
	return n
}

func (n *Node) getNode(path string) *Node {
walk: // outer loop for walking the tree
	for {
		if len(path) > len(n.path) {
			if path[:len(n.path)] == n.path {
				path = path[len(n.path):]
				// If this node does not have a wildcard (param or catchAll)
				// child,  we can just look up the next child node and continue
				// to walk down the tree
				if !n.wildChild {
					c := path[0]
					for i := 0; i < len(n.indices); i++ {
						if c == n.indices[i] {
							n = n.children[i]
							continue walk
						}
					}

					// Nothing found.
					// We can recommend to redirect to the same URL without a
					// trailing slash if a leaf exists for that path.
					return nil

				}

				// handle wildcard child
				n = n.children[0]
				switch n.nType {
				case param:
					// find param end (either '/' or path end)
					end := 0
					for end < len(path) && path[end] != '/' {
						end++
					}

					// we need to go deeper!
					if end < len(path) {
						if len(n.children) > 0 {
							path = path[end:]
							n = n.children[0]
							continue walk
						}

						// ... but we can't
						return nil
					}

					return n

				case catchAll:
					return n
				default:
					panic("invalid node type")
				}
			}
		} else if path == n.path {
			return n
		}
		return nil
	}
}

// Returns the handle registered with the given path (key). The values of
// wildcards are saved to a map.
// If no handle can be found, a TSR (trailing slash redirect) recommendation is
// made if a handle exists with an extra (without the) trailing slash for the
// given path.
func (n *Node) getValue(path string) (handle RouterHandle, p Params, outnode *Node, tsr bool) {
walk: // outer loop for walking the tree
	for {
		if len(path) > len(n.path) {
			if path[:len(n.path)] == n.path {
				path = path[len(n.path):]
				// If this node does not have a wildcard (param or catchAll)
				// child,  we can just look up the next child node and continue
				// to walk down the tree
				if !n.wildChild {
					c := path[0]
					for i := 0; i < len(n.indices); i++ {
						if c == n.indices[i] {
							n = n.children[i]
							continue walk
						}
					}

					// Nothing found.
					// We can recommend to redirect to the same URL without a
					// trailing slash if a leaf exists for that path.
					tsr = (path == "/" && n.handle != nil)
					return

				}

				// handle wildcard child
				outnode = n.children[0]
				n = outnode
				switch n.nType {
				case param:
					// find param end (either '/' or path end)
					end := 0
					for end < len(path) && path[end] != '/' {
						end++
					}

					// save param value
					if p == nil {
						// lazy allocation
						p = make(Params, 0, n.maxParams)
					}
					i := len(p)
					p = p[:i+1] // expand slice within preallocated capacity
					p[i].Key = n.path[1:]
					p[i].Value = path[:end]

					// we need to go deeper!
					if end < len(path) {
						if len(n.children) > 0 {
							path = path[end:]
							n = n.children[0]
							continue walk
						}

						// ... but we can't
						tsr = (len(path) == end+1)
						return
					}

					if handle = n.handle; handle != nil {
						return
					} else if len(n.children) == 1 {
						// No handle found. Check if a handle for this path + a
						// trailing slash exists for TSR recommendation
						n = n.children[0]
						tsr = (n.path == "/" && n.handle != nil)
					}

					return

				case catchAll:
					// save param value
					if p == nil {
						// lazy allocation
						p = make(Params, 0, n.maxParams)
					}
					i := len(p)
					p = p[:i+1] // expand slice within preallocated capacity
					p[i].Key = n.path[2:]
					p[i].Value = path

					handle = n.handle
					return

				default:
					panic("invalid node type")
				}
			}
		} else if path == n.path {
			outnode = n
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			if handle = n.handle; handle != nil {
				return
			}

			if path == "/" && n.wildChild && n.nType != root {
				tsr = true
				return
			}

			// No handle found. Check if a handle for this path + a
			// trailing slash exists for trailing slash recommendation
			for i := 0; i < len(n.indices); i++ {
				if n.indices[i] == '/' {
					n = n.children[i]
					tsr = (len(n.path) == 1 && n.handle != nil) ||
						(n.nType == catchAll && n.children[0].handle != nil)
					return
				}
			}

			return
		}

		// Nothing found. We can recommend to redirect to the same URL with an
		// extra trailing slash if a leaf exists for that path
		tsr = (path == "/") ||
			(len(n.path) == len(path)+1 && n.path[len(path)] == '/' &&
				path == n.path[:len(n.path)-1] && n.handle != nil)
		return
	}
}

// Makes a case-insensitive lookup of the given path and tries to find a handler.
// It can optionally also fix trailing slashes.
// It returns the case-corrected path and a bool indicating whether the lookup
// was successful.
func (n *Node) findCaseInsensitivePath(path string, fixTrailingSlash bool) (ciPath []byte, found bool) {
	return n.findCaseInsensitivePathRec(
		path,
		strings.ToLower(path),
		make([]byte, 0, len(path)+1), // preallocate enough memory for new path
		[4]byte{},                    // empty rune buffer
		fixTrailingSlash,
	)
}

// shift bytes in array by n bytes left
func shiftNRuneBytes(rb [4]byte, n int) [4]byte {
	switch n {
	case 0:
		return rb
	case 1:
		return [4]byte{rb[1], rb[2], rb[3], 0}
	case 2:
		return [4]byte{rb[2], rb[3]}
	case 3:
		return [4]byte{rb[3]}
	default:
		return [4]byte{}
	}
}

// recursive case-insensitive lookup function used by n.findCaseInsensitivePath
func (n *Node) findCaseInsensitivePathRec(path, loPath string, ciPath []byte, rb [4]byte, fixTrailingSlash bool) ([]byte, bool) {
	loNPath := strings.ToLower(n.path)

walk: // outer loop for walking the tree
	for len(loPath) >= len(loNPath) && (len(loNPath) == 0 || loPath[1:len(loNPath)] == loNPath[1:]) {
		// add common path to result
		ciPath = append(ciPath, n.path...)

		if path = path[len(n.path):]; len(path) > 0 {
			loOld := loPath
			loPath = loPath[len(loNPath):]

			// If this node does not have a wildcard (param or catchAll) child,
			// we can just look up the next child node and continue to walk down
			// the tree
			if !n.wildChild {
				// skip rune bytes already processed
				rb = shiftNRuneBytes(rb, len(loNPath))

				if rb[0] != 0 {
					// old rune not finished
					for i := 0; i < len(n.indices); i++ {
						if n.indices[i] == rb[0] {
							// continue with child node
							n = n.children[i]
							loNPath = strings.ToLower(n.path)
							continue walk
						}
					}
				} else {
					// process a new rune
					var rv rune

					// find rune start
					// runes are up to 4 byte long,
					// -4 would definitely be another rune
					var off int
					for max := min(len(loNPath), 3); off < max; off++ {
						if i := len(loNPath) - off; utf8.RuneStart(loOld[i]) {
							// read rune from cached lowercase path
							rv, _ = utf8.DecodeRuneInString(loOld[i:])
							break
						}
					}

					// calculate lowercase bytes of current rune
					utf8.EncodeRune(rb[:], rv)
					// skipp already processed bytes
					rb = shiftNRuneBytes(rb, off)

					for i := 0; i < len(n.indices); i++ {
						// lowercase matches
						if n.indices[i] == rb[0] {
							// must use a recursive approach since both the
							// uppercase byte and the lowercase byte might exist
							// as an index
							if out, found := n.children[i].findCaseInsensitivePathRec(
								path, loPath, ciPath, rb, fixTrailingSlash,
							); found {
								return out, true
							}
							break
						}
					}

					// same for uppercase rune, if it differs
					if up := unicode.ToUpper(rv); up != rv {
						utf8.EncodeRune(rb[:], up)
						rb = shiftNRuneBytes(rb, off)

						for i := 0; i < len(n.indices); i++ {
							// uppercase matches
							if n.indices[i] == rb[0] {
								// continue with child node
								n = n.children[i]
								loNPath = strings.ToLower(n.path)
								continue walk
							}
						}
					}
				}

				// Nothing found. We can recommend to redirect to the same URL
				// without a trailing slash if a leaf exists for that path
				return ciPath, (fixTrailingSlash && path == "/" && n.handle != nil)
			}

			n = n.children[0]
			switch n.nType {
			case param:
				// find param end (either '/' or path end)
				k := 0
				for k < len(path) && path[k] != '/' {
					k++
				}

				// add param value to case insensitive path
				ciPath = append(ciPath, path[:k]...)

				// we need to go deeper!
				if k < len(path) {
					if len(n.children) > 0 {
						// continue with child node
						n = n.children[0]
						loNPath = strings.ToLower(n.path)
						loPath = loPath[k:]
						path = path[k:]
						continue
					}

					// ... but we can't
					if fixTrailingSlash && len(path) == k+1 {
						return ciPath, true
					}
					return ciPath, false
				}

				if n.handle != nil {
					return ciPath, true
				} else if fixTrailingSlash && len(n.children) == 1 {
					// No handle found. Check if a handle for this path + a
					// trailing slash exists
					n = n.children[0]
					if n.path == "/" && n.handle != nil {
						return append(ciPath, '/'), true
					}
				}
				return ciPath, false

			case catchAll:
				return append(ciPath, path...), true

			default:
				panic("invalid node type")
			}
		} else {
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			if n.handle != nil {
				return ciPath, true
			}

			// No handle found.
			// Try to fix the path by adding a trailing slash
			if fixTrailingSlash {
				for i := 0; i < len(n.indices); i++ {
					if n.indices[i] == '/' {
						n = n.children[i]
						if (len(n.path) == 1 && n.handle != nil) ||
							(n.nType == catchAll && n.children[0].handle != nil) {
							return append(ciPath, '/'), true
						}
						return ciPath, false
					}
				}
			}
			return ciPath, false
		}
	}

	// Nothing found.
	// Try to fix the path by adding / removing a trailing slash
	if fixTrailingSlash {
		if path == "/" {
			return ciPath, true
		}
		if len(loPath)+1 == len(loNPath) && loNPath[len(loPath)] == '/' &&
			loPath[1:] == loNPath[1:len(loPath)] && n.handle != nil {
			return append(ciPath, n.path...), true
		}
	}
	return ciPath, false
}
