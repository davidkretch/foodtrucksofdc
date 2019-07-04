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


.PHONY: all backend frontend get_pdfs convert_pdfs load_db

help:
	@echo "  all               deploy the backend and frontend"
	@echo "  backend           deploy the backend"
	@echo "  frontend          deploy the frontend"
	@echo "  test              run tests"
	@echo "  get_pdfs          deploy DC gov PDF retrieval"
	@echo "  convert_pdfs      deploy DC gov PDF conversion"
	@echo "  load_db           deploy DC gov database load"

all: backend frontend

backend: get_pdfs convert_pdfs load_db

frontend:
	@echo -e "\nDeploying the frontend"
	cd frontend && \
	npm run build && \
	firebase deploy

buckets:
	@echo -e "\nCreating buckets"
	-${GSUTIL} mb gs://${BUCKET_CLOUD_FUNCTIONS}
	-${GSUTIL} mb gs://${BUCKET_OBJECTS}

get_pdfs_cron:
	@echo -e "\nSetting up a schedule for PDF conversion"
	-${SHELL} gcloud pubsub topics create get-pdfs
	-${SHELL} gcloud scheduler jobs delete get-pdfs --quiet
	-${SHELL} gcloud scheduler jobs create pubsub get-pdfs \
	--schedule="0 * * * *" \
	--topic=get-pdfs \
	--message-body="{}" \
	--time-zone=America/New_York

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

test:
	@echo -e "\nRunning tests"
