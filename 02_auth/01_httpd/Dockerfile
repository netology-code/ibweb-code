FROM httpd:2.4-alpine
RUN apk add --no-cache apr-util-dbd_pgsql postgresql-client
COPY httpd.conf /usr/local/apache2/conf/httpd.conf
COPY htdocs /usr/local/apache2/htdocs
COPY *.txt /usr/local/apache2/passwd/
COPY wait-for-postgres.sh wait-for-postgres.sh
RUN chmod +x wait-for-postgres.sh
CMD ["./wait-for-postgres.sh", "httpd-foreground"]

