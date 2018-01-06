package yaku

type Dictionary map[string]string

func AppendResources8(in Dictionary, str string, indexes []uint8) Dictionary {
	for i := 0; i < len(indexes)-1; i++ {
		in[str[indexes[i]:indexes[i+1]]] = ""
	}
	return in
}

func AppendResources16(in Dictionary, str string, indexes []uint16) Dictionary {
	for i := 0; i < len(indexes)-1; i++ {
		in[str[indexes[i]:indexes[i+1]]] = ""
	}
	return in
}

func GetResources() Dictionary {
	mp := make(Dictionary)
	mp = AppendResources16(mp, _Yaku_name, _Yaku_index[:])
	mp = AppendResources8(mp, _Yakuman_name, _Yakuman_index[:])
	mp = AppendResources8(mp, _Fu_name, _Fu_index[:])
	mp = AppendResources8(mp, _Limit_name, _Limit_index[:])
	return mp
}
