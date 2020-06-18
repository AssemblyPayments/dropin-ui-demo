# Drop-In UI Demo

A very small server written in Go that demonstrates the Assembly Payments Drop-In UI for secure credit card acquisition.

For more detailed information about Assembly Payments APIs, development guides and more, please see <https://developer.assemblypayments.com/docs>.

## License

Copyright 2020 Assembly Payments

The software in this repository is available as open source under the terms of the [Apache License](https://github.com/AssemblyPayments/dropin-ui-demo/blob/master/LICENSE).

## Installation

1. [Install Go](https://golang.org/doc/install).
2. Clone this repo. CD into the `dropin-ui-demo` directory and enter `go build server.go`.
3. If you are new to Go, you could check out [this](https://golang.org/doc/code.html) and [this](https://tour.golang.org/) and [this](https://golang.org/doc/effective_go.html).

Please note that knowledge of Go is not required to run this demo.

## Operation

1. You will need credentials to authorize the server to interact with the Prelive environment. If you haven't already, sign up [here](https://dashboard.prelive.assemblypayments.com/#/sign-up/prelive).
2. Add a file named `secret` to the repo directory. Write your prelive [compound key](https://developer.assemblypayments.com/docs/keys#compound) (`[username]:[API token]`) to this file. Get these values from the account page of the prelive [dashboard](https://dashboard.prelive.assemblypayments.com/#/accounts). **Don't use your production credentials with this demo.**
3. CD into the `dropin-ui-demo` directory and enter `./server`.
4. Open a web browser to <http://localhost:8081/dropin.html>.
5. Use the card number `4111 1111 1111 1111`.

## To get going with Docker & docker-compose

1. docker-compose build
2. Add creds (as per 2 above) to `secrets`
3. docker-compose up
4. Open a web browser to <http://localhost:8081/dropin.html>.
