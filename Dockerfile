## Currently this dockefile is used only for tests, no docker image is used in production ##
FROM golang:1.10-stretch

ENV PORT=8080
ENV DB_URL=postgres://postgres:postgres@postgres:5432/test?sslmode=disable
ENV DB_TYPE=postgres
ENV JWT_EXPIRATION_MINUTES=30
ENV JWT_SECRET=fd89asdufiasfjsbfaujsdjf[asf8ashf[asubfksjdfjsnfjkdsf]]
ENV DOMAIN=example.com
ENV PLATFORM_ENV=test
ENV TIMEOUT_SECONDS=30