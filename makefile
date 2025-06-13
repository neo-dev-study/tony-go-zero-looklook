

.PHONY: Docker_Mac_Env
Docker_Mac_Env:
	docker-compose -f docker-compose-env-mac.yml up -d

.PHONY: Docker_Mac_Start
Docker_Mac_Start:
	docker-compose up -d 

