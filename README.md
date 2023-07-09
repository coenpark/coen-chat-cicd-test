## 참고한 자료(도움된 자료)

[표준 Go 프로젝트 레이아웃](https://github.com/golang-standards/project-layout/blob/master/README_ko.md)

[gonic middleware로 사용하기](https://gin-gonic.com/ko-kr/docs/examples/using-middleware/)

[Docker redis-cli로 접근하기](https://jistol.github.io/docker/2017/09/01/docker-redis/)

[gonic middleware로 사용하기](https://gin-gonic.com/ko-kr/docs/examples/using-middleware/)

[local에서 GCP 인증 정보 설정](https://cloud.google.com/docs/authentication/provide-credentials-adc?hl=ko)

## docker cli
- docker build --tag chat .
- docker run -d --name chat [...moreOption?] chat

> 빌드시에 주의할점!
> 이미지 이름은 [도커허브 유저네임]/[레포지토리네임]:[테크네임]

- docker run -v /Users/coen/.config/gcloud/application_default_credentials.json:/path/in/container/application_default_credentials.json -d -p 8080:8080 --name=chat coenpark/chat

- docker run -v /credentials/application_default_credentials.json:/path/in/container/application_default_credentials.json

- docker buildx build --platform=linux/amd64/v3 -t coenpark/chat .
