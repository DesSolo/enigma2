package secrets

// OptionFunc ...
type OptionFunc func(*Provider)

// WithTokenLength ...
func WithTokenLength(length int) OptionFunc {
	return func(p *Provider) {
		p.tokenLength = length
	}
}

// WithTokenSaveRetries ...
func WithTokenSaveRetries(tokenSaveRetries int) OptionFunc {
	return func(p *Provider) {
		p.tokenSaveRetries = tokenSaveRetries
	}
}
