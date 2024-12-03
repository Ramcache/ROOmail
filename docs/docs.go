// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Поддержка API",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
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
        "/admin/file/upload": {
            "post": {
                "description": "Загрузка файла на сервер",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "файлы"
                ],
                "summary": "Загрузка файла",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Файл для загрузки",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"file_path\": \"uploaded/file/path\"}",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Ошибка разбора формы",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка чтения файла или сохранения файла",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/tasks/create": {
            "post": {
                "description": "Создает новую задачу с указанными данными",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "задачи"
                ],
                "summary": "Создание новой задачи",
                "parameters": [
                    {
                        "description": "Данные задачи",
                        "name": "task",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ROOmail_internal_models.Task"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Задача успешно создана",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Неверный JSON",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Неавторизован",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка создания задачи",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/tasks/delete/{id}": {
            "delete": {
                "description": "Удаляет задачу и все связанные с ней данные",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Задачи"
                ],
                "summary": "Удаление задачи",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID задачи",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Задача успешно удалена",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Некорректный идентификатор задачи",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Неавторизованный доступ",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Не удалось удалить задачу",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/tasks/update/{id}": {
            "put": {
                "description": "Обновление информации о задаче, такой как название, описание, срок выполнения, приоритет и список пользователей.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Задачи"
                ],
                "summary": "Обновить задачу",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор задачи",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Данные задачи для обновления",
                        "name": "task",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ROOmail_internal_models.Task"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"message\": \"Задача успешно обновлена\"}",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Неавторизованный доступ",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "patch": {
                "description": "Обновление одного или нескольких полей задачи по её идентификатору",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Задачи"
                ],
                "summary": "Частичное обновление задачи",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Идентификатор задачи",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Обновляемые поля задачи",
                        "name": "updates",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Задача успешно обновлена",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Некорректный идентификатор задачи или JSON",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Неавторизованный доступ",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/users_list": {
            "get": {
                "description": "Возвращает список пользователей с возможностью фильтрации по имени пользователя.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Получить список пользователей",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Фильтр по имени пользователя (поддерживает подстроку)",
                        "name": "username",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/ROOmail_internal_models.UsersList"
                            }
                        }
                    },
                    "401": {
                        "description": "Ошибка авторизации",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка получения пользователей",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/auth/login": {
            "post": {
                "description": "Аутентифицирует пользователя и возвращает jwt_token токен",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Вход пользователя",
                "parameters": [
                    {
                        "description": "Имя пользователя и пароль",
                        "name": "loginRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_handlers_auth.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешный вход",
                        "schema": {
                            "$ref": "#/definitions/internal_handlers_auth.LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Неверное имя пользователя или пароль",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка генерации токена",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "post": {
                "description": "Отзывает jwt_token токен пользователя",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Выход пользователя",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer \u003ctoken\u003e",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешный выход",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Требуется заголовок авторизации или он некорректен",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка отзыва токена",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "ROOmail_internal_models.Task": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "due_date": {
                    "type": "string"
                },
                "file_path": {
                    "description": "Путь к файлу, если есть",
                    "type": "string"
                },
                "priority": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "user_ids": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "ROOmail_internal_models.UsersList": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "internal_handlers_auth.LoginRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "internal_handlers_auth.LoginResponse": {
            "type": "object",
            "properties": {
                "role": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Документация ROOmail API",
	Description:      "Это документация API для проекта ROOmail.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
