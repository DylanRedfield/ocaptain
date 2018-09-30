from rasa_core.interpreter import RasaNLUInterpreter
from rasa_core.channels import HttpInputChannel
from rasa_core.channels.twilio import TwilioInput
from rasa_core.agent import Agent
from rasa_core.channels.console import ConsoleInputChannel
import json

def run():
    interpreter = RasaNLUInterpreter("models/current/nlu")

    while True:
        message = raw_input("Enter message: ")

        result = interpreter.parse(unicode(message))

        print(json.dumps(result, indent=2))



if __name__ == '__main__':
    run()
