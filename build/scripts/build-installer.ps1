param(
    [string]$Version,
    [switch]$Release
)

$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Resolve-Path (Join-Path $scriptDir "..\..")
$installerScript = Join-Path $projectRoot "build\installer\musicalc.iss"

Push-Location $projectRoot
try {
    $versionValue = $Version

    if ([string]::IsNullOrWhiteSpace($versionValue)) {
        $versionValue = git describe --tags --exact-match 2>$null
        if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrWhiteSpace($versionValue)) {
            if ($Release) {
                throw "Release installer builds require HEAD to be exactly on a version tag."
            }

            $latestTag = git describe --tags --abbrev=0 2>$null
            $shortCommit = git rev-parse --short HEAD
            if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrWhiteSpace($shortCommit)) {
                $versionValue = "dev"
            } elseif ([string]::IsNullOrWhiteSpace($latestTag)) {
                $versionValue = "dev-$shortCommit"
            } else {
                $versionValue = "$latestTag-dev-$shortCommit"
            }
        }
    } elseif ($Release) {
        $tagVersion = git describe --tags --exact-match 2>$null
        if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrWhiteSpace($tagVersion)) {
            throw "Release installer builds require HEAD to be exactly on a version tag."
        } elseif ($versionValue -ne $tagVersion) {
            throw "Release version '$versionValue' does not match exact Git tag '$tagVersion'."
        }
    }

    if ($Release -and ($versionValue -notmatch '^\d+\.\d+\.\d+$')) {
        throw "Release tag '$versionValue' is invalid. Expected MAJOR.MINOR.PATCH, for example 0.8.7."
    }

    if ($versionValue -notmatch '^[0-9A-Za-z][0-9A-Za-z.-]*$') {
        throw "Installer version '$versionValue' is invalid. Use only letters, numbers, dots, and hyphens."
    }

    iscc "/DMyAppVersion=$versionValue" $installerScript
    $exitCode = $LASTEXITCODE
} finally {
    Pop-Location
}

if ($exitCode -ne 0) {
    exit $exitCode
}
