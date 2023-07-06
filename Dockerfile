FROM golang:alpine 

WORKDIR /app

COPY . .

# WORKDIR /app

RUN go build -o main ./cmd/main.go

EXPOSE 8000

CMD [ "./main" ]
#CMD ["sh", "-c", "while true; do sleep 1000; done"]