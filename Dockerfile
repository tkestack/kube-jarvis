FROM ubuntu
ADD conf /conf
ADD kube-jarvis /
ADD translation /translation
CMD ["/kube-jarvis"]