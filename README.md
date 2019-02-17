# cwli

## Usage

```
$ cwli -h
Usage of cwli:
  -query string
      Query string
  -showptr
      Show @ptr in query result
  -version
      Print version and exit
```

```
$ cwli
ap-northeast-1> show groups;
/aws/kinesisfirehose/my-kinesis
/aws/lambda/my-lambda
/aws/rds/cluster/my-rds/slow-query

ap-northeast-1> head /aws/lambda/my-lambda limit 1;
*************************** 1. row ***************************
Timestamp: 1516438129835
Message: START RequestId: bbd0dbdd-fdbe-11e7-a41f-398ced7d9303 Version: $LATEST

ap-northeast-1> source /aws/lambda/my-lambda start=2018/11/19 end=2019/11/21 | field @timestamp, @message | limit 2;
{"@timestamp":"2019-04-09 09:38:19.455","@message":"END RequestId: ab83228c-6af1-4c0b-a304-e043aecbe84a\n"}
{"@timestamp":"2019-04-09 09:38:19.455","@message":"2019-04-09T09:38:19.455Z\tab83228c-6af1-4c0b-a304-e043aecbe84a\tLogEC2InstanceStateChange\n"}
// Status: Complete
// Statistics: BytesScanned=859801 RecordsMatched=2065 RecordsScanned=2065
```

### Passing to external command

```
ap-northeast-1> source /aws/lambda/my-lambda start=2018/11/19 end=2019/11/21 | field @timestamp, @message ! head -n 1 | tee output.jsonl;
{"@timestamp":"2019-04-09 09:38:19.455","@message":"2019-04-09T09:38:19.455Z\tab83228c-6af1-4c0b-a304-e043aecbe84a\tLogEC2InstanceStateChange\n"}
```

```
$ cat output.jsonl
{"@timestamp":"2019-04-09 09:38:19.455","@message":"2019-04-09T09:38:19.455Z\tab83228c-6af1-4c0b-a304-e043aecbe84a\tLogEC2InstanceStateChange\n"}
```
