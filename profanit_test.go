package src

import (
	"bytes"
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"
)

type Node struct {
	child map[rune]*Node
	end   bool
}

func newNode() *Node {
	return &Node{}
}

func (n *Node)checkAndFixMap(){
	if nil == n.child {
		n.child=map[rune]*Node{}
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
	node, ok := n.child[key]
	return node, ok
}

func (n *Node) setNodeEnd() {
	n.end = true
}

func (n *Node) paserWord(word []rune) {
	l := len(word)
	next := n
	for i := 0; i < l; i++ {
		next = next.setChild(word[i])
		// if next.end {
		// 	return
		// }
	}
	if(next != n){
		next.setNodeEnd()
	}
}

func (n *Node) checkWord(word []rune) bool {
	l := len(word)
	next := n
	ok := false
	for i := 0; i < l; i++ {
		next, ok = next.getChild(word[i])
		if ok {
			if next.end {
				return true
			}
		} else {
			return false
		}
	}
	return false
}

type findWordStat struct{
	index [][2]int
	pos int
	mark int
	isEnd bool
	root *Node
}

func newFindWordStat(n *Node)*findWordStat{
	return &findWordStat{
		index :[][2]int{},
		mark:-1,
		root:n,
	}
}

func (self *findWordStat)setRoot(n *Node){
	self.root = n
}
func (self *findWordStat)setPos(i int){
	self.pos = i
}
func (self *findWordStat)setMark(i int){
	self.mark = i
}

func (self *findWordStat)flush(){
	if self.mark >= self.pos{
		self.flushPushIndex()
	}else{
		self.flushNone()
	}
}

func (self *findWordStat)flushNone(){
	self.pos++
	self.mark = -1
}

func (self *findWordStat)reFind(n *Node){
	self.flush()
	self.setRoot(n)
}

func (self *findWordStat)findNextNode(key rune) (*Node, bool){
	return self.root.getChild(key)
}

func (self *findWordStat)flushPushIndex(){
	index := [2]int{self.pos,self.mark}
	self.index = append(self.index, index)
	self.pos = self.mark + 1
	self.mark = -1
}

func (self *findWordStat)getIndex() [][2]int{
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
	// fmt.Println(string(word), result)
	if len(result) > 0 {
		for _, v := range result {
			for i := v[0]; i <= v[1]; i++ {
				word[i] = mark
			}
		}
	}
	return word
}

func (n *Node) Printf(str string) {
	for k, v := range n.child {
		if v.end {
			v.Printf(str+"==>"+string([]rune{k})+"(end)")
		}else{
			v.Printf(str+"==>"+string([]rune{k}))
		}
		if v.child == nil {
			fmt.Println(str + "==>" + string([]rune{k}) + "(end)")
		}
	}
}

var data = [][]byte{
	[]byte("123"),
	[]byte("aba"),
	[]byte("abc"),
	[]byte("av123"),
	[]byte("你好"),
}

var root = initRoot()
var word = initWordData()

func initWordData() map[int][]rune{
	result := map[int][]rune{}
	// fd, _ := os.Open("./words.txt")
	fd, _ := os.Open("./2.txt")
	defer fd.Close()
	r := bufio.NewReader(fd)
	i := 0
	for {
		data, err := readLine(r)
		i++
		if nil == err {
			result[i] = data
		} else {
			return result
		}
	}
}

func initRoot() *Node{
	node := newNode()
	for _, data := range word {
		node.paserWord(data)
	}
	return node
}

func readLine(r *bufio.Reader) ([]rune, error) {
	line, err := r.ReadBytes('\n')
	if err == nil {
		data :=bytes.Runes(line[:len(line)-1])
		return data, err
	}else{
		return []rune{}, err
	}
}



func TestReplaceWord(t *testing.T) {
	// root.Printf("root")
	t1 := time.Now().UnixNano()
	sum := 0
	for _, v := range word {
		src := string(v)
		res := replace(v, '*')
		// replace(v, '*')
		fmt.Println("======")
		fmt.Println(src)
		fmt.Println(string(res))
		if src == string(res) {
			sum++
			fmt.Println(src)
		}
	}
	runeStr:= bytes.Runes([]byte("男男女女哈淫水哈东亚病夫爱液啊"))
	res := replace(runeStr, '*')
	fmt.Println("男男女女哈淫水哈东亚病夫爱液啊")
	fmt.Println(string(res))
	t2 := time.Now().UnixNano()
	fmt.Println("未过滤 :",sum, "总:", len(word), "耗时:", nsToMs(t2-t1),"ms")
}

func equalRune(a, b []rune) bool {
	if len(a) != len(b){
		return false
	}else{
		l := len(a)
		for i := 0; i < l; i++ {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
}

func nsToMs(ns int64) float64{
	return float64(ns)/1000000.0
}

func BenchmarkReplaceWord(b *testing.B){
	for i := 0; i < b.N; i++ {
		replace(word[100], '*')
	}
}
