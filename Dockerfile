FROM scratch
COPY bin/captcha.go /captcha
ENTRYPOINT ["captcha"]
