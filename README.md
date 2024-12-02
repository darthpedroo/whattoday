# Whattoday


Whattoday es una aplicación para subir posts sobre lo que hiciste en el día con la condición de que solo podes postear algo cada 24 horas. Esto fomenta que solo crees un post por día que esté enfocado en contar lo que hiciste.

Para crear un post tenes que tener un usuario pero si no iniciaste sesión podés ver los posteos del resto de la gente.

## Objetivo

El objetivo de este proyecto era aprender GO

## Demo
Una demo de la aplicación se encuentra [aca](!https://darthpedroo.github.io/whattoday-frontend/index.html)

### IMPORTANTE
El backend está hosteado en [render](!https://render.com/) y tardá 50 segundos en levantar el servidor para la primera request. 
Así que al visitar la página va a tardar un tiempo en cargar

 
## Features

- Endpoint para crear usuarios
- Endpoint para crear posteos
- Endpoiint para inicar sesión / registrarse
- Autorización con JWT ( Json Web Tokens )
- Implementación CORS para que el frontend pueda acceder a los endpoints.



## Instalación

### Clona el proyecto

#### Con ssh: 
```bash
git clone https://github.com/darthpedroo/whattoday.git 
```

#### Con https:
```bash
git clone git@github.com:darthpedroo/whattoday.git
```

### Instala las depencias

```bash
go mod tidy
```

### Configura las variables de entorno

Crea un archivo [.env](!https://towardsdatascience.com/use-environment-variable-in-your-next-golang-project-39e17c3aaa66) y configura las variables de entorno.

```bash
SECRET = "tu_clave_jwt" # Clave para registrar los tokens
PORT = "8080" #Puerto donde hostear el server
```

## Uso

```bash 
go run main.go
```
### Resultado Esperado:
```bash
Listening and serving HTTP on :6969
```

### Compilar el paquete 

#### Linux

```bash
go build -o whattoday-linux
```

#### Windows 
```bash
go build -o whattoday-windows.exe
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Tests aren't required but they are appreciated ! 

## License

[MIT](https://choosealicense.com/licenses/mit/)
