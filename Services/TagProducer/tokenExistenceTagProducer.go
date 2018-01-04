package TagProducer

type tokenExistenceTagProducer struct {
	tokens map[string]uint32
	tag    string
}

func (p tokenExistenceTagProducer) SatisfyTag(torrentTokens []string) bool {
	tokensLen := len(torrentTokens)
	for i := 0; i < tokensLen; i++ {
		if _, ok := p.tokens[torrentTokens[i]]; ok {
			return true
		}
	}
	return false
}

func (p tokenExistenceTagProducer) GetTag() string {
	return p.tag
}
