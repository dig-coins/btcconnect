package btctx

type SignHelper interface {
	GetPrivateKey(address string) (string, error)
}

func NewFixSignHelper(wifPrivateKey string) SignHelper {
	return &fixSignHelper{
		wifPrivateKey: wifPrivateKey,
	}
}

type fixSignHelper struct {
	wifPrivateKey string
}

func (impl *fixSignHelper) GetPrivateKey(_ string) (string, error) {
	return impl.wifPrivateKey, nil
}
