# Concurrency Approach in Library Management System

## Overview
The Library Management System has been enhanced to support concurrent book reservations using Go's concurrency primitives: Goroutines, Channels, and Mutexes. This ensures safe handling of multiple simultaneous reservation requests while preventing race conditions and implementing auto-cancellation for timed-out reservations.

## Key Components

### Channels
- A buffered channel (`reservationChan` in `Library`) queues incoming `ReservationRequest`s.
- Each request includes the book ID, member ID, and a response channel for returning errors synchronously to the caller.
- Buffering (size 100) allows non-blocking sends for high throughput.

### Goroutines
- **Worker Goroutines**: Multiple (configurable, default 4) worker goroutines consume from the reservation channel and process requests concurrently. This enables parallel handling of reservations.
- **Timer Goroutines**: For each successful reservation, a dedicated goroutine starts a 5-second timer. If the book remains reserved (not borrowed) after the timeout, it is automatically unreserved.

### Mutex
- A `sync.Mutex` (`mu` in `Library`) protects all shared state (e.g., `Books` and `Members` maps) during reads and writes.
- Critical sections (e.g., status updates in `processReservation` and `BorrowBook`) acquire the lock to ensure atomicity and prevent race conditions like double reservations.

## Flow for Reservations
1. **Request Submission** (`ReserveBook`): Creates a `ReservationRequest` with a temporary response channel and sends it to `reservationChan`.
2. **Concurrent Processing**:
   - A worker goroutine dequeues the request.
   - Acquires the mutex, checks book availability atomically.
   - If available: Updates status to "Reserved", sets `ReservedBy`, sends `nil` error via response channel, and spawns a timer goroutine.
   - If unavailable: Sends an appropriate error via response channel.
3. **Response**: The caller blocks briefly (`<-resp`) to receive the error synchronously.
4. **Auto-Cancellation**: The timer goroutine waits 5 seconds, then (under lock) checks and reverts the status to "Available" if still reserved by the original member.
5. **Borrow Integration**: `BorrowBook` (under lock) allows borrowing if available or reserved by the caller, clearing `ReservedBy` and setting status to "Borrowed".

## Handling Concurrent Requests
- **Simulation**: Menu option 8 launches 5 goroutines attempting to reserve the same book simultaneously (members 1-5).
- Workers process in parallel, but the mutex serializes updates—ensuring only one succeeds, others get "book not available".
- Outputs appear asynchronously in the console, demonstrating concurrency.
- Timers run independently, showcasing multiple goroutines per request.

## Error Handling
- **Double Reservations**: Prevented by atomic check-and-set under mutex.
- **Timeouts**: Graceful auto-unreserve with logging.
- **Invalid States**: Errors for non-existent books/members, borrowed/reserved-by-other books.
- Channel-based responses ensure callers always receive feedback, even under load.

## Code Structure
- **services/library_service.go**: Core logic, interface, and structs (e.g., `Library`, `ReservationRequest`).
- **concurrency/reservation_worker.go**: Worker setup (`StartWorkers`) and processing (`processReservation`).
- **controllers/library_controller.go**: Updated console with reservation (7) and simulation (8) options.
- **models/**: Extended `Book` with `ReservedBy`.
- **main.go**: Initializes extra members/book and starts workers.

This design scales efficiently for concurrent access while maintaining simplicity and safety.