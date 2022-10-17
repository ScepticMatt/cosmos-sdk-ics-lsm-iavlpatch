module cosmossdk.io/core

go 1.19

require (
	cosmossdk.io/api v0.2.1
	cosmossdk.io/depinject v1.0.0-alpha.3
	cosmossdk.io/math v1.0.0-beta.3
	github.com/cosmos/cosmos-proto v1.0.0-alpha8
	github.com/stretchr/testify v1.8.0
	google.golang.org/protobuf v1.28.1
	gotest.tools/v3 v3.4.0
	sigs.k8s.io/yaml v1.3.0
)

require (
	github.com/cosmos/gogoproto v1.4.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20220827204233-334a2380cb91 // indirect
	golang.org/x/net v0.0.0-20221004154528-8021a29435af // indirect
	golang.org/x/sys v0.0.0-20220928140112-f11e5e49a4ec // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220930163606-c98284e70a91 // indirect
	google.golang.org/grpc v1.50.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// temporary until we tag a new go module
replace cosmossdk.io/math => ../math
