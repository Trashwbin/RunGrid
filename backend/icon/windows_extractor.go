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
			"$icon=[System.Drawing.Icon]::ExtractAssociatedIcon($iconSource);"+
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
