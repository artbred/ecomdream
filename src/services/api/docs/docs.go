// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/payments/create": {
            "post": {
                "description": "Create payment link",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payments"
                ],
                "summary": "Create payment link",
                "parameters": [
                    {
                        "description": "Payment data",
                        "name": "payment_data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/payments.CreatePaymentLinkRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/payments.CreatePaymentLinkResponse"
                        }
                    }
                }
            }
        },
        "/v1/prompts/create/{id}": {
            "post": {
                "description": "Start prediction for prompt",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "prompts"
                ],
                "summary": "Start prediction for prompt",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Version ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Prompt data",
                        "name": "prompt_data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/prompts.CreatePromptRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/prompts.CreatePromptRequest"
                        }
                    }
                }
            }
        },
        "/v1/prompts/list/{id}": {
            "get": {
                "description": "Returns prompts and images for version",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "prompts"
                ],
                "summary": "Returns prompts and images for version",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Version ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/prompts.CreatePromptRequest"
                        }
                    }
                }
            }
        },
        "/v1/versions/info/{id}": {
            "get": {
                "description": "Get info about version",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "versions"
                ],
                "summary": "Get info about version",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Version ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/versions.VersionInfoResponse"
                        }
                    }
                }
            }
        },
        "/v1/versions/train/{id}": {
            "post": {
                "description": "Start training process",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "versions"
                ],
                "summary": "Start training process",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Class name",
                        "name": "class",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Payment ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "file"
                        },
                        "description": "Data",
                        "name": "data",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/versions.TrainVersionResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "payments.CreatePaymentLinkRequest": {
            "type": "object",
            "properties": {
                "plan_id": {
                    "type": "integer"
                },
                "promocode_id": {
                    "type": "string"
                },
                "version_id": {
                    "type": "string"
                }
            }
        },
        "payments.CreatePaymentLinkResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "prompts.CreatePromptRequest": {
            "type": "object",
            "properties": {
                "amount_images": {
                    "type": "integer"
                },
                "negative_prompt": {
                    "type": "string"
                },
                "prompt": {
                    "type": "string"
                }
            }
        },
        "versions.TrainVersionResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                },
                "version_id": {
                    "type": "string"
                }
            }
        },
        "versions.VersionInfoResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "info": {},
                "is_ready": {
                    "type": "boolean"
                },
                "time_training": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
