```
cd aws
docker build -t go-aws .
docker run -e AWS_REGION=REGION -e AWS_ACCESS_KEY_ID=ID -e AWS_SECRET_ACCESS_KEY=KEY -v FILE:/aws/tail go-aws
```
