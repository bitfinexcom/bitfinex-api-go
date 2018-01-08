package tests

type MockNonceGenerator struct {
	nonce string
}

func (m *MockNonceGenerator) Next(nonce string) {
	m.nonce = nonce
}

func (m *MockNonceGenerator) GetNonce() string {
	return m.nonce
}
