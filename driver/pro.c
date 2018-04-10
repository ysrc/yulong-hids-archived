#include "ntddk.h"
#include "windef.h"
#include "define.h"

#define SYSNAME "System"
#define VERSIONLEN 100

const WCHAR devLink[] = L"\\DosDevices\\MonitorProcess";
const WCHAR devName[] = L"\\Device\\MonitorProcess";
UNICODE_STRING          devNameUnicd;
UNICODE_STRING          devLinkUnicd;
PVOID                    gpEventObject = NULL;           
ULONG                    ProcessNameOffset = 0;
CHAR                    outBuf[255];
BOOL                    g_bMainThread;
ULONG                    g_dwParentId;
CHECKLIST                CheckList;
ULONG                    BuildNumber;                                        
ULONG                    SYSTEMID;                   
PWCHAR                    Version[VERSIONLEN];

NTSTATUS PsLookupProcessByProcessId(IN ULONG ulProcId, OUT PEPROCESS * pEProcess);

ULONG GetProcessNameOffset()
{
    PEPROCESS curproc;
    int i;

    curproc = PsGetCurrentProcess();

    for (i = 0; i < 3 * PAGE_SIZE; i++)
    {
        if (!strncmp(SYSNAME, (PCHAR)curproc + i, strlen(SYSNAME)))
        {
            return i;
        }
    }

    return 0;
}

NTSTATUS GetRegValue(PCWSTR RegPath, PCWSTR ValueName, PWCHAR Value)
{
    int ReturnValue = 0;
    NTSTATUS Status;
    OBJECT_ATTRIBUTES ObjectAttributes;
    HANDLE KeyHandle;
    PKEY_VALUE_PARTIAL_INFORMATION valueInfoP;
    ULONG valueInfoLength, returnLength;
    UNICODE_STRING UnicodeRegPath;
    UNICODE_STRING UnicodeValueName;

    RtlInitUnicodeString(&UnicodeRegPath, RegPath);
    RtlInitUnicodeString(&UnicodeValueName, ValueName);

    InitializeObjectAttributes(&ObjectAttributes,
        &UnicodeRegPath,
        OBJ_CASE_INSENSITIVE, 
        NULL, 
        NULL); 

    Status = ZwOpenKey(&KeyHandle,
        KEY_ALL_ACCESS,
        &ObjectAttributes);
    if (Status != STATUS_SUCCESS)
    {
        DbgPrint("zwopenkey error\n");
        return 0;
    }

    valueInfoLength = sizeof(KEY_VALUE_PARTIAL_INFORMATION)+VERSIONLEN;
    valueInfoP = (PKEY_VALUE_PARTIAL_INFORMATION)ExAllocatePool
        (NonPagedPool, valueInfoLength);
    Status = ZwQueryValueKey(KeyHandle,
        &UnicodeValueName,
        KeyValuePartialInformation,
        valueInfoP,
        valueInfoLength,
        &returnLength);

    if (!NT_SUCCESS(Status))
    {
        DbgPrint("zwqueryvaluekey error:%08x\n", Status);
        ExFreePool(valueInfoP);
        ZwClose(KeyHandle);
        return Status;
    }
    else
    {
        RtlCopyMemory((PCHAR)Value, (PCHAR)valueInfoP->Data, valueInfoP->DataLength);
        ReturnValue = 1;
    }

    ExFreePool(valueInfoP);
    ZwClose(KeyHandle);
    return ReturnValue;
}



