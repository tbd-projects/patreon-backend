# The image to pull the base configuration from
FROM nginx:mainline-alpine
# The directory where any additional files will be referenced
WORKDIR /app
# Copy the custom default.conf from the WORKDIR (.) and overwrite the existing internal configuration in the NGINX container
COPY ./configs/nginx.conf /etc/nginx/conf.d/default.conf