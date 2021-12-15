// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
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
        "/": {
            "get": {
                "description": "get the status of server.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "Show the status of server.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/admin/user/:user_id/devices": {
            "post": {
                "description": "meant for internal admin use - adds a device to a user. can add with only device_definition_id or with MMY, which will create a device_definition on the fly",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user-devices"
                ],
                "parameters": [
                    {
                        "description": "add device to user. either MMY or id are required",
                        "name": "user_device",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.AdminRegisterUserDevice"
                        }
                    },
                    {
                        "type": "string",
                        "description": "user id",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/controllers.RegisterUserDeviceResponse"
                        }
                    }
                }
            }
        },
        "/device-definitions": {
            "get": {
                "description": "gets a specific device definition by make model and year",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "device-definitions"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "make eg TESLA",
                        "name": "make",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "model eg MODEL Y",
                        "name": "model",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "year eg 2021",
                        "name": "year",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.DeviceDefinition"
                        }
                    }
                }
            }
        },
        "/device-definitions/all": {
            "get": {
                "description": "returns a json tree of Makes, models, and years",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "device-definitions"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/controllers.DeviceMMYRoot"
                            }
                        }
                    }
                }
            }
        },
        "/device-definitions/vin/{vin}": {
            "get": {
                "description": "decodes a VIN by first looking it up on our DB, and then calling out to external sources. If it does call out, it will backfill our DB",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "device-definitions"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "VIN eg. 5YJ3E1EA6MF873863",
                        "name": "vin",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.DeviceDefinition"
                        }
                    }
                }
            }
        },
        "/device-definitions/{id}": {
            "get": {
                "description": "gets a specific device definition by id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "device-definitions"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "device definition id, KSUID format",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.DeviceDefinition"
                        }
                    }
                }
            }
        },
        "/device-definitions/{id}/integrations": {
            "get": {
                "description": "gets all the available integrations for a device definition. Includes the capabilities of the device with the integration",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "device-definitions"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "device definition id, KSUID format",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/controllers.DeviceCompatibility"
                            }
                        }
                    }
                }
            }
        },
        "/user/devices": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    },
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "adds a device to a user. can add with only device_definition_id or with MMY, which will create a device_definition on the fly",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user-devices"
                ],
                "parameters": [
                    {
                        "description": "add device to user. either MMY or id are required",
                        "name": "user_device",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.RegisterUserDevice"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/controllers.RegisterUserDeviceResponse"
                        }
                    }
                }
            }
        },
        "/user/devices/:user_device_id/name": {
            "patch": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "updates the Name on the user device record",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user-devices"
                ],
                "parameters": [
                    {
                        "description": "Name",
                        "name": "name",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.UpdateNameReq"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": ""
                    }
                }
            }
        },
        "/user/devices/:user_device_id/vin": {
            "patch": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "updates the VIN on the user device record",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user-devices"
                ],
                "parameters": [
                    {
                        "description": "VIN",
                        "name": "vin",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.UpdateVINReq"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": ""
                    }
                }
            }
        },
        "/user/devices/me": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "gets all devices associated with current user - pulled from token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user-devices"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/controllers.UserDeviceFull"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.AdminRegisterUserDevice": {
            "type": "object",
            "properties": {
                "country_code": {
                    "type": "string"
                },
                "created_date": {
                    "description": "unix timestamp",
                    "type": "integer"
                },
                "device_definition_id": {
                    "type": "string"
                },
                "image_url": {
                    "type": "string"
                },
                "make": {
                    "type": "string"
                },
                "model": {
                    "type": "string"
                },
                "vehicle_name": {
                    "type": "string"
                },
                "verified": {
                    "type": "boolean"
                },
                "vin": {
                    "type": "string"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "controllers.DeviceCompatibility": {
            "type": "object",
            "properties": {
                "capabilities": {
                    "type": "string"
                },
                "country": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "style": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                },
                "vendor": {
                    "type": "string"
                }
            }
        },
        "controllers.DeviceDefinition": {
            "type": "object",
            "properties": {
                "compatible_integrations": {
                    "description": "CompatibleIntegrations has systems this vehicle can integrate with",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/controllers.DeviceCompatibility"
                    }
                },
                "device_definition_id": {
                    "type": "string"
                },
                "image_url": {
                    "type": "string"
                },
                "metadata": {},
                "name": {
                    "type": "string"
                },
                "type": {
                    "$ref": "#/definitions/controllers.DeviceType"
                },
                "vehicle_data": {
                    "description": "VehicleInfo will be empty if not a vehicle type",
                    "$ref": "#/definitions/services.DeviceVehicleInfo"
                }
            }
        },
        "controllers.DeviceMMYRoot": {
            "type": "object",
            "properties": {
                "make": {
                    "type": "string"
                },
                "models": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/controllers.DeviceModels"
                    }
                }
            }
        },
        "controllers.DeviceModelYear": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "controllers.DeviceModels": {
            "type": "object",
            "properties": {
                "model": {
                    "type": "string"
                },
                "years": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/controllers.DeviceModelYear"
                    }
                }
            }
        },
        "controllers.DeviceType": {
            "type": "object",
            "properties": {
                "make": {
                    "type": "string"
                },
                "model": {
                    "type": "string"
                },
                "sub_model": {
                    "type": "string"
                },
                "type": {
                    "description": "Type is eg. Vehicle, E-bike, roomba",
                    "type": "string"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "controllers.RegisterUserDevice": {
            "type": "object",
            "properties": {
                "country_code": {
                    "type": "string"
                },
                "device_definition_id": {
                    "type": "string"
                },
                "make": {
                    "type": "string"
                },
                "model": {
                    "type": "string"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "controllers.RegisterUserDeviceResponse": {
            "type": "object",
            "properties": {
                "device_definition_id": {
                    "type": "string"
                },
                "integration_capabilities": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/controllers.DeviceCompatibility"
                    }
                },
                "user_device_id": {
                    "type": "string"
                }
            }
        },
        "controllers.UpdateNameReq": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "controllers.UpdateVINReq": {
            "type": "object",
            "properties": {
                "vin": {
                    "type": "string"
                }
            }
        },
        "controllers.UserDeviceFull": {
            "type": "object",
            "properties": {
                "country_code": {
                    "type": "string"
                },
                "custom_image_url": {
                    "type": "string"
                },
                "device_definition": {
                    "$ref": "#/definitions/controllers.DeviceDefinition"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "vin": {
                    "type": "string"
                }
            }
        },
        "services.DeviceVehicleInfo": {
            "type": "object",
            "properties": {
                "base_msrp": {
                    "type": "integer"
                },
                "driven_wheels": {
                    "type": "string"
                },
                "epa_class": {
                    "type": "string"
                },
                "fuel_type": {
                    "type": "string"
                },
                "mpg_city": {
                    "type": "string"
                },
                "mpg_highway": {
                    "type": "string"
                },
                "number_of_doors": {
                    "type": "string"
                },
                "vehicle_type": {
                    "description": "VehicleType PASSENGER CAR, from NHTSA",
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "2.0",
	Host:        "",
	BasePath:    "/v1",
	Schemes:     []string{},
	Title:       "DIMO Devices API",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
