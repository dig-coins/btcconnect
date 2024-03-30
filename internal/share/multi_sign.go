package share

type MultiSignAddressInfo struct {
	PublicKeys []string `json:"public_keys" yaml:"public_keys"`
	MinSignNum int      `json:"min_sign_num" yaml:"min_sign_num"`
}
