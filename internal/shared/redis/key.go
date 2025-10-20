package redis

type KeyPrefixes struct {
	PrefixBillingInquiry     string
	PrefixSaveBillingInquiry string
}

func (k *KeyPrefixes) SetPrefixBillingInquiry(prefix string) {
	k.PrefixBillingInquiry = prefix
}

func (k *KeyPrefixes) SetPrefixSaveBilling(prefix string) {
	k.PrefixSaveBillingInquiry = prefix
}
