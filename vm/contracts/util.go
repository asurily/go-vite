package contracts

import (
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/vm/contracts/abi"
	"github.com/vitelabs/go-vite/vm/util"
	"github.com/vitelabs/go-vite/vm_context/vmctxt_interface"
)

func IsUserAccount(db vmctxt_interface.VmDatabase, addr types.Address) bool {
	// TODO merge method
	_, code := util.GetContractCode(db, &addr)
	return len(code) == 0
}

func IsExistGid(db vmctxt_interface.VmDatabase, gid types.Gid) bool {
	value := db.GetStorage(&types.AddressConsensusGroup, abi.GetConsensusGroupKey(gid))
	return len(value) > 0
}
