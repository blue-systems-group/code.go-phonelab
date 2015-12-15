package phonelab

type ByTime []Logline

func (a ByTime) Len() int      { return len(a) }
func (a ByTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool {
	if a[i].timestamp.Equal(a[j].timestamp) {
		return a[i].fileorder < a[j].fileorder
	} else {
		return a[i].timestamp.Before(a[j].timestamp)
	}
}
