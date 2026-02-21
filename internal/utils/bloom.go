package utils

import (
	"hash/fnv"
	"math"
)

// BloomFilter 是一个简单的布隆过滤器实现
// 用于快速判断一个元素是否“可能存在”于集合中
// 特点：
// 1. 如果返回 false，则元素一定不存在（100% 准确）
// 2. 如果返回 true，则元素可能存在（存在一定的误判率，但在海量去重场景可接受）
// 3. 内存占用极小，不用存储原始字符串
type BloomFilter struct {
	bitset []bool
	k      uint // 哈希函数的数量
	m      uint // 位数组的长度
}

// NewBloomFilter 创建一个新的布隆过滤器
// n: 预期插入的元素数量
// p: 期望的误判率 (例如 0.01 表示 1%)
func NewBloomFilter(n uint, p float64) *BloomFilter {
	// 计算最佳的位数组长度 m
	// m = - (n * ln(p)) / (ln(2)^2)
	m := uint(math.Ceil(float64(n) * math.Log(p) / math.Log(1.0/math.Pow(2.0, math.Log(2.0)))))

	// 计算最佳的哈希函数数量 k
	// k = (m / n) * ln(2)
	k := uint(math.Ceil((float64(m) / float64(n)) * math.Log(2.0)))

	return &BloomFilter{
		bitset: make([]bool, m),
		k:      k,
		m:      m,
	}
}

// Add 向过滤器中添加一个元素
func (bf *BloomFilter) Add(data []byte) {
	h := fnv.New64a()
	h.Write(data)
	hash1 := h.Sum64()

	h2 := fnv.New64()
	h2.Write(data)
	hash2 := h2.Sum64()

	for i := uint(0); i < bf.k; i++ {
		// 使用双重哈希模拟多个哈希函数
		// index = (hash1 + i * hash2) % m
		index := (hash1 + uint64(i)*hash2) % uint64(bf.m)
		bf.bitset[index] = true
	}
}

// AddString 向过滤器中添加一个字符串
func (bf *BloomFilter) AddString(s string) {
	bf.Add([]byte(s))
}

// Contains 检查元素是否可能存在
func (bf *BloomFilter) Contains(data []byte) bool {
	h := fnv.New64a()
	h.Write(data)
	hash1 := h.Sum64()

	h2 := fnv.New64()
	h2.Write(data)
	hash2 := h2.Sum64()

	for i := uint(0); i < bf.k; i++ {
		index := (hash1 + uint64(i)*hash2) % uint64(bf.m)
		if !bf.bitset[index] {
			return false // 只要有一位是 0，则一定不存在
		}
	}
	return true // 全是 1，可能存在
}

// ContainsString 检查字符串是否可能存在
func (bf *BloomFilter) ContainsString(s string) bool {
	return bf.Contains([]byte(s))
}

// Reset 重置过滤器
func (bf *BloomFilter) Reset() {
	for i := range bf.bitset {
		bf.bitset[i] = false
	}
}
