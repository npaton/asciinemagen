# asciinemagen

Create an asciinemagen cast file from a flat text file. By default lines
starting with `>` (configurable with `iprompt` flag) are considered command
lines and other lines are considered to be command output (stdout/err, no
distinction for asciinema).

You can control which timing, prompts, and asciinema metadata with flags, see
`asciinemagen --help` for the details.
