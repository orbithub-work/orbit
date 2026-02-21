package services

import (
	"container/list"
	"sync"

	"media-assistant-os/internal/models"
)

// AssetCache 提供内存查找以减少扫描期间的数据库压力。
// 实现了简单的 LRU (Least Recently Used) 淘汰策略。
// 通过 sync.RWMutex 实现线程安全。
type AssetCache struct {
	mu      sync.RWMutex
	
	// 数据存储
	byPath  map[string]*list.Element
	byID    map[string]*list.Element
	
	// LRU 链表：Front 是最近使用，Back 是最久未使用
	lruList *list.List
	
	// 指纹索引 (不做强一致性淘汰，只作为辅助索引)
	byFP    map[string][]string // fingerprint → []assetID
	
	maxSize int
}

// NewAssetCache 创建资产缓存实例
func NewAssetCache(maxSize int) *AssetCache {
	return &AssetCache{
		byPath:  make(map[string]*list.Element, maxSize),
		byID:    make(map[string]*list.Element, maxSize),
		lruList: list.New(),
		byFP:    make(map[string][]string, maxSize/10),
		maxSize: maxSize,
	}
}

// GetByPath 根据路径获取缓存资产
func (c *AssetCache) GetByPath(path string) (*models.Asset, bool) {
	c.mu.Lock() // 需要 Lock 因为会移动链表节点
	defer c.mu.Unlock()
	
	if elem, ok := c.byPath[path]; ok {
		c.lruList.MoveToFront(elem) // 标记为最近使用
		return elem.Value.(*models.Asset), true
	}
	return nil, false
}

// GetByID 根据ID获取缓存资产
func (c *AssetCache) GetByID(id string) (*models.Asset, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if elem, ok := c.byID[id]; ok {
		c.lruList.MoveToFront(elem) // 标记为最近使用
		return elem.Value.(*models.Asset), true
	}
	return nil, false
}

// GetByFingerprint 根据指纹获取缓存资产ID列表
// 注意：这个方法不会更新 LRU 状态，因为它返回的是 ID 列表
func (c *AssetCache) GetByFingerprint(fp string) ([]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ids, ok := c.byFP[fp]
	return ids, ok
}

// Put 将资产放入缓存
func (c *AssetCache) Put(a *models.Asset) {
	if a == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	// 1. 如果已存在，更新值并移到最前
	if elem, ok := c.byID[a.ID]; ok {
		c.lruList.MoveToFront(elem)
		elem.Value = a
		// 也要更新 byPath 映射，防止路径变了
		c.byPath[a.Path] = elem
		return
	}

	// 2. 如果不存在，检查容量
	if c.lruList.Len() >= c.maxSize {
		c.removeOldest()
	}

	// 3. 插入新元素
	elem := c.lruList.PushFront(a)
	c.byID[a.ID] = elem
	c.byPath[a.Path] = elem

	if a.Fingerprint != nil && *a.Fingerprint != "" {
		c.byFP[*a.Fingerprint] = appendUnique(c.byFP[*a.Fingerprint], a.ID)
	}
}

// removeOldest 移除最久未使用的元素 (需在锁内调用)
func (c *AssetCache) removeOldest() {
	elem := c.lruList.Back()
	if elem == nil {
		return
	}
	c.lruList.Remove(elem)
	asset := elem.Value.(*models.Asset)
	
	// 清理索引
	delete(c.byID, asset.ID)
	delete(c.byPath, asset.Path)
	
	// 注意：byFP 清理比较麻烦，这里暂时不做，因为 string list 占用内存很小
	// 且指纹索引是辅助性的。如果非要严谨，需要遍历 slice 删除 ID。
}

// UpdateFingerprint 更新缓存中资产的指纹
func (c *AssetCache) UpdateFingerprint(id string, fp string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if elem, ok := c.byID[id]; ok {
		// 移动到最近使用
		c.lruList.MoveToFront(elem)
		
		asset := elem.Value.(*models.Asset)
		asset.Fingerprint = &fp
		c.byFP[fp] = appendUnique(c.byFP[fp], id)
	}
}

// Invalidate 使缓存中的资产失效
func (c *AssetCache) Invalidate(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if elem, ok := c.byID[id]; ok {
		c.lruList.Remove(elem)
		asset := elem.Value.(*models.Asset)
		delete(c.byID, id)
		delete(c.byPath, asset.Path)
	}
}

// Clear 清空缓存
func (c *AssetCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.lruList.Init()
	c.byPath = make(map[string]*list.Element, c.maxSize)
	c.byID = make(map[string]*list.Element, c.maxSize)
	c.byFP = make(map[string][]string, c.maxSize/10)
}

// Size 返回缓存中的资产数量
func (c *AssetCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.byPath)
}

// appendUnique 向切片添加唯一值
func appendUnique(slice []string, val string) []string {
	for _, s := range slice {
		if s == val {
			return slice
		}
	}
	return append(slice, val)
}
