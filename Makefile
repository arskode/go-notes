.PHONY: test
test:
	docker-compose -f docker-compose.testing.yml run --rm backend 
	docker-compose stop