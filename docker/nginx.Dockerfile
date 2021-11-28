# The image to pull the base configuration from
FROM nginx:latest

RUN apt-get update
RUN apt-get install -y certbot python3-certbot-nginx

RUN mkdir -p /etc/letsencrypt/live/api.pyaterochka-team.site