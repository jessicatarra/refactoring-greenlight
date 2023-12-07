## Refactoring Greenlight application 

This project aims to refactor the Greenlight application by implementing clean architecture and SOLID principles. The goal is to improve the codebase's maintainability, scalability, and testability while increasing test coverage.

To understand the initial state of the project, please refer to the [previous README file](https://github.com/jessicatarra/greenlight/blob/main/README.md).

### API Docs from V1 application version
https://greenlight.tarralva.com/swagger/index.html

### TODO
- [ ] Refactor initial routes implementation into separate service handlers
  - [ ] Add an authentication module
    - [x] Implement the create user feature
    - [x] Implement the activate user feature
    - [ ] Implement the create authentication token feature
  - [ ] Add a movies module
  - [ ] Add a healthcheck module
- [ ] Refactor multiple functionalities into internal packages
  - [x] Add support for background tasks
  - [x] Add `log/slog` package
  - [x] Add `mailer` package
  - [x] Add utils package
    - [x] Add validator `v2` package
    - [x] Add validator helpers utilities
    - [x] Add general helpers utilities
  - [x] Add response package
  - [x] Add request package
  - [x] Add errors package
  - [x] Add config package
- [ ] Refactor initial API server implementation
- [ ] Relocate initial middleware implementation into the different modules based on the concerns of each one
- [ ] Implement an actual modular monolith

### References

- https://github.com/golang-standards/project-layout
- https://autostrada.dev/
- https://github.com/qiangxue/go-rest-api/tree/
- https://github.com/powerman/go-service-example/
- https://github.com/powerman/go-monolith-example
- https://github.com/amitshekhariitbhu/go-backend-clean-architecture
- https://github.com/evrone/go-clean-template/
- https://github.com/booscaaa/clean-go/
- https://github.com/bxcodec/go-clean-arch/
