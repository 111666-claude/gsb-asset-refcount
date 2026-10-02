package assets

import "testing"

func TestLoadReturnsHandle(t *testing.T) {
	if handle := New(2).Load("hero"); handle == 0 {
		t.Fatal("加载成功应该返回非零句柄")
	}
}

func TestUnknownRefCountIsZero(t *testing.T) {
	if New(2).RefCount("ghost") != 0 {
		t.Fatal("没有加载过的资源引用计数应该是 0")
	}
}

func TestUnknownIsNotLoaded(t *testing.T) {
	if New(2).Loaded("ghost") {
		t.Fatal("没有加载过的资源不应该算已加载")
	}
}

func TestReleaseValidHandle(t *testing.T) {
	store := New(2)
	handle := store.Load("hero")
	if !store.Release(handle) {
		t.Fatal("释放有效句柄应该成功")
	}
}

func TestCountersStartAtZero(t *testing.T) {
	store := New(2)
	if store.Loads() != 0 || store.Unloads() != 0 || store.Scanned() != 0 {
		t.Fatal("计数器初始都应该是 0")
	}
}
