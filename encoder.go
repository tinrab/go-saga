package saga

type Encoder interface {
	Encode() ([]byte, error)
}
