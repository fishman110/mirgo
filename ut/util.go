package ut

import (
	"math/rand"
	"os"
	"path"
	"path/filepath"
)

func AbsInt(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// 随机 [low, high]
func RandomInt(low int, high int) int {
	if low == high {
		return low
	}

	return rand.Intn(high-low+1) + low
}

// c# random.next [low, high)
func RandomNext2(low, high int) int {
	return RandomInt(low, high-1)
}

// c# random.next [0, high)
func RandomNext(high int) int {
	return RandomNext2(0, high)
}

func RandomString(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		b := rand.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

func GetFiles(dir string, allow []string) []string {

	allowMap := map[string]bool{}
	if allow != nil {
		for _, v := range allow {
			allowMap[v] = true
		}
	}

	ret := []string{}
	filepath.Walk(dir, func(fpath string, f os.FileInfo, err error) error {
		if f == nil || f.IsDir() {
			return nil
		}

		ext := path.Ext(fpath)
		if allowMap[ext] {
			ret = append(ret, filepath.ToSlash(fpath))
		}

		return nil
	})

	return ret
}