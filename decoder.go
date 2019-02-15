package saga

type Decoder interface {
	Decode(data []byte) error
}
