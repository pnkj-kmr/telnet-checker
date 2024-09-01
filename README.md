# telnet-checker

telnet-checker helps to run multiple commands on end router (switch) in bulk.

_Basically, if you want to run few commands on end router and you want to do for multiple routers
then telnet-checker helps to perform the action with input file and generate a output file as result_

### HOW TO USE

_Download the relevent os package from [here](https://github.com/pnkj-kmr/telnet-checker/releases)_

_create a **input.json** file_

```
[
    {
        "host": "192.168.1.1",
        "tag": "",
        "commands": [
            {
                "expect":"name: ",      # omitempty
                "command": "user1"
            },{
                "expect":"assword: ",   # omitempty
                "command": "pas@123"
            },{
                "expect":"#",           # omitempty
                "command": "show run",
                "eof":"#"               # omitempty
            }
        ],
        "port": 23,                     # omitempty default: 23
        "timeout": 30                   # omitempty default: 20
    },
    ...
]
```

_After creating the file run the executable binary as_

```
./telnetchecker
```

### OUTPUT


_As a result **output.json** file will be created after completion_

```
[
  {
    "input": {
       {
        "host": "192.168.1.1",
        "tag": "",
        "commands": [
            {
                "expect":"name: ",
                "command": "user1"
            },{
                "expect":"assword: ",
                "command": "pas@123"
            },{
                "expect":"#",
                "command": "show run",
                "eof":"#" 
            }
        ],
        "port": 23,
        "timeout": 30
    }
    },
    "error": ["telnet: handshake failed: EOF"],
    "output": []            # result success result if any
  }
]

```

### HELP

```
./telnetchecker --help

----------------------
Usage of ./telnetchecker:
  -ct int
        timeout [secs] (default 120)
  -f string
        input file name (default "input.json")
  -o string
        output file name (default "output.json")
  -p int
        generic default port (default 23)
  -t int
        timeout [secs] (default 20)
  -w int
        number of workers (default 4)

-------
Example:

./telnetchecker -f x.json -t 30 -w 20

```

:)
