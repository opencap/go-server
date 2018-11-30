## Currently this dockefile is used only for tests, no docker image is used in production
FROM golang:1.10-stretch

ENV TEST_PORT=8080
ENV DB_URL=test.db
ENV DB_TYPE=sqlite3
ENV JWT_EXPIRATION_MINUTES=30
ENV JWT_SECRET=DFUIJHSDFAJDLFHBSDFLSDFHJSALFIGHDSFKGHDFLKG
ENV PLATFORM_ENV=test
ENV CREATE_USER_PASSWORD=examplepassword
ENV DOMAIN_NAME=example.com