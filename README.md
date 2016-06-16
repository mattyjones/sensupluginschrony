## sensupluginschrony

## Commands
 * checkChronyStats

## Usage

### checkChronyStats
Check specific stats for checkChronyStats

Ex. `./checkChronyStats --checkKey foo --warn 10 --crit 20`

Available keys are

- ReferenceID is a straight shot, it is either critical or not
- Stratum is the number of hops away, critical and warning are that number
- ReferenceTime is the time the last measurement from a source was processed, critical and warning values represent the number of seconds deviation

## Installation

1. godep go build -o bin/sensupluginschrony
1. chmod +x sensupluginschrony (*nix only)
1. cp sensupluginschrony /usr/local/bin

## Notes  
