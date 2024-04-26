FROM testdocker-registry.web-site.com:5000/alpine:tzdata
COPY config/ /server/config/
ADD ./http-server /server/http-server
RUN  chmod u+x /server/http-server
WORKDIR /server
EXPOSE 50051
COPY run-http.sh /run-http.sh
RUN chmod u+x /run-http.sh
CMD ["/run-http.sh"]