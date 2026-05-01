param(
    [switch]$Release
)

$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Resolve-Path (Join-Path $scriptDir "..\..")
$installerScript = Join-Path $projectRoot "build\installer\musicalc.iss"

Push-Location $projectRoot
try {
    $version = git describe --tags --exact-match 2>$null
    if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrWhiteSpace($version)) {
        if ($Release) {
            throw "Release installer builds require HEAD to be exactly on a version tag."
        }

        $latestTag = git describe --tags --abbrev=0 2>$null
        $shortCommit = git rev-parse --short HEAD
        if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrWhiteSpace($shortCommit)) {
            $version = "dev"
        } elseif ([string]::IsNullOrWhiteSpace($latestTag)) {
            $version = "dev-$shortCommit"
        } else {
            $version = "$latestTag-dev-$shortCommit"
        }
    }

    if ($Release -and ($version -notmatch '^\d+\.\d+\.\d+$')) {
        throw "Release tag '$version' is invalid. Expected MAJOR.MINOR.PATCH, for example 0.8.7."
    }

    iscc "/DMyAppVersion=$version" $installerScript
    $exitCode = $LASTEXITCODE
} finally {
    Pop-Location
}

if ($exitCode -ne 0) {
    exit $exitCode
}
