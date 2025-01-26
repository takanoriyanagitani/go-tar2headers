#!/bin/sh

gensample(){
	pushd .
	cd sample.d

	touch empty.txt
	echo hw > hw.txt

	ln -f hw.txt hw2.txt

	ln -f -s hw.txt hw3.txt

	mkdir -p empty.d

	popd

	tar --create --verbose --file sample.tar sample.d/
}

gensample

cat ./sample.tar |
	./tar2headers |
	jq -c
