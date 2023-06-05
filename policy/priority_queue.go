package policy

type priorityqueue []*lfuEntry
type lfuEntry struct {
	e     entry
	count int
}

func (pq priorityqueue) Less(i, j int) bool {

	if pq[i].count == pq[j].count {
		return pq[i].e.expires.Before(*pq[j].e.expires)
	}

	return pq[i].count < pq[j].count
}

func (pq priorityqueue) Len() int {
	return len(pq)
}

func (pq priorityqueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
