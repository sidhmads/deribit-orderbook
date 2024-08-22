# Deribit Orderbook Service

This service provides real-time order book data for various instruments from Deribit via WebSocket and streams it into Kafka.

## Steps to Run

1. **Open terminal and navigate to the code directory**
   ```bash
   cd /path/to/code-directory
   ```
2. **Run Docker Compose to start kafka**
   ```bash
   docker-compose up -d
   ```
3. **Install the required Go modules**
   ```bash
   go mod tidy
   ```
4. **Build the project**
   ```bash
   go build -o ./bin/deribit-orderbook ./
   ```
5. **Run the orderbook consumer for a specific instrument**
   ```bash
   ./bin/deribit-orderbook orderbook-consumer btc option

   ```
4. **Open a new terminal in the same working directory and run the orderbook service**
   ```bash
   ./bin/deribit-orderbook orderbook btc option
   ```

### Improvements
User can get orderbook for multiple instruments

1. **Run the orderbook consumer for multiple instruments**
   ```bash
   ./bin/deribit-orderbook orderbook-consumer btc,eth,usdc,usdt,eurr,any option,spot,future,future_combo,option_combo
   ```
2. **Run the orderbook service for multiple instruments**
   ```bash
   ./bin/deribit-orderbook orderbook btc,eth,usdc,usdt,eurr,any option,spot,future,future_combo,option_combo
   ```

## General Considerations
- The service is designed to be easily scalable within a microservice architecture. The use of the Cobra CLI allows for flexible command handling, making the service adaptable to different commands and arguments.

- Since Deribit's instrument API requires both the currency and the instrument type, users can set these parameters via the CLI using Cobra.

- A base WebSocket struct was implemented with customizable OnOpenCallback and OnMessageCallback functions. This design promotes reusability, allowing the WebSocket code to be easily adapted for different use cases such as OHLCV and Order Updates.

- Given the high frequency of insertions, deletions, and searches, a Red-Black Tree data structure was chosen for managing the order book.

- The codebase includes test cases for the custom order book and other critical functions.
