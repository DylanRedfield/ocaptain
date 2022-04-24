from rasa.core.agent import Agent
from rasa.core.http_interpreter import RasaNLUHttpInterpreter
from textual_channel import TextualInput
#from rasa_core.channels.console import CmdlineInput
#from rasa_core.channels.channel import RestInput
from rasa.utils.endpoints import read_endpoint_config

print("-- Start interpreter --")
#interpreter = RasaNLUInterpreter("models/current/nlu")

print("-- Start Load --")
endpoint_conf = read_endpoint_config("endpoints.yml", endpoint_type="action_endpoint")
#agent = Agent.load("models/dialogue", interpreter = interpreter, action_endpoint = endpoint_conf)
agent = Agent.load("models/",  action_endpoint = endpoint_conf)

print("-- Make input channel --")
input_channel = TextualInput(agent)

print("-- Running --")
agent.handle_channels([input_channel], 5006, serve_forever = True)
