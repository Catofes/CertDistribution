all:
		mkdir -p build
		go build -o build/web -ldflags "-X 'main._version_=$(shell git log --pretty=format:"%h" -1)'"\
			                          	 github.com/Catofes/CertDistribution


