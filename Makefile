PROJECT := serene-foundry-234813
BUCKET_CLOUD_FUNCTIONS := ${PROJECT}-functions
BUCKET_OBJECTS := ${PROJECT}-objects
DC_GOV_URL := https://dcra.dc.gov/mrv

# Windows-specific settings.
ifdef OS
	SHELL := sh
	GSUTIL := gsutil.cmd
else
	GSUTIL := gsutil
endif


.PHONY: all backend frontend

help:
	@echo "  all               deploy the backend and frontend"
	@echo "  backend           deploy the backend"
	@echo "  frontend          deploy the frontend"
	@echo "  test              run tests"

all: frontend backend

frontend:
	@echo -e "\nDeploying the frontend"
	cd frontend && \
	npm run build && \
	firebase deploy

backend: db dcgov

db: db_rating db_trucks

db_rating:
	@echo -e "\nDeploying average ratings update"
	export TRIGGER_EVENT=$$(cat backend/db/rating/trigger_event); \
	export TRIGGER_RESOURCE=$$(cat backend/db/rating/trigger_resource); \
	${SHELL} gcloud functions deploy set-avg-rating \
	--entry-point=SetAvgRating \
	--runtime=go111 \
	--source=backend/db/rating \
	--stage-bucket=${BUCKET_CLOUD_FUNCTIONS} \
	--trigger-event=$${TRIGGER_EVENT} \
	--trigger-resource=projects/${PROJECT}/databases/\(default\)/documents/$${TRIGGER_RESOURCE}

db_trucks:
	@echo -e "\nUploading food truck data"
	export PROJECT=${PROJECT}; \
	cd backend/db/trucks && \
	go run trucks.go

dcgov: get_pdfs convert_pdfs load_db

get_pdfs: buckets get_pdfs_cron
	@echo -e "\nDeploying DC gov PDF retrieval"
	${SHELL} gcloud functions deploy get-pdfs \
	--entry-point=GetPDFs \
	--runtime=go111 \
	--source=backend/dcgov/get_pdfs \
	--stage-bucket=${BUCKET_CLOUD_FUNCTIONS} \
	--timeout=300 \
	--set-env-vars=URL=${DC_GOV_URL},BUCKET=${BUCKET_OBJECTS},PROJECT=${PROJECT} \
	--trigger-topic=get-pdfs

get_pdfs_cron:
	@echo -e "\nSetting up a schedule for PDF conversion"
	-${SHELL} gcloud pubsub topics create get-pdfs
	-${SHELL} gcloud scheduler jobs delete get-pdfs --quiet
	-${SHELL} gcloud scheduler jobs create pubsub get-pdfs \
	--schedule="0 * * * *" \
	--topic=get-pdfs \
	--message-body="{}" \
	--time-zone=America/New_York

convert_pdfs: buckets
	@echo -e "\nDeploying DC gov PDF to CSV conversion"
	${SHELL} gcloud functions deploy convert-pdfs \
	--entry-point=convert_pdf \
	--memory=512MB \
	--runtime=python37 \
	--source=backend/dcgov/convert_pdf \
	--stage-bucket=${BUCKET_CLOUD_FUNCTIONS} \
	--timeout=300 \
	--trigger-event=google.storage.object.finalize \
	--trigger-resource=${BUCKET_OBJECTS}

load_db: buckets
	@echo -e "\nDeploying DC gov database load"
	${SHELL} gcloud functions deploy load-db \
	--entry-point=LoadDB \
	--runtime=go111 \
	--source=backend/dcgov/load_db \
	--stage-bucket=${BUCKET_CLOUD_FUNCTIONS} \
	--set-env-vars=PROJECT=${PROJECT} \
	--trigger-event=google.storage.object.finalize \
	--trigger-resource=${BUCKET_OBJECTS}

buckets:
	@echo -e "\nCreating buckets"
	-${GSUTIL} mb gs://${BUCKET_CLOUD_FUNCTIONS}
	-${GSUTIL} mb gs://${BUCKET_OBJECTS}

test:
	@echo -e "\nRunning tests"
