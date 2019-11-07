module github.com/peterq/pan-light/server

go 1.13

// replace github.com/peterq/pan-light => ../

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/cache v6.4.0+incompatible
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/iris-contrib/middleware/jwt v0.0.0-20191028172159-41f72a73786a
	github.com/kataras/iris/v12 v12.0.1
	github.com/pkg/errors v0.8.1
	golang.org/x/net v0.0.0-20191105084925-a882066a44e0
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)
