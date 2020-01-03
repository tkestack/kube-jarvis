all:
	go build -o  bin/kube-jarvis cmd/kube-jarvis/*.go
release:all
	mkdir kube-jarvis
	cp -R conf kube-jarvis/
	cp -R translation kube-jarvis/
	cp bin/kube-jarvis kube-jarvis/
	cp -R manifests kube-jarvis/
	tar cf kube-jarvis.tar.gz kube-jarvis
	rm -rf kube-jarvis
clean:
	rm -f kube-jarvis.tar.gz
	rm -fr kube-jarvis
	rm -fr bin/
test:
	go test ./pkg/...
