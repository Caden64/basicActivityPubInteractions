package main

import (
	"fmt"
	"maps"
	"net/http"
	"slices"

	"github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
	"github.com/gorilla/mux"
)

var users map[string]*activitypub.Person
var activities map[string]map[string]*activitypub.Activity

func main() {
	users = make(map[string]*activitypub.Person)
	activities = make(map[string]map[string]*activitypub.Activity)
	activities["caden"] = make(map[string]*activitypub.Activity)
	// ActivityPub data setup
	// followers
	users["caden"] = activitypub.PersonNew("caden")
	followers := activitypub.OrderedCollectionNew("followers")
	testFollower := activitypub.PersonNew("john")
	testFollow := activitypub.FollowNew("1", testFollower)
	activities["caden"][testFollow.ID.String()] = testFollow
	err := followers.Append(testFollow)
	if err != nil {
		panic(err)
	}
	users["caden"].Followers = followers
	// outbox
	outbox := activitypub.OrderedCollectionNew("outbox")
	err = outbox.Append(followers)
	if err != nil {
		panic(err)
	}
	users["caden"].Outbox = outbox
	// inbox
	inbox := activitypub.OrderedCollectionNew("inbox")
	err = inbox.Append(followers)
	if err != nil {
		panic(err)
	}
	users["caden"].Inbox = inbox
	// gorilla/mux boilerplate
	r := mux.NewRouter()
	u := r.PathPrefix("/user").Subrouter()
	u.HandleFunc("/{username}", handleUser)
	u.HandleFunc("/{username}/activity/{id:[0-9]+}", activityHandler)
	u.HandleFunc("/{username}/inbox", handleUserInbox)
	u.HandleFunc("/{username}/outbox", handleUserOutbox)
	r.Handle("/user/", u)
	http.ListenAndServe(":8080", r)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["username"]
	if slices.Contains(maps.Keys(users), user) {
		data, err := jsonld.Marshal(users[user])
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "ld+json")
		w.Write(data)
	} else {
		fmt.Println(user)
		w.Write([]byte("You have requested: " + user + " unfortunately this user does not exist"))
	}

}
func activityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["username"]
	id := vars["id"]
	if slices.Contains(maps.Keys(users), user) && slices.Contains(maps.Keys(activities[user]), id) {
		data, err := jsonld.Marshal(activities[user][id])
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "ld+json")
		w.Write(data)
	} else {
		fmt.Println(user)
		w.Write([]byte("You have requested: " + user + " unfortunately this user does not exist"))
	}

}

func handleUserInbox(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["username"]
	if slices.Contains(maps.Keys(users), user) {
		data, err := jsonld.Marshal(users[user].Inbox)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "ld+json")
		w.Write(data)
	} else {
		fmt.Println(user)
		w.Write([]byte("You have requested: " + user + " unfortunately this user does not exist"))
	}
}
func handleUserOutbox(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["username"]
	if slices.Contains(maps.Keys(users), user) {
		data, err := jsonld.Marshal(users[user].Outbox)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "ld+json")
		w.Write(data)
	} else {
		fmt.Println(user)
		w.Write([]byte("You have requested: " + user + " unfortunately this user does not exist"))
	}
}
