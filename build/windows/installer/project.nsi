Unicode true

####
## Please note: Template replacements don't work in this file. They are provided with default defines like
## mentioned underneath.
## If the keyword is not defined, "wails_tools.nsh" will populate them with the values from ProjectInfo.
## If they are defined here, "wails_tools.nsh" will not touch them. This allows to use this project.nsi manually
## from outside of Wails for debugging and development of the installer.
##
## For development first make a wails nsis build to populate the "wails_tools.nsh":
## > wails build --target windows/amd64 --nsis
## Then you can call makensis on this file with specifying the path to your binary:
## For a AMD64 only installer:
## > makensis -DARG_WAILS_AMD64_BINARY=..\..\bin\app.exe
## For a ARM64 only installer:
## > makensis -DARG_WAILS_ARM64_BINARY=..\..\bin\app.exe
## For a installer with both architectures:
## > makensis -DARG_WAILS_AMD64_BINARY=..\..\bin\app-amd64.exe -DARG_WAILS_ARM64_BINARY=..\..\bin\app-arm64.exe
####
## The following information is taken from the ProjectInfo file, but they can be overwritten here.
####
## !define INFO_PROJECTNAME    "MyProject" # Default "{{.Name}}"
## !define INFO_COMPANYNAME    "MyCompany" # Default "{{.Info.CompanyName}}"
## !define INFO_PRODUCTNAME    "MyProduct" # Default "{{.Info.ProductName}}"
## !define INFO_PRODUCTVERSION "1.0.0"     # Default "{{.Info.ProductVersion}}"
## !define INFO_COPYRIGHT      "Copyright" # Default "{{.Info.Copyright}}"
###
## !define PRODUCT_EXECUTABLE  "Application.exe"      # Default "${INFO_PROJECTNAME}.exe"
## !define UNINST_KEY_NAME     "UninstKeyInRegistry"  # Default "${INFO_COMPANYNAME}${INFO_PRODUCTNAME}"
####
## !define REQUEST_EXECUTION_LEVEL "admin"            # Default "admin"  see also https://nsis.sourceforge.io/Docs/Chapter4.html
####
## Include the wails tools
####
!include "wails_tools.nsh"

# The version information for this two must consist of 4 parts
VIProductVersion "${INFO_PRODUCTVERSION}.0"
VIFileVersion    "${INFO_PRODUCTVERSION}.0"

VIAddVersionKey "CompanyName"     "${INFO_COMPANYNAME}"
VIAddVersionKey "FileDescription" "${INFO_PRODUCTNAME} Installer"
VIAddVersionKey "ProductVersion"  "${INFO_PRODUCTVERSION}"
VIAddVersionKey "FileVersion"     "${INFO_PRODUCTVERSION}"
VIAddVersionKey "LegalCopyright"  "${INFO_COPYRIGHT}"
VIAddVersionKey "ProductName"     "${INFO_PRODUCTNAME}"

# Enable HiDPI support. https://nsis.sourceforge.io/Reference/ManifestDPIAware
ManifestDPIAware true

!include "MUI.nsh"
!include "Sections.nsh"
!include "StrFunc.nsh"

${StrStr}
${UnStrStr}

!define MUI_ICON "..\icon.ico"
!define MUI_UNICON "..\icon.ico"
# !define MUI_WELCOMEFINISHPAGE_BITMAP "resources\leftimage.bmp" #Include this to add a bitmap on the left side of the Welcome Page. Must be a size of 164x314
!define MUI_FINISHPAGE_NOAUTOCLOSE # Wait on the INSTFILES page so the user can take a look into the details of the installation steps
!define MUI_FINISHPAGE_RUN "$INSTDIR\\${PRODUCT_EXECUTABLE}"
!define MUI_FINISHPAGE_RUN_TEXT "$(RUN_TEXT)"
!define MUI_ABORTWARNING # This will warn the user if they exit from the installer.
!define MUI_LANGDLL_REGISTRY_ROOT "HKCU"
!define MUI_LANGDLL_REGISTRY_KEY "Software\${INFO_COMPANYNAME}\${INFO_PRODUCTNAME}"
!define MUI_LANGDLL_REGISTRY_VALUENAME "InstallerLanguage"

