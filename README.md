# shopping-cart-service-go
> Sample shopping cart CRUD service with Go and Kubernetes.

## Requirements

- [GNU Make](https://www.gnu.org/software/make/)
- [Go](https://golang.org/)
- [Docker](https://hub.docker.com/)
- [DBMate](https://github.com/amacneil/dbmate)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

## Usage

To install project dependencies, run the following command:
```bash
make install
```

To generate mocks, run the following command:
```bash
make generate_mocks
```

To run all tests, run the following command:
```bash
make test
```

To make sure the database is running, run the following command:
```bash
make run_local_db
```

To start the app and database locally, run the following command:
```bash
make start
```

To debug the local database, run the following command:
```bash
make debug_local_db
```

To run migrations, run the following command:
```bash
make migrate
```

To generate seed data, run the following command:
```bash
make generate_seed_data
```

To seed the database, run the following command:
```bash
make seed_db
```

To build the docker image, run the following command:
```bash
make build_image
```

To run the docker image, run the following command:
```bash
make debug_image
```

To push the docker image to dockerhub, run the following command:
```bash
make push_image
```