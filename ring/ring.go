package ring

/* Ring signatures allow one to sign messages on behalf of a ring of public keys */

type Signature struct {
	/* The ring of public Keys */
	ring []string
}

/* Need to adjust sign and verification from to implement anonymity */
func Sign() {

}
func Verify(message []byte) bool {
	return false
}

