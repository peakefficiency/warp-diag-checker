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
    A[User] -- 1. Provides diag.zip --> B[Checker]
    B -- 2. Reads Logs --> C[Log Data<br/>- daemon<br/>- netstat<br/>- etc.] --> E
    C --> F
    B -- 3. Gets Config --> D[Config] 
    D -- 4a. Sets Tests --> E[Tests]
    D -- 4b. Sets Search --> F[Search]
    E --> G[5. Results]
    F --> G
    G -- 6. Shows Results --> A

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

[Requires Homebrew](https://brew.sh/)

```bash {#install-brew}
brew tap peakefficiency/releases
brew install warp-diag-checker
```

To update, run:

```
brew upgrade warp-diag-checker
```

### Windows  

Windows Chocolatey installation is no longer officially supported due to delays in approval from Chocolatey.

Please use the Go install method instead:

### Install issues

As a fallback, you can install from source using:

```
go install github.com/peakefficiency/warp-diag-checker@latest
```
You would need to have Go installed for this

## Usage


### Check for known issues:

```
warp-diag-checker check /path/to/diag.zip
```
![check-usage](/check-usage.gif)


### Get diag info summary: 

```
warp-diag-checker info /path/to/diag.zip
```
![info-usage](/info-usage.gif)

### Dump diag file to stdout

_This allows for processing with standard shell utilities_

```
warp-diag-checker dump /path/to/diag.zip [-f filename]
```
![dump-usage](/dump-usage.gif)

### List files in zip

```
warp-diag-checker list /path/to/diag.zip 
```
![list-usage](/list-usage.gif)