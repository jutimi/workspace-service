package utils

import "github.com/google/uuid"

func IsDiffSlice[T string | int | uuid.UUID](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return true
	}

	s1Memo := make(map[T]T)
	for _, v := range s1 {
		s1Memo[v] = v
	}

	for _, v := range s2 {
		if _, ok := s1Memo[v]; !ok {
			return true
		}
	}

	return false
}

func IsSliceContain[T string | int | uuid.UUID](s1 T, s2 []T) bool {
	for _, v := range s2 {
		if v == s1 {
			return true
		}
	}

	return false
}
