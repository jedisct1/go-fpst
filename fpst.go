package fpst

type FPST struct {
	children []FPST
	key      []byte
	idx      uint16
	bitmap   uint16
	val      interface{}
}

func popcount(w uint32) uint {
	w -= (w >> 1) & 0x55555555
	w = (w & 0x33333333) + ((w >> 2) & 0x33333333)
	w = (w + (w >> 4)) & 0x0f0f0f0f
	w = (w * 0x01010101) >> 24
	return uint(w)
}

func quadbitAt(str []byte, i uint) byte {
	if i/2 >= uint(len(str)) {
		return 0
	}
	c := str[i/2]
	if i&1 == 0 {
		c >>= 4
	}
	return 1 + (c & 0xf)
}

func (t *FPST) bitmapIsSet(bit uint) bool {
	return (t.bitmap & (uint16(1) << bit)) != 0
}

func (t *FPST) bitmapSet(bit uint) {
	t.bitmap |= (uint16(1) << bit)
}

func (t *FPST) actualIndex(i uint) uint {
	b := t.bitmap & ((uint16(1) << i) - 1)
	return popcount(uint32(b))
}

func (t *FPST) childGet(i uint) *FPST {
	if !t.bitmapIsSet(i) {
		return nil
	}
	return &t.children[t.actualIndex(i)]
}

func (t *FPST) childSet(v *FPST, i uint) {
	previous := t.childGet(i)
	if previous != nil {
		*previous = *v
		return
	}
	count := popcount(uint32(t.bitmap)) + 1
	t.children = append(t.children, FPST{})
	ri := t.actualIndex(i)
	rcount := count - ri - 1
	if rcount > 0 {
		for j := uint(0); j < rcount; j++ {
			t.children[ri+1+j] = t.children[ri+j]
		}
	}
	t.children[ri] = *v
	t.bitmapSet(i)
}

func New() *FPST {
	return nil
}

func (trie *FPST) Insert(key []byte, val interface{}) *FPST {
	keyLen := uint(len(key))
	if trie == nil {
		return &FPST{
			key:      key,
			val:      val,
			idx:      0,
			bitmap:   0,
			children: make([]FPST, 1),
		}
	}
	t := trie
	i := uint(0)
	j := uint(0)
	c := byte(0)
	for {
		lk := t.key
		lkLen := uint(len(lk))
		minKeyLen := lkLen
		if keyLen < minKeyLen {
			minKeyLen = keyLen
		}
		x := byte(0)
		for ; j < minKeyLen; j++ {
			x = lk[j] ^ key[j]
			if x != 0 {
				break
			}
		}
		if x == 0 {
			if keyLen == lkLen {
				t.val = val
				return trie
			}
			x = 0xff
		}
		i = j * 2
		if (x & 0xf0) == 0 {
			i++
		}
		if t.bitmap == 0 {
			/* keep index from the new key */
		} else if i >= uint(t.idx) {
			i = uint(t.idx)
			j = i / 2
		} else {
			savedNode := *t
			t.key = key
			t.val = val
			t.idx = uint16(i)
			t.bitmap = 0
			t.children = nil
			c := quadbitAt(lk, i)
			t.childSet(&savedNode, uint(c))
			return trie
		}
		c = quadbitAt(key, i)
		if !t.bitmapIsSet(uint(c)) {
			break
		}
		t = t.childGet(uint(c))
	}
	t.idx = uint16(i)
	newNode := FPST{
		key:      key,
		val:      val,
		idx:      0,
		bitmap:   0,
		children: nil,
	}
	t.childSet(&newNode, uint(c))
	return trie
}

func (trie *FPST) StartsWithExistingKey(str []byte) (foundKey *[]byte, foundVal interface{}) {
	if trie == nil {
		return nil, nil
	}
	strLen := uint(len(str))
	t := trie
	j := uint(0)
	for {
		lk := t.key
		lkLen := uint(len(lk))
		for ; j < strLen; j++ {
			if j >= lkLen {
				return &t.key, t.val
			}
			if lk[j] != str[j] {
				break
			}
		}
		if j >= lkLen {
			return &t.key, t.val
		}
		if t.bitmap == 0 {
			break
		}
		i := uint(t.idx)
		if i > strLen*2 {
			break
		}
		if j > i/2 {
			j = i / 2
		}
		c := quadbitAt(str, i)
		if !t.bitmapIsSet(uint(c)) {
			if t.bitmapIsSet(0) {
				c = 0
			} else {
				break
			}
		}
		t = t.childGet(uint(c))
	}
	return
}

func (trie *FPST) HasKey(key []byte) interface{} {
	foundKey, foundVal := trie.StartsWithExistingKey(key)
	if foundKey != nil && len(*foundKey) == len(key) {
		return foundVal
	}
	return nil
}
