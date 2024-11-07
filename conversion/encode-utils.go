package conversion

func EncodeSlice[Source, Target any, F func(*Source) *Target](source []*Source, convert F) []*Target {
	slice := make([]*Target, len(source))

	for i, item := range source {
		slice[i] = convert(item)
	}

	return slice
}
