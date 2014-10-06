package mstree

import (
	"io"
	"regexp"
	"strings"
)

type node struct {
	Children map[string]*node
}

func newNode() *node {
	return &node{make(map[string]*node)}
}

func (n *node) insert(tokens []string, inserted *bool) {
	if len(tokens) == 0 {
		return
	}
	first, tail := tokens[0], tokens[1:]
	child, ok := n.Children[first]
	if !ok {
		*inserted = true
		child = newNode()
		n.Children[first] = child
	}
	child.insert(tail, inserted)
}

func (n *node) traverseDump(prefix string, writer io.Writer) {
	if len(n.Children) == 0 {
		io.WriteString(writer, prefix+"\n")
	} else {
		for k, node := range n.Children {
			var nPref string
			if prefix == "" {
				nPref = k
			} else {
				nPref = prefix + "." + k
			}
			node.traverseDump(nPref, writer)
		}
	}
}

func (n *node) search(pattern string) map[string]*node {
	if pattern == "*" {
		return n.Children
	}

	results := make(map[string]*node)

	wcIndex := strings.Index(pattern, "*")
	qIndex := strings.Index(pattern, "?")

	if wcIndex == -1 && qIndex == -1 {
		if node, ok := n.Children[pattern]; ok {
			results[pattern] = node
		}
		return results
	}

	if qIndex == -1 {
		// Only *
		lwcIndex := strings.LastIndex(pattern, "*")

		if wcIndex != lwcIndex || (wcIndex != 0 && wcIndex != len(pattern)-1) {
			// more than one wildcard or one wildcard in the middle
			re := regexp.MustCompile(strings.Replace(pattern, "*", ".*", -1))
			for k, node := range n.Children {
				if re.MatchString(k) {
					results[k] = node
				}
			}
			return results
		}

		if wcIndex == len(pattern)-1 {
			// wildcard at the end
			partial := pattern[:len(pattern)-1]
			for k, node := range n.Children {
				if strings.HasPrefix(k, partial) {
					results[k] = node
				}
			}
		} else {
			// wildcard at the begining
			partial := pattern[1:]
			for k, node := range n.Children {
				if strings.HasSuffix(k, partial) {
					results[k] = node
				}
			}
		}
	} else if wcIndex == -1 {
		// Only ?
		lqIndex := strings.LastIndex(pattern, "?")
		if qIndex != lqIndex || (qIndex != 0 && qIndex != len(pattern)-1) {
			// more than one ? or one ? in the middle
			re := regexp.MustCompile(strings.Replace(pattern, "?", ".?", -1))
			for k, node := range n.Children {
				if re.MatchString(k) {
					results[k] = node
				}
			}
			return results
		}

		if qIndex == len(pattern)-1 {
			// ? at the end
			partial := pattern[:len(pattern)-1]
			for k, node := range n.Children {
				if k[:len(k)-1] == partial {
					results[k] = node
				}
			}
		} else {
			// ? at the begining
			partial := pattern[1:]
			for k, node := range n.Children {
				if k[1:] == partial {
					results[k] = node
				}
			}
		}

	} else {
		// * and ? presents
		rePattern := strings.Replace(strings.Replace(pattern, "*", ".*", -1), "?", ".?", -1)
		re := regexp.MustCompile(rePattern)
		for k, node := range n.Children {
			if re.MatchString(k) {
				results[k] = node
			}
		}
	}

	return results
}
