package lru

import (
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestAdd(t *testing.T) {
	lru := New(int64(16), nil)   
	lru.Add("key1", String("1234"))
	lru.Add("key2", String("5678"))
	t.Logf("链表的长度为%d\n", lru.Len())
 //注意lru的最大内存限制为16字节，key1和val1是字符串，用的uniocde编码，一个字符占4个字节，
 // 所以key1和val1的长度都是8字节，加起来正好是16字节，所以不会触发删除操作
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed\n")
	}

	if v, ok := lru.Get("key2"); !ok || string(v.(String)) != "5678" {
		t.Fatalf("cache hit key2=5678 failed\n")
	}
}

func TestRemoveOldest(t *testing.T) {
	lru := New(int64(10), nil)
	lru.Add("key1", String("123456"))
	lru.Add("key2", String("5678"))
	lru.Add("key3", String("9"))

	if _, ok := lru.Get("key1"); ok {
		t.Fatalf("RemoveOldest key1 failed")
	}
}

func TestOnClear(t *testing.T) {
	keys := make([]string, 0)
	onClear := func(key string, value Value) {
		keys = append(keys, key)
	}

	lru := New(int64(16), onClear)
	lru.Add("key1", String("123456"))
	lru.Add("key2", String("5678"))
	lru.Add("key3", String("9"))

	if len(keys) != 1 || keys[0] != "key1" {
		t.Fatalf("调用回调函数,后返回的结果有 %v", keys)
	}
}

func TestGet(t *testing.T) {
	lru := New(int64(10), nil)
	lru.Add("key1", String("1234"))

	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("没有命中缓存key1")
	}

	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("命中了缓存key2")
	}
}
