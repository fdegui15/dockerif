FROM ubuntu:15.04
MAINTAINER fdeguilhen@gmail.com

ENV PATH $PATH:/home/go/bin
ENV GOPATH /home/go
ENV MOUNTDIR /opt/dv

# All the installation here
# Reducing the image size with --no-install-recommends and rm /var/lib/apt/lists
RUN apt-get update && \
	apt-get install -y --no-install-recommends ca-certificates golang git clamav clamav-daemon python mediainfo exiftool && \
	rm -rf /var/lib/apt/lists

# Installation of siegfried 
# RUN apt-get install -y golang git
RUN mkdir /home/go ;\
	go get github.com/richardlehane/siegfried/cmd/sf ;\
	sf -update 

# Installation of ClamAV
# RUN apt-get install -y clamav clamav-daemon
RUN freshclam
RUN sed -i "s|/var/run/clamav/clamd.ctl|/tmp/clamd.ctl|" /etc/clamav/clamd.conf
RUN sed -i "s|User clamav|User root|" /etc/clamav/clamd.conf

# Installation of Fido
# RUN apt-get install -y python
RUN mkdir /home/fido
ADD fido /home/fido
RUN chmod +x /home/fido/fido.py

# Installation of mediainfo & exiftool
# RUN apt-get install -y mediainfo exiftool
RUN locale-gen fr_FR.UTF-8

# Add and Compile code inspectFile
ADD inspectfile /home/go/src/inspectfile
RUN go build -o /home/go/bin/inspectFile /home/go/src/inspectfile/*.go

#Add HTML file
ADD demo /home/demo

EXPOSE 8080


CMD ["inspectFile","--server","0.0.0.0:8080"]

# Compile = "docker build -t inspectfile ." (in the directory containing Dockerfile)
# Launch = "docker run -d -p 8080:8080 --name dv inspectfile"
# or "docker run -d -p 8080:8080 -v /media/sf_Temp:/opt/dv --name dv inspectfile" if you want to mount the directory /media/sf_Temp to inspect file in this directory.
# To see the output "docker logs dv"
# To stop "docker stop dv"
# To remove "docker rm dv"


# To inspect a file "curl localhost:8080/inspect -F file=@/path/and/file"
# To inspect a local file "curl localhost:8080/localinspect\&file=/path/to/file" /path/to/file is relative to /media/sf_Temp in the host !!
# To inspect all the files in a directory "curl localhost:8080/inspectpath/opt/dv/path/to/inspect"
# You can also launch it in your browser if you don't upload the file !!

# To update the virus database, "curl localhost:8080/avupdate"
# Get the tools version used (contains also the virus and pronom databases): "curl localhost:8080/gettoolsversion"

