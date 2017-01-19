FROM scratch
ADD kube-test-container /
ADD template/index.html /

EXPOSE 8000
ENTRYPOINT ["/kube-test-container"]
