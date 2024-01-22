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
    A[User] -- 1. Provides \n warp-diag.zip --> B[warp-diag-checker]
    B -- 2. Unzips and reads --> C[Ingested Log Data<br/>- daemon.log<br/>- netstat.txt<br/>- etc.] --> E
	C[Ingested Log Data] --> F
    B -- 3. Fetches config --> D[Remote Config] 
    D -- 4a. Defines<br/>built-in tests --> E[Built-in Tests]
    D -- 4b. Defines<br/>log search terms --> F[LogSearch]
    E --> G[5. Collect Results]
    F --> G
    G -- 6. Display results --> A[User]

style A fill:#4F772D,color:#fff,stroke:#333,stroke-width:2px
    style B fill:#1F618D,color:#fff,stroke:#333,stroke-width:2px  
    style C fill:#DDA15E,color:#fff,stroke:#333,stroke-width:2px
    style D fill:#BC4749,color:#fff,stroke:#333,stroke-width:2px
    style E fill:#80FF72,color:#000,stroke:#333,stroke-width:2px
    style F fill:#F0A500,color:#000,stroke:#333,stroke-width:2px
    style G fill:#6495ED,color:#fff,stroke:#333,stroke-width:2px
```

## Installation

### Mac and Linux

To install on Mac or Linux via Homebrew:

```bash {#install-brew}
brew tap peakefficiency/releases
brew install warp-diag-checker
```

To update, run:

```
brew upgrade warp-diag-checker
```

### Windows  

To install on Windows via Chocolatey:

```powershell {#install-choco}
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
