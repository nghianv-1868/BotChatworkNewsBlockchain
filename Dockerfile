# should match cloud.gov manifest GOVERSION
FROM golang:1.15.5-alpine


# install git
RUN apk add git

# install google.golang.org/api/sheets/v4
RUN go get -u google.golang.org/api/sheets/v4
# install golang.org/x/oauth2/google
RUN go get -u golang.org/x/oauth2/google
# install github.com/gorilla/mux
RUN go get -u github.com/gorilla/mux

# set working directory
WORKDIR /app

# default command that is executed if none specified
CMD ["sh","-c","make build"]
