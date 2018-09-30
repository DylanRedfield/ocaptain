from rasa_nlu.model import Interpreter
import json

interpreter = Interpreter.load("./models/current/nlu")
message = u"I'd like to place an order for delivery?"
result = interpreter.parse(message)
print(json.dumps(result, indent=2))
