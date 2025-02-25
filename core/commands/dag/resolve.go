package dagcmd

import (
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/xbradylee/ipfs-kubo/core/commands/cmdenv"

	cmds "github.com/ipfs/go-ipfs-cmds"
)

func dagResolve(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
	api, err := cmdenv.GetApi(env, req)
	if err != nil {
		return err
	}

	rp, err := api.ResolvePath(req.Context, path.New(req.Arguments[0]))
	if err != nil {
		return err
	}

	return cmds.EmitOnce(res, &ResolveOutput{
		Cid:     rp.Cid(),
		RemPath: rp.Remainder(),
	})
}
