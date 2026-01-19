package msdk

import (
	"fmt"
	"kernel/kenum"
)

var (
	MsdkUrl string
)

func SetMsdkUrl(isTest bool, appId string) {
	if isTest {
		MsdkUrl = kenum.Msdk_Test_Domain
	} else {
		MsdkUrl = fmt.Sprintf(kenum.Msdk_Dev_Domain, appId)
	}
}
