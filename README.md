```
cd journal
docker build -t go-journal .
docker run -v DIR:/journal/log:ro -v FILE:/journal/tail go-journal
```
```
cd aws
docker build -t go-aws .
docker run -e AWS_REGION=REGION -e AWS_ACCESS_KEY_ID=ID -e AWS_SECRET_ACCESS_KEY=KEY -e GROUP:AWS_LOG_GROUP -e STREAM:AWS_LOG_STREAM -v FILE:/aws/tail go-aws
```
