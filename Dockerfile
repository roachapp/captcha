FROM scratch
COPY bin/captcha /captcha
ENTRYPOINT ["/captcha"]
