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
61b7079e186a38a57302e0d9a0c7a2827e111f5a commit refs/heads/master
a1fb967d2e74e3f7df8837b4bf4ee1c0c58ad8b4 commit refs/heads/cockroachdb
881bb04ff00348e26f2a9cf3eae7722a427b7a34 commit refs/heads/makefile-tutorial
008fa359342b6e01ecc46c6358ad6ebcb4297007 commit refs/heads/http-proto
39fd42314993eac69dd8477fd34171a72803c715 commit refs/heads/helloworld
c3ce73205f0482b05e26d9511bdd05fa445d0298 commit refs/heads/jsonapi
2300e611f38b7a4c8252633f768e4b6e1ed94de0 commit refs/heads/dsa
4583db32ed200038d67f9e3d11409d41877d6d89 commit refs/heads/grpc/consul-sd-integration
d9577349c791102cc5a95f3da0010e0a703cbe6d commit refs/heads/pubsub
e98564d54a49ab7a80ebd5b33a05ff78616f61f0 commit refs/heads/http/multipart
c8c8be47c424bf92a25692826426815f3ef060e4 commit refs/heads/grpc/name_resolving
b0aa0a78ff054d96c4d2729b8a528156e910881a commit refs/heads/nats
11c6ff06edbb4ae036596ba3e413d27215acd81d commit refs/heads/zap-examples
d5f60f4b00c2c24d5e493bfa2b967b4041ee4575 commit refs/heads/go-ldflags
64c30111ace9fbe94b942b8e36f2cb724ecbc414 commit refs/heads/go-plugin
c94c1ab510d87ca6224c6d5a17d82d235c383a8f commit refs/heads/go-build-tags
f23bf6868a6122ce350bb46b3501530ff45d5a37 commit refs/heads/zookeeper-tutorial
790e6466fed0ff15b26c2ae340f4d1546f224afd commit refs/heads/opentracing-tutorial-1
51b4c419778bf3e102e39b019200cb1bd6bef891 commit refs/heads/http-tracing
9b48ead4d33918383c668dffa960af3f1894b1cb commit refs/heads/kafka-go
3e5185138b40228da5278e75d09069d8e71f16cb commit refs/heads/sort_transform
6762034f0b8a2576557edb10b43cbd473f10f9fe commit refs/heads/errgroup
751dc0b6554483279d98d60f6be0ccd4aa00575c commit refs/heads/rest-unit-testing
35b3c674fb17fa0530d1c1b6beaa9b24842b69a0 commit refs/heads/nats-tutorial
6c19a6548892050df473f5e137dee221d6e320a2 commit refs/heads/opendbservice
0f13b96a521a616686cac928876a3fe6734fbf8e commit refs/heads/opendbcli
90ff4f6e9a0d1bc08f45dd1c7daaf9edf3d99a81 commit refs/heads/routeguide
42f305a92ae3e2525f45638d89447fd38376d834 commit refs/heads/addsvc
0672152c4657fb285bf9ce37137a37a4b04e7576 commit refs/heads/tree-digesting-review
0672152c4657fb285bf9ce37137a37a4b04e7576 commit refs/heads/pipeline/tree-digesting-review
9aa8063dbdde94b855df7346b4afdbd0de72a858 commit refs/heads/review
9aa8063dbdde94b855df7346b4afdbd0de72a858 commit refs/heads/pipelines/review
f0eaf64accd35804b2b110b8fb3ef972273e09b6 commit refs/heads/basic
be4bf21d8ad01ebcdf3b6bc66cf53afb88f2e4b7 commit refs/heads/priority_queue_heap
8d861c493cbb16e0084264b32d54f6802867cc49 commit refs/heads/intheap
3db843b0fa7a90513db3d488e06ffd47d32064db commit refs/heads/postgres
15fa2c71218603ab2535bddb44d27bac2b99de20 commit refs/heads/mongodb-go/tutorial-1
3166d0b92760aae71c8a6b259013843a784c6ea8 commit refs/heads/simple-producer-consumer
a63b1352887a950e28bc7b7d5b08ec83bc1b8dd5 commit refs/heads/rpc
93c99374d1ee97b5d4b57035bb7c4f4832894b3f commit refs/heads/topics
a8c147d18bb7bd2cdab93f8c3c8c367e02310e34 commit refs/heads/routing
f0cbad7b2e8c0df53c2140173e630ce8a2e92412 commit refs/heads/workqueues
fe877a8e4d645df1197399ca89e1897edae5c1b1 commit refs/heads/workqueue
c895733b9add631a6d250335a62c72a7ac5a8c90 commit refs/heads/example-2
ae5d763e1f0dfe007d07295190a834491a3dcf4a commit refs/heads/example-1
0e96428d9af597119e2032ea5428af2ec3db4d44 commit refs/heads/organising-database-access-3
cd44d0503542966401f620e5da5b9d47f586869a commit refs/heads/organising-database-access-2
79055e17210974e55946ab4036d3c71ed391a090 commit refs/heads/organising-database-access-1
a8dbda410766a3f15c19ae322b3eab70c73a6579 commit refs/heads/auth
8a01d102645c42637a2da2cbd51c144b888bc68b commit refs/heads/profilesvc
9d4d6dbbedd7893b0e7d0c879ac560c46a40e94f commit refs/heads/stringsvc-3
f87299eeb9c0451d050e7d0a23b04bfa8dd65010 commit refs/heads/part-2
ce95e35fc26999ed8c8dd3b932581dac7a0bdff7 commit refs/heads/part-1
711366e85c4bd66d68ca182b48e0bcac20922a20 commit refs/heads/service
af237108061a7e98a26dfbf99ff6e3bbf6716d07 commit refs/heads/stringsvc-2
44fc6e50c0ca5717688df5f37cffbe41650813bb commit refs/heads/stringsvc-1
6b6b06a48b06dd6092b9c1b65a5b57b14f7ff14d commit refs/heads/basic_auth
c971d849bb0ad54ab3648c4d534d815336d52133 commit refs/heads/simple
3dbdc6cf419c44e6ab1973daf0ac025cf240ed35 commit refs/heads/gqlgen
af740a6c60c69465836a1f537d5f544f20fa4300 commit refs/heads/tunneling
328089181834c9b5132104fc21aade2e966819fa commit refs/heads/client-connection
f8d99e80fb3302d2694fa09a18fc835f8aeb0e71 commit refs/heads/tree-digest-parallel
413930b9317e077445e3cc6a549a86191827093e commit refs/heads/share_memory_by_communication
9f548997e67d571be562339f212275845fba503a commit refs/heads/multiple-contexts
1d2fd5e58792122f3575c478d670a979152d4752 commit refs/heads/basic-deploy
30eaa9cdd998e7d28aea5610969e3673204163d7 commit refs/heads/tutorrial-1
d3b98bc606c342b23313356c1bbfd61b9e53a480 commit refs/heads/mongo-rest
30a6f3dd41395d59a3f4f1de03f6c01af1d0b0a2 commit refs/heads/https
c09faec8b11aceb7bd9c86eb496908c50b058d4e commit refs/heads/chat
7231b9ed3eb0de3cb4cd818d1bb1f86088e361c9 commit refs/heads/parser-inspect
32feb8acde2492e469557ac5dc9e7a92d60ab354 commit refs/heads/parser
b2e7702a2220e41288b80284d8d17c3cef852c57 commit refs/heads/scanner
32f308a26601db472da63f18b81180bb7410ed67 commit refs/heads/multipart-response
36c8ff3c512f4db9a94af56f714978a63d6d4a75 commit refs/heads/multiple-files-upload-server
fe596ae802efb1bed8bb500185dc51a19e443f45 commit refs/heads/web-search
fe596ae802efb1bed8bb500185dc51a19e443f45 commit refs/heads/context/web-search
8c52bc5ee0bcaa87a1b87a3c6c47c9e261c060e7 commit refs/heads/tree-digest-bounded
ad86b72c387ed350e7d9f655b65ff7f8ab49ceb9 commit refs/heads/tree-digest-serial
1892d361dd9f56abae75feb22d8cceb259d76aff commit refs/heads/explicit-cancellation-2
d26dae380561ac69746d9771f2ef6bd7dc81a0f1 commit refs/heads/explicit-cancellation-1
611708ffd65348569f2c62ad475e40deaf4e3ada commit refs/heads/buffered-channel
e672d2ab600c0b59e1f57c43c60ebfd77120fff1 commit refs/heads/fan_out-fan_in
8ea32d7e2da1f19b2fb1c9eb7a77798cce8ad6b9 commit refs/heads/squaring-numbers
```

## Credits

All credits are to the original authors mentioned at comments with references at the top of every example.

## License

Under MIT license.
