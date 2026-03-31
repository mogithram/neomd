package ui

import (
	"regexp"
	"sort"
	"strings"

	"github.com/sspaeti/neomd/internal/imap"
)

// replyPrefixRe matches common reply/forward prefixes (Re:, Fwd:, Fw:, AW:, SV:, VS:).
var replyPrefixRe = regexp.MustCompile(`(?i)^(re|fwd?|aw|sv|vs)\s*(\[\d+\])?\s*:\s*`)

// normalizeSubject strips reply/forward prefixes and lowercases for thread grouping.
func normalizeSubject(subject string) string {
	s := strings.TrimSpace(subject)
	for {
		stripped := replyPrefixRe.ReplaceAllString(s, "")
		stripped = strings.TrimSpace(stripped)
		if stripped == s {
			break
		}
		s = stripped
	}
	return strings.ToLower(s)
}

// threadedEmail pairs an email with its tree-drawing prefix for the inbox list.
type threadedEmail struct {
	email        imap.Email
	threadPrefix string // e.g. "┌─>" or "  " or ""
}

// threadEmails groups and reorders emails into threaded display order.
// Each thread is sorted internally by date ascending (oldest = root at bottom,
// newest replies on top), matching neomutt's reverse-thread style.
// Threads are sorted relative to each other by the most recent email in each thread.
// Returns the reordered emails with tree prefixes for rendering.
func threadEmails(emails []imap.Email) []threadedEmail {
	if len(emails) == 0 {
		return nil
	}

	// Phase 1: Build groups using InReplyTo/MessageID + subject+participant fallback.
	parent := make([]int, len(emails))
	for i := range parent {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	union := func(a, b int) {
		ra, rb := find(a), find(b)
		if ra != rb {
			parent[ra] = rb
		}
	}

	// Connect by InReplyTo -> MessageID.
	byMsgID := make(map[string]int, len(emails))
	for i := range emails {
		if id := emails[i].MessageID; id != "" {
			byMsgID[id] = i
		}
	}
	for i := range emails {
		if replyTo := emails[i].InReplyTo; replyTo != "" {
			if j, ok := byMsgID[replyTo]; ok {
				union(i, j)
			}
		}
	}

	// Subject+participant fallback: group emails with same normalized subject
	// where participants overlap (From of one appears in From/To of another).
	type subjGroup struct {
		indices []int
		from    string // extractAddrFromField of first email
		to      string
	}
	bySubject := make(map[string]*subjGroup)
	for i := range emails {
		subj := normalizeSubject(emails[i].Subject)
		if subj == "" {
			continue
		}
		from := extractAddrFromField(emails[i].From)
		to := extractAddrFromField(emails[i].To)

		if g, ok := bySubject[subj]; ok {
			// Check participant overlap with existing group members.
			if from == g.from || from == g.to || to == g.from || to == g.to {
				union(i, g.indices[0])
				g.indices = append(g.indices, i)
			}
		} else {
			bySubject[subj] = &subjGroup{indices: []int{i}, from: from, to: to}
		}
	}

	// Phase 2: Collect threads.
	threadMap := make(map[int][]int) // root -> indices
	for i := range emails {
		root := find(i)
		threadMap[root] = append(threadMap[root], i)
	}

	// Phase 3: Sort each thread internally by date ascending (oldest first = root).
	type thread struct {
		indices   []int
		newestIdx int // index of most recent email (for inter-thread sorting)
	}
	var threads []thread
	for _, indices := range threadMap {
		sort.Slice(indices, func(a, b int) bool {
			return emails[indices[a]].Date.Before(emails[indices[b]].Date)
		})
		newest := indices[len(indices)-1]
		threads = append(threads, thread{indices: indices, newestIdx: newest})
	}

	// Phase 4: Sort threads by most recent email (matching the overall sort: newest first).
	sort.Slice(threads, func(i, j int) bool {
		return emails[threads[i].newestIdx].Date.After(emails[threads[j].newestIdx].Date)
	})

	// Phase 5: Build output with thread connector lines (Twitter-style).
	// Newest reply on top, root (oldest) at bottom.
	// │ = continuation (more thread below), ╰ = root/last in thread.
	result := make([]threadedEmail, 0, len(emails))
	for _, t := range threads {
		n := len(t.indices)
		if n == 1 {
			result = append(result, threadedEmail{email: emails[t.indices[0]]})
			continue
		}
		// Reverse order: newest first, oldest (root) last.
		for k := n - 1; k >= 0; k-- {
			prefix := "│"
			if k == 0 {
				prefix = "╰" // root = bottom of thread
			}
			result = append(result, threadedEmail{
				email:        emails[t.indices[k]],
				threadPrefix: prefix,
			})
		}
	}

	return result
}

// extractAddrFromField returns the bare email address from a "Name <addr>" or "addr" string.
func extractAddrFromField(s string) string {
	if i := strings.LastIndex(s, "<"); i >= 0 {
		if j := strings.Index(s[i:], ">"); j >= 0 {
			return strings.ToLower(s[i+1 : i+j])
		}
	}
	return strings.ToLower(strings.TrimSpace(s))
}
