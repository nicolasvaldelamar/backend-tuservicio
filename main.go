package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_"github.com/go-playground/validator/v10/translations/id"
	_ "github.com/go-sql-driver/mysql"
	_ "golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/appservices")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	// Configuración de CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8100"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	router.Use(cors.New(config))
	// Configurar los encabezados CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8100")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	router.POST("/register", register)
	router.POST("/login", login)
	router.GET("/searchuser", searchUser)
	router.PUT("/updatephoto", updatePhoto)

	router.Run(":8080")
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Name    string `json:"name"`
	Lastname string `json:"lastname"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SearchUser struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Id        int    `json:"id"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type UserCompleto struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Id       int    `json:"id"`
	Descripcion string `json:"descripcion"`
	Adress  string `json:"adress"`
	City string `json:"city"`
	Name    string `json:"name"`
	Lastname string `json:"lastname"`
	Profile_picture string `json:"profile_picture"`
	Ocupacion string `json:"ocupacion"`
}


func register(c *gin.Context) {
	var user User

	// Verificar si el cuerpo de la solicitud contiene un objeto JSON válido
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Acceder a los campos del objeto User
	username := user.Username
	password := user.Password
	email := user.Email
	role := user.Role
	name := user.Name
	lastname := user.Lastname

	log.Println("DATOS", username, password, email, role, name, lastname)
	// Realizar alguna lógica adicional con los datos del objeto User
	// Verificar si el usuario ya existe en la base de datos
	var count string = ""
	var validEmail string = ""
	err := db.QueryRow("SELECT username FROM users WHERE username=? OR email=?", username, email).Scan(&count)

	if err != nil {
		println(err)
	}

	err = db.QueryRow("SELECT email FROM users WHERE email=?", email).Scan(&validEmail)

	if err != nil {
		println(err)
	}

	if count != "" || validEmail != "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "El usuario y/o email ya existe"})
		return
	}

	// Insertar nuevo usuario en la base de datos
	_, err = db.Query("INSERT INTO `appservices`.`users` (`username`, `password`, `email`, `role`, nombre, apellido) VALUES (?, ?, ?, ?, ?, ?);", username, password, email, role, name, lastname)
	if err != nil {
		println(err)
	}

	// Devolver una respuesta
	if username != "" || password != "" || email != "" || role != "" {
		c.JSON(http.StatusOK, gin.H{
			"message":  "Registro exitoso",
			"username": username,
			"email":    email,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error al registrar el usuario"})
	}

}

func login(c *gin.Context) {
	// Obtener datos del formulario de inicio de sesión
	var user Login

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//Acceder a los campos del objeto User
	username := user.Username
	password := user.Password
	log.Println("DATOS", username, password)
	log.Println("DATOS", user)
	rows, err := db.Query("SELECT username, password, id FROM users WHERE username=? AND password=?", username, password)
    if err != nil {
        panic(err.Error())
    }
    defer rows.Close()
	var usernameop, passwordop string
	var id int
    for rows.Next() {
        
        if err := rows.Scan(&usernameop, &passwordop, &id); err != nil {
            panic(err.Error())
        }
        fmt.Println(usernameop, passwordop)
    }
	if usernameop != username || passwordop != password {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error al iniciar sesión"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Inicio de sesión exitoso", "user": user, "id": id})
	}

}

// Endpoint para buscar un usuario por nombre de usuario y contraseña
func searchUser(c *gin.Context) {
	// Obtener los parámetros de búsqueda del formulario
	id := c.Query("id")
	println(id)
	// Realizar la búsqueda en la base de datos
	
	rows, err := db.Query("SELECT * FROM users WHERE id=?", id)
    if err != nil {
        println(err)
    }
    defer rows.Close()
	var nameop, lastnameop, usernameop, passwordop, emailop, roleop, profile_picturesop, cover_photoop, descripcionop, adressop, city, ocupacionop string
	var idop int
    for rows.Next() {
        
        if err := rows.Scan(&idop, &usernameop, &passwordop, &emailop, &roleop,  &nameop, &lastnameop,&profile_picturesop, &cover_photoop, &descripcionop, &adressop, &city, &ocupacionop);
		err != nil {
            println(err)
        }
        fmt.Println("Hola imprimiendo" + nameop, lastnameop, usernameop, passwordop, emailop, roleop, idop, nameop, lastnameop, profile_picturesop, ocupacionop)
    }
	datos := UserCompleto{usernameop, passwordop, emailop, roleop, idop, descripcionop, adressop, city, nameop, lastnameop, profile_picturesop, ocupacionop}

	if usernameop == "" || passwordop == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error usuario no encontrado"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"results": datos})
	}
}

func updatePhoto(c *gin.Context) {
	// Obtener los parámetros de búsqueda del formulario
	username := c.Query("username")
	password := c.Query("password")
	profile_pictures := c.Query("newimage")
	println("DATOS", username, password, profile_pictures)
	// Realizar la búsqueda en la base de datos
	_, err := db.Query("UPDATE users SET profile_picture=? WHERE username=? AND password=?", profile_pictures, username, password)
	if err != nil {
		println(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Actualización exitosa"})
}

func updateInformacion(c *gin.Context){
	username := c.Query("username")
	//terminar funcion para actualizar informacion en la base de datos
	password := c.Query("password")
	id := c.Query("id")
	name := c.Query("name")
	lastname := c.Query("lastname")
	descripcion := c.Query("descripcion")
	dereccion := c.Query("direccion")
	ciudad := c.Query("ciudad")
	ocupacion := c.Query("ocupacion")


	_, err := db.Query("UPDATE users SET name=?, lastname=?, descripcion=?, adress=?, city=?, ocupacion=? WHERE username=? AND password=? AND id=?", name, lastname, descripcion, dereccion, ciudad, ocupacion, username, password, id)
	if err != nil {
		println(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Actualización exitosa"})


}