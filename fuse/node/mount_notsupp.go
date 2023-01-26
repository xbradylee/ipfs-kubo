//go:build (!nofuse && openbsd) || (!nofuse && netbsd) || (!nofuse && plan9)
// +build !nofuse,openbsd !nofuse,netbsd !nofuse,plan9

package node

import (
	"errors"

	core "github.com/xbradylee/ipfs-kubo/core"
)

func Mount(node *core.IpfsNode, fsdir, nsdir string) error {
	return errors.New("FUSE not supported on OpenBSD or NetBSD. See #5334 (https://github.com/xbradylee/ipfs-kubo/issues/5334).")
}
