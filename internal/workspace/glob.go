package workspace

import "path"

// matchAny reports whether rel (a root-relative, forward-slash path)
// matches at least one of globs. An empty globs list matches every path,
// so a caller with no declared globs still gets the default-exclude
// boundary rather than an empty collection.
func matchAny(globs []string, rel string) bool {
	if len(globs) == 0 {
		return true
	}
	for _, g := range globs {
		if matchGlob(g, rel) {
			return true
		}
	}
	return false
}

// matchGlob reports whether rel matches pattern, both given as forward-
// slash paths. "**" matches zero or more whole path segments; any other
// segment is matched with shell-style path.Match (requirement R3).
func matchGlob(pattern, rel string) bool {
	return matchSegments(splitPath(pattern), splitPath(rel))
}

func splitPath(p string) []string {
	if p == "" {
		return nil
	}
	var segs []string
	start := 0
	for i := 0; i < len(p); i++ {
		if p[i] == '/' {
			segs = append(segs, p[start:i])
			start = i + 1
		}
	}
	segs = append(segs, p[start:])
	return segs
}

func matchSegments(pat, seg []string) bool {
	if len(pat) == 0 {
		return len(seg) == 0
	}
	if pat[0] == "**" {
		if matchSegments(pat[1:], seg) {
			return true
		}
		if len(seg) == 0 {
			return false
		}
		return matchSegments(pat, seg[1:])
	}
	if len(seg) == 0 {
		return false
	}
	ok, err := path.Match(pat[0], seg[0])
	if err != nil || !ok {
		return false
	}
	return matchSegments(pat[1:], seg[1:])
}
