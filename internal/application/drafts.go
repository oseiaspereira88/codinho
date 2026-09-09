package application

import "github.com/oseiaspereira88/codinho/internal/drafts"

// Drafts shares the local event store but never mutates the published catalog.
func (s *SessionService) Drafts() *drafts.Service { return s.drafts }
