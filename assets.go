// Package assets 是资源引用计数：重复加载只算一次、归零才卸载、每帧卸载预算与 LRU。
// 缺陷：重复加载重复计数、释放一次就进卸载队列、卸载预算不生效、卸载按释放顺序
// 而不是最久未用、每次查询扫描全部资源。
package assets

// Store 是资源库。
type Store struct {
	budget  int
	refs    map[string]int
	handles map[int]string
	loaded  map[string]bool
	pending []string
	lastUse map[string]int64
	next    int
	loads   int
	unloads int
	scanned int
}

// New 建资源库，budget 是每帧最多卸载多少个。
func New(budget int) *Store {
	return &Store{budget: budget, refs: map[string]int{}, handles: map[int]string{},
		loaded: map[string]bool{}, lastUse: map[string]int64{}}
}

// Load 加载一个资源并返回句柄。
// 缺陷：每次调用都让 Loads 加一。
func (s *Store) Load(name string) int {
	s.loads++
	if !s.loaded[name] {
		s.loaded[name] = true
	}
	s.refs[name]++
	s.next++
	s.handles[s.next] = name
	return s.next
}

// Use 记录一次使用时间。
func (s *Store) Use(name string, nowMs int64) bool {
	if !s.loaded[name] {
		return false
	}
	s.lastUse[name] = nowMs
	return true
}

// Release 释放一个句柄。
// 缺陷：释放一次就进卸载队列，不看引用计数；无效句柄也返回 true。
func (s *Store) Release(handle int) bool {
	name, ok := s.handles[handle]
	if !ok {
		return true
	}
	delete(s.handles, handle)
	s.pending = append(s.pending, name)
	return true
}

// Tick 推进一帧的卸载。
// 缺陷：一次把待卸载的全部卸掉，而且按释放顺序而不是最久未用。
func (s *Store) Tick() int {
	count := len(s.pending)
	for _, name := range s.pending {
		s.unload(name)
	}
	s.pending = nil
	return count
}

func (s *Store) unload(name string) {
	if !s.loaded[name] {
		return
	}
	delete(s.loaded, name)
	delete(s.refs, name)
	s.unloads++
}

// RefCount 是资源当前引用计数。
// 缺陷：扫描全部资源。
func (s *Store) RefCount(name string) int {
	s.scanned += len(s.refs)
	return s.refs[name]
}

// Loaded 是资源当前是否在内存里。
func (s *Store) Loaded(name string) bool { return s.loaded[name] }

// Loads 是累计加载次数。
func (s *Store) Loads() int { return s.loads }

// Unloads 是累计卸载次数。
func (s *Store) Unloads() int { return s.unloads }

// Pending 是待卸载数量。
func (s *Store) Pending() int { return len(s.pending) }

// Scanned 是累计扫描量（规模观测）。
func (s *Store) Scanned() int { return s.scanned }
