package policy

type priorityqueue []*lfuEntry
type lfuEntry struct {
	index int
	entry entry
	count int
}

func (l *lfuEntry) referenced() {
	l.count++
	l.entry.touch()
}

func (pq priorityqueue) Less(i, j int) bool {

	if pq[i].count == pq[j].count {
		return pq[i].entry.updateAt.Before(*pq[j].entry.updateAt)
	}

	return pq[i].count < pq[j].count
}

func (pq priorityqueue) Len() int {
	return len(pq)
}

func (pq priorityqueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityqueue) Pop() interface{} {
	old := *pq
	n := len(old)
	entry := old[n-1]
	old[n-1] = nil //避免内存泄露
	new := old[0 : n-1]
	for i := 0; i < len(new); i++ {
		new[i].index = i
	}
	*pq = new
	return entry
}

func (pq *priorityqueue) Push(x interface{}) { // 绑定push方法，插入新元素
	entry := x.(*lfuEntry)
	entry.index = len(*pq)
	*pq = append(*pq, x.(*lfuEntry))
}
