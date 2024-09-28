# Usar un contenedor para compilar el binario
FROM golang:1.22-alpine as builder

RUN mkdir /app
WORKDIR /app

COPY . .


RUN go mod download
RUN go build -v -o reverseproxy ./cmd/main.go
RUN chmod +x /app/reverseproxy

# Crear una imagen más pequeña solo con el binario
FROM alpine:latest


ARG port
ARG mongo_user
ARG mongo_pwd
ARG mongo_url
ARG mongo_port
ARG mongo_db 

ENV port=$port
ENV mongo_user=$mongo_user
ENV mongo_pwd=$mongo_pwd
ENV mongo_url=$mongo_url
ENV mongo_port=$mongo_port
ENV mongo_db =$mongo_db 

# Variable de entorno para manejar la hora de Colombia
ENV TZ=America/Bogota

RUN mkdir /app
RUN apk update

COPY --from=builder /app/reverseproxy /app/reverseproxy
COPY /web/dist /app/web
COPY config.json /app/config.json

EXPOSE 80

# Run the web service on container startup.
CMD [ "/app/reverseproxy" ]
#ENTRYPOINT ["tail", "-f", "/dev/null"]