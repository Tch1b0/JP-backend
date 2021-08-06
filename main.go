package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Tch1b0/JP-backend/pkg"
	"github.com/gorilla/mux"
)

var (
	posts []pkg.Post
	port = 5007
	account pkg.Account
)

func main() {
	posts = pkg.ReadJson("./posts/posts.json")
	var err error
	account, err = pkg.CreateFromJsonFile("./account.json")
	if err != nil {
		panic("You forgot to create the account.json file!")
	}

	router := mux.NewRouter()

	// Default route
	router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		fmt.Fprint(res, "Ok")
	})

	// Get 'post' info
	router.HandleFunc("/posts", getAllPosts)
	router.HandleFunc("/posts/count", getPostCount)
	router.HandleFunc("/posts/titles", getPostTitles)

	// Get certain 'post'
	router.HandleFunc("/post/{title}", getPostByName)
	router.HandleFunc("/post/index/{index}", getPostByIndex)
	router.HandleFunc("/post/{title}/logo", getPostLogo)
	router.HandleFunc("/post/{title}/banner", getPostBanner)
	
	// Upload new post
	router.HandleFunc("/post", createPost)

	// Validate Credentials
	router.HandleFunc("/verify", verify)

	// Start the Server
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

// ROUTES

// Respond with the variable 'posts'
func getAllPosts(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
	json.NewEncoder(res).Encode(posts)
}

// Respond with the count of the posts
func getPostCount(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
	json.NewEncoder(res).Encode(
		map[string]int{
			"count": len(posts),
		})
}

// Respond with a certain post, by its name
func getPostByName(res http.ResponseWriter, req *http.Request) {
	title := mux.Vars(req)["title"]
	for _, post := range posts {
		if post.Title == title {
			if req.Method == "DELETE" && pkg.IsOwnerFromReq(req, account) {
				deletePost(res, req, post)
				return
			}
			res.WriteHeader(200)
			json.NewEncoder(res).Encode(post)
			return
		}
	}
	
	res.WriteHeader(204)
	json.NewEncoder(res).Encode(
		map[string]string{
			"error": "There is no post with that name",
		})
}

// Get all post-titles
func getPostTitles(res http.ResponseWriter, req *http.Request) {
	var titles []string

	for _, post := range posts {
		titles = append(titles, post.Title)
	}

	res.WriteHeader(200)
	json.NewEncoder(res).Encode(titles)
}

func getPostLogo(res http.ResponseWriter, req *http.Request) {
	getImage(res, req, "Logo")
}
func getPostBanner(res http.ResponseWriter, req *http.Request) {
	getImage(res, req, "Banner")
}
func getImage(res http.ResponseWriter, req *http.Request, file string) {
	var title = mux.Vars(req)["title"]

	for _, post := range posts {
		if post.Title == title {
			var filetype string
			if file == "Logo" { 
				filetype = post.LogoType 
			} else { 
				filetype = post.BannerType 
			}

			pkg.SendImage(res, fmt.Sprintf("./posts/media/%s/%s.%s", title, file, filetype), filetype)
			return
		}
	}

	res.WriteHeader(404)
	json.NewEncoder(res).Encode(
		map[string]string{
			"error": "There is no post with that name",
		})
}

func getPostByIndex(res http.ResponseWriter, req *http.Request) {
	indexString := mux.Vars(req)["index"]
	index, _ := strconv.Atoi(indexString)

	if index > len(posts) - 1 || index < 0 {
		res.WriteHeader(404)
		json.NewEncoder(res).Encode(
			map[string]string{
				"error": "Index out of range",
			})
		return
	}

	if req.Method == "DELETE" && pkg.IsOwnerFromReq(req, account) {
		deletePost(res, req, posts[index])
		return
	}

	res.WriteHeader(200)
	json.NewEncoder(res).Encode(posts[index])
}

func createPost(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.WriteHeader(405)
		fmt.Fprint(res, "Only the POST method is allowed")
		return
	}

	req.ParseForm()
	req.ParseMultipartForm(0)

	// Authorize Request
	if len(req.FormValue("username")) <= 1 {
		res.WriteHeader(400)
		fmt.Fprint(res, "Credentials missing")
		return
	}
	if req.FormValue("username") != account.Username || !pkg.CheckPassword(req.FormValue("password"), account.Password) {
		res.WriteHeader(401)
		fmt.Fprint(res, "Either the username and/or password is wrong")
		return
	}
	title := req.FormValue("title")
	abs, _ := filepath.Abs(fmt.Sprintf("./posts/media/%s", title))
	err := os.Mkdir(abs, 0755)
	if err != nil {
		fmt.Println(err)
	}
	logoExt, err := pkg.UploadFile(req, "logo", fmt.Sprintf("./posts/media/%s", title), "Logo")
	if err != nil {fmt.Println(err)}
	bannerExt, err := pkg.UploadFile(req, "banner", fmt.Sprintf("./posts/media/%s", title), "Banner")
	if err != nil {fmt.Println(err)}

	post := pkg.Post{
		Title: req.FormValue("title"),
		LogoType: logoExt,
		BannerType: bannerExt,
		Description: req.FormValue("description"),
		LongDescription: req.FormValue("long-description"),
	}

	posts = append(posts, post)

	pkg.WriteJson("./posts/posts.json", posts)

	res.WriteHeader(200)
	json.NewEncoder(res).Encode(post)
}

func deletePost(res http.ResponseWriter, req *http.Request, p pkg.Post) {
	for i, post := range posts {
		if post.Title == p.Title {
			posts = append(posts[:i], posts[i+1:]...)
		}
	}
	pkg.DeleteMediaDir("./posts/media", p)
	pkg.WriteJson("./posts/posts.json", posts)
	res.WriteHeader(200)
	fmt.Fprint(res, posts)
}

func verify(res http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(0)
	if req.FormValue("password") == "" || req.FormValue("username") == "" {
		fmt.Fprint(res, "No password or username")
		return
	}

	if pkg.IsOwnerFromReq(req, account){
		fmt.Fprint(res, "Logged in")
		return
	} else {
		fmt.Fprint(res, "Wrong password or username")
		fmt.Println(account.Password)
		fmt.Println(req.FormValue("password"), req.FormValue("username"))
		return
	}
}