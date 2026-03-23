# OCaptain

OCaptain is a conversational AI platform for businesses that handles customer interactions across multiple messaging channels. It is designed primarily for restaurants, enabling automated booking, order taking, and FAQ responses via SMS, WhatsApp, and Facebook Messenger.

## Architecture

The project consists of two main components:

### PizzaBot (Go)
A Go HTTP server that acts as the messaging gateway and business logic layer. It:
- Receives inbound messages from Twilio SMS, Twilio WhatsApp, and Facebook Messenger
- Looks up the relevant business and recipient records from Firebase Firestore
- Forwards requests to the Rasa NLU server for intent handling
- Sends outbound responses back to customers via the appropriate messaging platform
- Persists all messages and conversation state in Firestore

### Textual (Python / Rasa)
A Rasa NLU/Core chatbot that understands and responds to natural language. It:
- Handles restaurant-specific intents: making reservations, placing orders, greeting, and common FAQ questions
- Manages multi-turn conversations with slot filling (name, party size, date/time)
- Queries reservation availability and saves confirmed reservations
- Escalates to a human employee when needed

## Features

- **Multi-channel messaging**: Twilio SMS, Twilio WhatsApp, Facebook Messenger
- **Reservation management**: Natural language booking with availability checking and alternative time suggestions
- **Order management**: Delivery and pickup order intake
- **Restaurant FAQs**: Automated answers about hours, dietary options, parking, dress code, and more
- **Business caching**: In-memory caching of business and recipient data to reduce Firestore reads
- **Multi-environment support**: `dev_local`, `dev` (online), and `prod` configurations

## Project Structure

```
ocaptain/
├── PizzaBot/               # Go backend service
│   ├── main.go             # HTTP server and request routing
│   ├── pizzabot.go         # Core bot logic and action handling
│   ├── data_types.go       # Data models (Business, Recipient, Message, etc.)
│   ├── constants.go        # Constants for actions, slots, collections
│   ├── opentable.go        # Reservation platform integration
│   ├── twilio.go           # Twilio SMS/WhatsApp client
│   ├── facebook_messenger.go # Facebook Messenger client
│   ├── google_client.go    # Google client integration
│   ├── swift_sms.go        # Swift SMS client
│   ├── input_handler.go    # Incoming message processing
│   ├── message_client.go   # Outgoing message dispatch
│   └── util.go             # Utility functions
├── textual/                # Rasa NLU chatbot
│   ├── domain.yml          # Intents, actions, slots, and responses
│   ├── config.yml          # Rasa pipeline configuration
│   ├── data/               # Training data (NLU examples and stories)
│   ├── textual_channel.py  # Custom Rasa input/output channel
│   ├── run.py              # Rasa agent startup script
│   └── endpoints.yml       # Action server endpoint configuration
└── env_values.json         # Environment configuration
```

## Prerequisites

**PizzaBot (Go service)**
- Go 1.13+
- Firebase service account credentials (`dev-firebase-config.json` or `prod-firebase-config.json`)
- Twilio account (SID and auth token)
- Facebook Messenger page access token (optional)

**Textual (Rasa service)**
- Python 3.7+
- Rasa 2.x
- Firebase Admin SDK credentials (`firebase-config.json`)

## Configuration

Copy and update `env_values.json` with your environment values:

```json
{
  "name": "dev",
  "twilio_account_sid": "<YOUR_TWILIO_SID>",
  "twilio_auth_token": "<YOUR_TWILIO_AUTH_TOKEN>",
  "swift_account_key": "<YOUR_SWIFT_KEY>",
  "firestore_config_file_name": "dev_firestore_config.json",
  "pizza_url": "https://<YOUR_DOMAIN>/PizzaBot/outsideSmsInput",
  "pizza_port": "80",
  "messenger_verify": "<YOUR_MESSENGER_VERIFY_TOKEN>"
}
```

Place your Firebase service account JSON file in `PizzaBot/` (named `dev-firebase-config.json` for development or `prod-firebase-config.json` for production).

## Running the Services

### PizzaBot

```bash
cd PizzaBot
go build -o ocaptain .

# Development (local)
./ocaptain dev_local

# Development (online, uses HTTPS)
./ocaptain dev

# Production
./ocaptain prod
```

The server listens on port `:443` (HTTPS) in `dev` and `prod` modes, and port `:80` in `dev_local` mode. TLS certificates are automatically managed via Let's Encrypt.

### Textual (Rasa)

Install dependencies:

```bash
cd textual
pip install rasa firebase-admin twilio
```

Train the model:

```bash
rasa train
```

Start the action server and then the Rasa agent:

```bash
rasa run actions &
python run.py
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/PizzaBot/outsideSmsInput` | Inbound SMS from Twilio |
| POST | `/PizzaBot/outsideTwilioWhatsappInput` | Inbound WhatsApp from Twilio |
| POST | `/PizzaBot/outsideFacebookInput` | Inbound message from Facebook Messenger |
| POST | `/PizzaBot/outsideGoogleInput` | Google webhook verification |
| POST | `/PizzaBot/businessInput` | Message sent from the business dashboard |
| POST | `/PizzaBot/sendAndSave` | Send a message and persist it to Firestore |
| POST | `/ocaptain` | Rasa action server hook |
| POST | `/ocaptain/sendAndSave` | Send and save from the OCaptain action server |

## Data Model

- **Business** – A registered restaurant with messaging credentials, hours, and employee list
- **Recipient** – A customer contact associated with a business
- **Message** – An individual chat message, stored in Firestore under the business/recipient hierarchy
- **Order** – A delivery or pickup order
- **Reservation** – A table reservation with name, party size, and scheduled time
