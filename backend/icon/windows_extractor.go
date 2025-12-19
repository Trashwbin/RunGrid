//go:build windows

package icon

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

type WindowsExtractor struct{}

func NewDefaultExtractor() Extractor {
	return WindowsExtractor{}
}

func (WindowsExtractor) Extract(ctx context.Context, source string, dest string) error {
	if err := ValidateSource(source); err != nil {
		return err
	}

	ext := strings.ToLower(filepath.Ext(source))
	if ext == ".png" {
		return CopyFile(source, dest)
	}

	script := buildPowerShellScript(source, dest)
	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-Sta", "-Command", script)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

func buildPowerShellScript(source string, dest string) string {
	sourceLiteral := escapePowerShellLiteral(source)
	destLiteral := escapePowerShellLiteral(dest)
	return fmt.Sprintf(
		"$ErrorActionPreference='Stop';"+
			"Add-Type -AssemblyName System.Drawing;"+
			"Add-Type -Namespace RunGrid -Name IconUtil -MemberDefinition '[DllImport(\"shell32.dll\", CharSet=CharSet.Unicode)] public static extern int SHDefExtractIcon(string pszIconFile, int iIndex, uint uFlags, out IntPtr phiconLarge, out IntPtr phiconSmall, uint nIconSize); [DllImport(\"user32.dll\", SetLastError=true, CharSet=CharSet.Auto)] public static extern int PrivateExtractIcons(string lpszFile, int nIconIndex, int cxIcon, int cyIcon, IntPtr[] phicon, uint[] piconid, uint nIcons, int flags); [DllImport(\"user32.dll\", SetLastError=true)] public static extern bool DestroyIcon(IntPtr hIcon);';"+
			"$src='%s';"+
			"$dst='%s';"+
			"$iconSource=$src;"+
			"if([System.IO.Path]::GetExtension($src).ToLower() -eq '.lnk'){"+
			"try{$wsh=New-Object -ComObject WScript.Shell;"+
			"$shortcut=$wsh.CreateShortcut($src);"+
			"$iconLocation=$shortcut.IconLocation;"+
			"$targetPath=$shortcut.TargetPath;"+
			"$iconCandidate='';"+
			"if($iconLocation){$iconCandidate=($iconLocation -split ',')[0].Trim()}"+
			"if($iconCandidate -and [System.IO.Path]::GetExtension($iconCandidate).ToLower() -eq '.lnk'){$iconCandidate=''}"+
			"if(-not $iconCandidate -and $targetPath){$iconCandidate=$targetPath}"+
			"if($iconCandidate){$iconCandidate=$iconCandidate.Trim('\"').Trim(\"'\");"+
			"$iconCandidate=[Environment]::ExpandEnvironmentVariables($iconCandidate);"+
			"if(-not [System.IO.Path]::IsPathRooted($iconCandidate)){$iconCandidate=[System.IO.Path]::Combine([System.IO.Path]::GetDirectoryName($src),$iconCandidate)}"+
			"if(Test-Path $iconCandidate){$iconSource=$iconCandidate}}}catch{}}"+
			"if([string]::IsNullOrWhiteSpace($iconSource)){$iconSource=$src};"+
			"$sizes=@(256,128,96,64,48,32);"+
			"foreach($sz in $sizes){"+
			"$hLarge=[IntPtr]::Zero; $hSmall=[IntPtr]::Zero;"+
			"$res=[RunGrid.IconUtil]::SHDefExtractIcon($iconSource,0,0,[ref]$hLarge,[ref]$hSmall,$sz);"+
			"if($res -ge 0 -and $hLarge -ne [IntPtr]::Zero){$icon=[System.Drawing.Icon]::FromHandle($hLarge);$bmp=$icon.ToBitmap();$bmp.Save($dst,[System.Drawing.Imaging.ImageFormat]::Png);$bmp.Dispose();$icon.Dispose();[RunGrid.IconUtil]::DestroyIcon($hLarge);if($hSmall -ne [IntPtr]::Zero){[RunGrid.IconUtil]::DestroyIcon($hSmall)};return}}"+
			"$hicons=@([IntPtr]::Zero);$ids=@(0);"+
			"foreach($sz in $sizes){$count=[RunGrid.IconUtil]::PrivateExtractIcons($iconSource,0,$sz,$sz,$hicons,$ids,1,0);if($count -gt 0 -and $hicons[0] -ne [IntPtr]::Zero){$icon=[System.Drawing.Icon]::FromHandle($hicons[0]);$bmp=$icon.ToBitmap();$bmp.Save($dst,[System.Drawing.Imaging.ImageFormat]::Png);$bmp.Dispose();$icon.Dispose();[RunGrid.IconUtil]::DestroyIcon($hicons[0]);return}}"+
			"$icon=$null;"+
			"try{$icon=New-Object System.Drawing.Icon($iconSource)}catch{};"+
			"if($null -eq $icon){$icon=[System.Drawing.Icon]::ExtractAssociatedIcon($iconSource)};"+
			"if($null -eq $icon){throw 'icon not found'};"+
			"$bmp=$icon.ToBitmap();"+
			"$bmp.Save($dst,[System.Drawing.Imaging.ImageFormat]::Png);"+
			"$bmp.Dispose();$icon.Dispose();",
		sourceLiteral,
		destLiteral,
	)
}

func escapePowerShellLiteral(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}
