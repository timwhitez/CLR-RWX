package main

import (
	"bufio"
	"fmt"
	clr "github.com/ropnop/go-clr"
	"log"
	"os"
	"syscall"
)

func checkOK(hr uintptr, caller string) {
	if hr != 0x0 {
		log.Fatalf("%s returned 0x%08x", caller, hr)
	}
}

func main() {
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	var pMetaHost uintptr
	hr := clr.CLRCreateInstance(&clr.CLSID_CLRMetaHost, &clr.IID_ICLRMetaHost, &pMetaHost)
	checkOK(hr, "CLRCreateInstance")
	metaHost := clr.NewICLRMetaHostFromPtr(pMetaHost)

	versionString := "v4.0.30319"
	pwzVersion, _ := syscall.UTF16PtrFromString(versionString)
	var pRuntimeInfo uintptr
	hr = metaHost.GetRuntime(pwzVersion, &clr.IID_ICLRRuntimeInfo, &pRuntimeInfo)
	checkOK(hr, "metahost.GetRuntime")
	runtimeInfo := clr.NewICLRRuntimeInfoFromPtr(pRuntimeInfo)

	var isLoadable bool
	hr = runtimeInfo.IsLoadable(&isLoadable)
	checkOK(hr, "runtimeInfo.IsLoadable")
	if !isLoadable {
		log.Fatal("[!] IsLoadable returned false. Bailing...")
	}

	hr = runtimeInfo.BindAsLegacyV2Runtime()
	checkOK(hr, "runtimeInfo.BindAsLegacyV2Runtime")

	var pRuntimeHost uintptr
	hr = runtimeInfo.GetInterface(&clr.CLSID_CorRuntimeHost, &clr.IID_ICorRuntimeHost, &pRuntimeHost)
	runtimeHost := clr.NewICORRuntimeHostFromPtr(pRuntimeHost)
	hr = runtimeHost.Start()
	checkOK(hr, "runtimeHost.Start")

	fmt.Println("[+] Loaded CLR into this process")

	var pIUnknown uintptr
	runtimeHost.GetDefaultDomain(&pIUnknown)

	runtimeHost.Release()
	runtimeInfo.Release()
	metaHost.Release()

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

}
