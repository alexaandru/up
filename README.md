# Up

A tiny sftp UPloading client, so tiny - it doesn't even need a longer name :)

## Purpose

The purpose was to have an uploader that can perform extra verification steps.
For now, it only checks the size (locally and remote) but other checks will likely
be added in the future.

## Setup

No setup is needed, all params can be passed in commandline. Please see `./up -h` for help,
copied here for convenience:

    Usage of ./up:
      -addr string
            Address to connect to (default "172.17.0.2:22")
      -dst string
            Destination folder (default "/tmp")
      -pass string
            Password to connect with (default "1234")
      -user string
            Username to connect as (default "test")
      <filename> required
            The filename to be uploaded. There is no default value.

Optionally, if there is an `up.json` file in the current folder, the options will be read from it.
The commandline has precedence over the config file.

Sample config file:

```JSON
{
	"Addr":      "172.17.0.2:22",
	"User":      "test",
	"Password":  "1234",
	"DstFolder": "/tmp",
}
```

## Use

Please see `./up -h` for the latest usage instructions. Sample use cases:

```Bash
./up some_file.txt # uploaded using defaults or up.json if exists
./up -user foo -pass bar f.txt # uploaded using custom user/pass, defaults for rest
```
On error it will log a message to stderr and return with a non-zero exit code.

## Testing

This repo includes a Dockerfile, that will create a tiny Alpine Linux based sshd server which
can be used for testing:

```Bash
docker build -t somename/sshd .
docker run somename/sshd
# in another terminal
go build
./up some_file.txt # should work out of the box. The defaults in .go and Dockerfile match.
```

## Limitations

This is strictly an "one file at a time" uploader. It has no other functionality (upload recursively,
download, etc.) and there are no plans to add it.

Currently, it can only use password for authentication, but this may change in the future.

## License

See ./LICENSE for details.
