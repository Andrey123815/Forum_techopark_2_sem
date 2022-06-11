package main

import (
	"db-forum/internal/app/forum"
	"db-forum/internal/app/forum/forumRepo"
	"db-forum/internal/app/post"
	"db-forum/internal/app/post/postRepo"
	"db-forum/internal/app/service"
	"db-forum/internal/app/service/serviceRepo"
	"db-forum/internal/app/thread"
	"db-forum/internal/app/thread/threadRepo"
	"db-forum/internal/app/user"
	"db-forum/internal/app/user/userRepo"
	"fmt"
	"github.com/fasthttp/router"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"log"
)

type ServerConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	Username string
	DbName   string
	Password string
	Host     string
	Port     string
}

type MainAppConfig struct {
	Server   ServerConfig
	Database DatabaseConfig
}

func getPostgres(config DatabaseConfig) *pgx.ConnPool {
	configString := fmt.Sprintf(`user=%s dbname=%s password=%s host=%s port=%s`,
		config.Username, config.DbName, config.Password, config.Host, config.Port)

	conn, err := pgx.ParseConnectionString(configString)
	if err != nil {
		log.Fatalln("config can't be parsed", err)
	}

	fmt.Println("Connection to Database...")
	fmt.Println("Configuration: ", configString)

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     conn,
		MaxConnections: 1000,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	})
	if err != nil {
		log.Fatalf("Error during connection to database: %s", err)
	}

	fmt.Println("Successful connection to database!")

	return pool
}

func main() {
	viper.SetConfigName("../config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	var appConfiguration MainAppConfig
	if errorDBConfParse := viper.Unmarshal(&appConfiguration); errorDBConfParse != nil {
		log.Fatalln(errorDBConfParse)
	}

	//fileSuccessLog, err := os.OpenFile("../info.log", os.O_RDWR|os.O_CREATE, 0666)
	//InfoLog := log.New(fileSuccessLog, "INFO\t", log.Ldate|log.Ltime)
	//fileErrorLog, err := os.OpenFile("../error.log", os.O_RDWR|os.O_CREATE, 0666)
	//ErrorLog := log.New(fileErrorLog, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db := getPostgres(appConfiguration.Database)

	r := router.New()
	user.SetUserRouting(r, &user.Handlers{
		UserRepo: userRepo.CreateUserRepository(db),
	})
	forum.SetForumRouting(r, &forum.Handlers{
		ForumRepo: forumRepo.CreateForumRepository(db),
		UserRepo:  userRepo.CreateUserRepository(db),
	})
	thread.SetThreadRouting(r, &thread.Handlers{
		ThreadRepo: threadRepo.CreateThreadRepository(db),
		UserRepo:   userRepo.CreateUserRepository(db),
	})
	post.SetPostRouting(r, &post.Handlers{
		PostRepo: postRepo.CreatePostRepository(db),
	})
	service.SetServiceRouting(r, &service.Handlers{
		ServiceRepo: serviceRepo.CreateServiceRepository(db),
	})

	fmt.Println("\nServer successfully started at localhost:" + appConfiguration.Server.Port + "!")

	log.Fatal(fasthttp.ListenAndServe(":"+appConfiguration.Server.Port, r.Handler))
}
