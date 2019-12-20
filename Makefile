SRCDIR=src

.PHONY: all image functest test

all: image

image: build
	docker build -t applause/metronomikon:latest .

functest:
	# Functional tests
	test/run.sh

%:
	$(MAKE) -C $(SRCDIR) $@
