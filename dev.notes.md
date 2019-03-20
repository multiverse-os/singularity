

func FDZero(set *FdSet) {
	for i := range set.Bits {
		set.Bits[i] = 0
	}
}
