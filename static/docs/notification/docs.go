// Package notification GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package notification

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Iliyan Motovski",
            "email": "iliyan.motovski@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/notification/send-async": {
            "post": {
                "description": "Send notification async",
                "tags": [
                    "Notification"
                ],
                "parameters": [
                    {
                        "description": "Request in body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.RequestDTONotification"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.RequestDTOEmail": {
            "type": "object",
            "required": [
                "from",
                "from_name",
                "subject",
                "template_name",
                "to"
            ],
            "properties": {
                "from": {
                    "type": "string"
                },
                "from_name": {
                    "type": "string"
                },
                "reply_to": {
                    "type": "string"
                },
                "subject": {
                    "type": "string"
                },
                "template_data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.RequestDTOTemplateDate"
                    }
                },
                "template_name": {
                    "type": "string"
                },
                "to": {
                    "type": "string"
                }
            }
        },
        "dto.RequestDTONotification": {
            "type": "object",
            "properties": {
                "email": {
                    "$ref": "#/definitions/dto.RequestDTOEmail"
                },
                "slack_message": {
                    "$ref": "#/definitions/dto.RequestDTOSlackMessage"
                },
                "sms": {
                    "$ref": "#/definitions/dto.RequestDTOSMS"
                }
            }
        },
        "dto.RequestDTOSMS": {
            "type": "object",
            "required": [
                "mobile_number",
                "text"
            ],
            "properties": {
                "mobile_number": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                }
            }
        },
        "dto.RequestDTOSlackMessage": {
            "type": "object",
            "required": [
                "bot_name",
                "channel_name",
                "message"
            ],
            "properties": {
                "bot_name": {
                    "type": "string"
                },
                "channel_name": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "dto.RequestDTOTemplateDate": {
            "type": "object",
            "required": [
                "key",
                "value"
            ],
            "properties": {
                "key": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "server.Error": {
            "type": "object",
            "properties": {
                "fields_error": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "global_error": {
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
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "NOTIFICATION API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
