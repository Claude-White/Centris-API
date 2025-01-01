// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/brokers": {
            "get": {
                "description": "Retrieves a list of all brokers",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Brokers"
                ],
                "summary": "Get all brokers",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/server.Broker"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/properties": {
            "post": {
                "description": "Retrieves a list of properties with pagination",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Properties"
                ],
                "summary": "Get all properties",
                "parameters": [
                    {
                        "description": "Pagination parameters",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.RequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/repository.Property"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/properties/agency/{agencyName}": {
            "post": {
                "description": "Retrieves a list of properties for a specific agency with pagination",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Properties"
                ],
                "summary": "Get all properties by agency",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Agency Name",
                        "name": "agencyName",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Pagination parameters",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.RequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/repository.Property"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid agency name or request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Agency properties not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/properties/broker/{brokerId}": {
            "post": {
                "description": "Retrieves a list of properties for a specific broker with pagination",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Properties"
                ],
                "summary": "Get all properties by broker",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Broker ID",
                        "name": "brokerId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Pagination parameters",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.RequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/repository.Property"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid broker ID or request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Broker properties not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/properties/category/{categoryName}": {
            "post": {
                "description": "Retrieves a list of properties for a specific category with pagination",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Properties"
                ],
                "summary": "Get all properties by category",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Category Name",
                        "name": "categoryName",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Pagination parameters",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.RequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/repository.Property"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid category name or request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Category properties not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/properties/coordinates": {
            "post": {
                "description": "Retrieves a property based on latitude and longitude",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Properties"
                ],
                "summary": "Get a property by coordinates",
                "parameters": [
                    {
                        "description": "Coordinates parameters",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/server.Coordinates"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/repository.Property"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Property not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/properties/{mls}": {
            "get": {
                "description": "Retrieves property details using the provided MLS number",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Properties"
                ],
                "summary": "Get property by MLS number",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "MLS number",
                        "name": "mls",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/repository.Property"
                        }
                    },
                    "400": {
                        "description": "Invalid MLS number",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Property not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to get property",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "repository.Property": {
            "type": "object",
            "properties": {
                "apartment_number": {
                    "type": "string"
                },
                "bathroom_number": {
                    "type": "integer"
                },
                "bedroom_number": {
                    "type": "integer"
                },
                "category": {
                    "type": "string"
                },
                "city_name": {
                    "type": "string"
                },
                "civic_number": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "description": "MLS number",
                    "type": "integer"
                },
                "latitude": {
                    "type": "string"
                },
                "longitude": {
                    "type": "string"
                },
                "neighbourhood_name": {
                    "type": "string"
                },
                "price": {
                    "type": "string"
                },
                "room_number": {
                    "type": "integer"
                },
                "street_name": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "server.Broker": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "profilePhoto": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "server.Coordinates": {
            "type": "object",
            "properties": {
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                }
            }
        },
        "server.RequestBody": {
            "type": "object",
            "properties": {
                "numberOfItems": {
                    "type": "integer"
                },
                "startPosition": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Swagger Centris API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