VOID ThreadCreateMon(IN HANDLE PId, IN HANDLE TId, IN BOOLEAN  bCreate)
{

    PEPROCESS   EProcess, PEProcess;
    NTSTATUS    status;
    HANDLE        dwParentPID;

    status = PsLookupProcessByProcessId((ULONG)PId, &EProcess);
    if (!NT_SUCCESS(status))
    {
        DbgPrint("error\n");
        return;
    }

    if (bCreate)
    {
        dwParentPID = PsGetCurrentProcessId();
        status = PsLookupProcessByProcessId(
            (ULONG)dwParentPID,
            &PEProcess);
        if (!NT_SUCCESS(status))
        {
            ObDereferenceObject(EProcess);
            DbgPrint("error\n");
            return;
        }
        if (PId == 4) 
        {
            ObDereferenceObject(PEProcess);
            ObDereferenceObject(EProcess);
            return;
        }
        if ((g_bMainThread == TRUE)
            && (g_dwParentId != dwParentPID)
            && (dwParentPID != PId)
            )
        {
            g_bMainThread = FALSE;
            sprintf(outBuf, "r_thread|%s|%d|%s|%d\n",
                (char *)((char *)EProcess + ProcessNameOffset),
                PId,
                (char *)((char *)PEProcess + ProcessNameOffset), dwParentPID);
            if (gpEventObject != NULL) 
                KeSetEvent((PRKEVENT)gpEventObject, 0, FALSE);
        }
        if (CheckList.ONLYSHOWREMOTETHREAD) 
        {
            ObDereferenceObject(PEProcess);
            ObDereferenceObject(EProcess);
            return;
        }
        DbgPrint("thread|%s|%d|%s|%d\n",
            (char *)((char *)EProcess + ProcessNameOffset),
            PId,
            (char *)((char *)PEProcess + ProcessNameOffset), dwParentPID);
        sprintf(outBuf, "thread|%s|%d|%s|%d\n",
            (char *)((char *)EProcess + ProcessNameOffset),
            PId,
            (char *)((char *)PEProcess + ProcessNameOffset), dwParentPID);
        if (gpEventObject != NULL) 
            KeSetEvent((PRKEVENT)gpEventObject, 0, FALSE);
        
        ObDereferenceObject(PEProcess);
    }
    else if (CheckList.SHOWTERMINATETHREAD)
    {
        DbgPrint("thread_over|%d\n", TId);
        sprintf(outBuf, "thread_over|%d\n", TId);
        if (gpEventObject != NULL)
            KeSetEvent((PRKEVENT)gpEventObject, 0, FALSE);
    }
    ObDereferenceObject(EProcess);
}


VOID ProcessCreateMon(HANDLE hParentId, HANDLE PId, BOOLEAN bCreate)
{

    PEPROCESS        EProcess, PProcess;
    NTSTATUS        status;
    HANDLE            TId;

    g_dwParentId = hParentId;
    status = PsLookupProcessByProcessId((ULONG)PId, &EProcess);
    if (!NT_SUCCESS(status))
    {
        DbgPrint("error\n");
        return;
    }
    status = PsLookupProcessByProcessId((ULONG)hParentId, &PProcess);
    if (!NT_SUCCESS(status))
    {
        DbgPrint("error\n");
        ObDereferenceObject(EProcess);
        return;
    }

    if (bCreate)
    {
        g_bMainThread = TRUE;
        DbgPrint("process|%s|%d|%s|%d\n",
            (char *)((char *)EProcess + ProcessNameOffset),
            PId,
            (char *)((char *)PProcess + ProcessNameOffset),
            hParentId
            );
        sprintf(outBuf, "process|%s|%d|%s|%d\n",
            (char *)((char *)EProcess + ProcessNameOffset),
            PId,
            (char *)((char *)PProcess + ProcessNameOffset),
            hParentId
            );
        if (gpEventObject != NULL) 
            KeSetEvent((PRKEVENT)gpEventObject, 0, FALSE);
    }
    else if (CheckList.SHOWTERMINATEPROCESS)
    {
        DbgPrint("process_over|%d\n", PId);
        sprintf(outBuf, "process_over|%d\n", PId);
        if (gpEventObject != NULL)
            KeSetEvent((PRKEVENT)gpEventObject, 0, FALSE);
    }

    ObDereferenceObject(PProcess);
    ObDereferenceObject(EProcess);
}

NTSTATUS OnUnload(IN PDRIVER_OBJECT pDriverObject)
{
    NTSTATUS            status;
    
    DbgPrint("OnUnload called\n");
    
    if (gpEventObject)
    {
        ObDereferenceObject(gpEventObject);
    }

    PsSetCreateProcessNotifyRoutine(ProcessCreateMon, TRUE);
    PsRemoveCreateThreadNotifyRoutine(ThreadCreateMon);
    IoDeleteSymbolicLink(&devLinkUnicd);
    
    if (pDriverObject->DeviceObject != NULL)
    {
        IoDeleteDevice(pDriverObject->DeviceObject);
    }
    
    return STATUS_SUCCESS;
}

