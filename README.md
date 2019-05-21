<!--
 Copyright (c) 2019 William

 This software is released under the MIT License.
 https://opensource.org/licenses/MIT
-->

# Introduction

A golang examples collection with many classic tutorials from broad internet.

This repository is consisted of lots of branches but with empty master branch which only includes introductory README doc.

Generally, every branch is an independent example, examples with different stages are created as a new branch with the same example names suffixed by their stage names with slash separated. Please use `git branch` to list all and `git checkout` to navigate among all.

Current branches lists below(`git for-each-ref --sort=-committerdate refs/heads/`):

```bash
90ff4f6e9a0d1bc08f45dd1c7daaf9edf3d99a81 commit refs/heads/grpc/routeguide
42f305a92ae3e2525f45638d89447fd38376d834 commit refs/heads/go-kit/addsvc
78f78514ac2ddba04b36583ebd068b7242c20219 commit refs/heads/net/rpc/basic
0672152c4657fb285bf9ce37137a37a4b04e7576 commit refs/heads/pipeline/tree-digesting-review
9aa8063dbdde94b855df7346b4afdbd0de72a858 commit refs/heads/pipelines/review
1224979d7854fb2164dc2a8a5565624d1574dd20 commit refs/heads/container/ring/basic
f0eaf64accd35804b2b110b8fb3ef972273e09b6 commit refs/heads/container/list/basic
be4bf21d8ad01ebcdf3b6bc66cf53afb88f2e4b7 commit refs/heads/container/heap/priority_queue_heap
299836667b33a3baccaa713bf50241ecebc924fe commit refs/heads/master
8d861c493cbb16e0084264b32d54f6802867cc49 commit refs/heads/container/heap/intheap
3db843b0fa7a90513db3d488e06ffd47d32064db commit refs/heads/sqlx/postgres
15fa2c71218603ab2535bddb44d27bac2b99de20 commit refs/heads/mongodb-go/tutorial-1
98913f7db3c66ab882ac401c76fc986bd1d3ce65 commit refs/heads/grpc/helloworld
3166d0b92760aae71c8a6b259013843a784c6ea8 commit refs/heads/rabbitmq-go/simple-producer-consumer
a63b1352887a950e28bc7b7d5b08ec83bc1b8dd5 commit refs/heads/rabbitmq-go/rpc
4b11c2397a973d1ea239e58c14d121e797fc40f5 commit refs/heads/rabbitmq-go/pubsub
93c99374d1ee97b5d4b57035bb7c4f4832894b3f commit refs/heads/rabbitmq-go/topics
a8c147d18bb7bd2cdab93f8c3c8c367e02310e34 commit refs/heads/rabbitmq-go/routing
f0cbad7b2e8c0df53c2140173e630ce8a2e92412 commit refs/heads/rabbitmq-go/workqueues
05a28b92601a8cf16fd5d0b3ab56f744b9e9b1da commit refs/heads/rabbitmq-go/helloworld
c895733b9add631a6d250335a62c72a7ac5a8c90 commit refs/heads/mqtt-go/example-2
ae5d763e1f0dfe007d07295190a834491a3dcf4a commit refs/heads/mqtt-go/example-1
1c865ea9b9f5e71ab4580f5f77ec1a5e1ee262a2 commit refs/heads/redis/example-2
0e96428d9af597119e2032ea5428af2ec3db4d44 commit refs/heads/db/organising-database-access-3
cd44d0503542966401f620e5da5b9d47f586869a commit refs/heads/db/organising-database-access-2
79055e17210974e55946ab4036d3c71ed391a090 commit refs/heads/db/organising-database-access-1
a8dbda410766a3f15c19ae322b3eab70c73a6579 commit refs/heads/jwt/auth
8a01d102645c42637a2da2cbd51c144b888bc68b commit refs/heads/go-kit/profilesvc
9d4d6dbbedd7893b0e7d0c879ac560c46a40e94f commit refs/heads/go-kit/stringsvc-3
f87299eeb9c0451d050e7d0a23b04bfa8dd65010 commit refs/heads/shippy/part-2
ce95e35fc26999ed8c8dd3b932581dac7a0bdff7 commit refs/heads/shippy/part-1
711366e85c4bd66d68ca182b48e0bcac20922a20 commit refs/heads/go-micro/service
af237108061a7e98a26dfbf99ff6e3bbf6716d07 commit refs/heads/go-kit/stringsvc-2
44fc6e50c0ca5717688df5f37cffbe41650813bb commit refs/heads/go-kit/stringsvc-1
10312ce85b78fdf90809741df6ace602498f60d6 commit refs/heads/gin/basic
6b6b06a48b06dd6092b9c1b65a5b57b14f7ff14d commit refs/heads/httprouter/basic_auth
c9757e94e85f04b7140ffa924f3b7aca7c6d16b8 commit refs/heads/httprouter/basic
c971d849bb0ad54ab3648c4d534d815336d52133 commit refs/heads/jwt/simple
3dbdc6cf419c44e6ab1973daf0ac025cf240ed35 commit refs/heads/graphql/gqlgen
af740a6c60c69465836a1f537d5f544f20fa4300 commit refs/heads/ssh/tunneling
328089181834c9b5132104fc21aade2e966819fa commit refs/heads/ssh/client-connection
f8d99e80fb3302d2694fa09a18fc835f8aeb0e71 commit refs/heads/piplines/tree-digest-parallel
413930b9317e077445e3cc6a549a86191827093e commit refs/heads/channel/share_memory_by_communication
9f548997e67d571be562339f212275845fba503a commit refs/heads/context/multiple-contexts
1d2fd5e58792122f3575c478d670a979152d4752 commit refs/heads/docker/basic-deploy
85f1cc1101175ad7b56ef869c30068074703fb37 commit refs/heads/redis/example-1
30eaa9cdd998e7d28aea5610969e3673204163d7 commit refs/heads/protocol-buffers/tutorrial-1
d3b98bc606c342b23313356c1bbfd61b9e53a480 commit refs/heads/mongo-rest
30a6f3dd41395d59a3f4f1de03f6c01af1d0b0a2 commit refs/heads/docker/https
c09faec8b11aceb7bd9c86eb496908c50b058d4e commit refs/heads/websocket/chat
f08a4eca1c8973f4175dd848a3def3f620e3d6c7 commit refs/heads/appengine/helloworld
1b94f2e0aef8f68f4c4b774f727890beb2ec8bc5 commit refs/heads/jsonapi
7231b9ed3eb0de3cb4cd818d1bb1f86088e361c9 commit refs/heads/go-compiler/parser-inspect
32feb8acde2492e469557ac5dc9e7a92d60ab354 commit refs/heads/go-compiler/parser
b2e7702a2220e41288b80284d8d17c3cef852c57 commit refs/heads/go-compiler/scanner
32f308a26601db472da63f18b81180bb7410ed67 commit refs/heads/multipart-response
36c8ff3c512f4db9a94af56f714978a63d6d4a75 commit refs/heads/multiple-files-upload-server
fe596ae802efb1bed8bb500185dc51a19e443f45 commit refs/heads/context/web-search
8c52bc5ee0bcaa87a1b87a3c6c47c9e261c060e7 commit refs/heads/piplines/tree-digest-bounded
ad86b72c387ed350e7d9f655b65ff7f8ab49ceb9 commit refs/heads/piplines/tree-digest-serial
1892d361dd9f56abae75feb22d8cceb259d76aff commit refs/heads/piplines/explicit-cancellation-2
d26dae380561ac69746d9771f2ef6bd7dc81a0f1 commit refs/heads/piplines/explicit-cancellation-1
611708ffd65348569f2c62ad475e40deaf4e3ada commit refs/heads/piplines/buffered-channel
e672d2ab600c0b59e1f57c43c60ebfd77120fff1 commit refs/heads/piplines/fan_out-fan_in
8ea32d7e2da1f19b2fb1c9eb7a77798cce8ad6b9 commit refs/heads/piplines/squaring-numbers
```

## Credits

All credits are to the original authors mentioned at comments with references at the top of every example.

## License

Under MIT license.
