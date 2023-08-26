// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "GhostikGH",
            "url": "https://t.me/GhostikGH",
            "email": "feodor200@mail.ru"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/log/{user_id}": {
            "get": {
                "description": "Лог пользователя за все время",
                "produces": [
                    "application/octet-stream"
                ],
                "tags": [
                    "Log"
                ],
                "summary": "Лог пользователя",
                "operationId": "user-log",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "user id",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user_log.Response"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/user_log.Response"
                        }
                    }
                }
            }
        },
        "/segment": {
            "get": {
                "description": "Получения списка всех сегментов",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Segment"
                ],
                "summary": "Получения списка всех сегментов",
                "operationId": "segment-list",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/segment_list.Response"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/segment_list.Response"
                        }
                    }
                }
            }
        },
        "/segment/{segment}": {
            "get": {
                "description": "Получение всех пользователей в данном сегменте",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Segment"
                ],
                "summary": "Получение всех пользователей в данном сегменте",
                "operationId": "segment-user-list",
                "parameters": [
                    {
                        "type": "string",
                        "description": "segment name",
                        "name": "segment",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/segment_users.Response"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/segment_users.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "Удаление сегмента",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Segment"
                ],
                "summary": "Удаление сегмента",
                "operationId": "segment-deletion",
                "parameters": [
                    {
                        "type": "string",
                        "description": "segment name",
                        "name": "segment",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/delete_segment.Response"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/delete_segment.Response"
                        }
                    }
                }
            }
        },
        "/segment/{segment}/{percentage}": {
            "post": {
                "description": "Создание сегмента",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Segment"
                ],
                "summary": "Создание сегмента",
                "operationId": "segment-creation",
                "parameters": [
                    {
                        "type": "string",
                        "description": "segment name",
                        "name": "segment",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "percentage",
                        "name": "percentage",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/create_segment.Response"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/create_segment.Response"
                        }
                    }
                }
            }
        },
        "/user": {
            "get": {
                "description": "Получения списка всех пользователей",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Получения списка всех пользователей",
                "operationId": "user-list",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user_list.Response"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/user_list.Response"
                        }
                    }
                }
            }
        },
        "/user/{user_id}": {
            "get": {
                "description": "Получение всех сегментов данного пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Получение всех сегментов данного пользователя",
                "operationId": "user-segment-list",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "user id",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user_segments.Response"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/user_segments.Response"
                        }
                    }
                }
            },
            "post": {
                "description": "Изменяет состояние сегментов у пользователя, если пользователя нет, то он созадется",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Изменение сегментов у одного пользователя",
                "operationId": "create-user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "user id",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "user update data",
                        "name": "input",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/create_user.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/create_user.Response"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/create_user.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "Удаление пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Удаление пользователя",
                "operationId": "user-deletion",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "user id",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/delete_user.Response"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/delete_user.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "create_segment.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "create_user.Request": {
            "type": "object",
            "properties": {
                "addedSeg": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "removeSeg": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "create_user.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "delete_segment.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "delete_user.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "segment_list.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "segmentList": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "segment_users.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "userList": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "user_list.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "userList": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "user_log.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "user_segments.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "segments": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "status": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8002",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Avito Intership (Backend)",
	Description:      "Segmentation for A_B tests",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
