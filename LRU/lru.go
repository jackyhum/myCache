package lru

import "container/list"

type Value interface {
	Len() int
}

//LRU 缓存部分的代码，无 并发安全
type Cache struct {
	maxBytes int64
	nbytes   int64
	ll       *list.List    					//Go标准库中的双向链表
	cache    map[string]*list.Element   	//键是字符串，值是双向链表中对应节点的指针
	OnClear func(key string, value Value)   //清楚目标节点时的回调函数
}

//entry 双向链表节点存储了键+值
type entry struct {
	key   string
	value Value
}


// New 构造函数
func New(maxBytes int64, onClear func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnClear: onClear,
	}
}

// 添加值
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok { //已经存在，则更新值，并将该节点移到队首
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else { //不存在，则添加新节点
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	//每次添加时检查是否超过最大内存限制
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}


func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnClear != nil {
			c.OnClear(kv.key, kv.value)
		}
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}