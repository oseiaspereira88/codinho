package mcpserver

import "github.com/oseiaspereira88/codinho/internal/curriculum"

func subjectIDs(legacy string, plural []string) ([]string, error) {
	if legacy != "" {
		if plural != nil {
			return nil, curriculum.ErrInvalidSelection
		}
		plural = []string{legacy}
	}
	return curriculum.NormalizeSubjects(plural)
}
