#  `book-it` - Go Command Line Reservation Tool

Book-It is a command line tool written in Go that allows users to book reservations via the Resy API.
This tool aims to allow you to book at those hard-to-get restaurants!

## Installation

Currently, the program needs to be compiled locally
by cloning this repo and running `go build`

```shell
# Clone the repoository
git clone github.com/bthuilot/book-it
# cd into cloned repository
cd book-it
# Build the program
go build -o book-it
```

## Usage

To use book-it, run the `book-it` command with the required flags:


### Required Flags
*NOTE:* run with `--help` for most up-to-date docs

* `-d, --date string`: The date of when to book in 'mm/dd/yy' format.
* `-p, --party-size int`: The amount of people the reservation should be for.
* `-t, --time string`: The time to make the reservation for in kitchen time format (i.e. 3:04PM).
* `-v, --venue-id int`: ID of the venue you want to reserve from.

### Optional Flags

* `--include-types stringArray`: Filter for what types of reservations to include, i.e. Booth, Outdoor, etc. Must match the exact name shown on Resy.
* `-s, --time-spread int`: The 'spread' of acceptable times. When set, the program will consider times within this duration of seconds from the targeted time to be acceptable reservations.
* `-h, --help`: Show help message and exit.

## Examples

1. Running the tool from the CLI

```shell
# Book an indoor reservation at UVA Next door (ID 50830) in NYC 
# for July 15, 2023 any time from 10-11PM (spread of 30min from 10:30)
book-it -v 50830 -d 07/15/23 -t 10:30PM -s 1800 -p 4 --include-types Indoor 
```

2. Setting up a cronjob to run at midnight (for those "booked in a second" restaurants)

```text 
0 0 * * * /home/$USER/book-it -v 12345 -d 07/15/23 -t 9:00PM -p2 
```