!insertmacro MUI_PAGE_WELCOME # Welcome to the installer page.
!insertmacro MUI_PAGE_COMPONENTS # Optional components (e.g. auto start).
# !insertmacro MUI_PAGE_LICENSE "resources\eula.txt" # Adds a EULA page to the installer
!insertmacro MUI_PAGE_DIRECTORY # In which folder install page.
!insertmacro MUI_PAGE_INSTFILES # Installing page.
!insertmacro MUI_PAGE_FINISH # Finished installation page.

!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES # Uinstalling page

!insertmacro MUI_LANGUAGE "English"
!insertmacro MUI_LANGUAGE "SimpChinese"

!define DATA_DIR "$APPDATA\${INFO_PROJECTNAME}"
!define DATA_RULES_DIR "${DATA_DIR}\rules"
!define RULES_SOURCE_DIR "..\..\..\rules"

LangString MSG_UNINSTALL_RUNNING ${LANG_ENGLISH} "${INFO_PRODUCTNAME} is running and will be closed to continue uninstalling. Continue?"
LangString MSG_UNINSTALL_RUNNING ${LANG_SIMPCHINESE} "${INFO_PRODUCTNAME} is running and will be closed to continue uninstalling. Continue?"
LangString MSG_INSTALL_RUNNING ${LANG_ENGLISH} "${INFO_PRODUCTNAME} is running and will be closed to continue installing. Continue?"
LangString MSG_INSTALL_RUNNING ${LANG_SIMPCHINESE} "${INFO_PRODUCTNAME} 正在运行，将被关闭以继续安装。是否继续？"
LangString SEC_SHORTCUTS ${LANG_ENGLISH} "Create shortcuts"
LangString SEC_SHORTCUTS ${LANG_SIMPCHINESE} "创建快捷方式"
LangString SEC_AUTOSTART ${LANG_ENGLISH} "Start ${INFO_PRODUCTNAME} with Windows"
LangString SEC_AUTOSTART ${LANG_SIMPCHINESE} "开机自启"
LangString RUN_TEXT ${LANG_ENGLISH} "Launch ${INFO_PRODUCTNAME}"
LangString RUN_TEXT ${LANG_SIMPCHINESE} "立即启动 ${INFO_PRODUCTNAME}"
LangString DESC_SHORTCUTS ${LANG_ENGLISH} "Create Start Menu and Desktop shortcuts."
LangString DESC_SHORTCUTS ${LANG_SIMPCHINESE} "创建开始菜单和桌面快捷方式。"
LangString DESC_AUTOSTART ${LANG_ENGLISH} "Launch ${INFO_PRODUCTNAME} automatically when you sign in."
LangString DESC_AUTOSTART ${LANG_SIMPCHINESE} "登录后自动启动 ${INFO_PRODUCTNAME}。"

## The following two statements can be used to sign the installer and the uninstaller. The path to the binaries are provided in %1
#!uninstfinalize 'signtool --file "%1"'
#!finalize 'signtool --file "%1"'

Name "${INFO_PRODUCTNAME}"
OutFile "..\..\bin\${INFO_PROJECTNAME}-${ARCH}-installer.exe" # Name of the installer's file.
InstallDir "$PROGRAMFILES64\${INFO_PRODUCTNAME}" # Default installing folder ($PROGRAMFILES is Program Files folder).
InstallDirRegKey HKLM "${UNINST_KEY}" "InstallLocation"
ShowInstDetails show # This will always show the installation details.

Section "-Main" SEC_MAIN
    !insertmacro wails.setShellContext

    !insertmacro wails.webview2runtime

    SetOutPath $INSTDIR

    !insertmacro wails.files

    CreateDirectory "${DATA_DIR}"
    CreateDirectory "${DATA_RULES_DIR}"
    SetOutPath "${DATA_RULES_DIR}"
    File /r "${RULES_SOURCE_DIR}\*.json"
    SetOutPath $INSTDIR

    !insertmacro wails.associateFiles
    !insertmacro wails.associateCustomProtocols

    !insertmacro wails.writeUninstaller
    WriteRegStr HKLM "${UNINST_KEY}" "InstallLocation" "$INSTDIR"
SectionEnd

Section /o "$(SEC_SHORTCUTS)" SEC_SHORTCUTS
    CreateShortcut "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"
    CreateShortCut "$DESKTOP\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"
SectionEnd

Section /o "$(SEC_AUTOSTART)" SEC_AUTOSTART
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Run" "${INFO_PRODUCTNAME}" "$\"$INSTDIR\${PRODUCT_EXECUTABLE}$\""
SectionEnd

