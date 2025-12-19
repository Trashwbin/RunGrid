//go:build windows

package service

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

func resolveShortcutTarget(ctx context.Context, source string) (string, error) {
	script := buildShortcutResolveScript(source)
	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-Sta", "-Command", script)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func buildShortcutResolveScript(source string) string {
	sourceLiteral := escapePowerShellLiteral(source)
	return fmt.Sprintf(
		"$ErrorActionPreference='Stop';"+
			"$src='%s';"+
			"$wsh=New-Object -ComObject WScript.Shell;"+
			"$shortcut=$wsh.CreateShortcut($src);"+
			"$target=$shortcut.TargetPath;"+
			"if($target){$target=$target.Trim('\"');"+
			"$target=$target.Trim(\"'\");"+
			"$target=[Environment]::ExpandEnvironmentVariables($target);"+
			"if(-not [System.IO.Path]::IsPathRooted($target)){$target=[System.IO.Path]::Combine([System.IO.Path]::GetDirectoryName($src),$target)}}"+
			"Write-Output $target;",
		sourceLiteral,
	)
}

func escapePowerShellLiteral(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}
