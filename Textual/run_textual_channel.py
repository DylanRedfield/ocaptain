from rasa_core.agent import Agent
from rasa_core.interpreter import RasaNLUInterpreter
from textual_channel import TextualInput
from rasa_core.channels.console import CmdlineInput
from rasa_core.channels.channel import RestInput
from rasa_core import utils

print("Start interpreter")
interpreter = RasaNLUInterpreter("models/current/nlu")

print("start load")
endpoint_conf = utils.read_endpoint_config("endpoints.yml", endpoint_type="action_endpoint")
agent = Agent.load("models/current/dialogue", interpreter = interpreter, action_endpoint = endpoint_conf)

print("Make input channel")
input_channel = TextualInput(agent)

print("start handle")
agent.handle_channels([input_channel], 5005, serve_forever = True)
