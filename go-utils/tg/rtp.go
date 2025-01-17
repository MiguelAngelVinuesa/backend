package tg

import (
	"fmt"
)

// VerifyRTP verifies the given RTP if it is supported by our game engine.
// This verification is irrespective of the gameID.
func VerifyRTP(rtp uint8) error {
	switch rtp {
	case 92, 94, 96: // always good.
		return nil
	case 42: // DEBUG MODE ONLY!
		return nil
	default:
		return fmt.Errorf("VerifyRTP - invalid RTP [%d]", rtp)
	}
}
