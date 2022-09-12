TODO:
- profiling [https://making.pusher.com/go-tool-trace/ , https://habr.com/ru/company/badoo/blog/301990/]
- auther 2.0
  - decode jwt with public key (interesting). jwt rsa decoding/encoding will be nice if we want to divide services by repos. + any other service can decode tokens on its side
- tests (testify, integration tests) [https://github.com/stretchr/testify , https://github.com/bquenin/go-modern-rest-api-tutorial/blob/main/api/authors/authors_test.go] !!!
- settings like in pydantic. try read .env –> ignore fail, read from host env, search for config keys (search, or create this feature for GO community) [https://github.com/kelseyhightower/envconfig]