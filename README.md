# warp-diag-checker

<p>
  <img title="warp-diag-checker" src="https://raw.githubusercontent.com/peakefficiency/warp-diag-checker/main/assets/icon.png"" width="400" />
  <br><br>
  <a href="https://github.com/github/peakefficiency/warp-diag-checker/releases"><img src="https://img.shields.io/github/release/peakefficiency/warp-diag-checker.svg" alt="Latest Release"></a>

  <a href="https://github.com/peakefficiency/warp-diag-checker/actions"><img src="https://github.com/peakefficiency/warp-diag-checker/workflows/Build/badge.svg" alt="Build Status"></a>
</p>

## General Overview
```mermaid
flowchart LR
    A[User] -- Provides warp-diag file --> B[CLI Interface]
    B -- Ingests file --> C[Test Runner]
    C -- Fetches --> D[Remote Config File]
    D -- Defines tests --> C
    C -- Runs tests on --> E[Ingested Data]
    E -- Outputs --> F[Results]
    F --> A
```

## Installation

### Mac and Linux

To install on Mac or Linux via Homebrew:

```
brew install warp-diag-checker
```

To update, run:

```
brew upgrade warp-diag-checker
```

### Windows  

To install on Windows via Chocolatey:

```
choco install warp-diag-checker
```

If the latest version is still under review, you may need to specify the intended version found at https://community.chocolatey.org/packages/warp-diag-checker

To update, run:

```
choco upgrade warp-diag-checker
```

### Install issues

As a fallback, you can install from source using:

```
go install github.com/peakefficiency/warp-diag-checker@latest
```
You would need to have Go installed for this

## Usage


Basic usage:

```
warp-diag-checker /path/to/diag.zio
```
