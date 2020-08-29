<!--
 Copyright (c) 2019 William

 This software is released under the MIT License.
 https://opensource.org/licenses/MIT
-->

# Introduction

A golang examples collection with many classic tutorials from broad internet.

This repository is consisted of lots of branches but with empty master branch which only includes introductory README doc.

Generally, every branch is an independent example, examples with different stages are created as a new branch with the same example names suffixed by their stage names with slash separated. Please use `git branch` to list all and `git checkout` to navigate among all.

Current branches lists below(`git for-each-ref --sort=-committerdate refs/remotes/origin`):

```bash
86c27dc237d0b99619a51119c327de3eb4cad50e commit  refs/remotes/origin/kafka-tour
5a0b9011644e315a1167ec7c12b76a0a613d2785 commit  refs/remotes/origin/mongodb-go/tutorial-1.2
cf5a6e0a7ced8da63da055806c257d3b1e85e9e6 commit  refs/remotes/origin/mongodb-go/tutorial-1
ce5ef5bb5efb25ddc8ae413a7eda1d44425c6a0c commit  refs/remotes/origin/intheap
64e48cc7cdacbba68b0e4d9cb7bb1a7a4ea211b5 commit  refs/remotes/origin/http-tracing
427011a8f748601103664ff1db9a7e3a810eefab commit  refs/remotes/origin/grpc/consul-sd-integration
22968b80d641892fc992101eeebcd9a2432e232b commit  refs/remotes/origin/example-2
b0867e1137c01591122f67b86f2787ba139f3b9e commit  refs/remotes/origin/example-1
521941a4c99832f34c57bdc6b7478242b62ae12b commit  refs/remotes/origin/dsa
defb443f65dbcc96d58895cd146d865e91b4011e commit  refs/remotes/origin/basic
9aac2dc979188fd3c34ce23685f39ea190d71a91 commit  refs/remotes/origin/auth
edd613c74a01410808c1a29f525e4c4c7aa327ef commit  refs/remotes/origin/addsvc
3ce1ad5b207d2a9b3def17dfdefae4c0774c540c commit  refs/remotes/origin/master
3ce1ad5b207d2a9b3def17dfdefae4c0774c540c commit  refs/remotes/origin/HEAD
f462b3dfc9144085d342774e3dad489c1f7a63d3 commit  refs/remotes/origin/progress-bar
3b99a9035c37fc581519aa04c6a41b8dbdadda24 commit  refs/remotes/origin/multidomain-http-request
e9cfea93a5b20c5fa313680b9395e64248ad635d commit  refs/remotes/origin/cron
58edf3f53f9a634fcfdf9ae0155cb391b0f02793 commit  refs/remotes/origin/go-env-series-part1
73d13e84921384aa1744042232c56111e42df5e1 commit  refs/remotes/origin/tcp-server-shutdown
793d755e48d92fef5b1dff2ee226f4e7d68892a1 commit  refs/remotes/origin/faking-io
a1fb967d2e74e3f7df8837b4bf4ee1c0c58ad8b4 commit  refs/remotes/origin/cockroachdb
881bb04ff00348e26f2a9cf3eae7722a427b7a34 commit  refs/remotes/origin/makefile-tutorial
008fa359342b6e01ecc46c6358ad6ebcb4297007 commit  refs/remotes/origin/http-proto
39fd42314993eac69dd8477fd34171a72803c715 commit  refs/remotes/origin/helloworld
c3ce73205f0482b05e26d9511bdd05fa445d0298 commit  refs/remotes/origin/jsonapi
f5cf5b8734820048892bbdf6a680ea63943b8cc6 commit  refs/remotes/origin/cobra
d9577349c791102cc5a95f3da0010e0a703cbe6d commit  refs/remotes/origin/pubsub
e98564d54a49ab7a80ebd5b33a05ff78616f61f0 commit  refs/remotes/origin/http/multipart
e98564d54a49ab7a80ebd5b33a05ff78616f61f0 commit  refs/remotes/origin/http/body
c8c8be47c424bf92a25692826426815f3ef060e4 commit  refs/remotes/origin/grpc/name_resolving
b0aa0a78ff054d96c4d2729b8a528156e910881a commit  refs/remotes/origin/nats
11c6ff06edbb4ae036596ba3e413d27215acd81d commit  refs/remotes/origin/zap-examples
d5f60f4b00c2c24d5e493bfa2b967b4041ee4575 commit  refs/remotes/origin/go-ldflags
64c30111ace9fbe94b942b8e36f2cb724ecbc414 commit  refs/remotes/origin/go-plugin
c94c1ab510d87ca6224c6d5a17d82d235c383a8f commit  refs/remotes/origin/go-build-tags
f23bf6868a6122ce350bb46b3501530ff45d5a37 commit  refs/remotes/origin/zookeeper-tutorial
790e6466fed0ff15b26c2ae340f4d1546f224afd commit  refs/remotes/origin/opentracing-tutorial-1
51b4c419778bf3e102e39b019200cb1bd6bef891 commit  refs/remotes/origin/go-tour/http-tracing
9b48ead4d33918383c668dffa960af3f1894b1cb commit  refs/remotes/origin/kafka-go
3e5185138b40228da5278e75d09069d8e71f16cb commit  refs/remotes/origin/sort_transform
3e5185138b40228da5278e75d09069d8e71f16cb commit  refs/remotes/origin/algorithm/sort_transform
6762034f0b8a2576557edb10b43cbd473f10f9fe commit  refs/remotes/origin/sync/errgroup
6762034f0b8a2576557edb10b43cbd473f10f9fe commit  refs/remotes/origin/errgroup
751dc0b6554483279d98d60f6be0ccd4aa00575c commit  refs/remotes/origin/rest-unit-testing
35b3c674fb17fa0530d1c1b6beaa9b24842b69a0 commit  refs/remotes/origin/shijuvar/nats-tutorial
35b3c674fb17fa0530d1c1b6beaa9b24842b69a0 commit  refs/remotes/origin/nats-tutorial
6c19a6548892050df473f5e137dee221d6e320a2 commit  refs/remotes/origin/stdsql/opendbservice
6c19a6548892050df473f5e137dee221d6e320a2 commit  refs/remotes/origin/opendbservice
0f13b96a521a616686cac928876a3fe6734fbf8e commit  refs/remotes/origin/stdsql/opendbcli
0f13b96a521a616686cac928876a3fe6734fbf8e commit  refs/remotes/origin/opendbcli
1f51d45ac2d0b6ce468cf80d47bb507cab4b126a commit  refs/remotes/origin/go-kit/nats
90ff4f6e9a0d1bc08f45dd1c7daaf9edf3d99a81 commit  refs/remotes/origin/routeguide
90ff4f6e9a0d1bc08f45dd1c7daaf9edf3d99a81 commit  refs/remotes/origin/grpc/routeguide
42f305a92ae3e2525f45638d89447fd38376d834 commit  refs/remotes/origin/go-kit/addsvc
78f78514ac2ddba04b36583ebd068b7242c20219 commit  refs/remotes/origin/net/rpc/basic
0672152c4657fb285bf9ce37137a37a4b04e7576 commit  refs/remotes/origin/tree-digesting-review
0672152c4657fb285bf9ce37137a37a4b04e7576 commit  refs/remotes/origin/pipeline/tree-digesting-review
9aa8063dbdde94b855df7346b4afdbd0de72a858 commit  refs/remotes/origin/review
9aa8063dbdde94b855df7346b4afdbd0de72a858 commit  refs/remotes/origin/pipelines/review
1224979d7854fb2164dc2a8a5565624d1574dd20 commit  refs/remotes/origin/container/ring/basic
f0eaf64accd35804b2b110b8fb3ef972273e09b6 commit  refs/remotes/origin/container/list/basic
be4bf21d8ad01ebcdf3b6bc66cf53afb88f2e4b7 commit  refs/remotes/origin/priority_queue_heap
be4bf21d8ad01ebcdf3b6bc66cf53afb88f2e4b7 commit  refs/remotes/origin/container/heap/priority_queue_heap
8d861c493cbb16e0084264b32d54f6802867cc49 commit  refs/remotes/origin/container/heap/intheap
3db843b0fa7a90513db3d488e06ffd47d32064db commit  refs/remotes/origin/sqlx/postgres
3db843b0fa7a90513db3d488e06ffd47d32064db commit  refs/remotes/origin/postgres
98913f7db3c66ab882ac401c76fc986bd1d3ce65 commit  refs/remotes/origin/grpc/helloworld
3166d0b92760aae71c8a6b259013843a784c6ea8 commit  refs/remotes/origin/simple-producer-consumer
3166d0b92760aae71c8a6b259013843a784c6ea8 commit  refs/remotes/origin/rabbitmq-go/simple-producer-consumer
a63b1352887a950e28bc7b7d5b08ec83bc1b8dd5 commit  refs/remotes/origin/rpc
a63b1352887a950e28bc7b7d5b08ec83bc1b8dd5 commit  refs/remotes/origin/rabbitmq-go/rpc
4b11c2397a973d1ea239e58c14d121e797fc40f5 commit  refs/remotes/origin/rabbitmq-go/pubsub
93c99374d1ee97b5d4b57035bb7c4f4832894b3f commit  refs/remotes/origin/topics
93c99374d1ee97b5d4b57035bb7c4f4832894b3f commit  refs/remotes/origin/rabbitmq-go/topics
a8c147d18bb7bd2cdab93f8c3c8c367e02310e34 commit  refs/remotes/origin/routing
a8c147d18bb7bd2cdab93f8c3c8c367e02310e34 commit  refs/remotes/origin/rabbitmq-go/routing
f0cbad7b2e8c0df53c2140173e630ce8a2e92412 commit  refs/remotes/origin/workqueues
f0cbad7b2e8c0df53c2140173e630ce8a2e92412 commit  refs/remotes/origin/rabbitmq-go/workqueues
05a28b92601a8cf16fd5d0b3ab56f744b9e9b1da commit  refs/remotes/origin/rabbitmq-go/helloworld
fe877a8e4d645df1197399ca89e1897edae5c1b1 commit  refs/remotes/origin/workqueue
fe877a8e4d645df1197399ca89e1897edae5c1b1 commit  refs/remotes/origin/rabbitmq-go/workqueue
c895733b9add631a6d250335a62c72a7ac5a8c90 commit  refs/remotes/origin/mqtt-go/example-2
ae5d763e1f0dfe007d07295190a834491a3dcf4a commit  refs/remotes/origin/mqtt-go/example-1
1c865ea9b9f5e71ab4580f5f77ec1a5e1ee262a2 commit  refs/remotes/origin/redis/example-2
0e96428d9af597119e2032ea5428af2ec3db4d44 commit  refs/remotes/origin/organising-database-access-3
0e96428d9af597119e2032ea5428af2ec3db4d44 commit  refs/remotes/origin/db/organising-database-access-3
cd44d0503542966401f620e5da5b9d47f586869a commit  refs/remotes/origin/organising-database-access-2
cd44d0503542966401f620e5da5b9d47f586869a commit  refs/remotes/origin/db/organising-database-access-2
79055e17210974e55946ab4036d3c71ed391a090 commit  refs/remotes/origin/organising-database-access-1
79055e17210974e55946ab4036d3c71ed391a090 commit  refs/remotes/origin/db/organising-database-access-1
a8dbda410766a3f15c19ae322b3eab70c73a6579 commit  refs/remotes/origin/jwt/auth
8a01d102645c42637a2da2cbd51c144b888bc68b commit  refs/remotes/origin/profilesvc
8a01d102645c42637a2da2cbd51c144b888bc68b commit  refs/remotes/origin/go-kit/profilesvc
9d4d6dbbedd7893b0e7d0c879ac560c46a40e94f commit  refs/remotes/origin/stringsvc-3
9d4d6dbbedd7893b0e7d0c879ac560c46a40e94f commit  refs/remotes/origin/go-kit/stringsvc-3
f87299eeb9c0451d050e7d0a23b04bfa8dd65010 commit  refs/remotes/origin/shippy/part-2
f87299eeb9c0451d050e7d0a23b04bfa8dd65010 commit  refs/remotes/origin/part-2
ce95e35fc26999ed8c8dd3b932581dac7a0bdff7 commit  refs/remotes/origin/shippy/part-1
ce95e35fc26999ed8c8dd3b932581dac7a0bdff7 commit  refs/remotes/origin/part-1
711366e85c4bd66d68ca182b48e0bcac20922a20 commit  refs/remotes/origin/service
711366e85c4bd66d68ca182b48e0bcac20922a20 commit  refs/remotes/origin/go-micro/service
af237108061a7e98a26dfbf99ff6e3bbf6716d07 commit  refs/remotes/origin/stringsvc-2
af237108061a7e98a26dfbf99ff6e3bbf6716d07 commit  refs/remotes/origin/go-kit/stringsvc-2
44fc6e50c0ca5717688df5f37cffbe41650813bb commit  refs/remotes/origin/stringsvc-1
44fc6e50c0ca5717688df5f37cffbe41650813bb commit  refs/remotes/origin/go-kit/stringsvc-1
10312ce85b78fdf90809741df6ace602498f60d6 commit  refs/remotes/origin/gin/basic
6b6b06a48b06dd6092b9c1b65a5b57b14f7ff14d commit  refs/remotes/origin/httprouter/basic_auth
6b6b06a48b06dd6092b9c1b65a5b57b14f7ff14d commit  refs/remotes/origin/basic_auth
c9757e94e85f04b7140ffa924f3b7aca7c6d16b8 commit  refs/remotes/origin/httprouter/basic
c971d849bb0ad54ab3648c4d534d815336d52133 commit  refs/remotes/origin/simple
c971d849bb0ad54ab3648c4d534d815336d52133 commit  refs/remotes/origin/jwt/simple
3dbdc6cf419c44e6ab1973daf0ac025cf240ed35 commit  refs/remotes/origin/graphql/gqlgen
3dbdc6cf419c44e6ab1973daf0ac025cf240ed35 commit  refs/remotes/origin/gqlgen
af740a6c60c69465836a1f537d5f544f20fa4300 commit  refs/remotes/origin/tunneling
af740a6c60c69465836a1f537d5f544f20fa4300 commit  refs/remotes/origin/ssh/tunneling
328089181834c9b5132104fc21aade2e966819fa commit  refs/remotes/origin/ssh/client-connection
328089181834c9b5132104fc21aade2e966819fa commit  refs/remotes/origin/client-connection
f8d99e80fb3302d2694fa09a18fc835f8aeb0e71 commit  refs/remotes/origin/tree-digest-parallel
f8d99e80fb3302d2694fa09a18fc835f8aeb0e71 commit  refs/remotes/origin/piplines/tree-digest-parallel
413930b9317e077445e3cc6a549a86191827093e commit  refs/remotes/origin/share_memory_by_communication
413930b9317e077445e3cc6a549a86191827093e commit  refs/remotes/origin/channel/share_memory_by_communication
9f548997e67d571be562339f212275845fba503a commit  refs/remotes/origin/multiple-contexts
9f548997e67d571be562339f212275845fba503a commit  refs/remotes/origin/context/multiple-contexts
1d2fd5e58792122f3575c478d670a979152d4752 commit  refs/remotes/origin/docker/basic-deploy
1d2fd5e58792122f3575c478d670a979152d4752 commit  refs/remotes/origin/basic-deploy
85f1cc1101175ad7b56ef869c30068074703fb37 commit  refs/remotes/origin/redis/example-1
30eaa9cdd998e7d28aea5610969e3673204163d7 commit  refs/remotes/origin/tutorrial-1
30eaa9cdd998e7d28aea5610969e3673204163d7 commit  refs/remotes/origin/protocol-buffers/tutorrial-1
d3b98bc606c342b23313356c1bbfd61b9e53a480 commit  refs/remotes/origin/mongo-rest
30a6f3dd41395d59a3f4f1de03f6c01af1d0b0a2 commit  refs/remotes/origin/https
30a6f3dd41395d59a3f4f1de03f6c01af1d0b0a2 commit  refs/remotes/origin/docker/https
c09faec8b11aceb7bd9c86eb496908c50b058d4e commit  refs/remotes/origin/websocket/chat
c09faec8b11aceb7bd9c86eb496908c50b058d4e commit  refs/remotes/origin/chat
f08a4eca1c8973f4175dd848a3def3f620e3d6c7 commit  refs/remotes/origin/appengine/helloworld
7231b9ed3eb0de3cb4cd818d1bb1f86088e361c9 commit  refs/remotes/origin/parser-inspect
7231b9ed3eb0de3cb4cd818d1bb1f86088e361c9 commit  refs/remotes/origin/go-compiler/parser-inspect
32feb8acde2492e469557ac5dc9e7a92d60ab354 commit  refs/remotes/origin/parser
32feb8acde2492e469557ac5dc9e7a92d60ab354 commit  refs/remotes/origin/go-compiler/parser
b2e7702a2220e41288b80284d8d17c3cef852c57 commit  refs/remotes/origin/scanner
b2e7702a2220e41288b80284d8d17c3cef852c57 commit  refs/remotes/origin/go-compiler/scanner
32f308a26601db472da63f18b81180bb7410ed67 commit  refs/remotes/origin/multipart-response
36c8ff3c512f4db9a94af56f714978a63d6d4a75 commit  refs/remotes/origin/multiple-files-upload-server
fe596ae802efb1bed8bb500185dc51a19e443f45 commit  refs/remotes/origin/web-search
fe596ae802efb1bed8bb500185dc51a19e443f45 commit  refs/remotes/origin/context/web-search
8c52bc5ee0bcaa87a1b87a3c6c47c9e261c060e7 commit  refs/remotes/origin/tree-digest-bounded
8c52bc5ee0bcaa87a1b87a3c6c47c9e261c060e7 commit  refs/remotes/origin/piplines/tree-digest-bounded
ad86b72c387ed350e7d9f655b65ff7f8ab49ceb9 commit  refs/remotes/origin/tree-digest-serial
ad86b72c387ed350e7d9f655b65ff7f8ab49ceb9 commit  refs/remotes/origin/piplines/tree-digest-serial
1892d361dd9f56abae75feb22d8cceb259d76aff commit  refs/remotes/origin/piplines/explicit-cancellation-2
1892d361dd9f56abae75feb22d8cceb259d76aff commit  refs/remotes/origin/explicit-cancellation-2
d26dae380561ac69746d9771f2ef6bd7dc81a0f1 commit  refs/remotes/origin/piplines/explicit-cancellation-1
d26dae380561ac69746d9771f2ef6bd7dc81a0f1 commit  refs/remotes/origin/explicit-cancellation-1
611708ffd65348569f2c62ad475e40deaf4e3ada commit  refs/remotes/origin/piplines/buffered-channel
611708ffd65348569f2c62ad475e40deaf4e3ada commit  refs/remotes/origin/buffered-channel
e672d2ab600c0b59e1f57c43c60ebfd77120fff1 commit  refs/remotes/origin/piplines/fan_out-fan_in
e672d2ab600c0b59e1f57c43c60ebfd77120fff1 commit  refs/remotes/origin/fan_out-fan_in
8ea32d7e2da1f19b2fb1c9eb7a77798cce8ad6b9 commit  refs/remotes/origin/squaring-numbers
8ea32d7e2da1f19b2fb1c9eb7a77798cce8ad6b9 commit  refs/remotes/origin/piplines/squaring-numbers
```

## Credits

All credits are to the original authors mentioned at comments with references at the top of every example.

## License

Under MIT license.
