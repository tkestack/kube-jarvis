FROM ubuntu
ADD conf /conf
ADD cmd/kube-jarvis/kube-jarvis /
ADD translation /translation
CMD ["/kube-jarvis"]