TODO:
- tests (testify, integration tests) [https://github.com/stretchr/testify , https://github.com/bquenin/go-modern-rest-api-tutorial/blob/main/api/authors/authors_test.go]
- profiling [https://making.pusher.com/go-tool-trace/ , https://habr.com/ru/company/badoo/blog/301990/]
- Поиск по роут паттерну (regex :( )
- settings like in pydantic. try read .env –> ignore fail, read from host env, search for config keys (search, or create this feature for GO community) [https://github.com/kelseyhightower/envconfig]
- ci/cd
- auther 2.0
  - remove db, fill auth info on startup from auth-api. 
  - listen for changes in auth info from auth-api. 
  - decode jwt with public key (hz, interesting)