Section "uninstall"
    !insertmacro wails.setShellContext

    DeleteRegValue HKCU "Software\Microsoft\Windows\CurrentVersion\Run" "${INFO_PRODUCTNAME}"

    Delete /REBOOTOK "$INSTDIR\${PRODUCT_EXECUTABLE}"
    Delete /REBOOTOK "$INSTDIR\uninstall.exe"

    Delete "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk"
    Delete "$DESKTOP\${INFO_PRODUCTNAME}.lnk"

    !insertmacro wails.unassociateFiles
    !insertmacro wails.unassociateCustomProtocols

    RMDir /r "$APPDATA\${PRODUCT_EXECUTABLE}" # Remove WebView2 data path if present.
    RMDir /r "${DATA_DIR}"
    RMDir /r "$LOCALAPPDATA\${INFO_PROJECTNAME}"
    RMDir /r $INSTDIR

    !insertmacro wails.deleteUninstaller
SectionEnd

!insertmacro MUI_FUNCTION_DESCRIPTION_BEGIN
    !insertmacro MUI_DESCRIPTION_TEXT ${SEC_SHORTCUTS} "$(DESC_SHORTCUTS)"
    !insertmacro MUI_DESCRIPTION_TEXT ${SEC_AUTOSTART} "$(DESC_AUTOSTART)"
!insertmacro MUI_FUNCTION_DESCRIPTION_END

Function .onInit
   !insertmacro MUI_LANGDLL_DISPLAY
   !insertmacro wails.checkArchitecture
   StrCpy $R0 ""
   SetRegView 64
   ReadRegStr $R0 HKLM "${UNINST_KEY}" "InstallLocation"
   StrCmp $R0 "" 0 setInstallDir
   ReadRegStr $R0 HKCU "${UNINST_KEY}" "InstallLocation"
   StrCmp $R0 "" 0 setInstallDir
   SetRegView 32
   ReadRegStr $R0 HKLM "${UNINST_KEY}" "InstallLocation"
   StrCmp $R0 "" 0 setInstallDir
   ReadRegStr $R0 HKCU "${UNINST_KEY}" "InstallLocation"
   StrCmp $R0 "" 0 setInstallDir
   Goto doneInstallDir
setInstallDir:
   StrCpy $INSTDIR "$R0"
doneInstallDir:
   nsExec::ExecToStack '"$SYSDIR\\tasklist.exe" /FI "IMAGENAME eq ${PRODUCT_EXECUTABLE}" /FO CSV /NH'
   Pop $0
   Pop $1
   ${StrStr} $2 $1 "${PRODUCT_EXECUTABLE}"
   StrCmp $2 "" doneInstall
   IfSilent silentInstallKill
   MessageBox MB_OKCANCEL|MB_ICONEXCLAMATION "$(MSG_INSTALL_RUNNING)" IDOK +2
   Abort
silentInstallKill:
   ExecWait '"$SYSDIR\\taskkill.exe" /F /IM ${PRODUCT_EXECUTABLE}'
   Sleep 500
doneInstall:
   SectionGetFlags ${SEC_SHORTCUTS} $0
   IntOp $0 $0 | ${SF_SELECTED}
   SectionSetFlags ${SEC_SHORTCUTS} $0
   SectionGetFlags ${SEC_AUTOSTART} $1
   IntOp $1 $1 | ${SF_SELECTED}
   SectionSetFlags ${SEC_AUTOSTART} $1
FunctionEnd

Function un.onInit
    !insertmacro MUI_LANGDLL_DISPLAY
    nsExec::ExecToStack '"$SYSDIR\\tasklist.exe" /FI "IMAGENAME eq ${PRODUCT_EXECUTABLE}" /FO CSV /NH'
    Pop $0
    Pop $1
    ${UnStrStr} $2 $1 "${PRODUCT_EXECUTABLE}"
    StrCmp $2 "" done
    IfSilent silentKill
    MessageBox MB_OKCANCEL|MB_ICONEXCLAMATION "$(MSG_UNINSTALL_RUNNING)" IDOK +2
    Abort
silentKill:
    ExecWait '"$SYSDIR\\taskkill.exe" /F /IM ${PRODUCT_EXECUTABLE}'
    Sleep 500
done:
FunctionEnd
