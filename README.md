# albio

A CLI tool to gracefully manage EC2 instances on ALB/NLB.

## Usage

### Show loadbalancers status

```shell
$ albio status
[
    {
        "name": "albtest001",
        "dnsname": "albtest001-00000000.ap-northeast-1.elb.amazonaws.com",
        "arn": "arn:aws:elasticloadbalancing:ap-northeast-1:100000000000:loadbalancer/app/albtest001/11xxxxxxxxxxxxxxxxxx"
        "type": "application",
        "scheme": "internet-facing",
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

### Take an instance under maintenance

```shell
$ sudo mkdir -m 777 -p /var/lib/albio
$ albio detach --self > /var/lib/albio/loadbalancers.json
--> Detaching i-xxxxxxxxxxxxxxxxx from albtest001
$ # something to restart command
$ cat /var/lib/albio/loadbalancers.json | albio attach --self
--> Attaching i-0f5ffb9f0a75e6b85 to albtest001
```

## License

[The MIT License](./LICENSE).

## Author

[y_uuki](https://github.com/yuuki)
