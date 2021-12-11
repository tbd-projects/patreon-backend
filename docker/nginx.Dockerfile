# The image to pull the base configuration from
FROM nginx:latest

RUN apt-get update
RUN apt-get install -y certbot python3-certbot-nginx

RUN mkdir -p /etc/letsencrypt
RUN mkdir -p /static
RUN chmod 777 -R ./static
