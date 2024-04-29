module github.com/reviewpad/reviewpad/v4

go 1.20

require (
	github.com/aws/aws-secretsmanager-caching-go v1.1.2
	github.com/bmatcuk/doublestar/v4 v4.6.0
	github.com/bradleyfalzon/ghinstallation/v2 v2.7.0
	github.com/dukex/mixpanel v1.0.1
	github.com/golang/mock v1.6.0
	github.com/google/go-github/v52 v52.0.0
	github.com/google/uuid v1.3.1
	github.com/gorilla/mux v1.8.0
	github.com/graphql-go/graphql v0.8.1
	github.com/hasura/go-graphql-client v0.10.0
	github.com/jarcoal/httpmock v1.3.1
	github.com/jinzhu/copier v0.4.0
	github.com/libgit2/git2go/v31 v31.7.9
	github.com/mattn/go-shellwords v1.0.12
	github.com/migueleliasweb/go-github-mock v0.0.20
	github.com/mitchellh/mapstructure v1.5.0
	github.com/ohler55/ojg v1.19.4
	github.com/reviewpad/api/go v0.0.0-20231016095446-177fe259ce20
	github.com/reviewpad/go-conventionalcommits v0.10.0
	github.com/reviewpad/go-lib v0.0.0-20231016101521-daf31e8bf28f
	github.com/shurcooL/githubv4 be2daab69064
	github.com/shurcooL/graphql v0.0.0-20230722043721-ed46e5a46466
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.7.0
	github.com/stretchr/testify v1.8.4
	github.com/tomnomnom/linkheader v0.0.0-20180905144013-02ca5825eb80
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d
	golang.org/x/oauth2 v0.13.0
	google.golang.org/grpc v1.58.3
	google.golang.org/protobuf v1.31.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/ProtonMail/go-crypto v0.0.0-20230923063757-afb1ddc0824c // indirect
	github.com/aws/aws-sdk-go v1.45.25 // indirect
	github.com/cloudflare/circl v1.3.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-errors/errors v1.5.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-github/v55 v55.0.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.17.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231012201019-e917dd12ba7a // indirect
	nhooyr.io/websocket v1.8.7 // indirect
)

replace github.com/gin-gonic/gin v1.6.3 => github.com/gin-gonic/gin v1.7.7

exclude gopkg.in/yaml.v2 v2.2.2
