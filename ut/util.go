package ut

import (
	"bufio"
	"io"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"
)

const UintMax = ^uint(0)
const IntMax = int(^uint(0) >> 1)

func HasFlagUint16(a, b uint16) bool {
	return a&b != 0
}

func StringEqualFold(a string, b ...string) bool {
	for _, v := range b {
		if strings.EqualFold(a, v) {
			return true
		}
	}
	return false
}

func AbsInt(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
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

// 按空格拆分字符串。如果加了引号，那么认为是一个字符串
func SplitString(s string) []string {

	ret := []string{}

	start := 0
	var stat byte

	for i := 0; i < len(s); i++ {
		if unicode.IsSpace(rune(s[i])) {
			if stat == 1 {
				ret = append(ret, s[start:i])
				stat = 0
			}

		} else if s[i] == '\'' || s[i] == '"' {
			if stat == s[i] {
				ret = append(ret, s[start:i])
				stat = 0
			} else {
				if stat == 0 {
					stat = s[i]
					start = i + 1
				}
			}
		} else {
			if stat == 0 {
				start = i
				stat = 1
			}
		}

	}

	if stat != 0 {
		ret = append(ret, s[start:])
	}

	return ret
}

func ReadLines(filepath string) (lines []string, err error) {

	file, err := os.Open(filepath)
	if err != nil {
		return
	}

	return ReadLinesByReader(file), nil
}

func ReadLinesByReader(r io.Reader) []string {
	lines := []string{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}
