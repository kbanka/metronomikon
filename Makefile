SRCDIR=src

.PHONY: all image

all: image

image: build
	docker build -t applause/metronomikon:latest .

%:
	$(MAKE) -C $(SRCDIR) $@
