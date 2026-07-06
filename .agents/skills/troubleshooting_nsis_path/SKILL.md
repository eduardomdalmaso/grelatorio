---
name: Troubleshooting NSIS Path Build Error
description: Standard resolutions for 'makensis' not recognized errors when building Wails installers on Windows.
---

# Troubleshooting NSIS Path Build Error

This skill explains how to resolve the build error when executing `wails build -nsis` on Windows:
`makensis : The term 'makensis' is not recognized as the name of a cmdlet, function, script file, or operable program.`

## 🔴 The Problem
When running the Wails builder to package a Windows application into an NSIS installer, the build fails because the system cannot find the `makensis` executable, even if the Nullsoft Scriptable Install System (NSIS) has been successfully installed.

This occurs because:
1. NSIS is installed (typically at `C:\Program Files (x86)\NSIS\makensis.exe`).
2. The installation folder is not automatically added to the system's `PATH` environment variable, or the current terminal session has not yet loaded the updated `PATH`.

## 🟢 The Solution
To fix this, add the NSIS folder to the path in the current shell session before running the build command.

### In PowerShell
Add the directory to the `$env:Path` environment variable:
```powershell
$env:Path += ";C:\Program Files (x86)\NSIS"
wails build -nsis
```

### In Windows Command Prompt (cmd)
Add the directory to the `%PATH%` environment variable:
```cmd
set PATH=%PATH%;C:\Program Files (x86)\NSIS
wails build -nsis
```

To prevent this error in future sessions, you can permanently add `C:\Program Files (x86)\NSIS` to your system or user Environment Variables.
