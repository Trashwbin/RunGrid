//go:build !windows

package scanner

func lookupDisplayName(_ string) (string, bool) {
	return "", false
}
