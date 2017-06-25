package aho_corasick

import "container/list"

type ACAutomaton interface {
	Search(s string) <-chan Match
	SearchCallback(s string, callback func(start, end int))
}

type trie_node struct {
	flag         bool
	depth        int
	fail_link    *trie_node
	success_link *trie_node
	children     [](*trie_node)
}

func new_node(depth int) *trie_node {
	return &trie_node{false, depth, nil, nil, make([](*trie_node), (1 << 8))}
}

type ac_automaton struct {
	root *trie_node
}

func New(dict []string) ACAutomaton {
	ac := &ac_automaton{root: new_node(0)}
	ac.root.fail_link = ac.root
	for _, str := range dict {
		ac.add(str)
	}
	ac.build()
	return ac
}

func (self *ac_automaton) add(str string) {
	p := self.root
	for i := 0; i < len(str); i++ {
		b := str[i]
		if p.children[b] == nil {
			p.children[b] = new_node(i + 1)
		}
		p = p.children[b]
	}
	p.flag = true
}

func (self *ac_automaton) build() {
	queue := list.New()
	for _, n := range self.root.children {
		if n != nil {
			n.fail_link = self.root
			queue.PushBack(n)
		}
	}

	for queue.Len() > 0 {
		current := queue.Front()
		queue.Remove(current)

		crt_node := current.Value.(*trie_node)
		for i, next := range crt_node.children {
			if next == nil {
				continue
			}

			p := crt_node.fail_link
			for p != self.root {
				if p.children[i] != nil {
					break
				} else {
					p = p.fail_link
				}
			}

			if p.children[i] != nil {
				next.fail_link = p.children[i]
			} else {
				next.fail_link = self.root
			}

			queue.PushBack(next)
		}

		// End of word. Update success_link
		p := crt_node.fail_link
		for p != self.root && p.flag == false {
			p = p.fail_link
		}

		if p != self.root {
			crt_node.success_link = p
		}
	}
}

func (self *ac_automaton) Search(str string) <-chan Match {
	ret := make(chan Match)
	go func() {
		defer close(ret)
		self.SearchCallback(str, func(start, end int) {
			ret <- (&match{start, end})
		})
	}()
	return ret
}

func (self *ac_automaton) SearchCallback(str string, callback func(start, end int)) {
	p := self.root
	for i := 0; i < len(str); {
		ch := str[i]
		if p.children[ch] == nil {
			// Mismatch, follow fail_link
			if p == self.root {
				i++
			} else {
				p = p.fail_link
			}
			continue
		}

		// Matched, update state
		i++
		p = p.children[ch]

		// Check success link and output results
		succ := p
		for succ != nil {
			if succ.flag {
				callback(i-succ.depth, i)
			}
			succ = succ.success_link
		}
	}
}
