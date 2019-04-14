# cwli

cwli is [CloudWatch Logs Insights](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/AnalyzingLogData.html) Command-Line Client.

see https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/CWL_QuerySyntax.html

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
ap-northeast-1> help;
help
	Print a help message
show groups [prefix LOG-GROUP-PREFIX]
	Display Log Groups
head LOG-GROUP [limit N]
	Print first lines of a Log Group
source LOG-GROUP start=START-TIME end=END-TIME | field ... [! EXTERNAL-COMMAND]
	Perform CloudWatch Logs Insights query
```

```
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

### Print vertically

```
ap-northeast-1> source /aws/lambda/my-lambda start=2018/11/19 end=2019/11/21 | field @timestamp, @message | limit 3 | vertically;
*************************** 1. row ***************************
@timestamp: 2019-04-09 09:38:19.455
@message: REPORT RequestId: ab83228c-6af1-4c0b-a304-e043aecbe84a  Duration: 0.66 ms Billed Duration: 100 ms   Memory Size: 128 MB Max Memory Used: 48 MB

*************************** 2. row ***************************
@timestamp: 2019-04-09 09:38:19.455
@message: 2019-04-09T09:38:19.455Z  ab83228c-6af1-4c0b-a304-e043aecbe84a  LogEC2InstanceStateChange

*************************** 3. row ***************************
@timestamp: 2019-04-09 09:38:19.455
@message: 2019-04-09T09:38:19.455Z  ab83228c-6af1-4c0b-a304-e043aecbe84a  Received event: {
  "version": "0",
  "id": "4b22e040-e3dd-a0aa-84e0-a946526876a7",
  "detail-type": "AWS API Call via CloudTrail",
  "source": "aws.ec2",
  "account": "822997939312",
  "time": "2019-04-09T09:37:35Z",
  "region": "ap-northeast-1",
...
```