NTSTATUS DeviceIoControlDispatch(
    IN  PDEVICE_OBJECT  DeviceObject,
    IN  PIRP            pIrp
    )
{
    PIO_STACK_LOCATION              irpStack;
    NTSTATUS                        status;
    PVOID                           inputBuffer;
    ULONG                           inputLength;
    PVOID                           outputBuffer;
    ULONG                           outputLength;
    OBJECT_HANDLE_INFORMATION        objHandleInfo;

    status = STATUS_SUCCESS;
    irpStack = IoGetCurrentIrpStackLocation(pIrp);

    switch (irpStack->MajorFunction)
    {
    case IRP_MJ_CREATE:
        DbgPrint("Call IRP_MJ_CREATE\n");
        break;
    case IRP_MJ_CLOSE:
        DbgPrint("Call IRP_MJ_CLOSE\n");
        break;
    case IRP_MJ_DEVICE_CONTROL:
        DbgPrint("IRP_MJ_DEVICE_CONTROL\n");
        inputLength = irpStack->Parameters.DeviceIoControl.InputBufferLength;
        outputLength = irpStack->Parameters.DeviceIoControl.OutputBufferLength;
        switch (irpStack->Parameters.DeviceIoControl.IoControlCode)
        {
        case IOCTL_PASSEVENT: 
            inputBuffer = pIrp->AssociatedIrp.SystemBuffer;

            DbgPrint("inputBuffer:%08x\n", (HANDLE)inputBuffer);
            status = ObReferenceObjectByHandle(*(HANDLE *)inputBuffer,
                GENERIC_ALL,
                NULL,
                KernelMode,
                &gpEventObject,
                &objHandleInfo);
            
            if (status != STATUS_SUCCESS)
            {
                DbgPrint("wrong\n");
                break;
            }
            break;
        case IOCTL_UNPASSEVENT:
            if (gpEventObject)
                ObDereferenceObject(gpEventObject);
            DbgPrint("UNPASSEVENT called\n");
            break;
        case IOCTL_PASSBUF:
            RtlCopyMemory(pIrp->UserBuffer, outBuf, outputLength);
            break;
        case IOCTL_PASSEVSTRUCT:
            inputBuffer = pIrp->AssociatedIrp.SystemBuffer;
            memset(&CheckList, 0, sizeof(CheckList));
            RtlCopyMemory(&CheckList, inputBuffer, sizeof(CheckList));
            DbgPrint("%d:%d\n", CheckList.ONLYSHOWREMOTETHREAD, CheckList.SHOWTHREAD);
            break;
        default:
            break;
        }
        break;
    default:
        DbgPrint("Call IRP_MJ_UNKNOWN\n");
        break;
    }

    pIrp->IoStatus.Status = status;
    pIrp->IoStatus.Information = 0;
    IoCompleteRequest(pIrp, IO_NO_INCREMENT);
    return status;
}

NTSTATUS DriverEntry(IN PDRIVER_OBJECT pDriverObject, IN PUNICODE_STRING theRegistryPath)
{
    NTSTATUS                Status;
    PDEVICE_OBJECT            pDevice;

    DbgPrint("DriverEntry called!\n");
    g_bMainThread = FALSE;
	memset(outBuf, 0, 255);

    if (1 != GetRegValue(L"\\Registry\\Machine\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion", L"CSDVersion", Version))
    {
        DbgPrint("GetRegValueDword Wrong\n");
    }
    PsGetVersion(NULL, NULL, &BuildNumber, NULL);
    DbgPrint("[[[%d]]]:[[[%ws]]]", BuildNumber, Version);

    RtlInitUnicodeString(&devNameUnicd, devName);
    RtlInitUnicodeString(&devLinkUnicd, devLink);

    Status = IoCreateDevice(pDriverObject,
        0,
        &devNameUnicd,
        FILE_DEVICE_UNKNOWN,
        0,
        TRUE,
        &pDevice);
    if (!NT_SUCCESS(Status))
    {
        DbgPrint(("Can not create device.\n"));
        goto out;
    }

    Status = IoCreateSymbolicLink(&devLinkUnicd, &devNameUnicd);
    if (!NT_SUCCESS(Status))
    {
        DbgPrint(("Cannot create link.\n"));
        goto CleanDevice;
    }

    ProcessNameOffset = GetProcessNameOffset();
    pDriverObject->DriverUnload = OnUnload;
    pDriverObject->MajorFunction[IRP_MJ_CREATE] =
        pDriverObject->MajorFunction[IRP_MJ_CLOSE] =
        pDriverObject->MajorFunction[IRP_MJ_DEVICE_CONTROL] = DeviceIoControlDispatch;

    Status = PsSetCreateProcessNotifyRoutine(ProcessCreateMon, FALSE);
    if (!NT_SUCCESS(Status))
    {
        DbgPrint("PsSetCreateProcessNotifyRoutine error\n");
        goto CleanSymbolLink;
    }
    Status = PsSetCreateThreadNotifyRoutine(ThreadCreateMon);
    if (!NT_SUCCESS(Status))
    {
        DbgPrint("PsSetCreateThreadNotifyRoutine error\n");
        goto CleanProcessNotify;
    }

    return STATUS_SUCCESS;

CleanProcessNotify:
    PsSetCreateProcessNotifyRoutine(ProcessCreateMon, TRUE);
CleanSymbolLink:
    IoDeleteSymbolicLink(&devLinkUnicd);
CleanDevice:
    if (pDriverObject->DeviceObject != NULL) 
    {
        IoDeleteDevice(pDriverObject->DeviceObject);
    }
out:
    return Status;
}
