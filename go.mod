module github.com/CanalTP/gormungandr

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78
	github.com/CanalTP/gonavitia v0.0.0-20180817052458-0dcb887a472b
	github.com/CanalTP/gormungandr/kraken v0.1.0
	github.com/Microsoft/go-winio v0.4.7
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5
	github.com/beorn7/perks v1.0.0
	github.com/cenkalti/backoff v2.0.0+incompatible
	github.com/containerd/continuity v0.0.0-20180322171221-3e8f2ea4b190
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/go-connections v0.3.0
	github.com/docker/go-units v0.3.3
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gchaincl/sqlhooks v1.1.0
	github.com/gin-contrib/cors v0.0.0-20170318125340-cf4846e6a636
	github.com/gin-contrib/location v0.0.0-20180827025200-b7d60da6dc7c
	github.com/gin-contrib/sse v0.0.0-20170109093832-22d885f9ecc7
	github.com/gin-gonic/contrib v0.0.0-20180320084256-9b830a15f6ab
	github.com/gin-gonic/gin v0.0.0-20180329063307-6d913fc343cf
	github.com/golang/protobuf v1.3.2
	github.com/hashicorp/hcl v0.0.0-20180404174102-ef8a98b0bbce
	github.com/json-iterator/go v1.1.6
	github.com/lib/pq v0.0.0-20180327071824-d34b9ff171c2
	github.com/magiconair/properties v1.7.6
	github.com/mattn/go-isatty v0.0.3
	github.com/matttproud/golang_protobuf_extensions v1.0.1
	github.com/mitchellh/mapstructure v0.0.0-20180220230111-00c29f56e238
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
	github.com/modern-go/reflect2 v1.0.1
	github.com/newrelic/go-agent v1.11.0
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1
	github.com/opencontainers/runc v1.0.0-rc5
	github.com/ory/dockertest v0.0.0-20180810084858-272abdf5bd73
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pebbe/zmq4 v1.0.0
	github.com/pelletier/go-toml v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/prometheus/client_golang v1.0.0
	github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90
	github.com/prometheus/common v0.4.1
	github.com/prometheus/procfs v0.0.2
	github.com/rafaeljesus/rabbus v0.0.0-20180420204416-9b66eef60b25
	github.com/rafaeljesus/retry-go v0.0.0-20171214204623-5981a380a879
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/sony/gobreaker v0.4.1
	github.com/spf13/afero v1.1.0
	github.com/spf13/cast v1.2.0
	github.com/spf13/jwalterweatherman v0.0.0-20180109140146-7c0cea34c8ec
	github.com/spf13/pflag v0.0.0-20180403115518-1ce0cc6db402
	github.com/spf13/viper v1.0.2
	github.com/streadway/amqp v0.0.0-20180315184602-8e4aba63da9f
	github.com/stretchr/objx v0.1.1
	github.com/stretchr/testify v1.3.0
	github.com/ugorji/go v0.0.0-20180112141927-9831f2c3ac10
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/sys v0.0.0-20190422165155-953cdadca894
	golang.org/x/text v0.3.0
	golang.org/x/tools v0.0.0-20190628222527-fb37f6ba8261 // indirect
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2
	gopkg.in/yaml.v2 v2.2.1
)

replace github.com/CanalTP/gormungandr/kraken => ./kraken
