Name "glbchain-dev ${MAJORVERSION}.${MINORVERSION}.${BUILDVERSION}" # VERSION variables set through command line arguments
InstallDir "$InstDir"
OutFile "${OUTPUTFILE}" # set through command line arguments

# Links for "Add/Remove Programs"
!define HELPURL "https://github.com/lbchain-devchain/go-lbchain-dev/issues"
!define UPDATEURL "https://github.com/lbchain-devchain/go-lbchain-dev/releases"
!define ABOUTURL "https://github.com/lbchain-devchain/go-lbchain-dev#lbchain-devereum-go"
!define /date NOW "%Y%m%d"

PageEx license
  LicenseData {{.License}}
PageExEnd

# Install glbchain-dev binary
Section "Glbchain-dev" Glbchain-dev_IDX
  SetOutPath $INSTDIR
  file {{.Glbchain-dev}}

  # Create start menu launcher
  createDirectory "$SMPROGRAMS\${APPNAME}"
  createShortCut "$SMPROGRAMS\${APPNAME}\${APPNAME}.lnk" "$INSTDIR\glbchain-dev.exe" "--fast" "--cache=512"
  createShortCut "$SMPROGRAMS\${APPNAME}\Attach.lnk" "$INSTDIR\glbchain-dev.exe" "attach" "" ""
  createShortCut "$SMPROGRAMS\${APPNAME}\Uninstall.lnk" "$INSTDIR\uninstall.exe" "" "" ""

  # Firewall - remove rules (if exists)
  SimpleFC::AdvRemoveRule "Glbchain-dev incoming peers (TCP:30303)"
  SimpleFC::AdvRemoveRule "Glbchain-dev outgoing peers (TCP:30303)"
  SimpleFC::AdvRemoveRule "Glbchain-dev UDP discovery (UDP:30303)"

  # Firewall - add rules
  SimpleFC::AdvAddRule "Glbchain-dev incoming peers (TCP:30303)" ""  6 1 1 2147483647 1 "$INSTDIR\glbchain-dev.exe" "" "" "lbchain-devchain" 30303 "" "" ""
  SimpleFC::AdvAddRule "Glbchain-dev outgoing peers (TCP:30303)" ""  6 2 1 2147483647 1 "$INSTDIR\glbchain-dev.exe" "" "" "lbchain-devchain" "" 30303 "" ""
  SimpleFC::AdvAddRule "Glbchain-dev UDP discovery (UDP:30303)" "" 17 2 1 2147483647 1 "$INSTDIR\glbchain-dev.exe" "" "" "lbchain-devchain" "" 30303 "" ""

  # Set default IPC endpoint (https://github.com/lbchain-devereum/EIPs/issues/147)
  ${EnvVarUpdate} $0 "lbchain-devEREUM_SOCKET" "R" "HKLM" "\\.\pipe\glbchain-dev.ipc"
  ${EnvVarUpdate} $0 "lbchain-devEREUM_SOCKET" "A" "HKLM" "\\.\pipe\glbchain-dev.ipc"

  # Add instdir to PATH
  Push "$INSTDIR"
  Call AddToPath
SectionEnd

# Install optional develop tools.
Section /o "Development tools" DEV_TOOLS_IDX
  SetOutPath $INSTDIR
  {{range .DevTools}}file {{.}}
  {{end}}
SectionEnd

# Return on top of stack the total size (as DWORD) of the selected/installed sections.
Var GetInstalledSize.total
Function GetInstalledSize
  StrCpy $GetInstalledSize.total 0

  ${if} ${SectionIsSelected} ${Glbchain-dev_IDX}
    SectionGetSize ${Glbchain-dev_IDX} $0
    IntOp $GetInstalledSize.total $GetInstalledSize.total + $0
  ${endif}

  ${if} ${SectionIsSelected} ${DEV_TOOLS_IDX}
    SectionGetSize ${DEV_TOOLS_IDX} $0
    IntOp $GetInstalledSize.total $GetInstalledSize.total + $0
  ${endif}

  IntFmt $GetInstalledSize.total "0x%08X" $GetInstalledSize.total
  Push $GetInstalledSize.total
FunctionEnd

# Write registry, Windows uses these values in various tools such as add/remove program.
# PowerShell: Get-ItemProperty HKLM:\Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall\* | Select-Object DisplayName, InstallLocation, InstallDate | Format-Table –AutoSize
function .onInstSuccess
  # Save information in registry in HKEY_LOCAL_MACHINE branch, Windows add/remove functionality depends on this
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${GROUPNAME} ${APPNAME}" "DisplayName" "${GROUPNAME} - ${APPNAME} - ${DESCRIPTION}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${GROUPNAME} ${APPNAME}" "UninstallString" "$\"$INSTDIR\uninstall.exe$\""
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${GROUPNAME} ${APPNAME}" "QuietUninstallString" "$\"$INSTDIR\uninstall.exe$\" /S"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${GROUPNAME} ${APPNAME}" "InstallLocation" "$INSTDIR"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${GROUPNAME} ${APPNAME}" "InstallDate" "${NOW}"
  # Wait for Alex
  #WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${GROUPNAME} ${APPNAME}" "DisplayIcon" "$\"$INSTDIR\logo.ico$\""
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${GROUPNAME} ${APPNAME}" "Publisher" "${GROUPNAME}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${GROUPNAME} ${APPNAME}" "HelpLink" "${HELPURL}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${GROUPNAME} ${APPNAME}" "URLUpdateInfo" "${UPDATEURL}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${GROUPNAME} ${APPNAME}" "URLInfoAbout" "${ABOUTURL}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${GROUPNAME} ${APPNAME}" "DisplayVersion" "${MAJORVERSION}.${MINORVERSION}.${BUILDVERSION}"
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${GROUPNAME} ${APPNAME}" "VersionMajor" ${MAJORVERSION}
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${GROUPNAME} ${APPNAME}" "VersionMinor" ${MINORVERSION}
  # There is no option for modifying or repairing the install
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${GROUPNAME} ${APPNAME}" "NoModify" 1
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${GROUPNAME} ${APPNAME}" "NoRepair" 1

  Call GetInstalledSize
  Pop $0
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${GROUPNAME} ${APPNAME}" "EstimatedSize" "$0"

  # Create uninstaller
  writeUninstaller "$INSTDIR\uninstall.exe"
functionEnd

Page components
Page directory
Page instfiles
