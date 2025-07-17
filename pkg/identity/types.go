package identity

type DidDoc struct {
	Context             []string                   `json:"@context"`
	Id                  string                     `json:"id"`
	AlsoKnownAs         []string                   `json:"alsoKnownAs,omitempty"`
	VerificationMethods []DidDocVerificationMethod `json:"verificationMethods,omitempty"`
	Service             []DidDocService            `json:"service,omitempty"`
}

type DidDocVerificationMethod struct {
	Id                 string `json:"id"`
	Type               string `json:"type"`
	Controller         string `json:"controller"`
	PublicKeyMultibase string `json:"publicKeyMultibase"`
}

type DidDocService struct {
	Id              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}
