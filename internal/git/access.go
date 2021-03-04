package git

import "sync"

// Thread safe store of {[gitRepo]: token}
type accessTokenMap struct {
	sync.RWMutex
	tokens map[string]string
}

// Get is a thread safe getter to do a map lookup in a getAccessTokens
func (t *accessTokenMap) Get(repo string) (string, bool) {
	t.RLock()
	defer t.RUnlock()

	token, exists := t.tokens[repo]
	return token, exists
}

// Set is a thread safe setter method to modify a gitAccessTokenMap
func (t *accessTokenMap) Set(repo, token string) {
	t.Lock()
	defer t.Unlock()

	t.tokens[repo] = token
}

// AccessTokens is a thread-safe global store of Personal Access Tokens which
// is used to store PATs as they are discovered throughout the Install lifecycle
var AccessTokens = accessTokenMap{
	tokens: map[string]string{},
}
