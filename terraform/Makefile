.PHONY: default
default:
	@echo "make (dev|staging|production)"

.PHONY: dev
dev:
	cd dev ; make all

.PHONY: cfclean
cfclean:
	echo "------- CFCLEAN -------"
	for app in admin entree extract fetch pack serve walk ; do cf delete -f -r $${app} ; done
	for service in db_queues db_search db_work ; do cf delete-service -f -w $${service} ; done
	for service in s3_backup s3_fetch s3_extract s3_meta ; do cf delete-service -f -w $${service} ; done
	