albio
=====

A CLI tool to gracefully manage instances on ELB or ALB.

```shell
$ albio status
[
    {
        "name": "elbtest001",
        "dnsname": "elbtest001-00000000.ap-northeast-1.elb.amazonaws.com",
        "instances": [
            {
                "instance-id": "i-xxxxxxxxxxxxxxxxx"
            },
            {
                "instance-id": "i-yyyyyyyyyyyyyyyyy"
            }
        ]
    }
]
```

```shell
$ albio detach
$ # something restart command
$ albio attach
```

## License

[The MIT License](./LICENSE).

## Author

[y_uuki](https://github.com/yuuki)
