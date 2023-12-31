// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Jessica Tarra"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/movies": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Fetch a list of movies with server-side pagination",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Movies"
                ],
                "summary": "List movies with pagination",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Movie title",
                        "name": "title",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "Movie genres",
                        "name": "genres",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of movies per page",
                        "name": "page_size",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Sort order",
                        "name": "sort",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Movie list",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/database.Movie"
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create a new movie",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Movies"
                ],
                "summary": "Create a movie",
                "parameters": [
                    {
                        "description": "Request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.createMovieRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Movie created",
                        "schema": {
                            "$ref": "#/definitions/database.Movie"
                        }
                    }
                }
            }
        },
        "/movies/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Retrieve a movie by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Movies"
                ],
                "summary": "Get a movie by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Movie ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Movie details",
                        "schema": {
                            "$ref": "#/definitions/database.Movie"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update an existing movie",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Movies"
                ],
                "summary": "Update a movie by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Movie ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.updateMovieRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Movie updated",
                        "schema": {
                            "$ref": "#/definitions/database.Movie"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Delete a movie by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Movies"
                ],
                "summary": "Delete a movie by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Movie ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/tokens/authentication": {
            "post": {
                "description": "Creates an authentication token for a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Create authentication token",
                "parameters": [
                    {
                        "description": "Request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.CreateAuthTokenRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Authentication token",
                        "schema": {
                            "$ref": "#/definitions/domain.Token"
                        }
                    }
                }
            }
        },
        "/users": {
            "post": {
                "description": "Registers a new user.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Register User",
                "parameters": [
                    {
                        "description": "User registration data",
                        "name": "name",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.CreateUserRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/domain.User"
                        }
                    }
                }
            }
        },
        "/users/activated": {
            "put": {
                "description": "Activates a user account using a token that was previously sent when successfully register a new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Activate User",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Token for user activation",
                        "name": "token",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.User"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "database.Movie": {
            "type": "object",
            "properties": {
                "genres": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "runtime": {
                    "type": "string",
                    "example": "0"
                },
                "title": {
                    "type": "string"
                },
                "version": {
                    "type": "integer"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "domain.CreateAuthTokenRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "domain.CreateUserRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "domain.Token": {
            "type": "object",
            "properties": {
                "expiry": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "domain.User": {
            "type": "object",
            "properties": {
                "activated": {
                    "type": "boolean"
                },
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "main.createMovieRequest": {
            "type": "object",
            "properties": {
                "genres": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "runtime": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "main.updateMovieRequest": {
            "type": "object",
            "properties": {
                "genres": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "runtime": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                },
                "year": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "Greenlight API Docs",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
