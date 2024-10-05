FROM debian:latest
 
ARG APP_VERSION="v0.0.0"


ENV PORT=9001
EXPOSE ${PORT}

ENV DOMAIN=0.0.0.0
ENV USELETSENCRYPT=N

# RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections
WORKDIR /app
COPY ./bin/gomockapi_${APP_VERSION} ./gomockapi
 

RUN apt update && \
    chmod +x  ./gomockapi

CMD [ "./gomockapi" ]

 