package interfaces

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// Objecter2 extends Objecter with special zjson.Encode
type Objecter2 interface {
	Encode2(enc *zjson.Encoder)
	pool.Objecter
}
