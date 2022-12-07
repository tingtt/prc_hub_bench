OAS := https://raw.githubusercontent.com/tingtt/prc_hub_back/e-isucon/presentation/oas.yml
.PHONY: oas-code-gen
oas-code-gen:
	oapi-codegen --config infrastructure/externalapi/backend/client.yaml ${OAS}
	oapi-codegen --config infrastructure/externalapi/backend/models.yaml ${OAS}
