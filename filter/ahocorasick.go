package filter

import (
	"github.com/sgoware/ds/queue/arrayqueue"
)

type acNode struct {
	value    rune
	children map[rune]*acNode
	word     *string
	fail     *acNode
}

func newAcNode(r rune) *acNode {
	return &acNode{
		value:    r,
		children: make(map[rune]*acNode),
		word:     nil,
	}
}

type AcModel struct {
	root *acNode
}

func NewAcModel() *AcModel {
	return &AcModel{
		root: newAcNode(0),
	}
}

func (m *AcModel) AddWords(words ...string) {
	for _, word := range words {
		m.AddWord(word)
	}

	m.buildFailPointers()
}

func (m *AcModel) AddWord(word string) {
	now := m.root
	runes := []rune(word)

	for _, r := range runes {
		if next, ok := now.children[r]; ok {
			now = next
		} else {
			next = newAcNode(r)
			now.children[r] = next
			now = next
		}
	}

	now.word = new(string)
	*now.word = word
}

func (m *AcModel) DelWords(words ...string) {
	for _, word := range words {
		m.DelWord(word)
	}

	m.buildFailPointers()
}

func (m *AcModel) DelWord(word string) {
	var lastLeaf *acNode
	var lastLeafNextRune rune
	now := m.root
	runes := []rune(word)

	for _, r := range runes {
		if next, ok := now.children[r]; !ok {
			return
		} else {
			if now.word != nil {
				lastLeaf = now
				lastLeafNextRune = r
			}
			now = next
		}
	}

	delete(lastLeaf.children, lastLeafNextRune)
}

func (m *AcModel) buildFailPointers() {
	q := arrayqueue.New(m.root)

	for q.Len() > 0 {
		temp, _ := q.Top()
		q.Pop()
		for _, node := range temp.(*acNode).children {
			if temp.(*acNode) == m.root {
				node.fail = m.root
			} else {
				p := temp.(*acNode).fail
				for p != nil {
					if next, found := p.children[node.value]; found {
						node.fail = next
						break
					}
					p = p.fail
				}
				if p == nil {
					node.fail = m.root
				}
			}

			q.Push(node)
		}
	}
}

func (m *AcModel) Listen(addChan, delChan <-chan string) {
	go func() {
		var words []string

		for word := range addChan {
			words = append(words, word)
			if len(addChan) == 0 {
				m.AddWords(words...)
				word = word[:0]
			}
		}
	}()

	go func() {
		var words []string

		for word := range delChan {
			words = append(words, word)
			if len(delChan) == 0 {
				m.DelWords(words...)
				word = word[:0]
			}
		}
	}()
}

func (m *AcModel) FindAll(text string) []string {
	var matches []string
	var found bool

	now := m.root
	var temp *acNode
	runes := []rune(text)

	for pos := 0; pos < len(runes); pos++ {
		_, found = now.children[runes[pos]]
		if !found && now != m.root {
			now = now.fail
			for ; !found && now != m.root; now, found = now.children[runes[pos]] {
				now = now.fail
			}
		}

		// 若找到匹配成功的字符串结点, 则指向那个结点, 否则指向根结点
		if next, ok := now.children[runes[pos]]; ok {
			now = next
		} else {
			now = m.root
		}

		temp = now

		for temp != m.root {
			if temp.word != nil {
				matches = append(matches, *temp.word)
			}
			temp = temp.fail
		}
	}

	var res []string
	set := make(map[string]struct{})

	for _, word := range matches {
		if _, ok := set[word]; !ok {
			set[word] = struct{}{}
			res = append(res, word)
		}
	}

	return res
}

func (m *AcModel) FindAllCount(text string) map[string]int {
	res := make(map[string]int)
	var found bool
	var temp *acNode

	now := m.root
	runes := []rune(text)

	for pos := 0; pos < len(runes); pos++ {
		_, found = now.children[runes[pos]]
		if !found && now != m.root {
			now = now.fail
			for ; !found && now != m.root; now, found = now.children[runes[pos]] {
				now = now.fail
			}
		}

		// 若找到匹配成功的字符串结点, 则指向那个结点, 否则指向根结点
		if next, ok := now.children[runes[pos]]; ok {
			now = next
		} else {
			now = m.root
		}

		temp = now

		for temp != m.root {
			if temp.word != nil {
				res[*temp.word]++
			}
			temp = temp.fail
		}
	}

	return res
}

func (m *AcModel) FindOne(text string) string {
	var found bool
	var temp *acNode

	now := m.root
	runes := []rune(text)

	for pos := 0; pos < len(runes); pos++ {
		_, found = now.children[runes[pos]]
		if !found && now != m.root {
			now = now.fail
			for ; !found && now != m.root; now, found = now.children[runes[pos]] {
				now = now.fail
			}
		}

		// 若找到匹配成功的字符串结点, 则指向那个结点, 否则指向根结点
		if next, ok := now.children[runes[pos]]; ok {
			now = next
		} else {
			now = m.root
		}

		temp = now

		for temp != m.root {
			if temp.word != nil {
				return *temp.word
			}
			temp = temp.fail
		}
	}

	return ""
}

func (m *AcModel) IsSensitive(text string) bool {
	return m.FindOne(text) != ""
}

func (m *AcModel) Replace(text string, repl rune) string {
	var found bool
	var temp *acNode

	now := m.root
	runes := []rune(text)

	for pos := 0; pos < len(runes); pos++ {
		_, found = now.children[runes[pos]]
		if !found && now != m.root {
			now = now.fail
			for ; !found && now != m.root; now, found = now.children[runes[pos]] {
				now = now.fail
			}
		}

		// 若找到匹配成功的字符串结点, 则指向那个结点, 否则指向根结点
		if next, ok := now.children[runes[pos]]; ok {
			now = next
		} else {
			now = m.root
		}

		temp = now

		for temp != m.root {
			if temp.word != nil {
				for i := pos - len([]rune(*temp.word)) + 1; i <= pos; i++ {
					runes[i] = repl
				}
			}
			temp = temp.fail
		}
	}

	return string(runes)
}

func (m *AcModel) Remove(text string) string {
	var found bool
	var temp *acNode

	now := m.root
	runes := []rune(text)

	for pos := 0; pos < len(runes); pos++ {
		_, found = now.children[runes[pos]]
		if !found && now != m.root {
			now = now.fail
			for ; !found && now != m.root; now, found = now.children[runes[pos]] {
				now = now.fail
			}
		}

		// 若找到匹配成功的字符串结点, 则指向那个结点, 否则指向根结点
		if next, ok := now.children[runes[pos]]; ok {
			now = next
		} else {
			now = m.root
		}

		temp = now

		for temp != m.root {
			if temp.word != nil {
				runes = append(runes[:pos-len([]rune(*temp.word))+1], runes[pos+1:]...)
				pos -= len([]rune(*temp.word))
			}
			temp = temp.fail
		}
	}

	return string(runes)
}
