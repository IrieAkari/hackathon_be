package main

import (
	"log"
	"net/http"
	"os"
	_ "os/signal"
	_ "syscall"

	likeHandlers "hackathon/handlers/like"
	postHandlers "hackathon/handlers/post"
	replyHandlers "hackathon/handlers/reply"
	userHandlers "hackathon/handlers/user"
	"hackathon/utils"

	"github.com/rs/cors"
)

func main() {
	utils.InitDB()

	mux := http.NewServeMux()
	mux.HandleFunc("/username", userHandlers.UserNameGetHandler)      // GET /username?name=...
	mux.HandleFunc("/useremail", userHandlers.UserEmailGetHandler)    // GET /useremail?email=...
	mux.HandleFunc("/userdelete", userHandlers.UserDeleteHandler)     // GET /userdelete?email=...
	mux.HandleFunc("/userregister", userHandlers.UserRegisterHandler) // POST /userregister (name, email, password)
	mux.HandleFunc("/users", userHandlers.UsersGetHandler)            // GET /users
	mux.HandleFunc("/postcreate", postHandlers.PostCreateHandler)     // POST /postcreate (email, content)
	mux.HandleFunc("/postdelete", postHandlers.PostDeleteHandler)     // GET /postdelete?postid=...
	mux.HandleFunc("/posts", postHandlers.PostsGetHandler)            // GET /posts
	mux.HandleFunc("/postget", postHandlers.PostGetHandler)           // GET /postget?postid=...
	mux.HandleFunc("/replycreate", replyHandlers.ReplyCreateHandler)  // POST /replycreate (email, content, parent_id)
	mux.HandleFunc("/replydelete", replyHandlers.ReplyDeleteHandler)  // GET /replydelete?replyid=...
	mux.HandleFunc("/replys", replyHandlers.ReplysGetHandler)         // GET /replys?parentid=...
	mux.HandleFunc("/likecreate", likeHandlers.LikeCreateHandler)     // POST /likecreate (email, post_id)
	mux.HandleFunc("/likedelete", likeHandlers.LikeDeleteHandler)     // GET /likedelete?likeid=...
	mux.HandleFunc("/likeget", likeHandlers.LikeGetHandler)           // POST /likeget (user_id, post_id)

	handler := cors.Default().Handler(mux)

	utils.CloseDBWithSysCall()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("Server running on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
