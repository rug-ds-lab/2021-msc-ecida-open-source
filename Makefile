.ONESHELL:

images: 
	docker build -t celciusfarenheit:v1 src/examples/celcius-farenheit
	docker build -t logger:v2 src/examples/logger
