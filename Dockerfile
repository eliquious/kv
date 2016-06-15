FROM scratch

ADD main /

EXPOSE 9034 9034
ENTRYPOINT ["/main"]
