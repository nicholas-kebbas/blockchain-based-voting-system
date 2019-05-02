package signature

import "crypto/ecdsa"

/* The signature is composed of the content we're sending, and the Public Key of the sender */

/* So since node to node transactions won't necessarily be a thing we need,
we'll assume we push this up to the blockchain immediately once a vote has been cast.
The block of the blockchain then needs to have a signature parameter as well, which we derive
in handler.

We'll have a public key there so others can verify the integrity of the transaction
  */
type Signature struct {
	content string
	PublicKey ecdsa.PublicKey
}
