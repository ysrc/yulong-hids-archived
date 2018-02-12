// +build windows

package monitor

/*
#cgo windows LDFLAGS:-lWS2_32 -liphlpapi
#include <stdio.h>
#include <windows.h>
#include <tchar.h>
#include "accctrl.h"
#include "aclapi.h"
#include <tchar.h>
#pragma comment(lib, "advapi32.lib")


char * getprocessowner(char  *filename)
{
    static char outputfilename[256] = { 0 };
	memset(outputfilename, 0, 256);
	DWORD dwRtnCode = 0;
	PSID pSidOwner = NULL;
	BOOL bRtnBool = TRUE;
	LPTSTR AcctName = NULL;
	LPTSTR DomainName = NULL;
	DWORD dwAcctName = 1, dwDomainName = 1;
	SID_NAME_USE eUse = SidTypeUnknown;
	HANDLE hFile;
	PSECURITY_DESCRIPTOR pSD = NULL;
	hFile = CreateFileA(
		filename,
		GENERIC_READ,
		FILE_SHARE_READ,
		NULL,
		OPEN_EXISTING,
		FILE_ATTRIBUTE_NORMAL,
		NULL);

	if (hFile == INVALID_HANDLE_VALUE) {
		DWORD dwErrorCode = 0;
		dwErrorCode = GetLastError();
		return "";
	}

	dwRtnCode = GetSecurityInfo(
		hFile,
		SE_FILE_OBJECT,
		OWNER_SECURITY_INFORMATION,
		&pSidOwner,
		NULL,
		NULL,
		NULL,
		&pSD);
	if (dwRtnCode != ERROR_SUCCESS) {
		DWORD dwErrorCode = 0;
		dwErrorCode = GetLastError();
		CloseHandle(hFile);
		return "";
	}
	bRtnBool = LookupAccountSid(
		NULL,           // local computer
		pSidOwner,
		AcctName,
		(LPDWORD)&dwAcctName,
		DomainName,
		(LPDWORD)&dwDomainName,
		&eUse);
	AcctName = (LPTSTR)GlobalAlloc(
		GMEM_FIXED,
		dwAcctName*2+10);
	if (AcctName == NULL) {
		DWORD dwErrorCode = 0;
		dwErrorCode = GetLastError();
		CloseHandle(hFile);
		return "";
	}

	DomainName = (LPTSTR)GlobalAlloc(
		GMEM_FIXED,
		dwDomainName*2+10);
	if (DomainName == NULL) {
		DWORD dwErrorCode = 0;
		dwErrorCode = GetLastError();
		GlobalFree(AcctName);
		CloseHandle(hFile);
		return "";

	}
	bRtnBool = LookupAccountSid(
		NULL,                   // name of local or remote computer
		pSidOwner,              // security identifier
		AcctName,               // account name buffer
		(LPDWORD)&dwAcctName,   // size of account name buffer
		DomainName,             // domain name
		(LPDWORD)&dwDomainName, // size of domain name buffer
		&eUse);                 // SID type


	if (bRtnBool == FALSE) {
		DWORD dwErrorCode = 0;
		dwErrorCode = GetLastError();
		GlobalFree(AcctName);
		GlobalFree(DomainName);
		CloseHandle(hFile);
		return "";
	}
	else if (bRtnBool == TRUE)
	{
	     GlobalFree(DomainName);
		 CloseHandle(hFile);
		 memcpy(outputfilename,AcctName,wcslen(AcctName)*2);
		 GlobalFree(AcctName);
		 return outputfilename;
    }
	CloseHandle(hFile);
	GlobalFree(AcctName);
	GlobalFree(DomainName);
	return "";
}
*/
import "C"
import (
	"log"
	"runtime"
	"yulong-hids/agent/common"

	"os"
	"strings"

	"github.com/go-fsnotify/fsnotify"
)

// StartFileMonitor 开始文件行为监控
func StartFileMonitor(resultChan chan map[string]string) {
	log.Println("StartFileMonitor")
	var pathList []string
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()
	for _, path := range common.Config.MonitorPath {
		if strings.HasPrefix(path, "/") {
			continue
		}
		// web目录 循环添加子目录
		if path == "%web%" {
			iterationWatcher(common.ServerInfo.Path, watcher, pathList)
			continue
		}
		if strings.Contains(path, "%windows%") {
			path = strings.Replace(path, "%windows%", os.Getenv("SystemDrive")+`\windows`, 1)
		} else if strings.Contains(path, "%system32%") {
			if runtime.GOARCH == "386" {
				path = strings.Replace(path, "%system32%", os.Getenv("SystemDrive")+`\windows\SysNative`, 1)
			} else {
				path = strings.Replace(path, "%system32%", os.Getenv("SystemDrive")+`\windows\System32`, 1)
			}
		}
		pathList = append(pathList, strings.ToLower(path))
		if strings.HasSuffix(path, "*") {
			iterationWatcher([]string{strings.Replace(path, "*", "", 1)}, watcher, pathList)
		}
		watcher.Add(strings.ToLower(path))
	}
	var resultdata map[string]string
	for {
		select {
		case event := <-watcher.Events:
			resultdata = make(map[string]string)
			if common.InArray(filter.File, strings.ToLower(event.Name), false) ||
				common.InArray(pathList, strings.ToLower(event.Name), false) ||
				common.InArray(common.Config.Filter.File, strings.ToLower(event.Name), true) {
				continue
			}
			if len(event.Name) == 0 {
				continue
			}
			resultdata["source"] = "file"
			resultdata["action"] = event.Op.String()
			resultdata["path"] = event.Name
			resultdata["hash"] = ""
			resultdata["user"] = ""
			f, err := os.Stat(event.Name)
			if err == nil && !f.IsDir() {
				if f.Size() <= fileSize {
					if hash, err := getFileMD5(event.Name); err == nil {
						if common.InArray(common.Config.Filter.File, strings.ToLower(hash), false) {
							continue
						}
						resultdata["hash"] = hash
					}
				}
				user := C.getprocessowner(C.CString(event.Name))
				resultdata["user"] = C.GoString(user)
			}
			if isFileWhite(resultdata) {
				continue
			}
			resultChan <- resultdata
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}
