package src

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Node struct {
	child map[rune]*Node
	end   bool
}

func newNode() *Node {
	return &Node{}
}
func (n *Node) checkAndFixMap() {
	if nil == n.child {
		n.child = map[rune]*Node{}
	}
}
func (n *Node) setChild(key rune) *Node {
	n.checkAndFixMap()
	c, ok := n.child[key]
	if ok {
		return c
	} else {
		node := newNode()
		n.child[key] = node
		return node
	}
}
func (n *Node) getChild(key rune) (*Node, bool) {
	if n.child == nil || n == nil {
		return nil, false
	} else {
		node, ok := n.child[key]
		return node, ok
	}
}
func (n *Node) setNodeEnd() {
	n.end = true
}
func (n *Node) passerWord(word []rune) {
	l := len(word)
	next := n
	for i := 0; i < l; i++ {
		next = next.setChild(word[i])
	}
	if next != n {
		next.setNodeEnd()
	}
}
func (n *Node) checkWord(word []rune) bool {
	l := len(word)
	p := n
	for i := 0; i < l; i++ {
		next, ok := p.getChild(word[i])
		if ok {
			if next.end {
				return false
			} else {
				p = next
			}
		} else {
			if p != n {
				p = n
				i--
			}
		}
	}
	return true
}

type findWordStat struct {
	index [][2]int
	pos   int
	mark  int
	isEnd bool
	root  *Node
}

func newFindWordStat(n *Node) *findWordStat {
	return &findWordStat{
		index: [][2]int{},
		mark:  -1,
		root:  n,
	}
}
func (self *findWordStat) setRoot(n *Node) {
	self.root = n
}
func (self *findWordStat) setPos(i int) {
	self.pos = i
}
func (self *findWordStat) setMark(i int) {
	self.mark = i
}
func (self *findWordStat) flush() {
	if self.mark >= self.pos {
		self.flushPushIndex()
	} else {
		self.flushNone()
	}
}
func (self *findWordStat) flushNone() {
	self.pos++
	self.mark = -1
}
func (self *findWordStat) reFind(n *Node) {
	self.flush()
	self.setRoot(n)
}
func (self *findWordStat) findNextNode(key rune) (*Node, bool) {
	return self.root.getChild(key)
}
func (self *findWordStat) flushPushIndex() {
	index := [2]int{self.pos, self.mark}
	self.index = append(self.index, index)
	self.pos = self.mark + 1
	self.mark = -1
}
func (self *findWordStat) getIndex() [][2]int {
	return self.index
}
func (n *Node) findWord(word []rune) [][2]int {
	l := len(word)
	finder := newFindWordStat(n)
	for i := 0; i < l; i++ {
		next, ok := finder.findNextNode(word[i])
		if ok {
			finder.setRoot(next)
			if next.end {
				finder.setMark(i)
				if i == l-1 {
					finder.flush()
					return finder.getIndex()
				}
			}
		} else {
			i--
			finder.reFind(n)
		}
	}
	return finder.getIndex()
}
func replace(word []rune, mark rune) []rune {
	result := root.findWord(word)
	if len(result) > 0 {
		for _, v := range result {
			for i := v[0]; i <= v[1]; i++ {
				word[i] = mark
			}
		}
	}
	return word
}

var root = newNode()

func Configurate(fileName string) {
	word := loadWord(fileName)
	for _, data := range word {
		root.passerWord([]rune(data))
	}
}

func loadWord(filename string) []string {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0755)
	defer file.Close()
	if err != nil {
		log.Panicln("Open UserID File, Error:", err.Error())
	} else {
		data, err := ioutil.ReadAll(file)
		if nil != err {
			log.Panicln("Read UserID File, Error:", err.Error())
		}
		return readFileLine(string(data))
	}
	return []string{}
}

func readFileLine(str string) []string {
	win := strings.Split(str, "\r\n")
	unix := strings.Split(str, "\n")
	if len(win) == len(unix) {
		return win
	} else {
		return unix
	}
}

func ValidWord(plain string) bool {
	return root.checkWord([]rune(plain))
}
