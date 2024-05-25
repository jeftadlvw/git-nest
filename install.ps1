# PowerShell script to install git-nest on Windows
# Download and run locally with
# powershell -c "iex (Get-Content -Path "git-nest_install.ps1 -Raw)

$ErrorActionPreference = "Stop"

# variables
$repositoryName = "jeftadlvw/git-nest"
$repository = "https://github.com/$repositoryName"
$installDir = "$Env:UserProfile\AppData\Local\Programs\git-nest"
$binaryName = "git-nest.exe"

$arch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture
switch ($arch) {
    'X64'  { $arch = 'amd64' }
    'Arm64' { $arch = 'arm64' }
    default { Write-Error "No official binary for target '$arch'" }
}

Write-Host "Installing git-nest to $installDir\$binaryName."
Write-Host ""

# check if git is installed
if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
    Write-Host "Warning: git is not installed."
}

# get the latest release information from GitHub API
$latestTag = (Invoke-RestMethod -Uri "https://api.github.com/repos/$repositoryName/releases/latest").tag_name
if (-not $latestTag) {
    Write-Host "Error: unable to retrieve latest release tag"
    exit 1
}

# create installation directory if it doesn't exist
if (-not (Test-Path -Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir
}

# download binary to the installation directory
$downloadUrl = "$repository/releases/download/$latestTag/git-nest_windows-$arch"
$downloadPath = "$installDir\$binaryName"
Write-Host "Downloading from $downloadUrl"

Invoke-WebRequest -Uri $downloadUrl -OutFile $downloadPath
if (-not (Test-Path -Path $downloadPath)) {
    Write-Host "Error: unable to download binary from $downloadUrl and install it to $downloadPath"
    exit 1
}

# make binary executable (not necessary on Windows, but we can set it to run as executable)
icacls $downloadPath /grant Everyone:RX

Write-Host ""
Write-Host "git-nest has been successfully installed to $installDir"
Write-Host ""

# check if installation directory is in PATH
if (-not ($env:PATH -contains $installDir)) {
    Write-Host "Notice: $installDir is not in your PATH."
    Write-Host "Add it to your PATH and restart your shell in order to use git-nest. E.g:"
    Write-Host ""
    Write-Host "    [System.Environment]::SetEnvironmentVariable('Path', `$env:Path + ';$installDir', [System.EnvironmentVariableTarget]::User)"
    Write-Host ""
}
